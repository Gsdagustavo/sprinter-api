package usecases

import (
	"context"

	"github.com/VitorFranciscoDev/sprinter-api/domain/entities"
	"github.com/VitorFranciscoDev/sprinter-api/infrastructure/datastore"
)

type ActivityUseCase struct {
	r datastore.ActivityRepository
}

func NewActivityUseCase(r datastore.ActivityRepository) *ActivityUseCase {
	return &ActivityUseCase{r: r}
}

func (uc *ActivityUseCase) Save(ctx context.Context, activity *entities.Activity) error {
	return uc.r.Save(ctx, activity)
}

func (uc *ActivityUseCase) Get(ctx context.Context, id int64) (*entities.Activity, error) {
	return uc.r.Get(ctx, id)
}

func (uc *ActivityUseCase) GetAll(ctx context.Context) ([]*entities.Activity, error) {
	return uc.r.GetAll(ctx)
}
