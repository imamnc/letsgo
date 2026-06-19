package user

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	dbsql "letsgo/db/sqlc"
	"letsgo/internal/app"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already in use")
)

type Repository struct {
	app *app.Application
}

func NewRepository(application *app.Application) *Repository {
	return &Repository{app: application}
}

func (r *Repository) CreateUser(ctx context.Context, arg dbsql.CreateUserParams) (dbsql.User, error) {
	user, err := r.app.Queries.CreateUser(ctx, arg)
	if err != nil {
		return dbsql.User{}, translateUserError(err)
	}
	return user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id int32) (dbsql.User, error) {
	user, err := r.app.Queries.GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return dbsql.User{}, ErrUserNotFound
	}
	return user, err
}

func (r *Repository) ListUsers(ctx context.Context, limit, offset int32) ([]dbsql.User, error) {
	return r.app.Queries.ListUsers(ctx, dbsql.ListUsersParams{Limit: limit, Offset: offset})
}

func (r *Repository) UpdateUser(ctx context.Context, arg dbsql.UpdateUserParams) (dbsql.User, error) {
	user, err := r.app.Queries.UpdateUser(ctx, arg)
	if err != nil {
		return dbsql.User{}, translateUserError(err)
	}
	return user, nil
}

func (r *Repository) DeleteUser(ctx context.Context, id int32) error {
	err := r.app.Queries.DeleteUser(ctx, id)
	return translateUserError(err)
}

func (r *Repository) AssignPermissions(ctx context.Context, userID int32, permissionIDs []int32) error {
	return r.app.Queries.AssignUserPermissions(ctx, dbsql.AssignUserPermissionsParams{
		UserID:  userID,
		Column2: permissionIDs,
	})
}

func (r *Repository) DetachPermissions(ctx context.Context, userID int32, permissionIDs []int32) error {
	return r.app.Queries.DetachUserPermissions(ctx, dbsql.DetachUserPermissionsParams{
		UserID:  userID,
		Column2: permissionIDs,
	})
}

func (r *Repository) SyncPermissions(ctx context.Context, userID int32, permissionIDs []int32) error {
	return r.app.Queries.SyncUserPermissions(ctx, dbsql.SyncUserPermissionsParams{
		UserID:  userID,
		Column2: permissionIDs,
	})
}

func translateUserError(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ErrUserNotFound
	}

	if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
		return ErrEmailAlreadyUsed
	}

	return err
}
