package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"codemen.org/web/storage"
	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/justinas/nosurf"
)

func (h Handler) EditUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	editUser, err := h.storage.GetUserByID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	var form UserForm
	form.User = *editUser
	form.CSRFToken = nosurf.Token(r)
	h.pareseEditUserTemplate(w, form)
}

func (h Handler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	uID, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	var form UserForm
	user := storage.User{ID: uID}
	if err := h.decoder.Decode(&user, r.PostForm); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	form.User = user
	if err := user.Validate(); err != nil {
		if vErr, ok := err.(validation.Errors); ok {
			newErrs := make(map[string]error)
			for key, val := range vErr {
				newErrs[strings.Title(key)] = val
			}
			form.FormError = newErrs
		}
		h.pareseEditUserTemplate(w, form)
		return
	}

	updateUser, err := h.storage.UpdateUser(user)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	http.Redirect(w, r, fmt.Sprintf("/users/%v/edit", updateUser.ID), http.StatusSeeOther)
}

func (h Handler) pareseEditUserTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("edit-user.html")
	if t == nil {
		log.Println("unaable to lookup edit-user template")
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

}
