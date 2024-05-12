-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS players_statistic (
    id SERIAL PRIMARY KEY,
    player_id INTEGER,
    match_id INTEGER,
    game_date TIMESTAMP,
    opponent VARCHAR(255),
    score VARCHAR(255),
    fantasy_points NUMERIC,
    status VARCHAR(255),
    goals INTEGER,
    assists INTEGER,
    shots INTEGER,
    pims INTEGER,
    hits INTEGER,
    saves INTEGER,
    missed_goals INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players_statistic;
-- +goose StatementEnd
