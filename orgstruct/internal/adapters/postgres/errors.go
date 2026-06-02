package postgres

import (
	"fmt"

	"github.com/plopyblopy/orgstruct/internal/domain"
)

func HandleSQLError(code string, msg string) error {
	switch code {
	case domain.SelfReference:
		return domain.NewSelfReferenceError(msg)
	case domain.Cycle:
		return domain.NewCycleError(msg)
	case domain.ForeignKeyViolation:
		return domain.NewForeignKeyViolationError(msg)
	default:
		return domain.NewDefaultSqlError(fmt.Sprintf("code: %s, message: %s", code, msg))
	}
}
