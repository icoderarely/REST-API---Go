package sqlconnect

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"restapi/internal/models"
	"restapi/pkg/utils"
)

// ErrTeacherNotFound indicates that a teacher record does not exist.
var ErrTeacherNotFound = errors.New("teacher not found")

// SortOption defines a sortable column/order pair for teacher queries.
type SortOption struct {
	Field string
	Order string
}

// ListTeachers retrieves teachers applying optional filters and sort order.
func ListTeachers(filters map[string]string, sorts []SortOption) ([]models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect listing teachers to database")
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
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
		return nil, utils.ErrorHandler(err, "unable to fetch teacher list")
	}
	defer rows.Close()

	teachers := make([]models.Teacher, 0)
	for rows.Next() {
		var t models.Teacher
		if err := rows.Scan(&t.ID, &t.FirstName, &t.LastName, &t.Email, &t.Class, &t.Subject); err != nil {
			return nil, utils.ErrorHandler(err, "unable to parse teacher row")
		}
		teachers = append(teachers, t)
	}

	if err := rows.Err(); err != nil {
		return nil, utils.ErrorHandler(err, "error finishing teacher rows iteration")
	}
	return teachers, nil
}

// InsertTeachers stores the provided teachers and returns them with IDs populated.
func InsertTeachers(teachers []models.Teacher) ([]models.Teacher, error) {
	if len(teachers) == 0 {
		return []models.Teacher{}, nil
	}

	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect inserting teachers to database")
	}
	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to prepare teacher insert statement")
	}
	defer stmt.Close()

	inserted := make([]models.Teacher, 0, len(teachers))
	for _, t := range teachers {
		res, err := stmt.Exec(t.FirstName, t.LastName, t.Email, t.Class, t.Subject)
		if err != nil {
			return nil, utils.ErrorHandler(err, "unable to execute teacher insert")
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			return nil, utils.ErrorHandler(err, "unable to read inserted teacher id")
		}
		t.ID = int(lastID)
		inserted = append(inserted, t)
	}

	return inserted, nil
}

// GetTeacherByID fetches a teacher by ID.
func GetTeacherByID(id int) (*models.Teacher, error) {
	db, err := ConnectDB()
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to connect selecting teacher to database")
	}
	defer db.Close()

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).
		Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrTeacherNotFound
	}
	if err != nil {
		return nil, utils.ErrorHandler(err, "unable to scan teacher")
	}
	return &teacher, nil
}

// UpdateTeacher replaces all teacher fields.
func UpdateTeacher(teacher models.Teacher) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect updating teacher to database")
	}
	defer db.Close()

	result, err := db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		teacher.FirstName, teacher.LastName, teacher.Email, teacher.Class, teacher.Subject, teacher.ID,
	)
	if err != nil {
		return utils.ErrorHandler(err, "unable to execute teacher update")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "unable to check teacher update result")
	}
	if rowsAffected == 0 {
		return ErrTeacherNotFound
	}
	return nil
}

// DeleteTeacher removes a teacher by ID.
func DeleteTeacher(id int) error {
	db, err := ConnectDB()
	if err != nil {
		return utils.ErrorHandler(err, "unable to connect deleting teacher to database")
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		return utils.ErrorHandler(err, "unable to execute teacher delete")
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return utils.ErrorHandler(err, "unable to check teacher delete result")
	}
	if rowsAffected == 0 {
		return ErrTeacherNotFound
	}
	return nil
}
