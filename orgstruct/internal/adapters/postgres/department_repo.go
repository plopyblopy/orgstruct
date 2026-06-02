package postgres

import (
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
// Добавляет postgres.Department в БД.
func (r *DepartmentRepository) Add(model domain.Department) (*domain.Department, error) {
	entity := NewDepartment(model)
	result := r.db.Create(&entity)

	if result.Error != nil {
		var pgErr *pgconn.PgError
		if errors.As(result.Error, &pgErr) {
			return nil, HandleSQLError(pgErr.Code, pgErr.Message)
		} else {
			return nil, result.Error
		}
	}

	if result.RowsAffected == 0 {
		return nil, domain.ErrNotAdded
	}

	model.Id = entity.Id
	model.CreatedAt = entity.CreatedAt

	return &model, nil
}
