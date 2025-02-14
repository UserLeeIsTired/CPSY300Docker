package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	db, err := NewDatabase()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.Close()

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!!")
	})

	r.Post("/students", func(w http.ResponseWriter, r *http.Request) {
		PostStudent(w, r, db)
	})

	r.Get("/students", func(w http.ResponseWriter, r *http.Request) {
		GetAllStudents(w, r, db)
	})

	r.Get("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		GetStudentByID(w, r, db, id)
	})

	http.ListenAndServe("0.0.0.0:8080", r)

}
