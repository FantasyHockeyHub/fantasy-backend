-- +goose Up
-- +goose StatementBegin
CREATE TABLE fantasy_store
(
    id                 SERIAL PRIMARY KEY,
    product_name       VARCHAR(100) NOT NULL,
    price              INTEGER,
    league             SMALLINT,
    rarity             SMALLINT,
    player_cards_count INTEGER,
    photo_link         VARCHAR(255)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS fantasy_store;
-- +goose StatementEnd
