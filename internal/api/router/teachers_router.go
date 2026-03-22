package router

import (
	"net/http"

	"restapi/internal/api/handlers"
)

func registerTeacherRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /teachers/", handlers.GetTeachersHandler)
	mux.HandleFunc("POST /teachers/", handlers.AddTeachersHandler)

	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler)
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)

	mux.HandleFunc("GET /teachers/{id}/students", handlers.GetStudentsFromTeachersHandler)
	mux.HandleFunc("GET /teachers/{id}/studentcount", handlers.GetStudentCountHandler)
}
