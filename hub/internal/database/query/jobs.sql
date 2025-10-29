-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: GetJobByClusterAndJobId :one
SELECT * FROM jobs
WHERE cluster_id = $1 AND id = $2 LIMIT 1;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY started_at DESC, id DESC;

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
    status = COALESCE($4, status)
WHERE id = $1
RETURNING *;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = $1;

-- name: DeleteJobsByCluster :exec
DELETE FROM jobs
WHERE cluster_id = $1;

-- name: DeleteJobsByNode :exec
DELETE FROM jobs
WHERE node_id = $1;
