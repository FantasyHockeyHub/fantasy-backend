-- +goose Up
-- +goose StatementBegin
CREATE TABLE if NOT EXISTS tournaments (
    id BIGINT PRIMARY KEY,
    league SMALLINT,
    title VARCHAR(255),
    matches_ids INTEGER[],
    started_at BIGINT,
    end_at BIGINT,
    players_amount INTEGER DEFAULT 0::INTEGER,
    deposit INTEGER DEFAULT 100::INTEGER,
    prize_fond INTEGER DEFAULT 0::INTEGER,
    status_tournament VARCHAR(255) DEFAULT 'not_yet_started'::VARCHAR
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tournaments
-- +goose StatementEnd
