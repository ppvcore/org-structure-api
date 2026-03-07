package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"org-structure-api/internal/core/employee/model"
	empSvc "org-structure-api/internal/core/employee/service"

	"github.com/gorilla/mux"
)

type EmployeeHandler struct {
	svc empSvc.Interface
}

func NewEmployeeHandler(svc empSvc.Interface) *EmployeeHandler {
	return &EmployeeHandler{svc: svc}
}

func (h *EmployeeHandler) CreateEmployee(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	depIDStr, ok := vars["id"]
	if !ok {
		http.Error(w, "department id is required", http.StatusBadRequest)
		return
	}
	depID, err := strconv.ParseUint(depIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var emp model.Employee
	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	emp.DepartmentID = uint(depID)

	if err := h.svc.Create(r.Context(), &emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}
