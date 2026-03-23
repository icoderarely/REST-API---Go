package router

import (
	"net/http"

	"restapi/internal/api/handlers"
)

func registerExecsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /execs", handlers.GetExecsHandler)
	mux.HandleFunc("POST /execs", handlers.AddExecsHandler)
	mux.HandleFunc("PATCH /execs", handlers.PatchExecsHandler)

	mux.HandleFunc("GET /execs/{id}", handlers.GetExecHandler)
	mux.HandleFunc("PATCH /execs/{id}", handlers.PatchExecHandler)
	mux.HandleFunc("DELETE /execs/{id}", handlers.DeleteExecHandler)
	mux.HandleFunc("POST /execs/{id}/updatepassword", notImplementedHandler)
	mux.HandleFunc("POST /execs/login", notImplementedHandler)
	mux.HandleFunc("POST /execs/logout", notImplementedHandler)
	mux.HandleFunc("POST /execs/forgotpassword", notImplementedHandler)
	mux.HandleFunc("POST /execs/resetpassword/reset/{resetcode}", notImplementedHandler)
}

func notImplementedHandler(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}
