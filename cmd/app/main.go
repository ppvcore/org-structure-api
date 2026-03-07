package main

import (
	"log"
	cfg "org-structure-api/internal/config"
	depRepo "org-structure-api/internal/core/department/repository"
	depSvc "org-structure-api/internal/core/department/service"
	empRepo "org-structure-api/internal/core/employee/repository"
	empSvc "org-structure-api/internal/core/employee/service"
	db "org-structure-api/internal/database"
	srv "org-structure-api/internal/server"
)

func main() {
	/* if err := logger.InitZap(); err != nil {
		log.Fatalf("cannot initialize logger: %v", err)
	}
	defer logger.Close()
	zapLogs := logger.Logger */

	cfg, err := cfg.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config")
	}

	pg, err := db.NewPostgresClient(cfg.Postgres)
	if err != nil {
		log.Fatal("Failed to connect to Postgres")
	}

	empRepo := empRepo.NewEmployeeRepoGorm(pg)
	empSvc := empSvc.NewEmployeeSvc(empRepo)

	depRepo := depRepo.NewDepartmentRepoGorm(pg)
	depSvc := depSvc.NewDepartmentSvc(depRepo, empSvc)

	srv := srv.NewServer(cfg.Server, depSvc, empSvc)
	srv.Start()
}
