package depSvc

import (
	"context"
	"errors"
	depDtoResp "org-structure-api/internal/core/department/dto/response"
	"org-structure-api/internal/core/department/model"
	depRepo "org-structure-api/internal/core/department/repository"
	empSvc "org-structure-api/internal/core/employee/service"
	"strings"
)

type DepartmentSvc struct {
	repo   depRepo.Interface
	empSvc empSvc.Interface
}

func NewDepartmentSvc(repo depRepo.Interface, empSvc empSvc.Interface) *DepartmentSvc {
	return &DepartmentSvc{
		repo:   repo,
		empSvc: empSvc,
	}
}

func (s *DepartmentSvc) Create(ctx context.Context, dep *model.Department) error {
	dep.Name = trim(dep.Name)
	if dep.Name == "" || len(dep.Name) > 200 {
		return errors.New("invalid department name")
	}

	children, err := s.repo.ListByParentID(ctx, dep.ParentID)
	if err != nil {
		return err
	}
	for _, c := range children {
		if c.Name == dep.Name {
			return errors.New("department name must be unique in parent")
		}
	}

	return s.repo.Create(ctx, dep)
}

func (s *DepartmentSvc) GetByID(ctx context.Context, id uint, depth int, includeEmpl bool) (*depDtoResp.DepartmentResponse, error) {
	dep, err := s.repo.GetByID(ctx, id)
	if err != nil || dep == nil {
		return nil, err
	}

	full := &depDtoResp.DepartmentResponse{
		ID: dep.ID, Name: dep.Name, ParentID: dep.ParentID, CreatedAt: dep.CreatedAt,
	}

	if includeEmpl {
		emps, _ := s.empSvc.ListByDepartment(ctx, id)
		full.Employees = emps
	}

	if depth > 0 {
		children, _ := s.repo.ListByParentID(ctx, &id)
		for _, ch := range children {
			childFull, _ := s.GetByID(ctx, ch.ID, depth-1, includeEmpl)
			full.Children = append(full.Children, childFull)
		}
	}

	return full, nil
}

func (s *DepartmentSvc) Update(ctx context.Context, dep *model.Department) error {
	if dep.Name != "" {
		dep.Name = trim(dep.Name)
		if dep.Name == "" || len(dep.Name) > 200 {
			return errors.New("invalid department name")
		}
	}

	return s.repo.Update(ctx, dep)
}

func (s *DepartmentSvc) Delete(ctx context.Context, id uint, cascade bool, reassignTo *uint) error {
	if cascade {
		return s.repo.Delete(ctx, id)
	}

	if reassignTo != nil {
		return s.repo.Delete(ctx, id)
	}

	return errors.New("invalid delete parameters")
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
