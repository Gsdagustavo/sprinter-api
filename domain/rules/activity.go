package rules

import (
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
)

func ValidateActivityStart(activity *entities.Activity) error {
	switch activity.Type {
	case entities.CYCLING, entities.RUNNING, entities.WALKING:
		break
	default:
		return derr.InvalidActivityType
	}

	if activity.StartDate != time.Now() || activity.StartDate.IsZero() {
		return derr.InvalidActivityStartDate
	}

	return nil

}
func ValidateActivityFinish(activity *entities.Activity) error {
	if len(activity.Route) == 1 {
		return derr.InvalidActivityRoute
	}

	if activity.EndDate != time.Now() || activity.EndDate.IsZero() {
		return derr.InvalidActivityEndDate
	}
	return nil
}
