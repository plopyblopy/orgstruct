package usecases

import (
	"github.com/plopyblopy/orgstruct/internal/domain"
)

func PostDepartment() func(name string, parentId *int) (domain.Department, error) {
	return func(name string, parentId *int) (domain.Department, error) {
		return domain.Department{}, nil
	}
}
