package value_objects

import (
	"testing"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/stretchr/testify/assert"
)

func TestEmailValidation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{"test@example.com", "test@example.com", nil},
		{"  TEST@EXAMPLE.COM  ", "test@example.com", nil},
		{"invalid-email", "", exceptions.ErrInvalidEmail},
		{"", "", exceptions.ErrInvalidEmail},
		{"user.name+tag+sorting@example.com", "user.name+tag+sorting@example.com", nil},
		{"user@subdomain.example.com", "user@subdomain.example.com", nil},
		{"user@.com", "", exceptions.ErrInvalidEmail},
		{"user@com", "", exceptions.ErrInvalidEmail},
		{"user@exam_ple.com", "", exceptions.ErrInvalidEmail},
	}

	for _, test := range tests {
		email, err := NewEmail(test.input)
		if test.err != nil {
			assert.Error(t, err)
			assert.Equal(t, test.err, err)
		} else {
			assert.NoError(t, err)
			assert.Equal(t, test.expected, email.Value())
		}
	}
}
