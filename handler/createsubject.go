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

type SubjectForm struct {
	ListClass []storage.Class
	Subject   storage.Subject
	FormError map[string]error
	CSRFToken string
}

func (h Handler) CreateSubject(w http.ResponseWriter, r *http.Request) {
	listClass, err := h.storage.ListOfClass()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseCreateSubjectTemplate(w, SubjectForm{
		ListClass: listClass,
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreSubject(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := SubjectForm{}
	subject := storage.Subject{}
	if err := h.decoder.Decode(&subject, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.Subject = subject
	if err := subject.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			formErr := make(map[string]error)
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}
			form.FormError = formErr
			form.CSRFToken = nosurf.Token(r)
			fmt.Println(form.FormError)
			h.pareseCreateSubjectTemplate(w, form)
			return
		}
	}

	if err := h.storage.CreateSubject(subject); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/users/create-subject", http.StatusSeeOther)
}

func (h Handler) pareseCreateSubjectTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-subject.html")
	if t == nil {
		log.Println("unable to lookup create-user template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
