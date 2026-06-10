package health

import (
	"context"

	dbsql "letsgo/internal/database/sqlc"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	db      *pgxpool.Pool
	queries *dbsql.Queries
}

func NewRepository(db *pgxpool.Pool, queries *dbsql.Queries) *Repository {
	return &Repository{
		db:      db,
		queries: queries,
	}
}

func (r *Repository) CheckDatabase(ctx *fiber.Ctx) error {
	return r.db.Ping(ctx.Context())
}

func (r *Repository) CountUsers(ctx context.Context) (int64, error) {
	return r.queries.CountUsers(ctx)
}
