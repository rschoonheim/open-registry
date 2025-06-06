package config

import (
	"fmt"
	"os"
)

// EnvConfig holds configuration loaded from environment variables
type EnvConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	AppPort    string
}

// LoadEnvConfig loads configuration from environment variables
func LoadEnvConfig() *EnvConfig {
	config := &EnvConfig{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "postgres"),
		JWTSecret:  getEnv("JWT_SECRET", "your-256-bit-secret-key-change-this-in-production"),
		AppPort:    getEnv("PORT", "3000"),
	}

	return config
}

// getEnv gets an environment variable or returns a default value if not set
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetDSN returns the Data Source Name for the database
func (c *EnvConfig) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)
}
