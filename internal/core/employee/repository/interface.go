package empRepo

import (
	"context"
	"org-structure-api/internal/core/employee/model"
)

type Interface interface {
	Create(ctx context.Context, emp *model.Employee) error
	ListByDepartmentID(ctx context.Context, departmentID uint) ([]*model.Employee, error)
	UpdateDepartmentIDForAll(ctx context.Context, oldDepID uint, newDepID uint) error
	DeleteByDepartmentID(ctx context.Context, depID uint) error
}
