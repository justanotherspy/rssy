package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	Host                string
	DatabasePath        string
	FeedRefreshInterval time.Duration
	AllowedOrigins      []string
}

func Load() *Config {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	port := getEnv("PORT", "8080")
	host := getEnv("HOST", "localhost")
	dbPath := getEnv("DATABASE_PATH", "./rssy.db")

	refreshInterval := getEnvAsDuration("FEED_REFRESH_INTERVAL", "10m")
	allowedOrigins := getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:5173"})

	return &Config{
		Port:                port,
		Host:                host,
		DatabasePath:        dbPath,
		FeedRefreshInterval: refreshInterval,
		AllowedOrigins:      allowedOrigins,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key, defaultValue string) time.Duration {
	valueStr := getEnv(key, defaultValue)
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		log.Printf("Invalid duration for %s, using default: %s", key, defaultValue)
		duration, _ = time.ParseDuration(defaultValue)
	}
	return duration
}

func getEnvAsSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		// Support comma-separated values
		parts := strings.Split(value, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			trimmed := strings.TrimSpace(part)
			if trimmed != "" {
				result = append(result, trimmed)
			}
		}
		return result
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Invalid integer for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}
