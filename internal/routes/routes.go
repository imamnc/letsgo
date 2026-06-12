package routes

import (
	"letsgo/internal/app"
	"letsgo/internal/modules/auth"
	"letsgo/internal/modules/health"
	"letsgo/internal/modules/user"
)

func Register(application *app.Application) {
	// Create API group
	api := application.App.Group("/api")
	// Create v1 group
	v1 := api.Group("/v1")

	// Register health routes
	health.Register(application, v1)
	// Register auth routes
	auth.Register(application, v1)
	// Register user routes
	user.Register(application, v1)
}
