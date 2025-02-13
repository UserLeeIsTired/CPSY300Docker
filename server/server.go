package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
)

func PostStudent(w http.ResponseWriter, r *http.Request, db *Database) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if len(student.StudentID) != 9 {
		http.Error(w, "Student ID must be 9 character long", http.StatusBadRequest)
		return
	}

	for _, char := range student.StudentID {
		if !unicode.IsDigit(char) {
			http.Error(w, "Student ID should only contain numbers", http.StatusBadRequest)
			return
		}
	}

	if strings.TrimSpace(student.StudentName) == "" {
		http.Error(w, "Student name should not be space or empty", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(student.Course) == "" {
		http.Error(w, "Student Course should not be space or empty", http.StatusBadRequest)
		return
	}

	_, err = db.CreateStudent(student.StudentID, student.StudentName, student.Course)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"message": "Student created successfully",
		"user": map[string]string{
			"Student ID":   student.StudentID,
			"Student Name": student.StudentName,
			"Course":       student.Course,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func GetAllStudents(w http.ResponseWriter, r *http.Request, db *Database) {
	students, err := db.GetAllStudents()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(students)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Students retrieved successfully",
		"users":   json.RawMessage(responseJSON),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func GetStudentByID(w http.ResponseWriter, r *http.Request, db *Database, id string) {
	student, err := db.GetStudentByID(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(student)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"message": "Student retrieved successfully",
		"users":   json.RawMessage(responseJSON),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
