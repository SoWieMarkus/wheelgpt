package config

// Configuration for the database connection
type DatabaseConfig struct {
	HostName     string
	Port         int
	User         string
	Password     string
	DatabaseName string
}

func LoadDatabaseConfig() DatabaseConfig {
	// TODO load from environment variables
	return DatabaseConfig{
		HostName:     "postgres",
		Port:         5432,
		User:         "wheelgpt",
		Password:     "secret",
		DatabaseName: "wheelgpt",
	}
}
