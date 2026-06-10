package health

import (
	"letsgo/internal/app"

	"github.com/gofiber/fiber/v2"
)

func Register(application *app.Application, router fiber.Router) {
	// Intialize the repository
	repository := NewRepository(application.DB, application.Queries)
	// Initialize the service with the repository
	service := NewService(repository)
	// Initialize the handler with the application context
	handler := NewHandler(service)

	// Health check route
	router.Get("/health", handler.Check)

}
