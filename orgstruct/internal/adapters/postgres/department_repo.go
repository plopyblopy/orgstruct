package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/plopyblopy/orgstruct/internal/domain"
)

// DepartmentRepository встраивается структура Db.
type DepartmentRepository struct {
	Db
}

// NewDepartmentRepository конструктор для DepartmentRepository.
func NewDepartmentRepository(db *Db) *DepartmentRepository {
	return &DepartmentRepository{Db: *db}
}

// Add реализация интерфейса DepartamentAdder.
//
// Добавляет postgres.Department в БД.
func (r *DepartmentRepository) Add(ctx context.Context, model *domain.Department) error {
	entity := NewDepartment(*model)
	// result := r.db.Create(&entity)
	result := r.db.WithContext(context.Background()).Create(&entity)

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
