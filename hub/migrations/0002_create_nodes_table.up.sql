CREATE TABLE IF NOT EXISTS nodes (
    id BIGSERIAL PRIMARY KEY,
    cluster_id VARCHAR(64) NOT NULL REFERENCES clusters(id) ON DELETE CASCADE,
    node_name VARCHAR(255) NOT NULL,
    node_token_hash TEXT NOT NULL,
    current_job_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (cluster_id, node_name),
    UNIQUE (node_token_hash)
);

CREATE INDEX nodes_cluster_id_idx ON nodes (cluster_id);
