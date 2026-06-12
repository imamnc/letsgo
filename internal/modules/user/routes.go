package user

import (
	"letsgo/internal/app"

	"github.com/gofiber/fiber/v2"
)

func Register(application *app.Application, router fiber.Router) {
	module := NewModule(application)

	users := router.Group("/users")

	users.Get("/", module.handler.List)
	users.Post("/", module.handler.Create)
	users.Get("/:id", module.handler.Get)
	users.Put("/:id", module.handler.Update)
	users.Delete("/:id", module.handler.Delete)
}
