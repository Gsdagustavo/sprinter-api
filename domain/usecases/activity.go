package usecases

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// ActivityUseCase defines a use case interface with methods related to activity managing
type ActivityUseCase interface {
	//AddNewActivity attempts to add a new activity with the given information
	AddNewActivity(
		ctx context.Context,
		activity entities.Activity,
	) (int64, error)
}
