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
		id_user,
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
		return 0, derr.JoinError("failed to execute query", err)
	}

	activityID, err := res.LastInsertId()
	if err != nil {
		return 0, derr.JoinError("failed to get last inserted ID", err)
	}

	return activityID, nil
}

func (r activityRepository) FinishActivity(
	ctx context.Context,
	activity *entities.Activity,
) (int64, error) {
	const query = `
	UPDATE activities SET 
	    end_date = ?
	WHERE id = ?
	`

	tx, err := r.conn.BeginTx(ctx, nil)
	if err != nil {
		return 0, derr.JoinError("failed to begin transaction", err)
	}

	_, err = tx.ExecContext(
		ctx,
		query,
		activity.EndDate,
		activity.ID,
	)
	if err != nil {
		return 0, derr.JoinError("failed to execute the query", err)
	}

	err = r.saveActivityPoints(ctx, tx, activity.Route)
	if err != nil {
		return 0, derr.JoinError("failed to save activity points", err)
	}

	return activity.ID, nil
}
func (r activityRepository) saveActivityPoints(ctx context.Context, tx *sql.Tx, points []entities.Point) error {
	const query = `
	INSERT INTO points (
		id_activity,
		latitude,
	    longitude
	) VALUES (?, ?, ?)
	`

	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		return derr.JoinError("failed to prepare statement", err)
	}
	defer stmt.Close()

	for _, point := range points {
		_, err := stmt.ExecContext(
			ctx,
			query,
			point.ActivityID,
			point.Latitude,
			point.Longitude,
		)
		if err != nil {
			return derr.JoinError("failed to execute query", err)
		}
	}

	return nil
}
