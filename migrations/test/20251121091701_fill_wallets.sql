-- +goose Up
-- +goose StatementBegin
INSERT INTO app.wallets (id, balance) VALUES
('3f9a1b9e-2f64-4f42-9b4d-2d1c9a5ef901', 100),
('5d2c7e80-1a34-4b74-8cc2-9f0e4f3c2a12', 10),
('5d2c7e80-1a34-4b74-8cc2-9f0e4f3c2a13', 0),
('5d2c7e80-1a34-4b74-8cc2-9f0e4f3c2a14', 10000)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
TRUNCATE TABLE app.wallets;
-- +goose StatementEnd
