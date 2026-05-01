package modules

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/domain/usecases"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router/logger"
	"github.com/gorilla/mux"
)

func NewUserModule(userUseCases usecases.UserUseCase) router.Module {
	return userModule{
		userUseCases: userUseCases,
		name:         "User",
		path:         "/user",
	}
}

type userModule struct {
	userUseCases usecases.UserUseCase
	name         string
	path         string
}

func (m userModule) Name() string {
	return m.name
}

func (m userModule) Path() string {
	return m.path
}

func (m userModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{
			Path:        "/update-information",
			Description: "Update user information",
			Handler:     m.updateUserInformation,
			HttpMethods: []string{http.MethodPut},
			Public:      false,
		},
	}
}

func (m userModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}

func (m userModule) updateUserInformation(w http.ResponseWriter, r *http.Request) {
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

	var userInformation entities.UserInformation
	err = json.Unmarshal(body, &userInformation)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, derr.BadRequestError)
		return
	}

	userInformation.ID = user.ID
	response, err := m.userUseCases.UpdateUserInformation(ctx, userInformation)
	if err != nil {
		slog.ErrorContext(ctx, "failed to update user information", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}
