package handler

import (
	"log"
	"net/http"
)

func (h Handler) StudentsList(w http.ResponseWriter, r *http.Request) {
	listStudent, err := h.storage.ListOfStudents()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseClassListTemplate(w, listStudent)
}

func (h Handler) pareseClassListTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("student-list.html")
	if t == nil {
		log.Println("unable to lookup stuent-template template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
