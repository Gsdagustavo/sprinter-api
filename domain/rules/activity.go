package rules

import (
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
)

func ValidateActivity(activity *entities.Activity) error {
	if activity.Type != 1 && activity.Type != 2 && activity.Type != 3 {
		return derr.InvalidActivityType
	}
	if activity.StartTime == activity.EndTime {
		return derr.InvalidActivityDuration
	}
	return nil
}
