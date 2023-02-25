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

type ClassForm struct {
	Class     storage.Class
	FormError map[string]error
	CSRFToken string
}

func (h Handler) CreateClass(w http.ResponseWriter, r *http.Request) {
	h.pareseCreateClassTemplate(w, ClassForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreClass(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form := ClassForm{}
	class := storage.Class{}
	if err := h.decoder.Decode(&class, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.Class = class
	if err := class.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			formErr := make(map[string]error)
			for key, val := range vErr {
				formErr[strings.Title(key)] = val
			}

			form.FormError = formErr
			form.CSRFToken = nosurf.Token(r)
			fmt.Println(form.FormError)
			h.pareseCreateClassTemplate(w, form)
			return
		}
	}
	getclass, _ := h.storage.GetClassByClassName(class.ClassName)
	if getclass != nil {
		formErr := make(map[string]error)
		formErr["ClassName"] = fmt.Errorf("Class is axiset")
		form.FormError = formErr
		form.CSRFToken = nosurf.Token(r)
		h.pareseCreateClassTemplate(w, form)
		return
	}

	if err := h.storage.CreateClass(class); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/users/create-class", http.StatusSeeOther)
}

func (h Handler) pareseCreateClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-class.html")
	if t == nil {
		log.Println("unable to lookup create-user template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
