package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupEmployeeRoutes(r *mux.Router, handler *EmployeeHandler) {
	r.HandleFunc("/departments/{id:[0-9]+}/employees", handler.CreateEmployee).Methods(http.MethodPost)
}
