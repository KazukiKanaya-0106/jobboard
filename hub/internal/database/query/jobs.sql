-- name: GetJob :one
SELECT * FROM jobs
WHERE id = $1 LIMIT 1;

-- name: GetJobByClusterAndJobId :one
SELECT * FROM jobs
WHERE cluster_id = $1 AND job_number = $2 LIMIT 1;

-- name: ListJobs :many
SELECT * FROM jobs
ORDER BY created_at DESC;

-- name: ListJobsByCluster :many
SELECT * FROM jobs
WHERE cluster_id = $1
ORDER BY started_at DESC;

-- name: ListJobsByNode :many
SELECT * FROM jobs
WHERE node_id = $1
ORDER BY started_at DESC;

-- name: ListJobsByUser :many
SELECT * FROM jobs
WHERE "user" = $1
ORDER BY started_at DESC;

-- name: CreateJob :one
INSERT INTO jobs (
  node_id, cluster_id, job_number, tag, "user", started_at, finished_at, duration_hours
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8
)
RETURNING *;

-- name: UpdateJob :one
UPDATE jobs
SET tag = $2,
    "user" = $3,
    started_at = $4,
    finished_at = $5,
    duration_hours = $6
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
