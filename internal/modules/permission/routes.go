package permission

import (
	"letsgo/internal/app"
	"letsgo/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(application *app.Application, router fiber.Router) {
	module := NewModule(application)

	permissions := router.Group("/permissions")
	permissions.Use(middleware.Auth(application))

	permissions.Get("/", module.handler.List)
	permissions.Post("/", module.handler.Create)
	permissions.Get("/:id", module.handler.Get)
	permissions.Put("/:id", module.handler.Update)
	permissions.Delete("/:id", module.handler.Delete)
}
