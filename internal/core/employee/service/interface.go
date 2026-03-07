package empSvc

import (
	"context"
	"org-structure-api/internal/core/employee/model"
)

type Interface interface {
	Create(ctx context.Context, emp *model.Employee) error
	GetByID(ctx context.Context, id uint) (*model.Employee, error)
	ListByDepartment(ctx context.Context, departmentID uint) ([]*model.Employee, error)
	Update(ctx context.Context, emp *model.Employee) error
	Delete(ctx context.Context, id uint) error
}
