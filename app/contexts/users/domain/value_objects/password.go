package value_objects

import (
	"encoding/json"
	"regexp"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"golang.org/x/crypto/bcrypt"
)

type Password struct {
	hashedValue string
}

func NewPassword(plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, exceptions.ErrWeakPassword
	}

	if !hasUpperCase(plainPassword) || !hasLowerCase(plainPassword) || !hasDigit(plainPassword) {
		return Password{}, exceptions.ErrWeakPassword
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return Password{}, shared_exceptions.NewBadRequestError("failed to hash password", err.Error())
	}

	return Password{hashedValue: string(hashed)}, nil
}

func NewPasswordFromHash(hash string) Password {
	return Password{hashedValue: hash}
}

func (p Password) Hash() string {
	return p.hashedValue
}

func (p Password) Verify(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.hashedValue), []byte(plainPassword))
	return err == nil
}

func (p *Password) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	*p = NewPasswordFromHash(s)
	return nil
}

func hasUpperCase(s string) bool {
	return regexp.MustCompile(`[A-Z]`).MatchString(s)
}

func hasLowerCase(s string) bool {
	return regexp.MustCompile(`[a-z]`).MatchString(s)
}

func hasDigit(s string) bool {
	return regexp.MustCompile(`[0-9]`).MatchString(s)
}
