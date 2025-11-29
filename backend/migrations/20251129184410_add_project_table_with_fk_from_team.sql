-- +goose Up
-- +goose StatementBegin
CREATE TABLE projects (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
  description VARCHAR(500),
  team_id UUID NOT NULL,
  disabled BOOLEAN DEFAULT FALSE,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);

ALTER TABLE projects
ADD CONSTRAINT fk_team_id 
  FOREIGN KEY (team_id) 
  REFERENCES teams(id)
  ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE projects;
DROP CONSTRAINT fk_team_id;
-- +goose StatementEnd
