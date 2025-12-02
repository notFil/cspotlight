-- +goose Up
-- +goose StatementBegin
ALTER TABLE csp_reports RENAME COLUMN blocked_uri TO blocked_url;
ALTER TABLE csp_reports RENAME COLUMN document_uri TO document_url;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE csp_reports RENAME COLUMN blocked_url TO blocked_uri;
ALTER TABLE csp_reports RENAME COLUMN document_url TO document_uri;
-- +goose StatementEnd
