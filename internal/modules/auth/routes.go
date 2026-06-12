package auth

import (
	"letsgo/internal/app"
	"letsgo/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(application *app.Application, router fiber.Router) {
	module := NewModule(application)

	auth := router.Group("/auth")

	auth.Post("/access-token", module.handler.AccessToken)
	auth.Post("/refresh-token", module.handler.RefreshToken)
	auth.Get("/user", middleware.Auth(application), module.handler.FetchUser)
}
