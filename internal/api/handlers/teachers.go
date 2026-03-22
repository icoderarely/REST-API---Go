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

// ─── Helpers ─────────────────────────────────────────────────────────────────

var filterableTeacherFields = map[string]string{
	"first_name": "first_name",
	"last_name":  "last_name",
	"email":      "email",
	"class":      "class",
	"subject":    "subject",
}

func getIDParam(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func isValidSortOrder(order string) bool {
	return order == "asc" || order == "desc"
}

func isValidSortField(field string) bool {
	return filterableTeacherFields[field] != ""
}

func buildTeacherFilters(r *http.Request) map[string]string {
	filters := make(map[string]string)
	for param, column := range filterableTeacherFields {
		if value := r.URL.Query().Get(param); value != "" {
			filters[column] = value
		}
	}
	return filters
}

func applyTeacherUpdates(teacher *models.Teacher, updates map[string]interface{}) error {
	setters := map[string]func(*models.Teacher, string){
		"first_name": func(t *models.Teacher, v string) { t.FirstName = v },
		"last_name":  func(t *models.Teacher, v string) { t.LastName = v },
		"email":      func(t *models.Teacher, v string) { t.Email = v },
		"class":      func(t *models.Teacher, v string) { t.Class = v },
		"subject":    func(t *models.Teacher, v string) { t.Subject = v },
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
		setter(teacher, value)
	}
	return nil
}

// ─── Collection Handlers (no ID) ─────────────────────────────────────────────

// GET /teachers/
func GetTeachersHandler(w http.ResponseWriter, r *http.Request) {
	teacherList, err := sqlconnect.ListTeachers(buildTeacherFilters(r), sqlconnect.BuildSortOptions(r, filterableTeacherFields))
	if err != nil {
		http.Error(w, "Unable to load teachers", http.StatusInternalServerError)
		return
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

	var newTeachers []models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&newTeachers); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	addedTeachers, err := sqlconnect.InsertTeachers(newTeachers)
	if err != nil {
		http.Error(w, "Unable to add teachers", http.StatusInternalServerError)
		return
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

	teacher, err := sqlconnect.GetTeacherByID(id)
	if errors.Is(err, sqlconnect.ErrTeacherNotFound) {
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

	var updated models.Teacher
	if err := json.NewDecoder(r.Body).Decode(&updated); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	updated.ID = id

	if err := sqlconnect.UpdateTeacher(updated); errors.Is(err, sqlconnect.ErrTeacherNotFound) {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
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

	teacher, err := sqlconnect.GetTeacherByID(id)
	if errors.Is(err, sqlconnect.ErrTeacherNotFound) {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Unable to retrieve teacher", http.StatusInternalServerError)
		return
	}

	if err := applyTeacherUpdates(teacher, updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	teacher.ID = id
	if err := sqlconnect.UpdateTeacher(*teacher); err != nil {
		http.Error(w, "Error updating teacher", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(teacher)
}

// DELETE /teachers/{id}
func DeleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	id, err := getIDParam(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := sqlconnect.DeleteTeacher(id); errors.Is(err, sqlconnect.ErrTeacherNotFound) {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error deleting teacher", http.StatusInternalServerError)
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
