package config

import (
	"fmt"
	"net/url"
	"os"
)

type Config struct {
	// Application configuration
	Port string
	// Database configuration
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	DatabaseURL string
}

func Load() Config {
	cfg := Config{
		// Application configuration
		Port: getEnv("PORT", "3000"),
		// Database configuration
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USERNAME", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_DATABASE", "letsgo"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
	cfg.DatabaseURL = buildDatabaseURL(cfg)
	return cfg
}

func buildDatabaseURL(cfg Config) string {
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.DBUser, cfg.DBPassword),
		Host:   fmt.Sprintf("%s:%s", cfg.DBHost, cfg.DBPort),
		Path:   cfg.DBName,
	}

	q := u.Query()
	q.Set("sslmode", cfg.DBSSLMode)
	u.RawQuery = q.Encode()

	return u.String()
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
