package permission

import (
	"context"
	"errors"
	"strings"

	dbsql "letsgo/db/sqlc"
	"letsgo/shared/format"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) CreatePermission(ctx context.Context, request CreatePermissionRequest) (dbsql.Permission, error) {
	return s.repository.CreatePermission(ctx, dbsql.CreatePermissionParams{
		Name:        request.Name,
		Code:        request.Code,
		Description: format.ToNullString(request.Description),
		ParentID:    format.ToNullInt32(request.ParentID),
	})
}

func (s *Service) FindPermissionByID(ctx context.Context, id int32) (dbsql.Permission, error) {
	return s.repository.FindPermissionByID(ctx, id)
}

func (s *Service) ListPermissions(ctx context.Context, limit, offset int32) ([]dbsql.Permission, error) {
	return s.repository.ListPermissions(ctx, limit, offset)
}

func (s *Service) UpdatePermission(ctx context.Context, id int32, request UpdatePermissionRequest) (dbsql.Permission, error) {
	permission, err := s.repository.FindPermissionByID(ctx, id)
	if err != nil {
		return dbsql.Permission{}, err
	}

	name := permission.Name
	if request.Name != nil {
		trimmed := strings.TrimSpace(*request.Name)
		if trimmed == "" {
			return dbsql.Permission{}, errors.New("name cannot be empty")
		}
		name = trimmed
	}

	code := permission.Code
	if request.Code != nil {
		trimmed := strings.TrimSpace(*request.Code)
		if trimmed == "" {
			return dbsql.Permission{}, errors.New("code cannot be empty")
		}
		code = trimmed
	}

	description := permission.Description
	if request.Description != nil {
		description = format.ToNullString(request.Description)
	}

	parentID := permission.ParentID
	if request.ParentID != nil {
		parentID = format.ToNullInt32(request.ParentID)
	}

	return s.repository.UpdatePermission(ctx, dbsql.UpdatePermissionParams{
		ID:          id,
		Name:        name,
		Code:        code,
		Description: description,
		ParentID:    parentID,
	})
}

func (s *Service) DeletePermission(ctx context.Context, id int32) error {
	return s.repository.DeletePermission(ctx, id)
}
