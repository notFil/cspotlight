-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS timescaledb;

CREATE TABLE csp_reports (
    id UUID DEFAULT gen_random_uuid(),
    age INTEGER,
    body JSONB NOT NULL,
    url VARCHAR(255) NOT NULL,
    user_agent VARCHAR(255),
    source_ip VARCHAR(255),
    type VARCHAR(50) NOT NULL DEFAULT 'csp-violation',
    project_id UUID,
    directive VARCHAR(50),
    blocked_url TEXT,
    document_url TEXT,
    disposition VARCHAR(20),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    PRIMARY KEY (id, created_at)
); 

CREATE INDEX idx_csp_reports_project_id ON csp_reports(project_id);
CREATE INDEX idx_csp_reports_url ON csp_reports(url);
CREATE INDEX idx_csp_reports_directive ON csp_reports(directive);
CREATE INDEX idx_csp_reports_disposition ON csp_reports(disposition);
CREATE INDEX idx_csp_reports_created_at ON csp_reports(created_at DESC);

SELECT create_hypertable('csp_reports', 'created_at');

ALTER TABLE csp_reports SET (
  timescaledb.compress,
  timescaledb.compress_segmentby = 'directive',
  timescaledb.compress_orderby = 'created_at DESC'
);

SELECT add_compression_policy('csp_reports', INTERVAL '30 days');

ALTER TABLE csp_reports
  ADD CONSTRAINT csp_reports_project_id_fkey
  FOREIGN KEY (project_id)
  REFERENCES projects(id)
  ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS csp_reports;
-- +goose StatementEnd
