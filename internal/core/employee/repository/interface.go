package empRepo

import (
	"context"
	"org-structure-api/internal/core/employee/model"
)

type Interface interface {
	Create(ctx context.Context, emp *model.Employee) error
	GetByID(ctx context.Context, id uint) (*model.Employee, error)
	ListByDepartmentID(ctx context.Context, departmentID uint) ([]*model.Employee, error)
	Update(ctx context.Context, emp *model.Employee) error
	Delete(ctx context.Context, id uint) error
	UpdateDepartmentIDForAll(ctx context.Context, oldDepID uint, newDepID uint) error
	DeleteByDepartmentID(ctx context.Context, depID uint) error
}
