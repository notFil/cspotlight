-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN default_project_id DROP NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN default_project_id SET NOT NULL;
-- +goose StatementEnd
