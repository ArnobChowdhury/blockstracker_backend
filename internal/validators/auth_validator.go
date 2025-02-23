package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

// Custom password validation rule
func StrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Ensure at least 8 characters
	if len(password) < 8 {
		return false
	}

	// Ensure at least one lowercase letter
	if !regexp.MustCompile(`[a-z]`).MatchString(password) {
		return false
	}

	// Ensure at least one uppercase letter
	if !regexp.MustCompile(`[A-Z]`).MatchString(password) {
		return false
	}

	// Ensure at least one digit
	if !regexp.MustCompile(`[0-9]`).MatchString(password) {
		return false
	}

	return true
}
