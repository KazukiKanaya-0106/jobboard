-- name: GetNode :one
SELECT * FROM nodes
WHERE id = $1 LIMIT 1;

-- name: GetNodeByClusterAndNodeId :one
SELECT * FROM nodes
WHERE cluster_id = $1 AND node_name = $2 LIMIT 1;

-- name: ListNodes :many
SELECT * FROM nodes
ORDER BY created_at DESC;

-- name: ListNodesByCluster :many
SELECT * FROM nodes
WHERE cluster_id = $1
ORDER BY node_name ASC;

-- name: CreateNode :one
INSERT INTO nodes (
  cluster_id, node_name, current_job_number
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: UpdateNodeJobNumber :one
UPDATE nodes
SET current_job_number = $2
WHERE id = $1
RETURNING *;

-- name: DeleteNode :exec
DELETE FROM nodes
WHERE id = $1;

-- name: DeleteNodesByCluster :exec
DELETE FROM nodes
WHERE cluster_id = $1;
