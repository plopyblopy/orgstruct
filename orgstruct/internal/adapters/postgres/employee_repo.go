package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
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
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			return HandleSQLError(pgErr.Code, pgErr.Message)
		} else {
			return result.Error
		}
	}

	if result.RowsAffected == 0 {
		return domain.ErrNotAdded
	}

	model.Id = entity.Id
	model.CreatedAt = entity.CreatedAt

	return nil
}
