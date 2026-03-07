package depRepo

import (
	"context"
	"org-structure-api/internal/core/department/model"
)

type Interface interface {
	Create(ctx context.Context, dep *model.Department) error
	GetByID(ctx context.Context, id uint) (*model.Department, error)
	Update(ctx context.Context, dep *model.Department) error
	Delete(ctx context.Context, id uint) error
	ListByParentID(ctx context.Context, parentID *uint) ([]*model.Department, error)
}
