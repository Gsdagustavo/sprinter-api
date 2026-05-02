package rules

import (
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
)

func ValidateActivityStart(activity *entities.Activity) error {
	switch activity.Type {
	case entities.CYCLING:
		break
	case entities.RUNNING:
		break
	case entities.WALKING:
		break
	default:
		return derr.InvalidActivityType
	}

	if activity.StartDate == activity.EndDate {
		return derr.InvalidActivityStartDate
	}
	return nil

}
