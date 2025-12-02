-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN image VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN image;
-- +goose StatementEnd
