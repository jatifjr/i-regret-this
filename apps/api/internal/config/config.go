package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment    string
	ServerPort     string
	PostgresURL    string
	MaxConn        int
	JWTSecret      string
	AllowedOrigins []string
}

func Load() (*Config, error) {
	if os.Getenv("ENVIRONMENT") != "production" {
		if err := godotenv.Load(".env"); err != nil {
			return nil, err
		}
	}

	// Get allowed origins from environment variable, default to "*" if not set
	originsStr := getEnv("ALLOWED_ORIGINS", "*")
	origins := strings.Split(originsStr, ",")
	// Trim spaces from each origin
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}

	return &Config{
		Environment:    getEnv("ENVIRONMENT", "development"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		PostgresURL:    getEnv("POSTGRES_URL", "postgres://postgres:password@localhost:5432/postgres?sslmode=disable"),
		MaxConn:        getEnvInt("MAX_CONN", 100),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		AllowedOrigins: origins,
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return fallback
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return fallback
	}
	return value
}
