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

func (e *BadRequestError) Error() string {
	return e.message
}

func NewBadRequestError(err validator.FieldError) *BadRequestError {
	return &BadRequestError{
		field:   err.Field(),
		tag:     err.Tag(),
		value:   err.Value(),
		message: err.Error(),
	}
}

type NotFoundError struct {
	entity  string
	message string
}

func (e *NotFoundError) Error() string {
	return e.message
}

func NewNotFoundError(entity string) *NotFoundError {
	return &NotFoundError{
		entity:  entity,
		message: "data not found (" + entity + ")",
	}
}
