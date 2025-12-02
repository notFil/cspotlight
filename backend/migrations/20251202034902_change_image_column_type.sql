-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN image TYPE text;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN image TYPE VARCHAR(255);
-- +goose StatementEnd
