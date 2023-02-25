package handler

import (
	"log"
	"math"
	"net/http"
	"strconv"

	"codemen.org/web/storage"
)

const userListLimit = 2

type UserList struct {
	Users       []storage.User
	SearchTerm  string
	CurrentPage int
	Limit       int
	Total       int
	TotalPage   int
}

type UserForm struct {
	User      storage.User
	FormError map[string]error
	CSRFToken string
}



func (h Handler) ListUser(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	var err error
	CurrentPage := 1
	pn := r.FormValue("page")
	if pn != "" {
		CurrentPage, err = strconv.Atoi(pn)
		if err != nil {
			CurrentPage = 1
		}
	}

	offset := 0
	if CurrentPage > 1 {
		offset = (CurrentPage * userListLimit) - userListLimit
	}

	st := r.FormValue("SearchTerm")
	uf := storage.UserFilter{
		SearchTerm: st,
		Offset:     offset,
		Limit:      userListLimit,
	}
	listUser, err := h.storage.ListUser(uf)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	t := h.Templates.Lookup("list-user.html")
	if t == nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	total := 0
	if len(listUser) > 0 {
		total = listUser[0].Total
	}

	totalPage := int(math.Ceil(float64(total) / float64(userListLimit)))

	data := UserList{
		Users:       listUser,
		SearchTerm:  st,
		CurrentPage: CurrentPage,
		Limit:       userListLimit,
		Total:       total,
		TotalPage:   totalPage,
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
	}
}
