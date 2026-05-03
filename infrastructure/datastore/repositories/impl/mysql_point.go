package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
)

func NewPointRepository(settings repositories.SettingsRepository) repositories.PointRepository {
	return &pointRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type pointRepository struct {
	conn     *sql.DB
	settings repositories.SettingsRepository
}

func (p pointRepository) SaveActivityPoints(ctx context.Context, points []entities.Point) ([]int64, error) {
	const query = `
	INSERT INTO points (
		activity_id,
		latitude,
	    longitude)       
	VALUES (?, ?, ?)
`

	var pointsId []int64
	for _, point := range points {
		_, err := p.conn.ExecContext(
			ctx,
			query,
			point.ActivityID,
			point.Latitude,
			point.Longitude,
		)
		if err != nil {
			return []int64{}, derr.JoinError("failed to insert the point on the database", err)
		}
		pointsId = append(pointsId, point.ID)
	}

	return pointsId, nil
}
