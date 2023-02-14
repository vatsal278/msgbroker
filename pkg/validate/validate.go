package validate

import "github.com/go-playground/validator/v10"

// Validate uses the validator library to validate a given input struct.
func Validate(x interface{}) error {
	validate := validator.New()
	errs := validate.Struct(x)
	if errs != nil {
		return errs
	}
	return nil
}
