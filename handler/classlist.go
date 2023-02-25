package handler

import (
	"log"
	"net/http"
)

func (h Handler) ClassList(w http.ResponseWriter, r *http.Request) {
	listClass, err := h.storage.ListOfClass()
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	h.pareseStudentListTemplate(w, listClass)
}

func (h Handler) pareseStudentListTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("class-list.html")
	if t == nil {
		log.Println("unable to lookup stuent-template template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
