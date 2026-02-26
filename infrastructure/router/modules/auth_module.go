package modules

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/VitorFranciscoDev/sprinter-api/domain"
	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/domain/entities/derr"
	"github.com/VitorFranciscoDev/sprinter-api/domain/logger"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/router"
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

func (a authModule) Name() string {
	return a.name
}

func (a authModule) Path() string {
	return a.path
}

func (a authModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "/login",
			Description: "Attempt user login",
			Handler:     a.login,
			HttpMethods: []string{http.MethodPost},
		},
		{
			Path:        "/register",
			Description: "Attempt user register",
			Handler:     a.register,
			HttpMethods: []string{http.MethodPost},
		},
	}

	for _, d := range defs {
		r.HandleFunc(a.path+d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, r
}

func (a authModule) sessionMiddleware(next http.Handler) http.Handler {
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

			valid, err := a.authUseCases.CheckCredentials(ctx, credentials)
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

			user, err = a.authUseCases.GetUserByEmail(ctx, credentials.Email)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get user by document", logger.Err(err))
				_ = router.Write(w, unauthorizedBytes)
				return
			}

			if user == nil {
				slog.ErrorContext(ctx, "user not found")
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

			user, err = a.authUseCases.Get(ctx, token)
			if err != nil {
				slog.ErrorContext(ctx, "failed to get user from token", logger.Err(err))
				return
			}

			if user == nil {
				slog.ErrorContext(ctx, "user not found")
				_ = router.Write(w, unauthorizedBytes)
				return
			}
		}

		next.ServeHTTP(w, r.WithContext(router.WithUser(ctx, user)))
	})
}

func (a authModule) login(w http.ResponseWriter, r *http.Request) {
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

	response, err := a.authUseCases.AttemptLogin(ctx, credentials)
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

func (a authModule) register(w http.ResponseWriter, r *http.Request) {
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

	response, err := a.authUseCases.AttemptRegister(ctx, credentials)
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

func (a authModule) me(w http.ResponseWriter, r *http.Request) {
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

	response, err := a.authUseCases.AttemptRegister(ctx, credentials)
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
