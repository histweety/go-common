package utils

import "github.com/go-playground/validator/v10"

var validate *validator.Validate

func initValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		initValidator()
	}

	return validate.Struct(s)
}

func ValidateField(field interface{}, tag string) error {
	if validate == nil {
		initValidator()
	}

	return validate.Var(field, tag)
}
