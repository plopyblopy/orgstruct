package handlers

import "github.com/plopyblopy/orgstruct/internal/domain"

func ParseDepartamentModelToResponse(m domain.Department) domain.DepartmentResponse {
	return domain.DepartmentResponse{
		Id:        m.Id,
		Name:      m.Name,
		ParentId:  m.ParentId,
		CreatedAt: m.CreatedAt,
	}
}

func ParseEmployeeModelToResponse(m domain.Employee) domain.EmployeeResponse {
	return domain.EmployeeResponse{
		Id:           m.Id,
		FullName:     m.FullName,
		Position:     m.Position,
		DepartmentId: m.DepartmentId,
		HiredAt:      m.HiredAt,
		CreatedAt:    m.CreatedAt,
	}
}
