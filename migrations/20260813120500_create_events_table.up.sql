-- +migrate Up
CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY,
    script_id UUID NOT NULL REFERENCES scripts(id) ON DELETE CASCADE,
    user_name TEXT NOT NULL,
    script_path TEXT NOT NULL,
    action TEXT NOT NULL CHECK (action IN ('execute', 'modify')),
    event_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_events_script_id ON events(script_id);
CREATE INDEX IF NOT EXISTS idx_events_script_path ON events(script_path);
CREATE INDEX IF NOT EXISTS idx_events_event_time ON events(event_time);
CREATE INDEX IF NOT EXISTS idx_events_action ON events(action);