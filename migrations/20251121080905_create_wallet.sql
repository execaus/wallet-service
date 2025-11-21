-- +goose Up
-- +goose StatementBegin

CREATE TABLE app.wallets (
    id UUID PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wallets;

-- +goose StatementEnd
