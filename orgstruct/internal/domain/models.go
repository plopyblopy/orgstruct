package domain

import "time"

// Department описывает отдел
type Department struct {
	Id        int
	Name      string
	ParentId  *int
	CreatedAt time.Time
}

// NewDepartment конструктор возвращающий Department
func NewDepartment(name string, parentId int) Department {
	return Department{
		Name:     name,
		ParentId: &parentId,
	}
}

// Department описывает сотрудника
type Employee struct {
	Id           int
	DepartmentId int
	FullName     string
	HiredAt      *time.Time
	CreatedAt    time.Time
}

// NewEmployee конструктор возвращающий Employee
func NewEmployee(departmentId int, fullName string, hiredAt *time.Time) Employee {
	return Employee{
		DepartmentId: departmentId,
		FullName:     fullName,
		HiredAt:      hiredAt,
	}
}
