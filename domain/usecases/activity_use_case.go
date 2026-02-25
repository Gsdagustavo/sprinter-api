package usecases

import (
	"context"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore"
)

type ActivityUseCase struct {
	r datastore.ActivityRepository
}

func NewActivityUseCase(
	r datastore.ActivityRepository,
) ActivityUseCase {
	return ActivityUseCase{r: r}
}

func (u ActivityUseCase) Save(
	ctx context.Context,
	activity *entities.Activity,
) error {
	return u.r.Save(ctx, activity)
}

func (u ActivityUseCase) Get(
	ctx context.Context,
	id int64,
) (*entities.Activity, error) {
	return u.r.Get(ctx, id)
}

func (u ActivityUseCase) GetAll(
	ctx context.Context,
) ([]*entities.Activity, error) {
	return u.r.GetAll(ctx)
}
