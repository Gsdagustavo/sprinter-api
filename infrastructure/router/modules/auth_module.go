package modules

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/domain/logger"
	"github.com/VitorFranciscoDev/sprinter-api/domain/usecases"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/router"
	"github.com/gorilla/mux"
)

type authModule struct {
	authUseCases usecases.AuthenticationUseCase
	name         string
	path         string
}

func NewAuthModule(authUseCases usecases.AuthenticationUseCase) router.Module {
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
