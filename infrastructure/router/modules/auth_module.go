package modules

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/gorilla/mux"
)

type authModule struct {
	authUseCases domain.AuthenticationUseCase
	name         string
	path         string
}

func NewAuthModule(authUseCases domain.AuthenticationUseCase) router.Module {
	return authModule{
		authUseCases: authUseCases,
		name:         "Authentication",
		path:         "/auth",
	}
}

type AuthenticationResponse struct {
	Token string `json:"token"`
}

func (m authModule) Name() string {
	return m.name
}

func (m authModule) Path() string {
	return m.path
}

func (m authModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{
			Path:        "/login",
			Description: "Attempt user login",
			Handler:     m.login,
			HttpMethods: []string{http.MethodPost},
			Public: true,
		},
		{
			Path:        "/register",
			Description: "Attempt user register",
			Handler:     m.register,
			HttpMethods: []string{http.MethodPost},
			Public: true,
		},
		{
			Path:        "/completeRegistration",
			Description: "Attempt complete user registration",
			Handler:     m.completeRegistration,
			HttpMethods: []string{http.MethodPost},
		},
	}
}

func (m authModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{m.sessionMiddleware}
}

func (m authModule) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		authHeader := r.Header.Get("Authorization")

		var user *entities.User
		var err error

		// Test basic auth
		email, password, ok := r.BasicAuth()
		if ok {
			credentials := entities.UserCredentials{
				Email:    email,
				Password: password,
			}

			valid, err := m.authUseCases.CheckCredentials(ctx, credentials)
			if err != nil {
				slog.ErrorContext(ctx, "failed to check user credentials", "cause", err)
				router.HandleError(w, derr.UnauthorizedError)
				return
			}

			if !valid {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			user, err = m.authUseCases.GetUserByEmail(ctx, credentials.Email)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get user by email", "cause", err)
				router.HandleError(w, derr.UnauthorizedError)
				return
			}

			if user == nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}
		} else {
			var token string
			if authHeader != "" {
				token = strings.ReplaceAll(authHeader, "Bearer ", "")
			}

			if token == "" {
				slog.ErrorContext(ctx, "no token found in the request")
				router.HandleError(w, derr.UnauthorizedError)
				return
			}

			user, err = m.authUseCases.GetUserByToken(ctx, token)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get user from token", "cause", err)
				router.HandleError(w, derr.UnauthorizedError)
				return
			}

			if user == nil {
				http.Error(w, "user not found", http.StatusUnauthorized)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(router.WithUser(ctx, user)))
	})
}

func (m authModule) login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	var credentials entities.UserCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	token, err := m.authUseCases.AttemptLogin(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt login", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	response := AuthenticationResponse{Token: token}
	w.Header().Set("Authorization", "Bearer "+response.Token)
	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}

func (m authModule) register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	var credentials entities.UserCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	token, err := m.authUseCases.AttemptRegister(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt register", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	response := AuthenticationResponse{Token: token}
	w.Header().Set("Authorization", "Bearer "+token)
	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}

func (m authModule) me(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	var credentials entities.UserCredentials
	err = json.Unmarshal(body, &credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	token, err := m.authUseCases.AttemptRegister(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt register", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	response := AuthenticationResponse{Token: token}
	w.Header().Set("Authorization", "Bearer "+token)
	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}

func (m authModule) completeRegistration(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	var information entities.UserInformation
	err = json.Unmarshal(body, &information)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = m.authUseCases.AttemptCompleteRegistration(ctx, information)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt register", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, router.NewSuccessfulResponse())
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}
