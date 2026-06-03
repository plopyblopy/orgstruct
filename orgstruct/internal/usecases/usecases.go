package usecases

import (
	"context"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

// PostDepartment usecase для добавления нового Department.
func PostDepartment(r domain.DepartamentAdder) func(ctx context.Context, name string, parentId *int) (*domain.Department, error) {
	return func(ctx context.Context, name string, parentId *int) (*domain.Department, error) {
		model, err := domain.NewDepartment(name, parentId)
		if err != nil {
			return nil, err
		}

		err = r.Add(ctx, model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}
}

// PostEmployee usecase для добавления нового Employee.
func PostEmployee(r domain.EmployeeAdder) func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error) {
	return func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error) {

		return nil, nil
	}
}
