-- +goose Up
-- +goose StatementBegin
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
DROP CONSTRAINT fk_team_id;
DROP CONSTRAINT fk_default_project_id;
-- +goose StatementEnd
