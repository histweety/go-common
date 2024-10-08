package utils

import (
	"github.com/arthadede/common-package/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
)

var validate *validator.Validate

func initValidator() {
	validate = validator.New(validator.WithRequiredStructEnabled())
}

func ValidateStruct(s interface{}) error {
	if validate == nil {
		initValidator()
	}

	err := validate.Struct(s)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Errorf("Error: %s\n", err.Error())
		}

		return errors.ErrStructValidation
	}

	return nil
}

func ValidateField(field interface{}, tag string) error {
	if validate == nil {
		initValidator()
	}

	err := validate.Var(field, tag)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			log.Errorf("Error: %s\n", err.Error())
		}

		return errors.ErrStructValidation
	}

	return nil
}
