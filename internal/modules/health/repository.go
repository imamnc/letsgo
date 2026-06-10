package health

import (
	"context"

	"letsgo/internal/app"

	"github.com/gofiber/fiber/v2"
)

type Repository struct {
	app *app.Application
}

func NewRepository(application *app.Application) *Repository {
	return &Repository{
		app: application,
	}
}

func (r *Repository) CheckDatabase(ctx *fiber.Ctx) error {
	return r.app.DB.Ping(ctx.Context())
}

func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	return r.app.Queries.CountUsers(ctx)
}
