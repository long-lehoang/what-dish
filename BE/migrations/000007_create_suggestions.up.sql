CREATE TABLE suggestion_sessions (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id             UUID,
    session_type        VARCHAR(30) NOT NULL,
    input_params        JSONB NOT NULL DEFAULT '{}',
    result_recipe_ids   UUID[] NOT NULL DEFAULT '{}',
    total_calories      INT,
    created_at          TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_suggestion_sessions_user ON suggestion_sessions(user_id, created_at DESC);
CREATE INDEX idx_suggestion_sessions_type ON suggestion_sessions(session_type);
CREATE INDEX idx_suggestion_sessions_created ON suggestion_sessions(created_at);

CREATE TABLE suggestion_configs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    config_type     VARCHAR(30) NOT NULL,
    name            VARCHAR(100) NOT NULL,
    description     TEXT,
    params          JSONB NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    sort_order      INT DEFAULT 0,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    updated_at      TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE suggestion_exclusions (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id         UUID NOT NULL,
    recipe_id       UUID NOT NULL,
    excluded_until  TIMESTAMPTZ NOT NULL,
    created_at      TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(user_id, recipe_id)
);

CREATE INDEX idx_suggestion_exclusions_user ON suggestion_exclusions(user_id, excluded_until);
