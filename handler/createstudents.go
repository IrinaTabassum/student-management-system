package handler

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"codemen.org/web/storage"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

type StudentForm struct {
	ListClass []storage.Class
	Student   storage.Student
	FormError map[string]error
	CSRFToken string
}

func (h Handler) CreateStudent(w http.ResponseWriter, r *http.Request) {
	listClass, err := h.storage.ListOfClass()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseCreateStudentTemplate(w, StudentForm{
		ListClass: listClass,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreStuents(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := StudentForm{}
	student := storage.Student{}
	if err := h.decoder.Decode(&student, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.Student = student
	if err := student.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			formErr := make(map[string]error)
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
			form.FormError = formErr
			form.CSRFToken = nosurf.Token(r)
			fmt.Println(form.FormError)
			h.pareseCreateStudentTemplate(w, form)
			return
		}
	}
	fmt.Println(student)
	if _, err := h.storage.CreateStudent(student); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/users/create-student", http.StatusSeeOther)
}

func (h Handler) pareseCreateStudentTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-student.html")
	if t == nil {
		log.Println("unable to lookup create-user template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
