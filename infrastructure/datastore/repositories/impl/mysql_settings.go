package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore/repositories"
	_ "github.com/go-sql-driver/mysql"
)

type settingsRepository struct {
	connection *sql.DB
}

func NewSettingsRepository(config entities.Settings) (repositories.SettingsRepository, error) {
	db, err := setupConnection(config)
	if err != nil {
		return nil, err
	}

	return settingsRepository{
		connection: db,
	}, nil
}

func (s settingsRepository) Connection() *sql.DB {
	return s.connection
}

func (s settingsRepository) Dismount() error {
	err := s.connection.Close()
	if err != nil {
		return derr.JoinError("failed to close database connection", err)
	}

	return nil
}

func (s settingsRepository) ServerTime(
		ctx context.Context,
) (*time.Time, error) {
	//language=sql
	query := "SELECT CURRENT_TIMESTAMP"

	var serverTime time.Time
	err := s.connection.QueryRowContext(ctx, query).Scan(&serverTime)
	if err != nil {
		return nil, derr.JoinError("failed to get server time", err)
	}

	return &serverTime, nil
}

func setupConnection(config entities.Settings) (*sql.DB, error) {
	connection := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.RepositorySettings.User,
		config.RepositorySettings.Password,
		config.RepositorySettings.Host,
		config.RepositorySettings.Port,
		config.RepositorySettings.Name,
	)

	db, err := sql.Open("mysql", connection)
	if err != nil {
		return nil, derr.JoinError("failed to open database connection", err)
	}

	db.SetMaxOpenConns(800)
	db.SetConnMaxLifetime(20 * time.Minute)
	db.SetConnMaxIdleTime(20 * time.Minute)

	return db, nil
}
