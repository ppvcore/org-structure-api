package depSvc

import (
	"context"
	depDtoResp "org-structure-api/internal/core/department/dto/response"
	dep "org-structure-api/internal/core/department/model"
)

type Interface interface {
	Create(ctx context.Context, dep *dep.Department) error
	GetByID(ctx context.Context, id uint, depth int, includeEmployees bool) (*depDtoResp.Department, error)
	Update(ctx context.Context, dep *dep.Department) (*depDtoResp.Department, error)
	Delete(ctx context.Context, id uint, mode string, reassignTo *uint) error
}
