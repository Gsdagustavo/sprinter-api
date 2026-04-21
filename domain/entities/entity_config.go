package entities

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

	// Environment is the settings used to handle environment
	EnvironmentSettings EnvironmentSettings `yaml:"EnvironmentSettings"`

	// LogSettings is the settings used to handle logs
	LogSettings LogSettings `yaml:"LogSettings"`

	// PasetoSettings is the settings for data for encryption
	PasetoSettings PasetoSettings `yaml:"PasetoSettings"`

	// FileStorage is the settings related to file storage
	FileStorageSettings FileStorageSettings `yaml:"FileStorageSettings"`

	// SMTPSettings is the for the SMTP server connection
	SMTPSettings SMTPSettings `yaml:"SMTPSettings"`
}

// LogSettings is the settings used to handle logs
type LogSettings struct {
	// LogDir represents the directory path of the log file
	LogDir string `yaml:"LogDir"`
}

// IsProduction returns true if the application is running in a production environment
func (c Settings) IsProduction() bool {
	return c.EnvironmentSettings.EnvironmentType == Production
}

// IsLocal returns true if the application is running in the local environment
func (c Settings) IsLocal() bool {
	return c.EnvironmentSettings.EnvironmentType == Development
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
	Port string `yaml:"Port"`

	// User specifies the username for authentication with the SMTP server
	User string `yaml:"User"`

	// Password is the SMTP password credential
	Password string `yaml:"Password"`

	// From specifies the sender's email address for the SMTP server
	From string `yaml:"From"`
}

// Addr returns the SMTP server address in the format "host:port"
func (s *SMTPSettings) Addr() string {
	return s.Host + ":" + s.Port
}
