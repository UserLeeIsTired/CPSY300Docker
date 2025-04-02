package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

func CreateSessionHeader(w http.ResponseWriter, r *http.Request, redis *Redis, isRemove bool) error {
	if isRemove {
		redis.DeleteKey(r.Header.Get("X-CSRF-TOKEN"))
		w.Header().Del("X-CSRF-TOKEN")
	} else {
		contentToken, err := redis.CreateKeyWithExpiration()

		if err != nil {
			return err
		}

		w.Header().Set("X-CSRF-TOKEN", contentToken)
	}
	return nil
}

func CheckSessionHeader(r *http.Request, redis *Redis) (string, error) {
	contentToken := r.Header.Get("X-CSRF-TOKEN")

	value, err := redis.GetValueByKey(contentToken)

	if err != nil {
		return "", errors.New("Unauthorized")
	}

	return value, nil
}

func Login(w http.ResponseWriter, r *http.Request, redis *Redis) {
	err := CreateSessionHeader(w, r, redis, false)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Login successfully")
}

func Logout(w http.ResponseWriter, r *http.Request, redis *Redis) {
	err := CreateSessionHeader(w, r, redis, true)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Logout successfully")
}

func PostStudent(w http.ResponseWriter, r *http.Request, db *Database, redis *Redis) {

	_, err := CheckSessionHeader(r, redis)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var student Student
	err = json.NewDecoder(r.Body).Decode(&student)

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
		http.Error(w, err.Error(), http.StatusConflict)
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

func GetAllStudents(w http.ResponseWriter, r *http.Request, db *Database, redis *Redis) {

	_, err := CheckSessionHeader(r, redis)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	students, err := db.GetAllStudents()

	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(students)
	if err != nil {
		fmt.Println(err)
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

func GetStudentByID(w http.ResponseWriter, r *http.Request, db *Database, redis *Redis, id string) {

	_, err := CheckSessionHeader(r, redis)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

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

func UpdateStudentByID(w http.ResponseWriter, r *http.Request, db *Database, redis *Redis, id string) {

	_, err := CheckSessionHeader(r, redis)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	var student Student

	err = json.NewDecoder(r.Body).Decode(&student)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if student.StudentID != "" && len(student.StudentID) != 9 {
		http.Error(w, "Student ID must be 9 character long", http.StatusBadRequest)
		return
	}

	if student.StudentID != "" {
		for _, char := range student.StudentID {
			if !unicode.IsDigit(char) {
				http.Error(w, "Student ID should only contain numbers", http.StatusBadRequest)
				return
			}
		}
	}

	err = db.UpdateStudentById(id, student.StudentID, student.StudentName, student.Course)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Ok")
}

func DeleteStudentById(w http.ResponseWriter, r *http.Request, db *Database, redis *Redis, id string) {

	_, err := CheckSessionHeader(r, redis)

	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	err = db.DeleteStudentById(id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode("Ok")
}
