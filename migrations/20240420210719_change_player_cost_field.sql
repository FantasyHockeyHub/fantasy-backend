-- +goose Up
-- +goose StatementBegin
ALTER TABLE players
ALTER COLUMN player_cost TYPE NUMERIC(4,1)
USING player_cost::NUMERIC(4,1);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE players
ALTER COLUMN player_cost TYPE INTEGER
USING player_cost::INTEGER;
-- +goose StatementEnd
