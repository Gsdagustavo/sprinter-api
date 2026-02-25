package entities

// Server represents the configuration for a server.
type Server struct {
	// Port is the server host port
	Port int `toml:"port"`

	// Host is the server host address
	Host string `toml:"host"`
}

// Database contains configuration details for connecting to a database.
type Database struct {
	// Host is the database host address
	Host string `toml:"host"`

	// Port is the database port
	Port string `toml:"port"`

	// Name is the database name
	Name string `toml:"name"`

	// User is the username for database connection auth
	User string `toml:"user"`

	// Password is the password for database connection auth
	Password string `toml:"password"`
}

// Paseto represents the configuration for PASETO-based token management.
type Paseto struct {
	// SecurityKey is the key used for token generation/validation
	SecurityKey string `toml:"paseto_security_key"`
}

// FileStorage represents the configuration for storing files.
type FileStorage struct {
	// StorageFolder is the folder for storing files
	StorageFolder string `toml:"storage_folder"`
}

// SMTPConfig represents the configuration for connecting to an SMTP server.
type SMTPConfig struct {

	// Host specifies the hostname or IP address of the SMTP server.
	Host string `toml:"host"`

	// Port specifies the port number for the SMTP server connection.
	Port string `toml:"port"`

	// User specifies the username for authentication with the SMTP server.
	User string `toml:"user"`

	//
	Password string `toml:"password"`

	// From specifies the sender's email address for the SMTP server.
	From string `toml:"from"`
}

// Addr returns the SMTP server address in the format "host:port".
func (s *SMTPConfig) Addr() string {
	return s.Host + ":" + s.Port
}

// Config holds the configuration settings for the application.
type Config struct {
	// LogDir represents the directory path of the log file
	LogDir string

	// Environment is the environment where the application is running
	Environment string

	// Server is the data for server host
	Server Server

	// Database is the data for database connection
	Database Database

	// Paseto is the data for encryption
	Paseto Paseto

	// FileStorage contains data related to file storage
	FileStorage FileStorage

	// SMTPConfig holds configuration details for the SMTP server connection.
	SMTPConfig SMTPConfig
}

// IsProduction returns true if the application is running in production environment
func (c Config) IsProduction() bool {
	return c.Environment == "production" || c.Environment == "prod"
}

// IsHomolog returns true if the application is running in homolog environment
func (c Config) IsHomolog() bool {
	return c.Environment == "homolog" || c.Environment == "homo"
}

// IsLocal returns true if the application is running in local environment
func (c Config) IsLocal() bool {
	return c.Environment == "local"
}
