-- +goose Up
-- +goose StatementBegin
CREATE SCHEMA auth;
CREATE SCHEMA content;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP SCHEMA auth;
DROP SCHEMA content;
-- +goose StatementEnd
