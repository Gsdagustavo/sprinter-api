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

func (r activityRepository) AddNewActivity(
	ctx context.Context,
	activity entities.Activity,
) (int64, error) {
	const query = `
	INSERT INTO activities (
		uuid,
		user_id,
		type,
		start_time,
		end_time) 
	VALUES (?, ?, ?, ?, ?)
`

	res, err := r.conn.ExecContext(
		ctx,
		query,
		activity.UUID,
		activity.UserID,
		activity.Type,
		activity.StartTime,
		activity.EndTime,
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
