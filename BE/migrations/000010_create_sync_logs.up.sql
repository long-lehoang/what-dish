CREATE TABLE sync_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    started_at      TIMESTAMPTZ NOT NULL,
    finished_at     TIMESTAMPTZ,
    status          VARCHAR(20) NOT NULL DEFAULT 'RUNNING',
    recipes_added   INT DEFAULT 0,
    recipes_updated INT DEFAULT 0,
    recipes_deleted INT DEFAULT 0,
    errors          JSONB DEFAULT '[]',
    duration_ms     INT,
    created_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_sync_logs_status ON sync_logs(status, created_at DESC);
