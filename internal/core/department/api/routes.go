package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupDepartmentRoutes(r *mux.Router, handler *DepartmentHandler) {
	r.HandleFunc("/departments", handler.CreateDepartment).Methods(http.MethodPost)
	r.HandleFunc("/departments/{id:[0-9]+}", handler.GetDepartment).Methods(http.MethodGet)
	r.HandleFunc("/departments/{id:[0-9]+}", handler.UpdateDepartment).Methods(http.MethodPatch)
	r.HandleFunc("/departments/{id:[0-9]+}", handler.DeleteDepartment).Methods(http.MethodDelete)
}
