-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL,
    class_id BIGSERIAL,
    student_role Integer,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL,

	PRIMARY KEY(id),
	FOREIGN KEY(class_id)
      REFERENCES classes(id),
	UNIQUE(class_id, student_role)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students;
-- +goose StatementEnd
