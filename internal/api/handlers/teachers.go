package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
)

var (
	teachers = make(map[int]models.Teacher)
	// mutex    = &sync.Mutex{}
	nextID = 1
)

// initialize some data
func init() {
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "John",
		LastName:  "Doe",
		Class:     "9A",
		Subject:   "Math",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10A",
		Subject:   "Algebra",
	}
	nextID++
	teachers[nextID] = models.Teacher{
		ID:        nextID,
		FirstName: "Jane",
		LastName:  "Doe",
		Class:     "113",
		Subject:   "Computer",
	}
	nextID++
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getTeachersHandler(w, r)
	case http.MethodPost:
		addTeacherHandler(w, r)
	case http.MethodPut:
		_, err := w.Write([]byte("PUT: Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	case http.MethodDelete:
		_, err := w.Write([]byte("DELETE: Teachers Route"))
		if err != nil {
			log.Fatal("Error writing to server:", err)
		}
	}
}

func getTeachersHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	fname := r.URL.Query().Get("first_name")
	lname := r.URL.Query().Get("last_name")

	teacherList := make([]models.Teacher, 0, len(teachers))

	if idStr == "" {
		for _, teacher := range teachers {
			// filter based on first_name/last_name query param
			if (fname == "" || teacher.FirstName == fname) && (lname == "" || teacher.LastName == lname) {
				teacherList = append(teacherList, teacher)
			}
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
	} else {
		id, _ := strconv.Atoi(idStr)
		teacher, exists := teachers[id]

		if !exists {
			http.Error(w, "ID not found", http.StatusNotFound)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(teacher)
	}
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

	// addedTeachers := make([]models.Teacher, len(newTeachers))
	// for i, newTeacher := range newTeachers {
	// 	newTeacher.ID = nextID
	// 	teachers[nextID] = newTeacher
	// 	addedTeachers[i] = newTeacher
	// 	nextID++
	// }

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
