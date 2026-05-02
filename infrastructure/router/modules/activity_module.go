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

func NewActivityModule(activityUseCases usecases.ActivityUseCase) router.Module {
	return activityModule{
		activityUseCases: activityUseCases,
		name:             "Activity",
		path:             "/activity",
	}
}

type activityModule struct {
	activityUseCases usecases.ActivityUseCase
	name             string
	path             string
}

func (m activityModule) Name() string {
	return m.name
}

func (m activityModule) Path() string {
	return m.path
}

func (m activityModule) Routes() []router.RouteDefinition {
	return []router.RouteDefinition{
		{
			Path:        "",
			Description: "Start new activity",
			Handler:     m.startActivity,
			HttpMethods: []string{http.MethodPost},
			Public:      false,
		},
	}
}

func (m activityModule) Middlewares() []mux.MiddlewareFunc {
	return []mux.MiddlewareFunc{}
}
func (m activityModule) startActivity(w http.ResponseWriter, r *http.Request) {
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

	var activity entities.Activity
	err = json.Unmarshal(body, &activity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to unmarshal request body", logger.Err(err))
		router.HandleError(w, derr.BadRequestError)
		return
	}

	activity.UserID = user.ID
	response, err := m.activityUseCases.StartActivity(ctx, activity)
	if err != nil {
		slog.ErrorContext(ctx, "failed to add new activity", logger.Err(err))
		router.HandleError(w, err)
		return
	}

	err = router.Write(w, response)
	if err != nil {
		slog.ErrorContext(ctx, "Failed to write response", logger.Err(err))
	}
}
