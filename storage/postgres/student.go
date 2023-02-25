package postgres

import (
	"fmt"
	"log"

	"codemen.org/web/storage"
)

const insertStudentQuery = `
		INSERT INTO students(
			class_id,
			student_role,
			first_name,
			last_name
		) VALUES (
			:class_id,
			:student_role,
			:first_name,
			:last_name
		) RETURNING *;
	`

func (s PostgresStorage) CreateStudent(st storage.Student) (*storage.Student, error) {

	stmt, err := s.DB.PrepareNamed(insertStudentQuery)
	if err != nil {
		log.Fatalln(err)
	}
	if err := stmt.Get(&st, st); err != nil {
		return nil, err
	}

	if st.ID == 0 {
		return nil, fmt.Errorf("unable to insert user into db")
	}
	return &st, nil
}

const listStudentQuery = `SELECT * FROM students WHERE deleted_at IS NULL;`

func (s PostgresStorage) ListOfStudents() ([]storage.Student, error) {
	var listStudent []storage.Student
	if err := s.DB.Select(&listStudent, listStudentQuery); err != nil {
		log.Println(err)
		return nil, err
	}
	return listStudent, nil
}

const getStudentByIDQuery = `SELECT * FROM students WHERE id=$1 AND deleted_at IS NULL`

func (s PostgresStorage) GetStudentByID(id string) (*storage.Student, error) {
	var st storage.Student
	if err := s.DB.Get(&st, getStudentByIDQuery, id); err != nil {
		log.Println(err)
		return nil, err
	}

	return &st, nil
}

const listSubjectuery = `SELECT * FROM subjects WHERE class_id=$1 AND deleted_at IS NULL;`

func (s PostgresStorage) ListOfSubjects(class_id int) ([]storage.Subject, error) {
	var listSubject []storage.Subject
	if err := s.DB.Select(&listSubject, listSubjectuery, class_id); err != nil {
		log.Println(err)
		return nil, err
	}
	return listSubject, nil
}

const insertStudentSubjectQuery = `
		INSERT INTO students_subject(
			student_id,
			subject_id,
			marke
		) VALUES (
			:student_id,
			:subject_id,
			:marke
		) ON CONFLICT(student_id, subject_id) DO UPDATE SET 
		student_id = EXCLUDED.student_id, 
		subject_id = EXCLUDED.subject_id, 
		marke = EXCLUDED.marke
		RETURNING *;
	`

func (s PostgresStorage) CreateStudentSubject(ss storage.StudentSubject) (*storage.StudentSubject, error) {

	stmt, err := s.DB.PrepareNamed(insertStudentSubjectQuery)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("preparemarks")
	if err := stmt.Get(&ss, ss); err != nil {
		return nil, err
	}

	if ss.ID == 0 {
		return nil, fmt.Errorf("unable to insert user into db")
	}
	return &ss, nil
}

const listStudentbyClassQuery = `SELECT students.id, students.first_name, sum(students_subject.marke) as total
FROM students
INNER JOIN students_subject
ON students.id = students_subject.student_id 
where students.class_id = $1 
group by students.id, students.first_name;`

func (s PostgresStorage) ListStudentbyClass(class_id string) ([]storage.AllStudentInAClass, error) {
	var liststudent []storage.AllStudentInAClass
	if err := s.DB.Select(&liststudent, listStudentbyClassQuery, class_id); err != nil {
		log.Println(err)
		return nil, err
	}
	return liststudent, nil
}

const viewStudentResultQuery = `SELECT students.id, students.first_name , subjects.subject_name, students_subject.marke
FROM students
INNER JOIN students_subject
ON students.id = students_subject.student_id
INNER JOIN subjects
ON students_subject.subject_id= subjects.id where students.id = $1;`

func (s PostgresStorage) ViewDetail(student_id string) ([]storage.StudentResult, error) {
	var viewstudents []storage.StudentResult
	if err := s.DB.Select(&viewstudents, viewStudentResultQuery, student_id); err != nil {
		log.Println(err)
		return nil, err
	}
	return viewstudents, nil
}
