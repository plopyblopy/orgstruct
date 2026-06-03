package domain

import "time"

type AddEmployeeRequest struct {
	FullName     string
	Position     string
	DepartmentId int
	HiredAt      *time.Time
}
