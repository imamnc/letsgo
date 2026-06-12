package env

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func Load() error {
	return godotenv.Load()
}

func Get(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func GetInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return fallback
}
