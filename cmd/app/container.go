package main

import (
	cfg "org-structure-api/internal/config"
	depRepo "org-structure-api/internal/core/department/repository"
	depSvc "org-structure-api/internal/core/department/service"
	empRepo "org-structure-api/internal/core/employee/repository"
	empSvc "org-structure-api/internal/core/employee/service"
	srv "org-structure-api/internal/server"

	"gorm.io/gorm"
)

func container(cfg cfg.ServerConfig, pg *gorm.DB) (*srv.Server, error) {
	empRepo, err := empRepo.NewEmployeeRepoGorm(pg)
	if err != nil {
		return nil, err
	}

	empSvc, err := empSvc.NewEmployeeSvc(empRepo)
	if err != nil {
		return nil, err
	}

	depRepo, err := depRepo.NewDepartmentRepoGorm(pg)
	if err != nil {
		return nil, err
	}

	depSvc, err := depSvc.NewDepartmentSvc(depRepo, empSvc)
	if err != nil {
		return nil, err
	}

	srv, err := srv.NewServer(cfg, depSvc, empSvc)
	if err != nil {
		return nil, err
	}

	return srv, nil
}
