package seeds

import (
	"context"
	"fmt"

	dbsql "letsgo/db/sqlc"
)

func SeedUserPermissions(ctx context.Context, queries *dbsql.Queries) error {
	if err := queries.SeedAllUserPermissions(ctx); err != nil {
		return fmt.Errorf("failed to seed user permissions: %w", err)
	}

	fmt.Println("All permissions assigned to all users")
	return nil
}
