package usecases

import (
	"context"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

// PostDepartment usecase для добавления нового Department.
func PostDepartment(repo domain.DepartamentAdder) func(ctx context.Context, name string, parentId *int) (*domain.Department, error) {
	return func(ctx context.Context, name string, parentId *int) (*domain.Department, error) {
		model, err := domain.NewDepartment(name, parentId)
		if err != nil {
			return nil, err
		}

		err = repo.Add(ctx, model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}
}

// PostEmployee usecase для добавления нового Employee.
func PostEmployee(repo domain.EmployeeAdder) func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error) {
	return func(ctx context.Context, r domain.AddEmployeeRequest) (*domain.Employee, error) {
		model, err := domain.NewEmployee(r.DepartmentId, r.FullName, r.Position, r.HiredAt)
		if err != nil {
			return nil, err
		}

		err = repo.Add(ctx, model)
		if err != nil {
			return nil, err
		}

		return model, nil
	}
}
