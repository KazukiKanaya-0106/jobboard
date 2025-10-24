CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT NOT NULL,
    cluster_id VARCHAR(255) NOT NULL,
    job_number INTEGER NOT NULL,
    tag VARCHAR(255),
    "user" VARCHAR(255),
    started_at TIMESTAMP,
    finished_at TIMESTAMP,
    duration_hours DECIMAL(10, 2),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE CASCADE,
    FOREIGN KEY (cluster_id) REFERENCES clusters(id) ON DELETE CASCADE,
    UNIQUE(cluster_id, job_number)
);