package modules

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/gorilla/mux"
)

func NewUserModule(userUsecase domain.UserUseCase) router.Module {
	return userModule{
		userUseCases: userUsecase,
		name:         "User",
		path:         "/user",
	}
}

type userModule struct {
	userUseCases domain.UserUseCase
	name         string
	path         string
}

func (m userModule) Name() string {
	return m.name
}

func (m userModule) Path() string {
	return m.path
}

func (m userModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "",
			Description: "Edit user profile",
			Handler:     m.editUser,
			HttpMethods: []string{http.MethodPut},
		},
	}

	for _, d := range defs {
		r.HandleFunc(m.path+d.Path, d.Handler).Methods(d.HttpMethods...)
	}

	return defs, r
}

func (m userModule) editUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := router.GetUser(r)
	if err != nil {
		slog.ErrorContext(ctx, "failed to get user from context", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	var userInformation entities.AccountInformation
	err = json.Unmarshal(body, &userInformation)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	userInformation.ID = user.ID
	response, err := m.userUseCases.UpdateUserProfile(ctx, userInformation)
	if err != nil {
		slog.ErrorContext(ctx, "failed to edit user", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}
