package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (h Handler) AllStudentInClass(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	listStudentbyClass, err := h.storage.ListStudentbyClass(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	fmt.Println(listStudentbyClass)

	h.pareseStudentListnyClassTemplate(w, listStudentbyClass)
}

func (h Handler) pareseStudentListnyClassTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("student-inclass.html")
	if t == nil {
		log.Println("unable to lookup stuent-template template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
