package empSvc_test

import (
	"context"
	"errors"
	"testing"

	"org-structure-api/internal/core/employee/model"
	empSvc "org-structure-api/internal/core/employee/service"

	"github.com/stretchr/testify/assert"
)

type fakeEmpRepo struct {
	data   map[uint]*model.Employee
	byDept map[uint][]uint
	lastID uint
}

func (r *fakeEmpRepo) Create(_ context.Context, emp *model.Employee) error {
	if emp.ID == 0 {
		r.lastID++
		emp.ID = r.lastID
	}
	r.data[emp.ID] = emp

	if emp.DepartmentID != 0 {
		r.byDept[emp.DepartmentID] = append(r.byDept[emp.DepartmentID], emp.ID)
	}
	return nil
}

func (r *fakeEmpRepo) GetByID(_ context.Context, id uint) (*model.Employee, error) {
	return r.data[id], nil
}

func (r *fakeEmpRepo) ListByDepartmentID(_ context.Context, deptID uint) ([]*model.Employee, error) {
	var list []*model.Employee
	for _, eid := range r.byDept[deptID] {
		if emp, ok := r.data[eid]; ok {
			list = append(list, emp)
		}
	}
	return list, nil
}

func (r *fakeEmpRepo) Update(_ context.Context, emp *model.Employee) error {
	if _, exists := r.data[emp.ID]; !exists {
		return errors.New("employee not found")
	}
	r.data[emp.ID] = emp

	return nil
}

func (r *fakeEmpRepo) Delete(_ context.Context, id uint) error {
	delete(r.data, id)
	return nil
}

func (r *fakeEmpRepo) UpdateDepartmentIDForAll(_ context.Context, oldDepID, newDepID uint) error {
	ids := r.byDept[oldDepID]
	delete(r.byDept, oldDepID)
	for _, eid := range ids {
		if emp, ok := r.data[eid]; ok {
			emp.DepartmentID = newDepID
			r.byDept[newDepID] = append(r.byDept[newDepID], eid)
		}
	}
	return nil
}

func (r *fakeEmpRepo) DeleteByDepartmentID(_ context.Context, depID uint) error {
	ids := r.byDept[depID]
	delete(r.byDept, depID)
	for _, id := range ids {
		delete(r.data, id)
	}
	return nil
}

func ptr[T any](v T) *T { return &v }

func TestEmployeeSvc_Create_InvalidFullName(t *testing.T) {
	repo := &fakeEmpRepo{
		data:   make(map[uint]*model.Employee),
		byDept: make(map[uint][]uint),
	}
	svc := empSvc.NewEmployeeSvc(repo)

	emp := &model.Employee{
		FullName:     "   ",
		Position:     "Developer",
		DepartmentID: 1,
	}

	err := svc.Create(context.Background(), emp)
	assert.Error(t, err)
	assert.Equal(t, "invalid employee full_name", err.Error())
}

func TestEmployeeSvc_Create_InvalidPosition(t *testing.T) {
	repo := &fakeEmpRepo{
		data:   make(map[uint]*model.Employee),
		byDept: make(map[uint][]uint),
	}
	svc := empSvc.NewEmployeeSvc(repo)

	emp := &model.Employee{
		FullName:     "John Doe",
		Position:     "",
		DepartmentID: 1,
	}

	err := svc.Create(context.Background(), emp)
	assert.Error(t, err)
	assert.Equal(t, "invalid employee position", err.Error())
}

func TestEmployeeSvc_Create_Success(t *testing.T) {
	repo := &fakeEmpRepo{
		data:   make(map[uint]*model.Employee),
		byDept: make(map[uint][]uint),
	}
	svc := empSvc.NewEmployeeSvc(repo)

	emp := &model.Employee{
		FullName:     "Anna Smith",
		Position:     "QA Engineer",
		DepartmentID: 5,
	}

	err := svc.Create(context.Background(), emp)
	assert.NoError(t, err)
	assert.NotZero(t, emp.ID)

	saved, _ := repo.GetByID(context.Background(), emp.ID)
	assert.Equal(t, "Anna Smith", saved.FullName)
	assert.Equal(t, "QA Engineer", saved.Position)
}

func TestEmployeeSvc_ReassignToDepartment(t *testing.T) {
	repo := &fakeEmpRepo{
		data:   make(map[uint]*model.Employee),
		byDept: make(map[uint][]uint),
	}
	svc := empSvc.NewEmployeeSvc(repo)

	repo.Create(context.Background(), &model.Employee{ID: 1, FullName: "A", Position: "x", DepartmentID: 10})
	repo.Create(context.Background(), &model.Employee{ID: 2, FullName: "B", Position: "y", DepartmentID: 10})

	err := svc.ReassignToDepartment(context.Background(), 10, 20)
	assert.NoError(t, err)

	emps, _ := repo.ListByDepartmentID(context.Background(), 10)
	assert.Empty(t, emps)

	emps, _ = repo.ListByDepartmentID(context.Background(), 20)
	assert.Len(t, emps, 2)
}

func TestEmployeeSvc_DeleteByDepartmentID(t *testing.T) {
	repo := &fakeEmpRepo{
		data:   make(map[uint]*model.Employee),
		byDept: make(map[uint][]uint),
	}
	svc := empSvc.NewEmployeeSvc(repo)

	repo.Create(context.Background(), &model.Employee{ID: 10, FullName: "X", Position: "Dev", DepartmentID: 7})
	repo.Create(context.Background(), &model.Employee{ID: 11, FullName: "Y", Position: "PM", DepartmentID: 7})

	err := svc.DeleteByDepartmentID(context.Background(), 7)
	assert.NoError(t, err)

	list, _ := repo.ListByDepartmentID(context.Background(), 7)
	assert.Empty(t, list)

	_, err = repo.GetByID(context.Background(), 10)
	assert.Nil(t, err)
}
