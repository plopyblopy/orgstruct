package domain

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	NotAdded     = "It was not added."
	NotFound     = "It was not found."
	AlreadyExist = "It already exists."
	Validation   = "Validation error."
)

var (
	ErrNotAdded     = errors.New(NotAdded)
	ErrNotFound     = errors.New(NotFound)
	ErrAlreadyExist = errors.New(AlreadyExist)
	ErrValidation   = errors.New(Validation)
)

type SQLSTATE = string

const (
	SelfReference       SQLSTATE = "P0002"
	Cycle               SQLSTATE = "P0003"
	ForeignKeyViolation SQLSTATE = "23503"
)

// HTTPStatusCoder для получения http кода от доменной ошибки.
type HTTPStatusCoder interface {
	Code() int
}

// MultiRow для определения является ли ошибка multi-row.
type MultiRow interface {
	// IsRow для оперделения является тип multi-row.
	IsRow() bool
	// Rows для получения строк от multi-row ошибки.
	Rows() any
}

// DefaultSqlError ошибка для неопределенных ошибок.
type DefaultSqlError struct {
	msg string
}

func NewDefaultSqlError(msg string) *DefaultSqlError { return &DefaultSqlError{msg: msg} }
func (e *DefaultSqlError) Error() string             { return e.msg }
func (e *DefaultSqlError) Code() int                 { return http.StatusInternalServerError }

// PanicError ошибка оборачивание panic в ошибку определенного типа.
type PanicError struct{ msg string }

func NewPanicError(msg string) *PanicError { return &PanicError{msg: msg} }
func (e PanicError) Error() string         { return e.msg }
func (e PanicError) Code() int             { return http.StatusInternalServerError }

// ValidationError multi-row ошибки валидации полей.
type ValidationError struct{ Fields []FieldError }

func NewValidationError(errs []FieldError) *ValidationError { return &ValidationError{Fields: errs} }
func (e *ValidationError) Error() string                    { return Validation }
func (e *ValidationError) Code() int                        { return http.StatusBadRequest }
func (e *ValidationError) IsRow() bool                      { return false }
func (e *ValidationError) Rows() any                        { return e.Fields }

// FieldError ошибка для конкретного поля.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func NewFieldError(field, msg string) FieldError { return FieldError{Field: field, Message: msg} }

// ErrLengthOutOfRange ошибка выхода за диапазон длины.
func ErrLengthOutOfRange(min, max, got int) string {
	return fmt.Sprintf("length %d is not within allowed range [%d, %d]", got, min, max)
}

//////////////
// SQL errors

// SelfReferenceError ошибка ссылка на саму себя.
type SelfReferenceError struct {
	msg string
}

func NewSelfReferenceError(msg string) *SelfReferenceError { return &SelfReferenceError{msg: msg} }
func (e *SelfReferenceError) Error() string                { return e.msg }
func (e *SelfReferenceError) Code() int                    { return http.StatusConflict }

// CycleError ошибка цикла в последовательности наследования.
type CycleError struct {
	msg string
}

func NewCycleError(msg string) *CycleError { return &CycleError{msg: msg} }
func (e *CycleError) Error() string        { return e.msg }
func (e *CycleError) Code() int            { return http.StatusConflict }

// ForeignKeyViolationError ошибка внешнего ключа.
type ForeignKeyViolationError struct {
	msg string
}

func NewForeignKeyViolationError(msg string) *ForeignKeyViolationError {
	return &ForeignKeyViolationError{msg: msg}
}
func (e *ForeignKeyViolationError) Error() string { return e.msg }
func (e *ForeignKeyViolationError) Code() int     { return http.StatusConflict }
