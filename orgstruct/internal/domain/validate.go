package domain

// Validator хранит ошибки валидации
type Validator struct {
	errs []FieldError
}

// NewValidator конструктор возвращает указатель на экземпляр Validator
func NewValidator() *Validator { return &Validator{errs: []FieldError{}} }

func (v *Validator) Validate() error {
	if len(v.errs) != 0 {
		return v.Error()
	}
	return nil
}

// HasErrors проверка наличия ошибок
func (v *Validator) HasErrors() bool { return len(v.errs) == 0 }

// Error реализует интерфейс error возвращая доменную ошибку ValidationError, включающая ошибки полей
func (v *Validator) Error() error {
	return NewValidationError(v.errs)
}

// MinMax проверка диапазона минимального и максимального значения
func (v *Validator) MinMax(fieldName string, min, max, act int) {
	if act < min || act > max {
		v.errs = append(v.errs, NewFieldError(fieldName, ErrLengthOutOfRange(1, 200, act)))
	}
}

func (v *Validator) NotNil(fieldName string, val any) {
	if val == nil {
		v.errs = append(v.errs, NewFieldError(fieldName, ErrNotNil()))
	}
}
