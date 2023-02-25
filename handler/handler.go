package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"codemen.org/web/storage"
	"github.com/Masterminds/sprig"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/form"
)

type Handler struct {
	sessionManager *scs.SessionManager
	decoder        *form.Decoder
	storage        dbStorage
	Templates      *template.Template
}

type dbStorage interface {
	ListUser(storage.UserFilter) ([]storage.User, error)
	CreateUser(storage.User) (*storage.User, error)
	UpdateUser(storage.User) (*storage.User, error)
	GetUserByID(string) (*storage.User, error)
	GetUserByUsername(string) (*storage.User, error)
	DeleteUserByID(string) error
	CreateClass(storage.Class) error
	ListOfClass() ([]storage.Class, error)
	CreateSubject(storage.Subject) error
	CreateStudent(storage.Student) (*storage.Student, error)
	ListOfStudents() ([]storage.Student, error)
	GetStudentByID(string) (*storage.Student, error)
	ListOfSubjects(int) ([]storage.Subject, error)
	CreateStudentSubject(storage.StudentSubject) (*storage.StudentSubject, error)
	ListStudentbyClass(string) ([]storage.AllStudentInAClass, error)
	ViewDetail(string) ([]storage.StudentResult, error)
	GetClassByClassName(string) (*storage.Class, error)
}

type ErrorPage struct {
	Code    int
	Message string
}

func (h Handler) Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	ep := ErrorPage{
		Code:    code,
		Message: error,
	}

	tf := "default"
	switch code {
	case 400, 401, 402, 403, 404:
		tf = "4xx"
	case 500, 501, 503:
		tf = "5xx"
	}

	tpl := fmt.Sprintf("%s.html", tf)
	t := h.Templates.Lookup(tpl)
	if t == nil {
		log.Fatalln("unable to find template")
	}

	if err := t.Execute(w, ep); err != nil {
		log.Fatalln(err)
	}
}

func NewHandler(sm *scs.SessionManager, formDecoder *form.Decoder, storage dbStorage) *chi.Mux {
	h := &Handler{
		sessionManager: sm,
		decoder:        formDecoder,
		storage:        storage,
	}

	h.ParseTemplates()
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(Method)

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Get("/", h.Home)
		r.Get("/login", h.Login)
		r.Post("/login", h.LoginPostHandler)
		r.Get("/registration", h.CreateUser)
		r.Post("/store", h.StoreUser)
	})

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "assets/src"))
	r.Handle("/static/*", http.StripPrefix("/static", http.FileServer(filesDir)))

	r.Group(func(r chi.Router) {
		r.Use(sm.LoadAndSave)
		r.Use(h.Authentication)

		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.ListUser)

			r.Get("/create-class", h.CreateClass)

			r.Post("/store-class", h.StoreClass)

			r.Get("/create-subject", h.CreateSubject)

			r.Post("/store-subject", h.StoreSubject)

			r.Get("/create-student", h.CreateStudent)

			r.Post("/store-student", h.StoreStuents)

			r.Get("/student-list", h.StudentsList)

			r.Get("/insetr-result/{id:[0-9]+}", h.InsertResult)

			r.Post("/insetr-result/{id:[0-9]+}/store", h.StoreResult)

			r.Get("/class-list", h.ClassList)

			r.Get("/all-student/class-{id:[0-9]+}", h.AllStudentInClass)

			r.Get("/{id:[0-9]+}/result-detail", h.StudentResultView)

			r.Get("/{id:[0-9]+}/edit", h.EditUser)

			r.Put("/{id:[0-9]+}/update", h.UpdateUser)

			r.Get("/{id:[0-9]+}/delete", h.DeleteUser)
		})

		r.Get("/logout", h.LogoutHandler)
	})

	return r
}

func Method(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			switch strings.ToLower(r.PostFormValue("_method")) {
			case "put":
				r.Method = http.MethodPut
			case "patch":
				r.Method = http.MethodPatch
			case "delete":
				r.Method = http.MethodDelete
			default:
			}
		}
		next.ServeHTTP(w, r)
	})
}

func (h Handler) Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := h.sessionManager.GetString(r.Context(), "userID")
		uID, err := strconv.Atoi(userID)
		if err != nil {
			log.Println(err)
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		if uID <= 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) ParseTemplates() error {
	templates := template.New("web-templates").Funcs(template.FuncMap{
		"calculatePreviousPage": func(currentPageNumber int) int {
			if currentPageNumber == 1 {
				return 0
			}

			return currentPageNumber - 1
		},

		"calculateNextPage": func(currentPageNumber, totalPage int) int {
			if currentPageNumber == totalPage {
				return 0
			}

			return currentPageNumber + 1
		},
	}).Funcs(sprig.FuncMap())

	newFS := os.DirFS("assets/templates")
	tmpl := template.Must(templates.ParseFS(newFS, "*/*/*.html", "*/*.html", "*.html"))
	if tmpl == nil {
		log.Fatalln("unable to parse templates")
	}

	h.Templates = tmpl
	return nil
}
