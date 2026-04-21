package entities

import "time"

// EnvironmentType represents the type of environment
type EnvironmentType string

const (
	// Development represents the environment type for development configurations
	Development EnvironmentType = "development"

	// Production represents the production environment configuration setting
	Production EnvironmentType = "production"
)

// Settings hold the configuration settings for the application
type Settings struct {
	// Server is the data for server host
	Server ServerSettings `yaml:"ServerSettings"`

	// RepositorySettings is the settings for database connection
	RepositorySettings RepositorySettings `yaml:"RepositorySettings"`

	// ServiceConfig represents a service with its name, display name, and description configuration
	ServiceConfig ServiceConfig `yaml:"ServiceConfig"`

	// Environment is the settings used to handle environment
	EnvironmentSettings EnvironmentSettings `yaml:"EnvironmentSettings"`

	// LogSettings is the settings used to handle logs
	LogSettings LogSettings `yaml:"LogSettings"`

	// PasetoSettings is the settings for data for encryption
	PasetoSettings PasetoSettings `yaml:"PasetoSettings"`

	// CORSConfig defines the configuration for Cross-Origin Resource Sharing (CORS) settings
	CORSConfig CORSConfig `json:"corsConfig" toml:"CORSConfig"`

	// FileStorage is the settings related to file storage
	FileStorageSettings FileStorageSettings `yaml:"FileStorageSettings"`

	// SMTPSettings is the for the SMTP server connection
	SMTPSettings SMTPSettings `yaml:"SMTPSettings"`
}

// ServiceConfig stores the name, display name, and description of a service
type ServiceConfig struct {
	// Name is the name of the service
	Name string `yaml:"Name"`

	// Display is the display name of the service
	Display string `yaml:"Display"`

	// Description is the description of the service
	Description string `yaml:"Description"`
}

// LogSettings is the settings used to handle logs
type LogSettings struct {
	// LogDir represents the directory path of the log file
	LogDir string `yaml:"LogDir"`
}

// IsProduction returns true if the application is running in a production environment
func (s Settings) IsProduction() bool {
	return s.EnvironmentSettings.EnvironmentType == Production
}

// IsLocal returns true if the application is running in the local environment
func (s Settings) IsLocal() bool {
	return s.EnvironmentSettings.EnvironmentType == Development
}

// ServerSettings represents the configuration for a server
type ServerSettings struct {
	// Port is the server host port
	Port int `yaml:"Port"`

	// Host is the server host address
	Host string `yaml:"Host"`

	// Domain is the server domain
	Domain string `yaml:"Domain"`
}

// RepositorySettings contains configuration details to connect to a database
type RepositorySettings struct {
	// Host is the database host address
	Host string `yaml:"Host"`

	// Port is the database port
	Port string `yaml:"Port"`

	// Name is the database name
	Name string `yaml:"Name"`

	// User is the username for database connection auth
	User string `yaml:"User"`

	// Password is the password for database connection auth
	Password string `yaml:"Password"`
}

// EnvironmentSettings represents the configuration for environment
type EnvironmentSettings struct {
	// EnvironmentType is the type of environment. Could be either "development" or "production"
	EnvironmentType EnvironmentType `yaml:"EnvironmentType"`
}

// PasetoSettings represents the configuration for PASETO-based token management
type PasetoSettings struct {
	// SecurityKey is the key used for token generation/validation
	SecurityKey string `yaml:"PasetoSecurityKey"`
}

// CORSConfig stores CORS-related configuration
type CORSConfig struct {
	// CORSOrigins is a list of allowed origins for CORS requests
	CORSOrigins []string `yaml:"CORSOrigins"`
}

// FileStorageSettings represents the configuration for storing files.
type FileStorageSettings struct {
	// StorageFolder is the folder for storing files
	StorageFolder string `yaml:"StorageFolder"`
}

// SMTPSettings represents the configuration for connecting to an SMTP server
type SMTPSettings struct {
	// Host specifies the hostname or IP address of the SMTP server
	Host string `yaml:"Host"`

	// Port specifies the port number for the SMTP server connection
	Port int `yaml:"Port"`

	// User specifies the username for authentication with the SMTP server
	User string `yaml:"User"`

	// Password is the SMTP password credential
	Password string `yaml:"Password"`

	// From specifies the sender's email address for the SMTP server
	From string `yaml:"From"`

	// MaxConnections specifies the maximum number of connections in the email connection pool
	MaxConnections int `yaml:"MaxConnections"`

	// IdleTimeout is the maximum time to wait for new activity on a connection
	// before closing it and removing it from the pool.
	IdleTimeout time.Duration `yaml:"IdleTimeout"`

	// PoolWaitTimeout is the maximum time to wait to obtain a connection from
	// a pool before timing out. This may happen when all open connections are
	// busy sending e-mails, and they're not returning to the pool fast enough.
	// This is also the timeout used when creating new SMTP connections.
	PoolWaitTimeout time.Duration `yaml:"PoolWaitTimeout"`
}
