package handler

import (
	"log"
	"net/http"
	"strings"

	"codemen.org/web/storage"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)


func (h Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	h.pareseCreateUserTemplate(w, UserForm{
		CSRFToken: nosurf.Token(r),
	})
}

func (h Handler) StoreUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	form := UserForm{}
	user := storage.User{}
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.User = user
	if err := user.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			for key, val := range vErr {
				form.FormError[strings.Title(key)] = val
			}
		}
		h.pareseCreateUserTemplate(w, form)
		return
	}

	_, err := h.storage.CreateUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (h Handler) pareseCreateUserTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("create-user.html")
	if t == nil {
		log.Println("unable to lookup create-user template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
