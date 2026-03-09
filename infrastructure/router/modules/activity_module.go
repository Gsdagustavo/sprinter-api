package modules

import (
	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/router"
	"github.com/gorilla/mux"
)

type activityModule struct {
	activityUseCases domain.ActivityUseCase
	name             string
	path             string
}

func NewActivityModule(activityUseCases domain.ActivityUseCase) router.Module {
	return activityModule{
		activityUseCases: activityUseCases,
		name:             "Activity",
		path:             "/activity",
	}
}

func (a activityModule) Name() string {
	return a.name
}

func (a activityModule) Path() string {
	return a.path
}

func (a activityModule) Setup(r *mux.Router) ([]router.RouteDefinition, *mux.Router) {
	return nil, nil
}
