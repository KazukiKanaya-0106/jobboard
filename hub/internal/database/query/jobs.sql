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
    cluster_id,
    node_id,
    started_at,
    status,
    tag
) VALUES (
    $1,
    $2,
    COALESCE($3, NOW()),
    COALESCE($4, 'running'),
    $5
)
RETURNING *;

-- name: UpdateJob :one
UPDATE jobs
SET started_at  = COALESCE($2, started_at),
    finished_at = COALESCE($3, finished_at),
    status      = COALESCE($4, status),
    tag         = COALESCE($5, tag),
    duration_hours = COALESCE($6, duration_hours),
    error_text = $7
WHERE id = $1
RETURNING *;
