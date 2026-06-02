package usecases

import (
	"github.com/plopyblopy/orgstruct/internal/domain"
)

func PostDepartment(r domain.DepartamentAdder) func(name string, parentId *int) (*domain.Department, error) {
	return func(name string, parentId *int) (*domain.Department, error) {
		dept, err := domain.NewDepartment(name, parentId)
		if err != nil {
			return nil, err
		}

		createdDept, err := r.Add(dept)
		if err != nil {
			return nil, err
		}

		return createdDept, nil
	}
}
