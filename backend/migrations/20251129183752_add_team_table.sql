-- +goose Up
-- +goose StatementBegin
CREATE TABLE teams (
	id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
	name VARCHAR(100) NOT NULL,
  description TEXT,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE teams;
-- +goose StatementEnd
