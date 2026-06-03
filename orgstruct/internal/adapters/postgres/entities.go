package postgres

import (
	"time"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

// Department GORM адаптированный тип.
type Department struct {
	Id               int         `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	Name             string      `gorm:"column:name;type:varchar(200);not null;check:length(name) >= 1 AND length(name) <= 200;check:name = trim(name);uniqueIndex:uq_child_name,priority:2;uniqueIndex:uq_root_name,where:parent_id IS NULL"`
	ParentId         *int        `gorm:"column:parent_id;type:int;null;uniqueIndex:uq_child_name,priority:1"`
	CreatedAt        time.Time   `gorm:"column:created_at;type:timestamp;not null"`
	ParentDepartment *Department `gorm:"foreignKey:ParentId;references:Id;constraint:OnDelete:CASCADE"`
}

// NewDepartment конструктор Department.
func NewDepartment(model domain.Department) *Department {
	return &Department{
		Name:     model.Name,
		ParentId: model.ParentId,
	}
}

// Employee GORM адаптированный тип.
type Employee struct {
	Id           int        `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	DepartmentId int        `gorm:"column:department_id;type:int;not null"`
	FullName     string     `gorm:"column:full_name;type:varchar(200);not null;check:length(full_name) >= 1 AND length(full_name) <= 200;check:full_name = trim(full_name)"`
	Position     string     `gorm:"column:position;type:varchar(200);not null;check:length(position) >= 1 AND length(position) <= 200;check:position = trim(position)"`
	HiredAt      *time.Time `gorm:"column:hired_at;type:timestamp;null"`
	CreatedAt    time.Time  `gorm:"column:created_at;type:timestamp;not null"`
	Department   Department `gorm:"foreignKey:DepartmentId;references:Id;constraint:OnDelete:CASCADE;constraint:fk_employees_department"`
}

// NewEmployee конструктор Employee.
func NewEmployee(model domain.Employee) *Employee {
	return &Employee{
		DepartmentId: model.DepartmentId,
		FullName:     model.FullName,
		Position:     model.Position,
		HiredAt:      model.HiredAt,
	}
}
