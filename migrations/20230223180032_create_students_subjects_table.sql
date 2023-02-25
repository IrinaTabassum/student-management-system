-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS students_subject (
    id BIGSERIAL,
    student_id BIGSERIAL,
    subject_id BIGSERIAL,
    marke Integer,
	created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	deleted_at TIMESTAMP DEFAULT NULL,

	PRIMARY KEY(id),
    UNIQUE(student_id, subject_id)


);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS students_subject;
-- +goose StatementEnd
