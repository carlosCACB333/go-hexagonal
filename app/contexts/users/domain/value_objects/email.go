package value_objects

import (
	"encoding/json"
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

func (e *Email) UnmarshalJSON(b []byte) error {

	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	email, err := NewEmail(s)
	if err != nil {
		return err
	}
	*e = email
	return nil

}
