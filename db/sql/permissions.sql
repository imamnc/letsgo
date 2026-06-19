-- name: CreatePermission :one
INSERT INTO permissions (name, code, description, parent_id)
VALUES ($1, $2, $3, $4)
RETURNING id, name, code, description, parent_id, created_at, updated_at;

-- name: GetPermissionByID :one
SELECT id, name, code, description, parent_id, created_at, updated_at
FROM permissions
WHERE id = $1;

-- name: ListPermissions :many
SELECT id, name, code, description, parent_id, created_at, updated_at
FROM permissions
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdatePermission :one
UPDATE permissions
SET name = $2,
    code = $3,
    description = $4,
    parent_id = $5,
    updated_at = now()
WHERE id = $1
RETURNING id, name, code, description, parent_id, created_at, updated_at;

-- name: DeletePermission :exec
DELETE FROM permissions
WHERE id = $1;
