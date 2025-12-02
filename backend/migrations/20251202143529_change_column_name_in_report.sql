-- +goose Up
-- +goose StatementBegin
ALTER TABLE csp_reports 
RENAME COLUMN report_body TO body;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE csp_reports 
RENAME COLUMN body TO report_body;
-- +goose StatementEnd
