package value_objects

import (
	"regexp"
	"strings"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
)

type Email struct {
	value string
}

func NewEmail(email string) (Email, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return Email{}, exceptions.ErrInvalidEmail
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(email) {
		return Email{}, exceptions.ErrInvalidEmail
	}

	return Email{value: email}, nil
}

func (e Email) Value() string {
	return e.value
}
