package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
)

// ─── Helpers ─────────────────────────────────────────────────────────────────

func connectDB(w http.ResponseWriter) (*sql.DB, error) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return nil, err
	}
	return db, nil
}

func getIDParam(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	validFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"email":      true,
		"class":      true,
		"subject":    true,
	}
	return validFields[field]
}

func addSort(r *http.Request, query string) string {
	sortParams := r.URL.Query()["sortBy"]
	if len(sortParams) == 0 {
		return query
	}
	query += " ORDER BY"
	for i, param := range sortParams {
		parts := strings.Split(param, ":")
		if len(parts) != 2 {
			continue
		}
		field, order := parts[0], parts[1]
		if !isValidSortField(field) || !isValidSortOrder(order) {
			continue
		}
		if i > 0 {
			query += ","
		}
		query += " " + field + " " + order
	}
	return query
}

func addFilters(r *http.Request, query string, args []interface{}) (string, []interface{}) {
	params := map[string]string{
		"first_name": "first_name",
		"last_name":  "last_name",
		"email":      "email",
		"class":      "class",
		"subject":    "subject",
	}
	for param, dbField := range params {
		if value := r.URL.Query().Get(param); value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

func scanTeacher(row *sql.Row, teacher *models.Teacher) error {
	return row.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
}

// ─── Collection Handlers (no ID) ─────────────────────────────────────────────

// GET /teachers/
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB(w)
	if err != nil {
		return
	}
	defer db.Close()

	query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
	var args []interface{}
	query, args = addFilters(r, query, args)
	query = addSort(r, query)

	rows, err := db.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	teacherList := make([]models.Teacher, 0)
	for rows.Next() {
		var teacher models.Teacher
		if err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject); err != nil {
			http.Error(w, "Error scanning database results", http.StatusInternalServerError)
			return
		}
		teacherList = append(teacherList, teacher)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teacherList),
		Data:   teacherList,
	})
}

// POST /teachers/
func AddTeachersHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	db, err := connectDB(w)
	if err != nil {
		return
	}
	defer db.Close()

	var newTeachers []models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&newTeachers); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	stmt, err := db.Prepare("INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		http.Error(w, "Error preparing SQL statement", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, t := range newTeachers {
		res, err := stmt.Exec(t.FirstName, t.LastName, t.Email, t.Class, t.Subject)
		if err != nil {
			http.Error(w, "Error inserting teacher", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error retrieving last insert ID", http.StatusInternalServerError)
			return
		}
		t.ID = int(lastID)
		addedTeachers = append(addedTeachers, t)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	})
}

func PatchTeachersHandler(w http.ResponseWriter, r *http.Request) {}

func DeleteTeachersHandler(w http.ResponseWriter, r *http.Request) {}

// ─── Single Resource Handlers (with ID) ──────────────────────────────────────

// GET /teachers/{id}
func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB(w)
	if err != nil {
		return
	}
	defer db.Close()

	var teacher models.Teacher
	err = scanTeacher(
		db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id),
		&teacher,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// PUT /teachers/{id}
func UpdateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB(w)
	if err != nil {
		return
	}
	defer db.Close()

	var existing models.Teacher
	err = scanTeacher(
		db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id),
		&existing,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve teacher", http.StatusInternalServerError)
		return
	}

	var updated models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updated.ID = id

	_, err = db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		updated.FirstName, updated.LastName, updated.Email, updated.Class, updated.Subject, updated.ID,
	)
	if err != nil {
		log.Println(err)
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

// PATCH /teachers/{id}
func PatchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	db, err := connectDB(w)
	if err != nil {
		return
	}
	defer db.Close()

	var existing models.Teacher
	err = scanTeacher(
		db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id),
		&existing,
	)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve teacher", http.StatusInternalServerError)
		return
	}

	teacherVal := reflect.ValueOf(&existing).Elem()
	teacherType := teacherVal.Type()
	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" && teacherVal.Field(i).CanSet() {
				teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
			}
		}
	}

	existing.ID = id
	_, err = db.Exec(
		"UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existing.FirstName, existing.LastName, existing.Email, existing.Class, existing.Subject, existing.ID,
	)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existing)
}

// DELETE /teachers/{id}
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := connectDB(w)
	if err != nil {
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving rows affected", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher successfully deleted",
		ID:     id,
	})
}
