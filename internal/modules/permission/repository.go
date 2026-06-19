package permission

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	dbsql "letsgo/db/sqlc"
	"letsgo/internal/app"
)

var (
	ErrPermissionNotFound    = errors.New("permission not found")
	ErrPermissionAlreadyUsed = errors.New("permission code or name already in use")
)

type Repository struct {
	app *app.Application
}

func NewRepository(application *app.Application) *Repository {
	return &Repository{app: application}
}

func (r *Repository) CreatePermission(ctx context.Context, arg dbsql.CreatePermissionParams) (dbsql.Permission, error) {
	permission, err := r.app.Queries.CreatePermission(ctx, arg)
	if err != nil {
		return dbsql.Permission{}, translatePermissionError(err)
	}
	return permission, nil
}

func (r *Repository) FindPermissionByID(ctx context.Context, id int32) (dbsql.Permission, error) {
	permission, err := r.app.Queries.GetPermissionByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return dbsql.Permission{}, ErrPermissionNotFound
	}
	return permission, err
}

func (r *Repository) ListPermissions(ctx context.Context, limit, offset int32) ([]dbsql.Permission, error) {
	return r.app.Queries.ListPermissions(ctx, dbsql.ListPermissionsParams{Limit: limit, Offset: offset})
}

func (r *Repository) UpdatePermission(ctx context.Context, arg dbsql.UpdatePermissionParams) (dbsql.Permission, error) {
	permission, err := r.app.Queries.UpdatePermission(ctx, arg)
	if err != nil {
		return dbsql.Permission{}, translatePermissionError(err)
	}
	return permission, nil
}

func (r *Repository) DeletePermission(ctx context.Context, id int32) error {
	err := r.app.Queries.DeletePermission(ctx, id)
	return translatePermissionError(err)
}

func translatePermissionError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrPermissionNotFound
	}

	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return ErrPermissionAlreadyUsed
	}

	return err
}
