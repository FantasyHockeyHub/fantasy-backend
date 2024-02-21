-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT EXISTS teams
(
    team_id         SERIAL PRIMARY KEY,
    team_abbrev VARCHAR(255) DEFAULT ''::VARCHAR,
    team_name VARCHAR(255) NOT NULL,
    team_logo TEXT DEFAULT ''::TEXT,
    league SMALLINT,
    conference_name VARCHAR(255),
    division VARCHAR(255),
    api_id INT DEFAULT 0::INT
);

CREATE INDEX if NOT EXISTS idx_teams
    ON teams (team_id, team_abbrev);

CREATE UNIQUE INDEX if NOT EXISTS uniq_idx_teams
    ON teams (team_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS teams;
DROP INDEX IF EXISTS idx_teams;
DROP INDEX IF EXISTS uniq_idx_teams;
-- +goose StatementEnd
