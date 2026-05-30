package domain

import (
	"errors"
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

type PanicError struct {
	msg string
}

func NewPanicError(msg string) error {
	return &PanicError{msg: msg}
}

func (e *PanicError) Error() string {
	return e.msg
}
