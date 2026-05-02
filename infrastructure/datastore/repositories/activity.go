package repositories

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// ActivityRepository defines methods for managing activities
type ActivityRepository interface {
	// AddNewActivity tries to add a new activity with the given information
	AddNewActivity(
		ctx context.Context,
		activity entities.Activity,
	) (int64, error)
}
