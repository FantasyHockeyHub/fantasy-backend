-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_roster
(
    id              SERIAL PRIMARY KEY,
    tournament_id   BIGINT REFERENCES tournaments (id) ON DELETE CASCADE,
    user_id         UUID,
    roster          INTEGER[],
    current_balance NUMERIC(4, 1) DEFAULT 100.0,
    points          NUMERIC(4, 1) DEFAULT 0.0,
    coins           INTEGER DEFAULT 0,
    place           INTEGER DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_roster;
-- +goose StatementEnd
