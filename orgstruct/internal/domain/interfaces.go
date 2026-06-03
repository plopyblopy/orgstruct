package domain

import "context"

// DepartamentAdder для добавления Department.
type DepartamentAdder interface {
	Add(ctx context.Context, model *Department) error
}
