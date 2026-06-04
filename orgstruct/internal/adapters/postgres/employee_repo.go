package postgres

import (
	"context"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

// EmployeeRepository встраивается структура Db.
type EmployeeRepository struct {
	Db
}

// NewEmployeeRepository конструктор для EmployeeRepository.
func NewEmployeeRepository(db *Db) *EmployeeRepository {
	return &EmployeeRepository{Db: *db}
}

// Add реализация интерфейса EmployeeAdder.
//
// Добавляет postgres.Employee в БД.
func (r *EmployeeRepository) Add(ctx context.Context, model *domain.Employee) error {
	entity := NewEmployee(*model)

	result := r.db.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		return HandleSQLError(result.Error)
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotAdded
	}

	model.Id = entity.Id
	model.CreatedAt = entity.CreatedAt

	return nil
}

func (r *EmployeeRepository) GetByDepartmentIds(ctx context.Context, deptIds ...int) ([]*domain.EmployeeResponse, error) {
	var emps []*Employee
	res := r.db.WithContext(ctx).Where("department_id IN ?", deptIds).Order("created_at ASC, full_name ASC").Find(&emps)
	if res.Error != nil {
		return nil, HandleSQLError(res.Error)
	}
	if res.RowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	resEmps := make([]*domain.EmployeeResponse, len(emps))
	for i := range emps {
		emp := emps[i]
		resEmps[i] = &domain.EmployeeResponse{
			Id:           emp.Id,
			FullName:     emp.FullName,
			DepartmentId: emp.DepartmentId,
			Position:     emp.Position,
			HiredAt:      emp.HiredAt,
			CreatedAt:    emp.CreatedAt,
		}
	}

	return resEmps, nil
}
