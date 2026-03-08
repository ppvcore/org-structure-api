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

func NewDepartmentSvc(repo depRepo.Interface, empSvc empSvc.Interface) (*DepartmentSvc, error) {
	if repo == nil {
		return nil, errors.New("nil department repository")
	}

	if empSvc == nil {
		return nil, errors.New("nil employee service")
	}

	return &DepartmentSvc{
		repo:   repo,
		empSvc: empSvc,
	}, nil
}

func (s *DepartmentSvc) Create(ctx context.Context, dep *model.Department) error {
	dep.Name = strings.TrimSpace(dep.Name)
	if dep.Name == "" || len(dep.Name) > 200 {
		return errors.New("invalid department name")
	}

	if dep.ParentID != nil {
		parent, _ := s.repo.GetByID(ctx, *dep.ParentID)
		if parent == nil {
			return errors.New("parent department not found")
		}
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

func (s *DepartmentSvc) GetByID(ctx context.Context, id uint, depth int, includeEmpl bool) (*depDtoResp.Department, error) {
	dep, err := s.repo.GetByID(ctx, id)
	if err != nil || dep == nil {
		return nil, err
	}

	full := &depDtoResp.Department{
		ID: dep.ID, Name: dep.Name, ParentID: dep.ParentID, CreatedAt: dep.CreatedAt,
	}

	if includeEmpl {
		emps, _ := s.empSvc.ListByDepartmentID(ctx, id)
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

func (s *DepartmentSvc) Update(ctx context.Context, dep *model.Department) (*depDtoResp.Department, error) {
	existing, err := s.repo.GetByID(ctx, dep.ID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, errors.New("department not found")
	}

	if dep.Name != "" {
		dep.Name = strings.TrimSpace(dep.Name)
		if dep.Name == "" || len(dep.Name) > 200 {
			return nil, errors.New("invalid department name")
		}

		children, _ := s.repo.ListByParentID(ctx, existing.ParentID)
		for _, c := range children {
			if c.ID != dep.ID && c.Name == dep.Name {
				return nil, errors.New("department name must be unique within parent")
			}
		}
	}

	newParentID := dep.ParentID
	if newParentID != nil {
		if *newParentID == dep.ID {
			return nil, errors.New("cannot set department as its own parent")
		}

		if s.wouldCreateCycle(ctx, dep.ID, *newParentID) {
			return nil, errors.New("would create cycle in department tree")
		}

		parent, _ := s.repo.GetByID(ctx, *newParentID)
		if parent == nil {
			return nil, errors.New("parent department not found")
		}
	}

	if err := s.repo.Update(ctx, dep); err != nil {
		return nil, err
	}

	return s.GetByID(ctx, dep.ID, 0, false)
}

func (s *DepartmentSvc) wouldCreateCycle(ctx context.Context, depID, newParentID uint) bool {
	current := newParentID
	visited := make(map[uint]bool)

	for current != 0 {
		if current == depID {
			return true
		}
		if visited[current] {
			return true
		}
		visited[current] = true

		p, err := s.repo.GetByID(ctx, current)
		if err != nil || p == nil || p.ParentID == nil {
			return false
		}
		current = *p.ParentID
	}
	return false
}

func (s *DepartmentSvc) Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	switch mode {
	case "cascade":
		return s.deleteCascade(ctx, id)

	case "reassign":
		if reassignTo == nil {
			return errors.New("reassign_to_department_id required for mode=reassign")
		}
		target, _ := s.repo.GetByID(ctx, *reassignTo)
		if target == nil {
			return errors.New("target department not found")
		}
		if *reassignTo == id {
			return errors.New("cannot reassign to the same department")
		}

		err = s.empSvc.ReassignToDepartment(ctx, id, *reassignTo)
		if err != nil {
			return err
		}
		return s.repo.Delete(ctx, id)

	default:
		return errors.New("mode must be 'cascade' or 'reassign'")
	}
}

func (s *DepartmentSvc) deleteCascade(ctx context.Context, id uint) error {
	children, _ := s.repo.ListByParentID(ctx, &id)
	for _, ch := range children {
		if err := s.deleteCascade(ctx, ch.ID); err != nil {
			return err
		}
	}
	_ = s.empSvc.DeleteByDepartmentID(ctx, id)
	return s.repo.Delete(ctx, id)
}
