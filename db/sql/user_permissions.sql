-- name: AssignUserPermissions :exec
INSERT INTO user_permissions (user_id, permission_id)
SELECT $1, unnest($2::int[])
ON CONFLICT DO NOTHING;

-- name: DetachUserPermissions :exec
DELETE FROM user_permissions
WHERE user_id = $1
  AND permission_id = ANY($2::int[]);

-- name: SyncUserPermissions :exec
WITH desired AS (
  SELECT unnest($2::int[]) AS permission_id
)
DELETE FROM user_permissions
WHERE user_id = $1
  AND permission_id NOT IN (SELECT permission_id FROM desired);

INSERT INTO user_permissions (user_id, permission_id)
SELECT $1, permission_id
FROM desired
ON CONFLICT DO NOTHING;

-- name: GetUserPermissions :many
SELECT p.id, p.name, p.code, p.description, p.parent_id, p.created_at, p.updated_at
FROM permissions p
JOIN user_permissions up ON p.id = up.permission_id
WHERE up.user_id = $1
ORDER BY p.id;
