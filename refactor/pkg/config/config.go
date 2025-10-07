package config

// Configuration for the database connection
type DatabaseConfig struct {
	HostName     string
	Port         int
	User         string
	Password     string
	DatabaseName string
}
