package handler

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"codemen.org/web/storage"
	"github.com/go-chi/chi"
	"github.com/justinas/nosurf"
)

type StudentSubjectData struct {
	Student    storage.Student
	AllSubject []storage.Subject
	CSRFToken  string
}
type StudentSubjectForm struct {
	SubjectMarks map[int]int
}

func (h Handler) InsertResult(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	indStudent, err := h.storage.GetStudentByID(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	tempdata := StudentSubjectData{}
	tempdata.Student = *indStudent

	listSubject, err := h.storage.ListOfSubjects(indStudent.Class_ID)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	tempdata.AllSubject = listSubject
	tempdata.CSRFToken = nosurf.Token(r)
	h.pareseInsertResultTemplate(w, tempdata)
}
func (h Handler) StoreResult(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
	var sb StudentSubjectForm

	err := h.decoder.Decode(&sb, r.PostForm)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	id := chi.URLParam(r, "id")
	Sid, err := strconv.Atoi(id)
	if err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	studentsubje := storage.StudentSubject{}
	for key, val := range sb.SubjectMarks {
		studentsubje.Student_id = Sid
		studentsubje.Subject_id = key
		studentsubje.Marke = val
		fmt.Println(studentsubje)
		if _, err := h.storage.CreateStudentSubject(studentsubje); err != nil {
			log.Println(err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		
	}
	http.Redirect(w, r, "/users/student-list", http.StatusSeeOther)

}
func (h Handler) pareseInsertResultTemplate(w http.ResponseWriter, data any) {
	t := h.Templates.Lookup("insert-result.html")
	if t == nil {
		log.Println("unable to lookup stuent-template template")
		h.Error(w, "internal server error", http.StatusInternalServerError)
	}

	if err := t.Execute(w, data); err != nil {
		log.Println(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
