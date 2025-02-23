package validators

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func RegisterCustomValidators() {
	Validate.RegisterValidation("strongpassword", StrongPassword)
}
