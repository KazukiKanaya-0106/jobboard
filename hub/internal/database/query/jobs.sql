-- name: GetJobByClusterAndJobID :one
SELECT * FROM jobs
WHERE cluster_id = $1 AND id = $2 LIMIT 1;

-- name: ListJobsByCluster :many
SELECT * FROM jobs
WHERE cluster_id = $1
ORDER BY started_at DESC, id DESC;

-- name: ListJobsByNode :many
SELECT * FROM jobs
WHERE node_id = $1
ORDER BY started_at DESC, id DESC;

-- name: CreateJob :one
INSERT INTO jobs (
  cluster_id, node_id
) VALUES (
  $1, $2
)
RETURNING *;

-- name: UpdateJob :one
UPDATE jobs
SET started_at = COALESCE($2, started_at),
    finished_at = COALESCE($3, finished_at),
    status = COALESCE($4, status),
    tag = COALESCE($5, tag)
WHERE id = $1
RETURNING *;
