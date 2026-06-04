package postgres

import (
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/plopyblopy/orgstruct/internal/domain"
)

// HandleSQLError преобразует ошибку от базы данных в доменную приложения.
func HandleSQLError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return domain.NewDefaultSqlError(err.Error())
	}

	switch pgErr.Code {
	case domain.SelfReference:
		return domain.NewSelfReferenceError(pgErr.Message)
	case domain.Cycle:
		return domain.NewCycleError(pgErr.Message)
	case domain.ForeignKeyViolation:
		return domain.NewForeignKeyViolationError(pgErr.Message)
	default:
		return domain.NewDefaultSqlError(fmt.Sprintf("code: %s, message: %s", pgErr.Code, pgErr.Message))
	}
}
