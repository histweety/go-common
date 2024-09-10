package errors

import (
	"github.com/go-playground/validator/v10"
)

type BadRequestError struct {
	field   string
	tag     string
	value   interface{}
	message string
}

func NewBadRequestError(err validator.FieldError) *BadRequestError {
	return &BadRequestError{
		field:   err.Field(),
		tag:     err.Tag(),
		value:   err.Value(),
		message: err.Error(),
	}
}
