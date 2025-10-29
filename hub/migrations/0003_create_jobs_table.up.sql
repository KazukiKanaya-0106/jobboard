CREATE TABLE IF NOT EXISTS jobs (
    id BIGSERIAL PRIMARY KEY,
    cluster_id VARCHAR(64) NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    node_id BIGINT NOT NULL REFERENCES nodes(id) ON DELETE CASCADE,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    status VARCHAR(16) NOT NULL DEFAULT 'running',
    tag VARCHAR(128),
    CONSTRAINT jobs_status_check CHECK (status IN ('running', 'completed', 'failed'))
);

CREATE INDEX jobs_cluster_id_idx ON jobs (cluster_id);
CREATE INDEX jobs_node_id_idx ON jobs (node_id);

CREATE UNIQUE INDEX jobs_one_active_per_node_idx
ON jobs (node_id)
WHERE finished_at IS NULL;

ALTER TABLE nodes
    ADD CONSTRAINT nodes_current_job_id_fk
    FOREIGN KEY (current_job_id)
    REFERENCES jobs(id)
    ON DELETE SET NULL
    DEFERRABLE INITIALLY IMMEDIATE;
