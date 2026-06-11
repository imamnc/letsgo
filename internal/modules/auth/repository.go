package auth

import (
	"context"
	"database/sql"
	"errors"

	dbsql "letsgo/db/sqlc"
	"letsgo/internal/app"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	app *app.Application
}

func NewRepository(application *app.Application) *Repository {
	return &Repository{app: application}
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (dbsql.User, error) {
	user, err := r.app.Queries.GetUserByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return dbsql.User{}, ErrUserNotFound
	}
	return user, err
}

func (r *Repository) FindUserByID(ctx context.Context, id int64) (dbsql.User, error) {
	user, err := r.app.Queries.GetUserByID(ctx, int32(id))
	if errors.Is(err, sql.ErrNoRows) {
		return dbsql.User{}, ErrUserNotFound
	}
	return user, err
}
