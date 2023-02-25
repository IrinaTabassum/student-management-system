-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS subjects (
    id BIGSERIAL,
    class_id BIGSERIAL,
    subject_name TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL,

	PRIMARY KEY(id),
	FOREIGN KEY(class_id)
      REFERENCES classes(id),
	UNIQUE(class_id, subject_name)  

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS subjects ;
-- +goose StatementEnd
