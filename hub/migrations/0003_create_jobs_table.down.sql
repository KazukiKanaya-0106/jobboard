ALTER TABLE nodes DROP CONSTRAINT IF EXISTS nodes_current_job_id_fk;

DROP INDEX IF EXISTS jobs_one_active_per_node_idx;
DROP INDEX IF EXISTS jobs_node_id_idx;
DROP INDEX IF EXISTS jobs_cluster_id_idx;

DROP TABLE IF EXISTS jobs;
