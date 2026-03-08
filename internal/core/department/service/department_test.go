package depSvc_test

import (
	"context"
	"errors"
	"testing"

	"org-structure-api/internal/core/department/model"
	depSvc "org-structure-api/internal/core/department/service"

	"github.com/stretchr/testify/assert"
)

type fakeRepo struct {
	data map[uint]*model.Department
}

func (r *fakeRepo) GetByID(_ context.Context, id uint) (*model.Department, error) {
	return r.data[id], nil
}

func (r *fakeRepo) ListByParentID(_ context.Context, parentID *uint) ([]*model.Department, error) {
	var list []*model.Department
	for _, d := range r.data {
		if (parentID == nil && d.ParentID == nil) ||
			(parentID != nil && d.ParentID != nil && *d.ParentID == *parentID) {
			list = append(list, d)
		}
	}
	return list, nil
}

func (r *fakeRepo) Create(_ context.Context, dep *model.Department) error {
	if dep.ID == 0 {
		dep.ID = uint(len(r.data) + 1)
	}
	r.data[dep.ID] = dep
	return nil
}

func (r *fakeRepo) Update(_ context.Context, dep *model.Department) error {
	if _, exists := r.data[dep.ID]; !exists {
		return errors.New("department not found")
	}
	r.data[dep.ID] = dep
	return nil
}

func (r *fakeRepo) Delete(_ context.Context, id uint) error {
	delete(r.data, id)
	return nil
}
func ptr[T any](v T) *T { return &v }

func TestDepartmentSvc_Create_InvalidName(t *testing.T) {
	repo := &fakeRepo{data: make(map[uint]*model.Department)}
	svc := depSvc.NewDepartmentSvc(repo, nil)

	dep := &model.Department{Name: "  "}

	err := svc.Create(context.Background(), dep)
	assert.Error(t, err)
	assert.Equal(t, "invalid department name", err.Error())
}

func TestDepartmentSvc_Create_DuplicateName(t *testing.T) {
	repo := &fakeRepo{data: make(map[uint]*model.Department)}

	repo.data[1] = &model.Department{ID: 1, Name: "Development", ParentID: nil}

	repo.data[10] = &model.Department{ID: 10, Name: "Backend", ParentID: ptr(uint(1))}

	svc := depSvc.NewDepartmentSvc(repo, nil)

	dep := &model.Department{Name: "Backend", ParentID: ptr(uint(1))}

	err := svc.Create(context.Background(), dep)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unique in parent")
}

func TestDepartmentSvc_Update_SelfParent(t *testing.T) {
	repo := &fakeRepo{data: make(map[uint]*model.Department)}
	repo.data[5] = &model.Department{ID: 5}

	svc := depSvc.NewDepartmentSvc(repo, nil)

	dep := &model.Department{ID: 5, ParentID: ptr(uint(5))}

	err := svc.Update(context.Background(), dep)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "own parent")
}
