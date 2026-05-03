package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
)

func NewActivityRepository(settings repositories.SettingsRepository) repositories.ActivityRepository {
	return &activityRepository{
		conn:     settings.Connection(),
		settings: settings,
	}
}

type activityRepository struct {
	conn     *sql.DB
	settings repositories.SettingsRepository
}

func (r activityRepository) StartActivity(
		ctx context.Context,
		activity *entities.Activity,
) (int64, error) {
	const query = `
	INSERT INTO activities (
		user_id,
		type,
	    start_date,
	VALUES (?, ?, ?)
`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		activity.UserID,
		activity.Type,
		activity.StartDate,
	)
	if err != nil {
		return -1, derr.JoinError("failed to execute query", err)
	}

	activityID, err := res.LastInsertId()
	if err != nil {
		return -1, derr.JoinError("failed to get last inserted ID", err)
	}

	return activityID, nil
}

func (r activityRepository) FinishActivity(
		ctx context.Context,
		activity entities.Activity,
) (int64, error) {
	const query = `
	UPDATE activities 
	SET end_date = ?,
	WHERE id = ?
`
	_, err := r.conn.ExecContext(
		ctx,
		query,
		activity.EndDate,
		activity.ID,
	)
	if err != nil {
		return -1, derr.JoinError("failed to execute the query", err)
	}

	_, err = r.saveActivityPoints(ctx, activity.Route)
	if err != nil {
		return -1, err
	}

	return activity.ID, nil
}
func (r activityRepository) saveActivityPoints(ctx context.Context, points []entities.Point) ([]int64, error) {
	const query = `
	INSERT INTO points (
		activity_id,
		latitude,
	    longitude)       
	VALUES (?, ?, ?)
`

	var pointsId []int64
	for _, point := range points {
		_, err := r.conn.ExecContext(
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
