-- 0001_create_jobs_table.down.sql

DROP INDEX IF EXISTS idx_jobs_created_at;
DROP TABLE IF EXISTS jobs;
