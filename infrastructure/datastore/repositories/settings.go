package repositories

import (
	"context"
	"database/sql"
	"time"
)

// SettingsRepository creates and manages the access for the database
type SettingsRepository interface {
	// Connection returns a database connection
	Connection() *sql.DB

	// Dismount closes all connections with the database
	Dismount() error

	// ServerTime returns the current time on the server
	ServerTime(ctx context.Context) (*time.Time, error)
}