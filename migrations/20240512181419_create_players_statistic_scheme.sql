-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS players_statistic (
    id SERIAL PRIMARY KEY,
    player_id INTEGER NOT NULL,
    match_id INTEGER NOT NULL,
    game_date TIMESTAMP,
    opponent VARCHAR(255),
    fantasy_points NUMERIC(4,1) DEFAULT 0.0::NUMERIC,
    goals INTEGER DEFAULT 0::INTEGER,
    assists INTEGER DEFAULT 0::INTEGER,
    shots INTEGER DEFAULT 0::INTEGER,
    pims INTEGER DEFAULT 0::INTEGER,
    hits INTEGER DEFAULT 0::INTEGER,
    saves INTEGER DEFAULT 0::INTEGER,
    missed_goals INTEGER DEFAULT 0::INTEGER,
    shutout BOOLEAN DEFAULT false::BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players_statistic;
-- +goose StatementEnd
