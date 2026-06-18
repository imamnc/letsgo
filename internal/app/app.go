package app

import (
	database "letsgo/db"
	dbsql "letsgo/db/sqlc"
	"letsgo/internal/config"
	"letsgo/shared/jwt"
	"letsgo/shared/redis"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Application struct {
	App     *fiber.App
	Config  config.Config
	DB      *pgxpool.Pool
	Queries *dbsql.Queries
	Jwt     *jwt.Provider
	Redis   *redis.Service
}

func New() *Application {
	// Load configuration
	cfg := config.Load()

	// Connect to the database
	db := database.Connect(cfg)

	// Initialize SQLC queries using database/sql compatibility layer
	sqlDB := stdlib.OpenDBFromPool(db)
	queries := dbsql.New(sqlDB)

	// Initialize Redis client
	redisClient := redis.New(cfg.RedisHost+":"+cfg.RedisPort, "", 0)

	// Initialize Fiber app
	fiberApp := fiber.New()

	// Init jwt provider
	jwtProvider := jwt.New(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTExpiry)

	// Create the application instance
	app := &Application{
		App:     fiberApp,
		Config:  cfg,
		DB:      db,
		Queries: queries,
		Jwt:     jwtProvider,
		Redis:   redis.NewService(redisClient),
	}

	// Return the application instance
	return app
}

func (a *Application) Run() error {
	// Start the Fiber app
	return a.App.Listen(":" + a.Config.AppPort)
}
