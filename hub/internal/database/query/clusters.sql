-- name: GetCluster :one
SELECT * FROM clusters
WHERE id = $1 LIMIT 1;

-- name: ListClusters :many
SELECT * FROM clusters
ORDER BY created_at DESC;

-- name: CreateCluster :one
INSERT INTO clusters (
  id, password_hash
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateCluster :one
UPDATE clusters
SET password_hash = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCluster :exec
DELETE FROM clusters
WHERE id = $1;
