package router

import (
	"net/http"

	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.RootHandler)

	registerTeacherRoutes(mux)
	registerStudentRoutes(mux)
	registerExecsRoutes(mux)

	// Executives handler
	mux.HandleFunc("/execs/", handlers.ExecsHandler)

	return mux
}
