package configs

import (
	"os"
)

// Config holds all application configuration.
type Config struct {
	DBDriver   string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	WebPort    string
	GRPCPort   string
	GraphQLPort string
}

// LoadConfig reads configuration from environment variables with defaults.
func LoadConfig() *Config {
	return &Config{
		DBDriver:    getEnv("DB_DRIVER", "mysql"),
		DBHost:      getEnv("DB_HOST", "localhost"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "root"),
		DBName:      getEnv("DB_NAME", "orders"),
		WebPort:     getEnv("WEB_SERVER_PORT", ":8000"),
		GRPCPort:    getEnv("GRPC_SERVER_PORT", "50051"),
		GraphQLPort: getEnv("GRAPHQL_SERVER_PORT", ":8080"),
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
