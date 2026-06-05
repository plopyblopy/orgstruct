package domain

import (
	"encoding/json"
	"time"
)

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

type AddEmployeeRequest struct {
	FullName     string
	Position     string
	DepartmentId int
	HiredAt      *time.Time
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

type UpdateDepartment struct {
	Id        int                     `json:"id"`
	Name      UpdateField[*string]    `json:"name"`
	ParentId  UpdateField[*int]       `json:"parent_id"`
	CreatedAt UpdateField[*time.Time] `json:"created_at"`
}

// UpdateField - это тип для обновления полей, который позволяет вам определить, следует ли обновлять поле, поскольку оно может быть необязательным и вообще не передаваться.
//
// Value: если было передано ненулевое значение, оно сохраняется здесь.
//
// Valid: false - поле не было передано, true - есть null или Value.
//
// IsNull: false != null, true == null.
type UpdateField[T any] struct {
	Value  T
	Valid  bool
	IsNull bool
}

// UnmarshalJSON реализует интерфей Unmarshaler.
//
// Проверяет, было ли передано поле, а также наличие значения или null.
func (f *UpdateField[T]) UnmarshalJSON(data []byte) error {
	// если null.
	if string(data) == "null" {
		f.Valid = true
		f.IsNull = true
		return nil
	}

	// если не null значение.
	if err := json.Unmarshal(data, &f.Value); err != nil {
		return err
	}

	f.Valid = true

	return nil
}
