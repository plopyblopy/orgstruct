package domain

import (
	"strings"
	"time"
)

// Department описывает отдел.
type Department struct {
	Id        int
	Name      string
	ParentId  *int
	CreatedAt time.Time
}

// NewDepartment конструктор валидирует и возвращает Department.
func NewDepartment(name string, parentId *int) (Department, error) {
	verr := []FieldError{}

	if len(name) < 1 || len(name) > 200 {
		verr = append(verr, NewFieldError("name", ErrLengthOutOfRange(1, 200, len(name))))
	}

	if len(verr) != 0 {
		return Department{}, NewValidationError(verr)
	}

	trimName := strings.TrimSpace(name)

	return Department{
		Name:     trimName,
		ParentId: parentId,
	}, nil
}

// Department описывает сотрудника.
type Employee struct {
	Id           int
	DepartmentId int
	FullName     string
	Position     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}

// NewEmployee конструктор валидирует и возвращает Employee.
func NewEmployee(departmentId int, fullName string, position string, hiredAt *time.Time) (Employee, error) {
	verr := []FieldError{}

	if len(fullName) < 1 || len(fullName) > 200 {
		verr = append(verr, NewFieldError("fullName", ErrLengthOutOfRange(1, 200, len(fullName))))
	}

	if len(position) < 1 || len(position) > 200 {
		verr = append(verr, NewFieldError("position", ErrLengthOutOfRange(1, 200, len(position))))
	}

	if len(verr) != 0 {
		return Employee{}, NewValidationError(verr)
	}

	trimFullName := strings.TrimSpace(fullName)
	trimPosition := strings.TrimSpace(fullName)

	return Employee{
		DepartmentId: departmentId,
		FullName:     trimFullName,
		Position:     trimPosition,
		HiredAt:      hiredAt,
	}, nil
}
