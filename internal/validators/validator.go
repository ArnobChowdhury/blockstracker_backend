package validators

import (
	responsemsg "blockstracker_backend/constants"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func GetCustomMessage(err validator.FieldError, req any) string {
	field, _ := reflect.TypeOf(req).FieldByName(err.Field())

	switch err.Tag() {
	case "required":
		return field.Name + " is required"
	case "email":
		return responsemsg.InvalidEmail
	case "strongpassword":
		return responsemsg.NotStrongPassword
	default:
		return "Validation failed for " + field.Name
	}
}

var Validate = validator.New()

func RegisterCustomValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("strongpassword", StrongPassword)
	}

}
