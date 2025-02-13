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

	http.ListenAndServe(":3000", r)

}
