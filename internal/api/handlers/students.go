package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
)

var filterableStudentFields = map[string]string{
	"first_name": "first_name",
	"last_name":  "last_name",
	"email":      "email",
	"class":      "class",
}

func getStudentIDParam(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func buildStudentFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)
	for param, column := range filterableStudentFields {
		if value := r.URL.Query().Get(param); value != "" {
			filters[column] = value
		}
	}
	return filters
}

func applyStudentUpdates(student *models.Student, updates map[string]interface{}) error {
	setters := map[string]func(*models.Student, string){
		"first_name": func(s *models.Student, v string) { s.FirstName = v },
		"last_name":  func(s *models.Student, v string) { s.LastName = v },
		"email":      func(s *models.Student, v string) { s.Email = v },
		"class":      func(s *models.Student, v string) { s.Class = v },
	}
	for field, raw := range updates {
		setter, ok := setters[field]
		if !ok {
			return fmt.Errorf("invalid field: %s", field)
		}
		value, ok := raw.(string)
		if !ok {
			return fmt.Errorf("invalid value for %s", field)
		}
		setter(student, value)
	}
	return nil
}

func GetStudentsHandler(w http.ResponseWriter, r *http.Request) {
	studentList, err := sqlconnect.ListStudents(buildStudentFilters(r), sqlconnect.BuildSortOptions(r, filterableStudentFields))
	if err != nil {
		http.Error(w, "Unable to load students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(studentList),
		Data:   studentList,
	})
}

func AddStudentsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var newStudents []models.Student
	if err := json.NewDecoder(r.Body).Decode(&newStudents); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	addedStudents, err := sqlconnect.InsertStudents(newStudents)
	if err != nil {
		http.Error(w, "Unable to add students", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Student `json:"data"`
	}{
		Status: "success",
		Count:  len(addedStudents),
		Data:   addedStudents,
	})
}

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	id, err := getStudentIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	student, err := sqlconnect.GetStudentByID(id)
	if errors.Is(err, sqlconnect.ErrStudentNotFound) {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Database query error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func UpdateStudentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getStudentIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var updated models.Student
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updated.ID = id

	if err := sqlconnect.UpdateStudent(updated); errors.Is(err, sqlconnect.ErrStudentNotFound) {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error updating student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func PatchStudentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getStudentIDParam(r)
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

	student, err := sqlconnect.GetStudentByID(id)
	if errors.Is(err, sqlconnect.ErrStudentNotFound) {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve student", http.StatusInternalServerError)
		return
	}

	if err := applyStudentUpdates(student, updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	student.ID = id
	if err := sqlconnect.UpdateStudent(*student); err != nil {
		http.Error(w, "Error updating student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

func DeleteStudentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getStudentIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := sqlconnect.DeleteStudent(id); errors.Is(err, sqlconnect.ErrStudentNotFound) {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error deleting student", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Student successfully deleted",
		ID:     id,
	})
}
