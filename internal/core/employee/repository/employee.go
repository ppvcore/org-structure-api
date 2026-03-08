package empRepo

import (
	"context"
	"errors"
	"org-structure-api/internal/core/employee/model"

	"gorm.io/gorm"
)

type EmployeeRepoGorm struct {
	db *gorm.DB
}

func NewEmployeeRepoGorm(db *gorm.DB) (*EmployeeRepoGorm, error) {
	if db == nil {
		return nil, errors.New("nil gorm DB")
	}

	return &EmployeeRepoGorm{db: db}, nil
}

func (r *EmployeeRepoGorm) Create(ctx context.Context, emp *model.Employee) error {
	return r.db.WithContext(ctx).Create(emp).Error
}

func (r *EmployeeRepoGorm) ListByDepartmentID(ctx context.Context, departmentID uint) ([]*model.Employee, error) {
	var employees []*model.Employee
	err := r.db.WithContext(ctx).
		Where("department_id = ?", departmentID).
		Order("created_at DESC").
		Find(&employees).Error
	return employees, err
}

func (r *EmployeeRepoGorm) UpdateDepartmentIDForAll(ctx context.Context, oldDepID, newDepID uint) error {
	return r.db.WithContext(ctx).
		Model(&model.Employee{}).
		Where("department_id = ?", oldDepID).
		Update("department_id", newDepID).Error
}

func (r *EmployeeRepoGorm) DeleteByDepartmentID(ctx context.Context, depID uint) error {
	return r.db.WithContext(ctx).
		Where("department_id = ?", depID).
		Delete(&model.Employee{}).Error
}
