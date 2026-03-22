package router

import (
	"net/http"

	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	// Teachers handler
	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	// Students handler
	mux.HandleFunc("GET /students/", handlers.GetStudentsHandler)
	mux.HandleFunc("POST /students/", handlers.AddStudentsHandler)

	mux.HandleFunc("GET /students/{id}", handlers.GetStudentHandler)
	mux.HandleFunc("PUT /students/{id}", handlers.UpdateStudentHandler)
	mux.HandleFunc("PATCH /students/{id}", handlers.PatchStudentHandler)
	mux.HandleFunc("DELETE /students/{id}", handlers.DeleteStudentHandler)

	// Executives handler
	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
