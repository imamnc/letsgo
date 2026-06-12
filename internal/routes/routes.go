package routes

import (
	docs "letsgo/docs"
	"letsgo/internal/app"
	"letsgo/internal/modules/auth"
	"letsgo/internal/modules/health"
	"letsgo/internal/modules/user"

	"github.com/gofiber/swagger"
)

func Register(application *app.Application) {
	// Create API group
	api := application.App.Group("/api")

	// Set Swagger info
	SetupDocsInfo(application)
	api.Get("/docs/*", swagger.HandlerDefault)

	// Create v1 group
	v1 := api.Group("/v1")

	// Register health routes
	health.Register(application, v1)
	// Register auth routes
	auth.Register(application, v1)
	// Register user routes
	user.Register(application, v1)
}

func SetupDocsInfo(application *app.Application) {
	// Set Swagger info
	docs.SwaggerInfo.Title = "LetsGO API Boilerplate"
	docs.SwaggerInfo.Description = "This is a REST API Boilerplate built with Go and Fiber. Battery included, Ready to use, and Good Developer Experience. Just clone and LetsGO!"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = application.Config.AppHost
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
	// Set the base path for the API
	docs.SwaggerInfo.BasePath = "/api/v1"
}
