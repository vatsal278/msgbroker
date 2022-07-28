package validate

import "github.com/go-playground/validator/v10"

func Validate(x interface{}) error {
	validate := validator.New()
	errs := validate.Struct(x)
	if errs != nil {
		return errs
	}
	return nil
}
