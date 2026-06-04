package domain

import (
	"context"
)

// DepartamentAdder для добавления Department.
type DepartamentAdder interface {
	Add(ctx context.Context, model *Department) error
}

// EmployeeAdder для добавления Employee.
type EmployeeAdder interface {
	Add(ctx context.Context, model *Employee) error
}

// WithDepthHierarchyGetter для получения среза Department с добавленным полем depth.
type WithDepthHierarchyGetter interface {
	GetWithDepthHierarchy(ctx context.Context, deptId int, depth int) ([]*FlatDepartments, error)
}

// EmployeesByDepartmentIdsGetter для получения среза Employee по срезу Department id.
type EmployeesByDepartmentIdsGetter interface {
	GetByDepartmentIds(ctx context.Context, ids ...int) ([]*EmployeeResponse, error)
}
