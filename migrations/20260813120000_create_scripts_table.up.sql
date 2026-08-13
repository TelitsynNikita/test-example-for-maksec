-- +migrate Up
CREATE TABLE IF NOT EXISTS scripts (
    id UUID PRIMARY KEY,
    host TEXT NOT NULL,
    user_name TEXT NOT NULL,
    template TEXT NOT NULL,
    path TEXT NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_scripts_host ON scripts(host);
CREATE INDEX IF NOT EXISTS idx_scripts_path ON scripts(path);
CREATE INDEX IF NOT EXISTS idx_scripts_created_at ON scripts(created_at);