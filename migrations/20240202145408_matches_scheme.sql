-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT EXISTS matches
(
    id         SERIAL PRIMARY KEY,
    home_team_id INT,
    home_team_score SMALLINT DEFAULT 0::SMALLINT,
    away_team_id INT,
    away_team_score SMALLINT DEFAULT 0::SMALLINT,
    start_at BIGINT,
    end_at BIGINT,
    event_id INT,
    status VARCHAR(255) DEFAULT 'not_yet_started'::VARCHAR,
    league SMALLINT
);

CREATE INDEX if NOT EXISTS idx_matches
    ON matches (id, event_id);

CREATE UNIQUE INDEX if NOT EXISTS uniq_idx_matches
    ON matches (id);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS matches;
DROP INDEX IF EXISTS idx_matches;
DROP INDEX IF EXISTS uniq_idx_matches;
-- +goose StatementEnd
