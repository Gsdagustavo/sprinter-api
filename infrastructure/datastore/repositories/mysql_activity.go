package repositories

import (
	"context"
	"database/sql"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
)

type activityRepository struct {
	conn     *sql.DB
	settings datastore.RepositorySettings
}

func NewActivityRepository(settings datastore.RepositorySettings) datastore.ActivityRepository {
	return activityRepository{conn: settings.Connection(), settings: settings}
}

func (a activityRepository) StartActivity(ctx context.Context, activityType entities.ActivityType) (int64, error) {
	const query = `
		INSERT INTO activities (type) VALUES (?)
	`

	result, err := a.conn.ExecContext(ctx, query, activityType)
	if err != nil {
		return -1, derr.JoinInternalError(err, "failed to execute query")
	}

	ID, err := result.LastInsertId()
	if err != nil {
		return -1, derr.JoinInternalError(err, "failed to get last inserted ID")
	}

	return ID, nil
}

func (a activityRepository) EndActivity(ctx context.Context, activityID int64) error {
	// TODO: Implement activity end
	const query = `
		
	`

	_, err := a.conn.ExecContext(ctx, query, activityID)
	if err != nil {
		return derr.JoinInternalError(err, "failed to execute query")
	}

	return nil
}
