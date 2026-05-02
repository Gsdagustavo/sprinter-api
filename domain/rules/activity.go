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
