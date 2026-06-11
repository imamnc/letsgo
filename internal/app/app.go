package app

import (
	dbsql "letsgo/db/sqlc"
	"letsgo/internal/config"
	"letsgo/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type Application struct {
	App     *fiber.App
	Config  config.Config
	DB      *pgxpool.Pool
	Queries *dbsql.Queries
}

func New() *Application {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()
	// Connect to the database
	db := database.Connect(cfg)
	// Initialize SQLC queries using database/sql compatibility layer
	sqlDB := stdlib.OpenDBFromPool(db)
	queries := dbsql.New(sqlDB)
	// Initialize Fiber app
	fiberApp := fiber.New()

	// Create the application instance
	app := &Application{
		App:     fiberApp,
		Config:  cfg,
		DB:      db,
		Queries: queries,
	}
	// Return the application instance
	return app
}

func (a *Application) Run() error {
	// Start the Fiber app
	return a.App.Listen(":" + a.Config.Port)
}
