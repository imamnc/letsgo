package health

import (
	"letsgo/internal/app"

	"github.com/gofiber/fiber/v2"
)

func Register(application *app.Application, router fiber.Router) {
	// Intialize the repository
	module := NewModule(application)

	// Health check route
	router.Get("/health", module.handler.Check)
}
