package domain

import "time"

type AddEmployeeRequest struct {
	FullName     string
	Position     string
	DepartmentId int
	HiredAt      *time.Time
}

type DepartmentResponse struct {
	Id        int       `json:"id"`
	Name      string    `json:"name"`
	ParentId  *int      `json:"parent_id"`
	CreatedAt time.Time `json:"created_at"`
}

type EmployeeResponse struct {
	Id           int        `json:"id"`
	DepartmentId int        `json:"department_id"`
	FullName     string     `json:"full_name"`
	Position     string     `json:"position"`
	HiredAt      *time.Time `json:"hired_at"`
	CreatedAt    time.Time  `json:"created_at"`
}

type DepartmentWithChildResponse struct {
	DepartmentResponse
	Employees []EmployeeResponse             `json:"employees,omitempty"`
	Children  []*DepartmentWithChildResponse `json:"children,omitempty"`
}

type FlatDepartments struct {
	Department
	Depth int
}
