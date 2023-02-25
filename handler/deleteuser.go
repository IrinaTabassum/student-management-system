package handler

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.storage.DeleteUserByID(id); err != nil {
		h.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}
