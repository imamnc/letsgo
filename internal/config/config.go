package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
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
	// JWT configuration
	JWTSecret string
	JWTIssuer string
	JWTExpiry time.Duration
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
		// JWT configuration
		JWTSecret: getEnv("JWT_SECRET", "secret"),
		JWTIssuer: getEnv("JWT_ISSUER", "letsgo"),
		JWTExpiry: time.Duration(getEnvInt64("JWT_EXPIRY", 3600)) * time.Second,
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

func getEnvInt64(key string, fallback int64) int64 {
	if value := os.Getenv(key); value != "" {
		if v, err := strconv.ParseInt(value, 10, 64); err == nil {
			return v
		}
	}
	return fallback
}
