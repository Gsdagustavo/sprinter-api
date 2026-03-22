package entities

// Settings holds the configuration settings for the application
type Settings struct {
	// Server is the data for server host
	Server Server `toml:"Server"`

	// RepositorySettings is the settings for database connection
	RepositorySettings RepositorySettings `toml:"RepositorySettings"`

	// Environment is the environment where the application is running
	Environment string `toml:"Environment"`

	// LogSettings is the settings used to handle logs
	LogSettings LogSettings `toml:"LogSettings"`

	// Paseto is the settings for data for encryption
	Paseto Paseto `toml:"Paseto"`

	// FileStorage is the settings related to file storage
	FileStorage FileStorage `toml:"FileStorage"`

	// SMTPConfig is the for the SMTP server connection
	SMTPConfig SMTPConfig `toml:"STMPConfig"`
}

// LogSettings is the settings used to handle logs
type LogSettings struct {
	// LogDir represents the directory path of the log file
	LogDir string `toml:"LogDir"`
}

// IsProduction returns true if the application is running in production environment
func (c Settings) IsProduction() bool {
	return c.Environment == "production"
}

// IsLocal returns true if the application is running in local environment
func (c Settings) IsLocal() bool {
	return c.Environment == "dev" || c.Environment == "development" || c.Environment == "local"
}

// Server represents the configuration for a server
type Server struct {
	// Port is the server host port
	Port int `toml:"Port"`

	// Host is the server host address
	Host string `toml:"Host"`

	// Domain is the server domain
	Domain string `toml:"Domain"`
}

// RepositorySettings contains configuration details to connect to a database
type RepositorySettings struct {
	// Host is the database host address
	Host string `toml:"Host"`

	// Port is the database port
	Port string `toml:"Port"`

	// Name is the database name
	Name string `toml:"Name"`

	// User is the username for database connection auth
	User string `toml:"User"`

	// Password is the password for database connection auth
	Password string `toml:"Password"`
}

// Paseto represents the configuration for PASETO-based token management
type Paseto struct {
	// SecurityKey is the key used for token generation/validation
	SecurityKey string `toml:"PasetoSecurityKey"`
}

// FileStorage represents the configuration for storing files.
type FileStorage struct {
	// StorageFolder is the folder for storing files
	StorageFolder string `toml:"StorageFolder"`
}

// SMTPConfig represents the configuration for connecting to an SMTP server
type SMTPConfig struct {
	// Host specifies the hostname or IP address of the SMTP server
	Host string `toml:"Host"`

	// Port specifies the port number for the SMTP server connection
	Port string `toml:"Port"`

	// User specifies the username for authentication with the SMTP server
	User string `toml:"User"`

	// Password is the SMTP password credential
	Password string `toml:"Password"`

	// From specifies the sender's email address for the SMTP server
	From string `toml:"From"`
}

// Addr returns the SMTP server address in the format "host:port"
func (s *SMTPConfig) Addr() string {
	return s.Host + ":" + s.Port
}
