package seeds

import (
	"context"
	"fmt"
	"strings"

	dbsql "letsgo/db/sqlc"
	"letsgo/shared/format"
)

type SeedPermissionParams struct {
	Name        string
	Code        string
	Description *string
	ParentID    *int32
}

func PreparePermissions() []SeedPermissionParams {
	readOnly := "Read-only permission"
	execute := "Execute permission"

	return []SeedPermissionParams{
		{
			Name:        "Read Permission",
			Code:        "permission.read",
			Description: &readOnly,
			ParentID:    nil,
		},
		{
			Name:        "Write Permission",
			Code:        "permission.write",
			Description: nil,
			ParentID:    nil,
		},
		{
			Name:        "Execute Permission",
			Code:        "permission.execute",
			Description: &execute,
			ParentID:    nil,
		},
	}
}

func SeedPermission(ctx context.Context, queries *dbsql.Queries) error {
	permissions := PreparePermissions()

	for _, permission := range permissions {
		_, err := queries.CreatePermission(ctx, dbsql.CreatePermissionParams{
			Name:        permission.Name,
			Code:        permission.Code,
			Description: format.ToNullString(permission.Description),
			ParentID:    format.ToNullInt32(permission.ParentID),
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				fmt.Printf("Permission with code %s already exists, skipping...\n", permission.Code)
				continue
			}
			return fmt.Errorf("failed to create permission %s: %w", permission.Code, err)
		}
		fmt.Printf("Permission with code %s created successfully\n", permission.Code)
	}

	return nil
}
