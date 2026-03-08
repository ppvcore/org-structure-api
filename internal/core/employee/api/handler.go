package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	empDtoReq "org-structure-api/internal/core/employee/dto/request"
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
	idStr := vars["id"]
	deptID, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid department id", http.StatusBadRequest)
		return
	}

	var req empDtoReq.CreateEmployee
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var hiredAt *time.Time
	if req.HiredAt != "" {
		t, err := time.Parse("2006-01-02", req.HiredAt)
		if err != nil {
			http.Error(w, "invalid hired_at format, must be YYYY-MM-DD", http.StatusBadRequest)
			return
		}
		hiredAt = &t
	}

	emp := &model.Employee{
		DepartmentID: uint(deptID),
		FullName:     req.FullName,
		Position:     req.Position,
		HiredAt:      hiredAt,
	}

	if err := h.svc.Create(r.Context(), emp); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}
