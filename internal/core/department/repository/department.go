package depRepo

import (
	"context"
	"errors"
	"org-structure-api/internal/core/department/model"

	"gorm.io/gorm"
)

type DepartmentRepoGorm struct {
	db *gorm.DB
}

func NewDepartmentRepoGorm(db *gorm.DB) *DepartmentRepoGorm {
	return &DepartmentRepoGorm{db: db}
}

func (r *DepartmentRepoGorm) Create(ctx context.Context, dep *model.Department) error {
	return r.db.WithContext(ctx).Create(dep).Error
}

func (r *DepartmentRepoGorm) GetByID(ctx context.Context, id uint) (*model.Department, error) {
	var dep model.Department
	err := r.db.WithContext(ctx).First(&dep, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &dep, err
}

func (r *DepartmentRepoGorm) Update(ctx context.Context, dep *model.Department) error {
	return r.db.WithContext(ctx).Save(dep).Error
}

func (r *DepartmentRepoGorm) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&model.Department{}, id).Error
}

func (r *DepartmentRepoGorm) ListByParentID(ctx context.Context, parentID *uint) ([]*model.Department, error) {
	var deps []*model.Department
	query := r.db.WithContext(ctx).Order("created_at ASC")
	if parentID == nil {
		query = query.Where("parent_id IS NULL")
	} else {
		query = query.Where("parent_id = ?", *parentID)
	}
	err := query.Find(&deps).Error
	return deps, err
}
