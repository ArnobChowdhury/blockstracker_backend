package validators_test

import (
	"blockstracker_backend/internal/validators"

	"testing"

	"github.com/stretchr/testify/assert"
)

func TestS(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Valid Password", "abcd1Abcd", true},
		{"No uppercase", "abcd1abcd", false},
		{"No lowercase", "ABCD1ABCD", false},
		{"No number", "ABCDBABCD", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validators.StrongPassword(tt.password)
			assert.Equal(t, tt.expected, result)
		})
	}

}
