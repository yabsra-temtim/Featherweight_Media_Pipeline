CREATE TABLE IF NOT EXISTS jobs (
    id                   TEXT PRIMARY KEY,
    original_name        TEXT NOT NULL,
    original_path        TEXT NOT NULL,
    original_size        BIGINT NOT NULL DEFAULT 0,
    width                INTEGER NOT NULL DEFAULT 0,
    quality              INTEGER NOT NULL DEFAULT 0,
    formats              JSONB NOT NULL DEFAULT '[]',
    status               TEXT NOT NULL,
    outputs              JSONB NOT NULL DEFAULT '[]',
    error                TEXT NOT NULL DEFAULT '',
    created_at           TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at           TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_jobs_created_at ON jobs (created_at);