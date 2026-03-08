// internal/server/api/routes.go (общий)
package api

import (
	depApi "org-structure-api/internal/core/department/api"
	empApi "org-structure-api/internal/core/employee/api"

	"github.com/gorilla/mux"
)

func SetupAllRoutes(r *mux.Router, depHand *depApi.DepartmentHandler, empHand *empApi.EmployeeHandler) {
	depApi.SetupDepartmentRoutes(r, depHand)
	empApi.SetupEmployeeRoutes(r, empHand)
}
