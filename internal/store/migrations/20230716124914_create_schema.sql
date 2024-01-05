-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA auth;
CREATE SCHEMA content;
CREATE EXTENSION "pgcrypto";
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA auth;
DROP SCHEMA content;
-- +goose StatementEnd
