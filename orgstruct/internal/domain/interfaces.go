package domain

// DepartamentAdder для добавления Department.
type DepartamentAdder interface {
	Add(model Department) (*Department, error)
}
