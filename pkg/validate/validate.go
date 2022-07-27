package validate

import "github.com/go-playground/validator"

func ValidateRequest(x interface{}) error {
	validate := validator.New()
	errs := validate.Struct(x)
	if errs != nil {
		return errs
	}
	return nil
}
