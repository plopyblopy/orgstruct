package postgres

import (
	"context"

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

	result := r.db.WithContext(context.Background()).Create(&entity)
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

func (r *DepartmentRepository) GetWithDepthHierarchy(ctx context.Context, deptId int, depth int) ([]*domain.FlatDepartments, error) {
	var flatDepts []*domain.FlatDepartments

	sql := `
		WITH RECURSIVE dept_tree AS (
			SELECT id, name, parent_id, created_at, 0 AS depth
			FROM departments
			WHERE id = ?

			UNION ALL

			SELECT d.id, d.name, d.parent_id, d.created_at, dt.depth + 1
			FROM departments d
			JOIN dept_tree dt ON d.parent_id = dt.id
			WHERE dt.depth < ?
		)
		SELECT id, name, parent_id, created_at, depth
		FROM dept_tree
	 	ORDER BY depth, id
	`

	result := r.db.WithContext(ctx).Raw(sql, deptId, depth).Scan(&flatDepts)
	if result.Error != nil {
		return nil, HandleSQLError(result.Error)
	}

	if result.RowsAffected == 0 {
		return nil, domain.ErrNotFound
	}

	return flatDepts, nil
}
