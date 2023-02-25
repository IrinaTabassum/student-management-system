package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func (h Handler) StudentResultView(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	studentResult, err := h.storage.ViewDetail(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	fmt.Println(studentResult)
	h.pareseResultVievTemplate(w, studentResult)
}

func (h Handler) pareseResultVievTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("view-result.html")
	if t == nil {
		log.Println("unable to lookup stuent-template template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
