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
	"github.com/Gsdagustavo/sprinter-api/domain/logger"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
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

func (m authModule) Name() string {
	return m.name
}

func (m authModule) Path() string {
	return m.path
}

func (m authModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "/login",
			Description: "Attempt user login",
			Handler:     m.login,
			HttpMethods: []string{http.MethodPost},
		},
		{
			Path:        "/register",
			Description: "Attempt user register",
			Handler:     m.register,
			HttpMethods: []string{http.MethodPost},
		},
		{
			Path:        "/completeRegistration",
			Description: "Attempt complete user registration",
			Handler:     m.completeRegistration,
			HttpMethods: []string{http.MethodPost},
		},
	}

	for _, d := range defs {
		r.HandleFunc(m.path+d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, r
}

func (m authModule) sessionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authHeader := r.Header.Get("Authorization")

		var user *entities.User

		unauthorizedBytes, err := json.Marshal(derr.UnauthorizedError)
		if err != nil {
			slog.ErrorContext(ctx, "failed to marshal error response", logger.Err(err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Test basic auth
		email, password, ok := r.BasicAuth()
		if ok {
			credentials := entities.UserCredentials{
				Email:    email,
				Password: password,
			}

			valid, err := m.authUseCases.CheckCredentials(ctx, credentials)
			if err != nil {
				slog.ErrorContext(ctx, "failed to check credentials", logger.Err(err))
				_ = router.Write(w, unauthorizedBytes)
				return
			}

			if !valid {
				slog.ErrorContext(ctx, "invalid credentials")
				_ = router.Write(w, unauthorizedBytes)
				return
			}

			user, err = m.authUseCases.GetUserByEmail(ctx, credentials.Email)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get user by document", logger.Err(err))
				_ = router.Write(w, unauthorizedBytes)
				return
			}

			if user == nil {
				slog.ErrorContext(ctx, "user not found in context")
				_ = router.Write(w, unauthorizedBytes)
				return
			}
		} else {
			var token string
			var err error

			if authHeader != "" {
				token = strings.ReplaceAll(authHeader, "Bearer ", "")
			}

			if token == "" {
				slog.ErrorContext(ctx, "no token found in the request")
				_ = router.Write(w, unauthorizedBytes)
				return
			}

			user, err = m.authUseCases.GetUserByToken(ctx, token)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get user from token", logger.Err(err))
				return
			}

			if user == nil {
				slog.ErrorContext(ctx, "user not found in context")
				_ = router.Write(w, unauthorizedBytes)
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

	response, err := m.authUseCases.AttemptLogin(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt login", logger.Err(err))
		router.HandleError(w, err)
		return
	}

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

	response, err := m.authUseCases.AttemptRegister(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt register", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+response.Token)
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

	response, err := m.authUseCases.AttemptRegister(ctx, credentials)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt register", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	w.Header().Set("Authorization", "Bearer "+response.Token)
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

	var information entities.AccountInformation
	err = json.Unmarshal(body, &information)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	response, err := m.authUseCases.AttemptCompleteRegistration(ctx, information)
	if err != nil {
		slog.ErrorContext(ctx, "failed to attempt register", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}
