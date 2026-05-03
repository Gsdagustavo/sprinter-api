package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/rules"
	"github.com/Gsdagustavo/sprinter-api/domain/usecases"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
)

func NewActivityUseCases(repository repositories.ActivityRepository) usecases.ActivityUseCase {
	return activityUseCases{
		repository: repository,
	}
}

type activityUseCases struct {
	repository repositories.ActivityRepository
}

func (u activityUseCases) StartActivity(ctx context.Context, activity *entities.Activity) (int64, error) {
	err := rules.ValidateActivityStart(activity)
	if err != nil {
		return 0, err
	}

	return u.repository.StartActivity(ctx, activity)
}

func (u activityUseCases) FinishActivity(ctx context.Context, activity *entities.Activity) (int64, error) {
	err := rules.ValidateActivityFinish(activity)
	if err != nil {
		return 0, err
	}

	err = rules.ValidateActivityPoints(&activity.Route)
	if err != nil {
		return 0, err
	}

	return u.repository.FinishActivity(ctx, *activity)
}
