-- +goose Up
-- +goose StatementBegin
CREATE TABLE coin_transactions
(
    id                  SERIAL PRIMARY KEY,
    profile_id          UUID REFERENCES user_profile (id) ON DELETE CASCADE,
    transaction_details VARCHAR(200) NOT NULL,
    amount              INTEGER,
    transaction_date    TIMESTAMP    NOT NULL,
    status              VARCHAR(30)  NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS coin_transactions CASCADE;
-- +goose StatementEnd
