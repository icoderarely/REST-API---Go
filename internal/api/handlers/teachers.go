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

var teachers = make(map[int]models.Teacher)

// mutex    = &sync.Mutex{}
// nextID = 1

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		updateTeacherHandler(w, r)
	case http.MethodPatch:
		patchTeacherHandler(w, r)
	case http.MethodDelete:
		deleteTeacherHandler(w, r)
	}
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortFiled(field string) bool {
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
	if len(sortParams) > 0 {
		query += " ORDER BY"
		for i, param := range sortParams {
			parts := strings.Split(param, ":")
			if len(parts) != 2 {
				continue
			}
			field, order := parts[0], parts[1]
			if !isValidSortFiled(field) || !isValidSortOrder(order) {
				continue
			}
			if i > 0 {
				query += ","
			}
			query += " " + field + " " + order
		}
	}
	return query
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
		var args []interface{}

		query, args = addFilters(r, query, args)

		query = addSort(r, query)

		rows, err := db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "db query error", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		teacherList := make([]models.Teacher, 0, len(teachers))

		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
			if err != nil {
				http.Error(w, "Error scanning db results", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	id, _ := strconv.Atoi(idStr)

	var teacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
	if err == sql.ErrNoRows {
		http.Error(w, "Teacher Not Found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
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
		value := r.URL.Query().Get(param)
		if value != "" {
			query += " AND " + dbField + " = ?"
			args = append(args, value)
		}
	}
	return query, args
}

func addTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to database", http.StatusInternalServerError)
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
		http.Error(w, "Error preparing SQL query", http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	addedTeachers := make([]models.Teacher, len(newTeachers))
	for i, newTeacher := range newTeachers {
		res, err := stmt.Exec(newTeacher.FirstName, newTeacher.LastName, newTeacher.Email, newTeacher.Class, newTeacher.Subject)
		if err != nil {
			http.Error(w, "Error inserting data into db", http.StatusInternalServerError)
			return
		}
		lastID, err := res.LastInsertId()
		if err != nil {
			http.Error(w, "Error getting the last id", http.StatusInternalServerError)
			return
		}
		newTeacher.ID = int(lastID)
		addedTeachers[i] = newTeacher
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	resp := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	json.NewEncoder(w).Encode(resp)
}

func updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println(err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updatedTeacher models.Teacher
	_ = json.NewDecoder(r.Body).Decode(&updatedTeacher)

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	updatedTeacher.ID = id
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, updatedTeacher.ID)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedTeacher)
}

func patchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
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

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		log.Println(err)
		http.Error(w, "Unable to connect to db", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(&existingTeacher.ID, &existingTeacher.FirstName, &existingTeacher.LastName, &existingTeacher.Email, &existingTeacher.Class, &existingTeacher.Subject)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Unable to retrieve data", http.StatusInternalServerError)
		return
	}

	// apply updates using reflect
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updates {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k+",omitempty" {
				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}
			}
		}
	}

	existingTeacher.ID = id
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?", existingTeacher.FirstName, existingTeacher.LastName, existingTeacher.Email, existingTeacher.Class, existingTeacher.Subject, existingTeacher.ID)
	if err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(existingTeacher)
}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	db, err := sqlconnect.ConnectDB()
	if err != nil {
		http.Error(w, "Error connecting to DB", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error retrieving data", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Error deleting teacher", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	resp := struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Teacher successfully deleted",
		ID:     id,
	}

	json.NewEncoder(w).Encode(resp)
}
