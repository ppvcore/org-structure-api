package api

import (
	"net/http"

	"github.com/gorilla/mux"
)

func SetupEmployeeRoutes(r *mux.Router, h *EmployeeHandler) {
	r.HandleFunc("/departments/{id:[0-9]+}/employees", h.CreateEmployee).Methods(http.MethodPost)
}
