package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"restapi/pkg/utils"
)

var errInvalidExecID = errors.New("invalid exec id")

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

func parseExecID(r *http.Request) (int, error) {
	return strconv.Atoi(r.PathValue("id"))
}

func writeExecError(w http.ResponseWriter, err error) {
	if errors.Is(err, errInvalidExecID) {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}
	if errors.Is(err, sqlconnect.ErrExecNotFound) {
		http.Error(w, "Exec not found", http.StatusNotFound)
		return
	}
	http.Error(w, "Internal server error", http.StatusInternalServerError)
}

func writeInvalidLoginError(w http.ResponseWriter) {
	http.Error(w, "Invalid username or password", http.StatusUnauthorized)
}

func loadExecFromRequest(r *http.Request) (*models.Exec, error) {
	id, err := parseExecID(r)
	if err != nil {
		return nil, errInvalidExecID
	}

	exec, err := sqlconnect.GetExecByID(id)
	if err != nil {
		return nil, err
	}
	return exec, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req models.Exec
	// data validation
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid req body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username & password required", http.StatusBadRequest)
		return
	}

	user, err := sqlconnect.GetExecByUsername(req.Username)
	if err != nil {
		if errors.Is(err, sqlconnect.ErrExecNotFound) {
			writeInvalidLoginError(w)
			return
		}
		http.Error(w, "Unable to process login", http.StatusInternalServerError)
		return
	}

	// is user active
	if user.InactiveStatus.Valid && user.InactiveStatus.Bool {
		http.Error(w, "Account is inactive", http.StatusForbidden)
		return
	}

	if err := utils.VerifyPassword(req.Password, user.Password); err != nil {
		switch {
		case errors.Is(err, utils.ErrPasswordMismatch):
			writeInvalidLoginError(w)
		case errors.Is(err, utils.ErrInvalidPasswordEncoding):
			http.Error(w, "Unable to process login", http.StatusInternalServerError)
		default:
			http.Error(w, "Unable to process login", http.StatusInternalServerError)
		}
		return
	}

	// generate jwt token
	tokenString, err := utils.SignToken(user.ID, req.Username, user.Role)
	if err != nil {
		http.Error(w, "Could not create token", http.StatusInternalServerError)
		return
	}

	// set cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    tokenString,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Now().Add(24 * time.Hour),
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Token   string `json:"token"`
	}{
		Status:  "success",
		Message: "Login verified",
		Token:   tokenString,
	})
}

var (
	execFieldIndexOnce sync.Once
	execFieldIndex     map[string]int
)

func applyExecUpdates(exec *models.Exec, updates map[string]interface{}) {
	if exec == nil || len(updates) == 0 {
		return
	}

	fieldIndex := getExecFieldIndex()
	value := reflect.ValueOf(exec).Elem()
	for field, raw := range updates {
		idx, ok := fieldIndex[field]
		if !ok {
			continue
		}
		structField := value.Field(idx)
		if !structField.CanSet() {
			continue
		}

		switch structField.Kind() {
		case reflect.String:
			if v, ok := raw.(string); ok {
				structField.SetString(v)
			}
		}
	}
}

func getExecFieldIndex() map[string]int {
	execFieldIndexOnce.Do(func() {
		typ := reflect.TypeOf(models.Exec{})
		execFieldIndex = make(map[string]int, typ.NumField())
		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)
			jsonTag := field.Tag.Get("json")
			if jsonTag == "-" {
				continue
			}
			name := strings.Split(jsonTag, ",")[0]
			if name == "" {
				name = field.Name
			}
			execFieldIndex[name] = i
		}
	})
	return execFieldIndex
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// clear jwt token, blacklist of invalidated jwt
	http.SetCookie(w, &http.Cookie{
		Name:     "Bearer",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		Expires:  time.Unix(0, 0),
		SameSite: http.SameSiteStrictMode,
	})

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"message": "Logged out"}`))
}
