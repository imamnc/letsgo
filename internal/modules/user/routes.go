package user

import (
	"letsgo/internal/app"
	"letsgo/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func Register(application *app.Application, router fiber.Router) {
	module := NewModule(application)

	users := router.Group("/users")
	users.Use(middleware.Auth(application))

	users.Get("/", module.handler.List)
	users.Post("/", module.handler.Create)
	users.Get("/:id", module.handler.Get)
	users.Put("/:id", module.handler.Update)
	users.Delete("/:id", module.handler.Delete)
	users.Post("/:id/permissions", module.handler.AssignPermissions)
	users.Delete("/:id/permissions", module.handler.DetachPermissions)
	users.Put("/:id/permissions", module.handler.SyncPermissions)
}
