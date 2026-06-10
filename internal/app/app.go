package app

import (
	"letsgo/internal/config"
	"letsgo/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Application struct {
	App    *fiber.App
	Config config.Config
	DB     *pgxpool.Pool
}

func New() *Application {
	// Load environment variables from .env if present
	_ = godotenv.Load()

	// Load configuration
	cfg := config.Load()
	// Connect to the database
	db := database.Connect(cfg)
	// Initialize Fiber app
	fiberApp := fiber.New()

	// Create the application instance
	app := &Application{
		App:    fiberApp,
		Config: cfg,
		DB:     db,
	}
	// Return the application instance
	return app
}

func (a *Application) Run() error {
	// Start the Fiber app
	return a.App.Listen(":" + a.Config.Port)
}
