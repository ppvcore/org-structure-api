package empSvc

import (
	"context"
	"errors"
	"strings"

	"org-structure-api/internal/core/employee/model"
	empRepo "org-structure-api/internal/core/employee/repository"
)

type EmployeeSvc struct {
	repo empRepo.Interface
}

func NewEmployeeSvc(repo empRepo.Interface) (*EmployeeSvc, error) {
	if repo == nil {
		return nil, errors.New("nil employee repository")
	}

	return &EmployeeSvc{repo: repo}, nil
}

func (s *EmployeeSvc) Create(ctx context.Context, emp *model.Employee) error {
	emp.FullName = strings.TrimSpace(emp.FullName)
	emp.Position = strings.TrimSpace(emp.Position)

	if emp.FullName == "" || len(emp.FullName) > 200 {
		return errors.New("invalid employee full_name")
	}
	if emp.Position == "" || len(emp.Position) > 200 {
		return errors.New("invalid employee position")
	}

	return s.repo.Create(ctx, emp)
}

func (s *EmployeeSvc) ListByDepartmentID(ctx context.Context, departmentID uint) ([]*model.Employee, error) {
	return s.repo.ListByDepartmentID(ctx, departmentID)
}

func (s *EmployeeSvc) ReassignToDepartment(ctx context.Context, oldDepID, newDepID uint) error {
	return s.repo.UpdateDepartmentIDForAll(ctx, oldDepID, newDepID)
}

func (s *EmployeeSvc) DeleteByDepartmentID(ctx context.Context, depID uint) error {
	return s.repo.DeleteByDepartmentID(ctx, depID)
}
