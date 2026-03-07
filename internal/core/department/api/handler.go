package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"org-structure-api/internal/core/department/model"
	depSvc "org-structure-api/internal/core/department/service"

	"github.com/gorilla/mux"
)

type DepartmentHandler struct {
	svc depSvc.Interface
}

func NewDepartmentHandler(svc depSvc.Interface) *DepartmentHandler {
	return &DepartmentHandler{svc: svc}
}

func (h *DepartmentHandler) CreateDepartment(w http.ResponseWriter, r *http.Request) {
	var dep model.Department
	if err := json.NewDecoder(r.Body).Decode(&dep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.svc.Create(r.Context(), &dep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dep)
}

func (h *DepartmentHandler) GetDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	depth := 1
	includeEmployees := true

	if d := r.URL.Query().Get("depth"); d != "" {
		if val, err := strconv.Atoi(d); err == nil && val >= 1 && val <= 5 {
			depth = val
		}
	}

	if ie := r.URL.Query().Get("include_employees"); ie == "false" {
		includeEmployees = false
	}

	dep, err := h.svc.GetByID(r.Context(), uint(id), depth, includeEmployees)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(dep)
}

func (h *DepartmentHandler) UpdateDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	var dep model.Department
	if err := json.NewDecoder(r.Body).Decode(&dep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	dep.ID = uint(id)

	if err := h.svc.Update(r.Context(), &dep); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(dep)
}

func (h *DepartmentHandler) DeleteDepartment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		http.Error(w, "id is required", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid id", http.StatusBadRequest)
		return
	}

	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "cascade"
	}
	reassignIDStr := r.URL.Query().Get("reassign_to_department_id")

	var cascade bool
	var reassignTo *uint

	switch mode {
	case "cascade":
		cascade = true
	case "reassign":
		if reassignIDStr == "" {
			http.Error(w, "reassign_to_department_id is required", http.StatusBadRequest)
			return
		}
		val, err := strconv.ParseUint(reassignIDStr, 10, 64)
		if err != nil {
			http.Error(w, "invalid reassign_to_department_id", http.StatusBadRequest)
			return
		}
		tmp := uint(val)
		reassignTo = &tmp
	default:
		http.Error(w, "invalid mode", http.StatusBadRequest)
		return
	}

	if err := h.svc.Delete(r.Context(), uint(id), cascade, reassignTo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
