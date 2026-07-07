CREATE TABLE sessions (
    token       TEXT PRIMARY KEY,
    player_id   UUID NOT NULL REFERENCES players(id) ON DELETE CASCADE,
    expires_at  TIMESTAMPTZ NOT NULL,
    ip          TEXT NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_player_id ON sessions (player_id);
CREATE INDEX idx_sessions_expires_at ON sessions (expires_at);
