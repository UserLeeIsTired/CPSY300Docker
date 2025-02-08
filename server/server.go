package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"
)

func PostUser(w http.ResponseWriter, r *http.Request, db *Database) {
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
		"message": "User created successfully",
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
