package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Gsdagustavo/sprinter-api/domain/entities"
	"github.com/Gsdagustavo/sprinter-api/domain/entities/derr"
	"github.com/Gsdagustavo/sprinter-api/infrastructure/datastore"
	_ "github.com/go-sql-driver/mysql"
)

type settingsRepository struct {
	connection *sql.DB
}

func NewSettingsRepository(config entities.Settings) (datastore.RepositorySettings, error) {
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
		return derr.JoinInternalError(err, "failed to close database connection")
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
		return nil, derr.JoinInternalError(err, "failed to get server time")
	}

	return &serverTime, nil
}

func setupConnection(config entities.Settings) (*sql.DB, error) {
	//err := migrateDatabase(config)
	//if err != nil {
	//	return nil, err
	//}

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
		return nil, derr.JoinInternalError(err, "failed to open database connection")
	}

	db.SetMaxOpenConns(800)
	db.SetConnMaxLifetime(20 * time.Minute)
	db.SetConnMaxIdleTime(20 * time.Minute)

	return db, nil
}

//func migrateDatabase(config entities.Settings) error {
//	connection := fmt.Sprintf(
//		"mysql://%s:%s@tcp(%s:%s)/%s?parseTime=true",
//		config.RepositorySettings.User,
//		config.RepositorySettings.Password,
//		config.RepositorySettings.Host,
//		config.RepositorySettings.Port,
//		config.RepositorySettings.Name,
//	)
//
//	fsMigrations, err := iofs.New(fs, "migrations")
//	if err != nil {
//		return errors.Join(errors.New("failed to create migration using fs"), err)
//	}
//
//	// Defining the settings used by the migration service
//	migration, err := migrate.NewWithSourceInstance(
//		"iofs",
//		fsMigrations,
//		connection,
//	)
//	if err != nil {
//		return errors.Join(errors.New("failed to create migration with iofs and connection"), err)
//	}
//	defer migration.Close()
//
//	currentVersion, isDirty, err := migration.Version()
//	if err != nil && err.Error() != "no migration" {
//		return errors.Join(errors.New("failed to load the database version"), err)
//	}
//
//	if isDirty {
//		return errors.New("the database is in a dirty state, clear the errors and restart the system")
//	}
//
//	// Checks the need to apply a database migration
//	if DatabaseVersion > currentVersion {
//		slog.Info(
//			"apply database migrations",
//			slog.int("from", int(currentVersion)),
//			slog.int("to", int(DatabaseVersion)),
//		)
//		err = migration.Up()
//		if err != nil {
//			slog.Error(
//				"error running migration",
//				slog.String("error", err.Error()),
//				slog.int("from", int(currentVersion)),
//				slog.int("to", int(DatabaseVersion)),
//			)
//			return errors.Join(errors.New("failed to run up migrations"), err)
//		}
//	} else if DatabaseVersion < currentVersion {
//		slog.Info(
//			"downgrade database",
//			slog.int("from", int(currentVersion)),
//			slog.int("to", int(DatabaseVersion)),
//		)
//		err = migration.Down()
//		if err != nil {
//			slog.Error(
//				"error running downgrade",
//				slog.String("error", err.Error()),
//				slog.int("from", int(currentVersion)),
//				slog.int("to", int(DatabaseVersion)),
//			)
//			return errors.Join(errors.New("failed to run down migrations"), err)
//		}
//	} else {
//		slog.Info("database is up to date", slog.int("version", int(DatabaseVersion)))
//	}
//
//	return nil
//}
