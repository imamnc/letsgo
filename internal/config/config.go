package config

import (
	"fmt"
	"net/url"
	"time"

	"letsgo/shared/env"
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
	_ = env.Load()

	cfg := Config{
		// Application configuration
		Port: env.Get("PORT", "3000"),
		// Database configuration
		DBHost:     env.Get("DB_HOST", "127.0.0.1"),
		DBPort:     env.Get("DB_PORT", "5432"),
		DBUser:     env.Get("DB_USERNAME", "postgres"),
		DBPassword: env.Get("DB_PASSWORD", "password"),
		DBName:     env.Get("DB_DATABASE", "letsgo"),
		DBSSLMode:  env.Get("DB_SSLMODE", "disable"),
		// JWT configuration
		JWTSecret: env.Get("JWT_SECRET", "secret"),
		JWTIssuer: env.Get("JWT_ISSUER", "letsgo"),
		JWTExpiry: time.Duration(env.GetInt64("JWT_EXPIRY", 3600)) * time.Second,
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
