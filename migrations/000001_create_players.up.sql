CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE players (
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email           TEXT NOT NULL UNIQUE,
    username        TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    gender          TEXT NOT NULL DEFAULT 'other',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_active_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    level           INT NOT NULL DEFAULT 1,
    prestige        INT NOT NULL DEFAULT 0,
    exp             INT NOT NULL DEFAULT 0,
    exp_max         INT NOT NULL DEFAULT 100,

    hp              INT NOT NULL DEFAULT 100,
    hp_max          INT NOT NULL DEFAULT 100,
    energy          INT NOT NULL DEFAULT 100,
    energy_max      INT NOT NULL DEFAULT 100,
    nerve           INT NOT NULL DEFAULT 50,
    nerve_max       INT NOT NULL DEFAULT 50,
    awake           INT NOT NULL DEFAULT 100,
    awake_max       INT NOT NULL DEFAULT 100,

    strength        INT NOT NULL DEFAULT 10,
    defense         INT NOT NULL DEFAULT 10,
    speed           INT NOT NULL DEFAULT 10,
    agility         INT NOT NULL DEFAULT 10,

    cash            BIGINT NOT NULL DEFAULT 0,
    bank            BIGINT NOT NULL DEFAULT 0,
    points          BIGINT NOT NULL DEFAULT 0,
    credits         BIGINT NOT NULL DEFAULT 0,

    gang_id         UUID,

    hospital_time   TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01',
    jail_time       TIMESTAMPTZ NOT NULL DEFAULT '1970-01-01'
);

CREATE INDEX idx_players_email ON players (email);
CREATE INDEX idx_players_username ON players (username);
CREATE INDEX idx_players_level ON players (level DESC);
CREATE INDEX idx_players_gang_id ON players (gang_id);
