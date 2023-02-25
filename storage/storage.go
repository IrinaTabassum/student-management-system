package storage

import (
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type UserFilter struct {
	SearchTerm string
	Offset     int
	Limit      int
}

type User struct {
	ID        int          `json:"id" form:"-" db:"id"`
	FirstName string       `json:"first_name" db:"first_name"`
	LastName  string       `json:"last_name" db:"last_name"`
	Email     string       `json:"email" db:"email"`
	Username  string       `json:"username" db:"username"`
	Password  string       `json:"password" db:"password"`
	Status    bool         `json:"status" db:"status"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
	Total     int          `json:"-" db:"total"`
}

type Class struct {
	ID        int          ` form:"-" db:"id"`
	ClassName string       `db:"class_name"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type Subject struct {
	ID          int          `json:"id" form:"-" db:"id"`
	Class_ID    int          `json:"class_id" db:"class_id"`
	SubjectName string       `json:"subject_name" db:"subject_name"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

type Student struct {
	ID        int          `form:"-" db:"id"`
	Class_ID  int          `db:"class_id"`
	Role      int          `db:"student_role"`
	FirstName string       `db:"first_name"`
	LastName  string       `db:"last_name"`
	TotalMark int          `db:"_"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt time.Time    `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}
type StudentSubject struct {
	ID         int          `form:"-" db:"id"`
	Student_id int          `form:"-" db:"student_id"`
	Subject_id int          `form:"-" db:"subject_id"`
	Marke      int          `form:"-" db:"marke"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  time.Time    `db:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
}

type StudentResult struct {
	Sid         int    `db:"id"`
	SFirst_Name string `db:"first_name"`
	Sub_Name    string `db:"subject_name"`
	Mark        int    `db:"marke"`
}
type AllStudentInAClass struct {
	Sid         int    `db:"id"`
	SFirst_Name string `db:"first_name"`
	Total       int    `db:"total"`
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.FirstName,
			validation.Required.Error("The first name field is required."),
		),
		validation.Field(&u.LastName,
			validation.Required.Error("The last name field is required."),
		),
		validation.Field(&u.Username,
			validation.Required.When(u.ID == 0).Error("The username field is required."),
		),
		validation.Field(&u.Email,
			validation.Required.When(u.ID == 0).Error("The email field is required."),
			is.Email.Error("The email field must be a valid email."),
		),
		validation.Field(&u.Password,
			validation.Required.When(u.ID == 0).Error("The password field is required."),
		),
	)
}

func (c Class) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.ClassName,
			validation.Required.Error("The first name field is required."),
		),
	)
}

func (s Subject) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Class_ID,
			validation.Required.Error("The first name field is required."),
		),
		validation.Field(&s.SubjectName,
			validation.Required.Error("The first name field is required."),
		),
	)
}
func (st Student) Validate() error {
	return validation.ValidateStruct(&st,
		validation.Field(&st.Class_ID,
			validation.Required.Error("The first name field is required."),
		),
		validation.Field(&st.Role,
			validation.Required.Error("The first name field is required."),
		),
		validation.Field(&st.FirstName,
			validation.Required.Error("The first name field is required."),
		),
		validation.Field(&st.LastName,
			validation.Required.Error("The first name field is required."),
		),
	)
}
