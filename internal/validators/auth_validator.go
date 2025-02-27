package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var lowercaseRegexp = regexp.MustCompile(`[a-z]`)
var uppercaseRegexp = regexp.MustCompile(`[A-Z]`)
var digitRegexp = regexp.MustCompile(`[0-9]`)

// Wrapper function for validator
func StrongPasswordValidator(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return StrongPassword(password)
}

func StrongPassword(password string) bool {
	if len(password) < 8 {
		return false
	}
	if !lowercaseRegexp.MatchString(password) {
		return false
	}
	if !uppercaseRegexp.MatchString(password) {
		return false
	}
	if !digitRegexp.MatchString(password) {
		return false
	}
	return true
}
