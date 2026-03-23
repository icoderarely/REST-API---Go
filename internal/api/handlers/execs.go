package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
)

func GetExecsHandler(w http.ResponseWriter, _ *http.Request) {
	execs, err := sqlconnect.ListExecs()
	if err != nil {
		http.Error(w, "Unable to load execs", http.StatusInternalServerError)
		return
	}
	writeExecCollection(w, execs)
}

func AddExecsHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var execs []models.Exec
	if err := json.NewDecoder(r.Body).Decode(&execs); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	insertedExecs, err := sqlconnect.InsertExecs(execs)
	if err != nil {
		http.Error(w, "Unable to create execs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeExecCollection(w, insertedExecs)
}

func PatchExecsHandler(w http.ResponseWriter, _ *http.Request) {}

func GetExecHandler(w http.ResponseWriter, r *http.Request) {
	exec, err := loadExecFromRequest(r)
	if err != nil {
		writeExecError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func PatchExecHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	exec, err := loadExecFromRequest(r)
	if err != nil {
		writeExecError(w, err)
		return
	}

	updates := make(map[string]interface{})
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	applyExecUpdates(exec, updates)
	if err := sqlconnect.UpdateExec(*exec); err != nil {
		writeExecError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(exec)
}

func DeleteExecHandler(w http.ResponseWriter, r *http.Request) {
	id, err := parseExecID(r)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := sqlconnect.DeleteExec(id); err != nil {
		writeExecError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
		ID     int    `json:"id"`
	}{
		Status: "Exec successfully deleted",
		ID:     id,
	})
}

func writeExecCollection(w http.ResponseWriter, execs []models.Exec) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string        `json:"status"`
		Count  int           `json:"count"`
		Data   []models.Exec `json:"data"`
	}{
		Status: "success",
		Count:  len(execs),
		Data:   execs,
	})
}

func loadExecFromRequest(r *http.Request) (*models.Exec, error) {
	id, err := parseExecID(r)
	if err != nil {
		return nil, err
	}
	return sqlconnect.GetExecByID(id)
}

func parseExecID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func writeExecError(w http.ResponseWriter, err error) {
	if errors.Is(err, sqlconnect.ErrExecNotFound) {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func applyExecUpdates(exec *models.Exec, updates map[string]interface{}) {
	for field, value := range updates {
		switch field {
		case "first_name":
			if v, ok := value.(string); ok {
				exec.FirstName = v
			}
		case "last_name":
			if v, ok := value.(string); ok {
				exec.LastName = v
			}
		case "email":
			if v, ok := value.(string); ok {
				exec.Email = v
			}
		case "username":
			if v, ok := value.(string); ok {
				exec.Username = v
			}
		case "password":
			if v, ok := value.(string); ok {
				exec.Password = v
			}
		case "role":
			if v, ok := value.(string); ok {
				exec.Role = v
			}
		}
	}
}
