-- +goose Up
-- +goose StatementBegin
CREATE TABLE players
(
    id             SERIAL PRIMARY KEY,
    api_id         INTEGER,
    position       SMALLINT,
    name           VARCHAR(255),
    team_id        INTEGER REFERENCES teams (team_id),
    sweater_number INTEGER,
    photo_link     VARCHAR(255),
    league         SMALLINT,
    player_cost    INTEGER DEFAULT 0:: INTEGER
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS players;
-- +goose StatementEnd
