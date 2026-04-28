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

type userModule struct {
	userUseCases domain.UserUseCase
	name         string
	path         string
}

// Name implements [router.Module].
func (u userModule) Name() string {
	return u.name
}

// Path implements [router.Module].
func (u userModule) Path() string {
	return u.path
}

// Setup implements [router.Module].
func (u userModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	defs := []router.RouteDefinition{
		{
			Path:        "/{user_id}",
			Description: "Edit user profile",
			Handler:     u.editUser,
			HttpMethods: []string{http.MethodPatch},
		},
	}
	for _, d := range defs {
		r.HandleFunc(u.path+d.Path, d.Handler).Methods(d.HttpMethods...)
	}
	return defs, r
}

func NewUserModule(userUsecase domain.UserUseCase) router.Module {
	return userModule{
		userUseCases: userUsecase,
		name:         "User",
		path:         "/user",
	}
}
func (u userModule) editUser(w http.ResponseWriter, r *http.Request) {
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

	var editedUser entities.EditUserProfileDTO
	err = json.Unmarshal(body, &editedUser)
	if err != nil {
		slog.ErrorContext(ctx, "failed to read body", logger.Err(err))
		router.HandleError(w, err)
		return
	}
	editedUser.ID = user.ID
	response, err := u.userUseCases.EditUserProfile(ctx, editedUser)
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
