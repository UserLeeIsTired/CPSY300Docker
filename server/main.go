package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func main() {
	db, err := NewDatabase()
	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer db.Close()

	r := chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "X-CSRF-TOKEN"},
		ExposedHeaders:   []string{"Link", "X-CSRF-TOKEN"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum age for the preflight OPTIONS request
	})
	r.Use(cors.Handler)

	redis := NewRedis()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!!")
	})

	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		Login(w, r, redis)
	})

	r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
		Logout(w, r, redis)
	})

	r.Post("/students", func(w http.ResponseWriter, r *http.Request) {
		PostStudent(w, r, db, redis)
	})

	r.Get("/students", func(w http.ResponseWriter, r *http.Request) {
		GetAllStudents(w, r, db, redis)
	})

	r.Get("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		GetStudentByID(w, r, db, redis, id)
	})

	r.Put("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		UpdateStudentByID(w, r, db, redis, id)
	})

	r.Delete("/students/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		DeleteStudentById(w, r, db, redis, id)
	})

	http.ListenAndServe("0.0.0.0:8080", r)
}
