-- +goose Up
-- +goose StatementBegin
CREATE TABLE csp_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    age INTEGER,
    report_body JSONB NOT NULL,
    url VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255),
    source_ip VARCHAR(255),
    type VARCHAR(50) NOT NULL DEFAULT 'csp-violation',
    project_id UUID REFERENCES projects(id) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL
); 

CREATE INDEX idx_csp_reports_project_id ON csp_reports(project_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS csp_reports;

DROP INDEX IF EXISTS idx_csp_reports_project_id;
-- +goose StatementEnd
