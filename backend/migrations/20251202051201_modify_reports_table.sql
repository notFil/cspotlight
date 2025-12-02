-- +goose Up
-- +goose StatementBegin
ALTER TABLE csp_reports 
ADD COLUMN directive VARCHAR(50),
ADD COLUMN blocked_uri TEXT,
ADD COLUMN document_uri TEXT,
ADD COLUMN disposition VARCHAR(20);

CREATE INDEX IF NOT EXISTS idx_csp_reports_url ON csp_reports(url);
CREATE INDEX IF NOT EXISTS idx_csp_reports_directive ON csp_reports(directive);
CREATE INDEX IF NOT EXISTS idx_csp_reports_project_id ON csp_reports(project_id);
CREATE INDEX IF NOT EXISTS idx_csp_reports_disposition ON csp_reports(disposition);
CREATE INDEX IF NOT EXISTS idx_csp_reports_created_at ON csp_reports(created_at DESC);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE csp_reports 
DROP COLUMN directive,
DROP COLUMN blocked_uri,
DROP COLUMN document_uri,
DROP COLUMN disposition;

DROP INDEX IF EXISTS idx_csp_reports_url;
DROP INDEX IF EXISTS idx_csp_reports_directive;
DROP INDEX IF EXISTS idx_csp_reports_project_id;
DROP INDEX IF EXISTS idx_csp_reports_disposition;
DROP INDEX IF EXISTS idx_csp_reports_created_at;
-- +goose StatementEnd
