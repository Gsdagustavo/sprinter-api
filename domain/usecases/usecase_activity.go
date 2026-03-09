package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain"
	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

type activityUseCases struct {
	repository datastore.ActivityRepository
}

func NewActivityUseCases(repository datastore.ActivityRepository) domain.ActivityUseCase {
	return activityUseCases{repository: repository}
}

func (a activityUseCases) StartActivity(ctx context.Context, activityType entities.ActivityType) (int64, error) {
	return a.repository.StartActivity(ctx, activityType)
}

func (a activityUseCases) EndActivity(ctx context.Context, activityID int64) error {
	return a.repository.EndActivity(ctx, activityID)
}
