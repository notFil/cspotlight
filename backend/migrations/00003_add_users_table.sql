-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(10) NOT NULL,
    default_project_id UUID,
    disabled BOOLEAN DEFAULT FALSE,
    team_id UUID,
    image TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE users
    ADD CONSTRAINT fk_team_id 
    FOREIGN KEY (team_id) 
    REFERENCES teams(id)
    ON DELETE CASCADE;

ALTER TABLE users
    ADD CONSTRAINT fk_default_project_id 
    FOREIGN KEY (default_project_id) 
    REFERENCES projects(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
