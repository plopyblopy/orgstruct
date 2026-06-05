package domain

import (
	"strings"
	"time"
)

// Department ограничения.
const (
	DeptNameMinLen = 1
	DeptNameMaxLen = 200
)

// Department описывает отдел.
type Department struct {
	Id        int
	Name      string
	ParentId  *int
	CreatedAt time.Time
}

// NewDepartment конструктор валидирует и возвращает Department.
func NewDepartment(name string, parentId *int) (*Department, error) {
	v := NewValidator()

	v.MinMax("name", DeptNameMinLen, DeptNameMaxLen, len(name))

	if err := v.Validate(); err != nil {
		return nil, err
	}

	trimName := strings.TrimSpace(name)

	return &Department{
		Name:     trimName,
		ParentId: parentId,
	}, nil
}

// Employee ограничения.
const (
	EmpFullNameMinLen = 1
	EmpFullNameMaxLen = 200
	EmpPositionMinLen = 1
	EmpPositionMaxLen = 200
)

// Employee описывает сотрудника.
type Employee struct {
	Id           int
	DepartmentId int
	FullName     string
	Position     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}

// NewEmployee конструктор валидирует и возвращает Employee.
func NewEmployee(departmentId int, fullName string, position string, hiredAt *time.Time) (*Employee, error) {
	v := NewValidator()

	v.MinMax("fullName", EmpFullNameMinLen, EmpFullNameMaxLen, len(fullName))
	v.MinMax("position", EmpPositionMinLen, EmpPositionMaxLen, len(position))

	if err := v.Validate(); err != nil {
		return nil, err
	}

	trimFullName := strings.TrimSpace(fullName)
	trimPosition := strings.TrimSpace(position)

	return &Employee{
		DepartmentId: departmentId,
		FullName:     trimFullName,
		Position:     trimPosition,
		HiredAt:      hiredAt,
	}, nil
}
