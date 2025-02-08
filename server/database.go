package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Student struct {
	StudentID   string
	StudentName string
	Course      string
	CreatedDate string
}

type Database struct {
	db *sql.DB
}

func NewDatabase() (*Database, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	var (
		host     = os.Getenv("HOST")
		port     = os.Getenv("PORT")
		user     = os.Getenv("USER")
		password = os.Getenv("PASSWORD")
		dbname   = os.Getenv("DBNAME")
	)

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	db, err := sql.Open("postgres", connectionString)

	return &Database{db: db}, nil
}

func (d *Database) Close() {
	if d != nil {
		d.db.Close()
	}
}

func (d *Database) CreateStudent(studentID string, studentName string, course string) (string, error) {
	stmt, err := d.db.Prepare("INSERT INTO my_student (id, name, course) VALUES ($1, $2, $3)")

	if err != nil {
		return "", err
	}

	defer stmt.Close()

	_, err = stmt.Exec(studentID, studentName, course)

	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("Student %s is created successfully\n", studentName)

	return result, nil
}

func (d *Database) GetAllStudents() ([]Student, error) {
	rows, err := d.db.Query("SELECT id, name, course, created_date FROM my_student")

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var students []Student

	for rows.Next() {
		var student Student
		if err := rows.Scan(&student.StudentID, &student.StudentName, &student.Course, &student.CreatedDate); err != nil {
			return nil, err
		}
		students = append(students, student)
	}

	return students, nil

}

func (d *Database) GetStudentByID(id string) (*Student, error) {
	stmt, err := d.db.Prepare("SELECT id, name, course, created_date FROM my_student WHERE id = $1")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var student Student

	row := stmt.QueryRow(id)
	if err := row.Scan(&student.StudentID, &student.StudentName, &student.Course, &student.CreatedDate); err != nil {
		return nil, err
	}

	return &student, nil
}
