-- name: GetNodeByWebhookSecretHash :one
SELECT id, cluster_id, node_name, webhook_secret_hash, current_job_id, created_at
FROM nodes
WHERE webhook_secret_hash = $1
LIMIT 1;

-- name: ListNodesByCluster :many
SELECT id, cluster_id, node_name, webhook_secret_hash, current_job_id, created_at
FROM nodes
WHERE cluster_id = $1
ORDER BY node_name ASC;

-- name: CreateNode :one
INSERT INTO nodes (
  cluster_id, node_name, webhook_secret_hash
) VALUES (
  $1, $2, $3
)
RETURNING id, cluster_id, node_name, webhook_secret_hash, current_job_id, created_at;

-- name: UpdateNodeCurrentJob :one
UPDATE nodes
SET current_job_id = $2
WHERE id = $1
RETURNING id, cluster_id, node_name, webhook_secret_hash, current_job_id, created_at;

-- name: DeleteNodeByCluster :execrows
DELETE FROM nodes
WHERE id = $1 AND cluster_id = $2;