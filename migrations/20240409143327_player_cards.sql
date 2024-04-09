-- +goose Up
-- +goose StatementBegin
CREATE TABLE player_cards
(
    id           SERIAL PRIMARY KEY,
    profile_id   UUID REFERENCES user_profile (id) ON DELETE CASCADE,
    player_id    INTEGER REFERENCES players (id) ON DELETE CASCADE,
    rarity       SMALLINT,
    multiply     NUMERIC(3, 2),
    bonus_metric SMALLINT,
    unpacked     BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS player_cards;
-- +goose StatementEnd
