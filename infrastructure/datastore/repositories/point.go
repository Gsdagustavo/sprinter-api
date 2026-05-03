package repositories

import (
	"context"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
)

// PointRepository defines methods for managing the points
type PointRepository interface {
	// SaveActivityPoints tries to save the activity points with the given information
	SaveActivityPoints(
		ctx context.Context,
		points []entities.Point,
	) ([]int64, error)
}
