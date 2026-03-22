package sqlconnect

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"restapi/internal/models"
	"restapi/pkg/utils"
)

// ErrStudentNotFound indicates that a student record does not exist.
var ErrStudentNotFound = errors.New("student not found")

// ListStudents retrieves students applying optional filters and sort order.
func ListStudents(filters map[string]string, sorts []SortOption) ([]models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect listing students to database")
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class FROM students WHERE 1=1"
	args := make([]interface{}, 0, len(filters))
	for field, value := range filters {
		query += fmt.Sprintf(" AND %s = ?", field)
		args = append(args, value)
	}

	if len(sorts) > 0 {
		parts := make([]string, 0, len(sorts))
		for _, sort := range sorts {
			parts = append(parts, fmt.Sprintf("%s %s", sort.Field, sort.Order))
		}
		query += " ORDER BY " + strings.Join(parts, ", ")
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to fetch student list")
	}
	defer rows.Close()

	students := make([]models.Student, 0)
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(&s.ID, &s.FirstName, &s.LastName, &s.Email, &s.Class); err != nil {
			return nil, utils.ErrorHandler(err, "unable to parse student row")
		}
		students = append(students, s)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorHandler(err, "error finishing student rows iteration")
	}
	return students, nil
}

// InsertStudents stores the provided students and returns them with IDs populated.
func InsertStudents(students []models.Student) ([]models.Student, error) {
	if len(students) == 0 {
		return []models.Student{}, nil
	}

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect inserting students to database")
	}
	defer db.Close()

	query := GenerateInsertQuery(models.Student{})
	if query == "" {
		return nil, utils.ErrorHandler(errors.New("invalid insert query"), "unable to build student insert query")
	}

	stmt, err := db.Prepare(query)
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to prepare student insert statement")
	}
	defer stmt.Close()

	inserted := make([]models.Student, 0, len(students))
	for _, s := range students {
		values := GenerateStructValues(s)
		res, err := stmt.Exec(values...)
		if err != nil {
			return nil, utils.ErrorHandler(err, "unable to execute student insert")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "unable to read inserted student id")
		}
		s.ID = int(lastID)
		inserted = append(inserted, s)
	}

	return inserted, nil
}

// GetStudentByID fetches a student by ID.
func GetStudentByID(id int) (*models.Student, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect selecting student to database")
	}
	defer db.Close()

	var student models.Student
	err = db.QueryRow("SELECT id, first_name, last_name, email, class FROM students WHERE id = ?", id).
		Scan(&student.ID, &student.FirstName, &student.LastName, &student.Email, &student.Class)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrStudentNotFound
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to scan student")
	}
	return &student, nil
}

// UpdateStudent replaces all student fields.
func UpdateStudent(student models.Student) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect updating student to database")
	}
	defer db.Close()

	result, err := db.Exec(
		"UPDATE students SET first_name = ?, last_name = ?, email = ?, class = ? WHERE id = ?",
		student.FirstName, student.LastName, student.Email, student.Class, student.ID,
	)
	if err != nil {
		return utils.ErrorHandler(err, "unable to execute student update")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "unable to check student update result")
	}
	if rowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}

// DeleteStudent removes a student by ID.
func DeleteStudent(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect deleting student to database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM students WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "unable to execute student delete")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "unable to check student delete result")
	}
	if rowsAffected == 0 {
		return ErrStudentNotFound
	}
	return nil
}

// CountStudents returns the total number of students matching filters.
func CountStudents(filters map[string]string) (int, error) {
	db, err := ConnectDB()
	if err != nil {
		return 0, utils.ErrorHandler(err, "unable to connect counting students to database")
	}
	defer db.Close()

	query := "SELECT COUNT(*) FROM students WHERE 1=1"
	args := make([]interface{}, 0, len(filters))
	for field, value := range filters {
		query += fmt.Sprintf(" AND %s = ?", field)
		args = append(args, value)
	}

	var count int
	if err := db.QueryRow(query, args...).Scan(&count); err != nil {
		return 0, utils.ErrorHandler(err, "unable to count students")
	}

	return count, nil
}
