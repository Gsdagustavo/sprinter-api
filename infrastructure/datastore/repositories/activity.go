package repositories

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// ActivityRepository defines methods for managing activities
type ActivityRepository interface {
	// StartActivity tries to start a new activity with the given information
	StartActivity(
		ctx context.Context,
		activity *entities.Activity,
	) (int64, error)
}
