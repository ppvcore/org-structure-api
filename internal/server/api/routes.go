package api

import (
	depApi "org-structure-api/internal/core/department/api"
	empApi "org-structure-api/internal/core/employee/api"

	"github.com/gorilla/mux"
)

func SetupRoutes(r *mux.Router, depHandler *depApi.DepartmentHandler, empHandler *empApi.EmployeeHandler) {
	depApi.SetupDepartmentRoutes(r, depHandler)
	empApi.SetupEmployeeRoutes(r, empHandler)
}
