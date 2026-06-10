package routes

import (
	"letsgo/internal/app"
	"letsgo/internal/modules/health"
)

func Register(application *app.Application) {
	// Create API group
	api := application.App.Group("/api")
	// Create v1 group
	v1 := api.Group("/v1")

	// Register health routes
	health.Register(application, v1)
}
