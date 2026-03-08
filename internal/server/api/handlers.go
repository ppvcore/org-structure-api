package api

import (
	depApi "org-structure-api/internal/core/department/api"
	depSvc "org-structure-api/internal/core/department/service"
	empApi "org-structure-api/internal/core/employee/api"
	empSvc "org-structure-api/internal/core/employee/service"
)

func SetupAllHandlers(depSvc depSvc.Interface, empSvc empSvc.Interface) (depHand *depApi.DepartmentHandler, empHand *empApi.EmployeeHandler) {
	depHand = depApi.NewDepartmentHandler(depSvc)
	empHand = empApi.NewEmployeeHandler(empSvc)
	return
}
