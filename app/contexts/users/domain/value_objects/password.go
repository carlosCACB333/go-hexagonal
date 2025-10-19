package value_objects

import (
	"regexp"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	shared_ports "github.com/carloscacb333/go-hexagonal/app/shared/domain/ports"
)

type Password struct {
	hashedValue string
	hasher      shared_ports.Hasher
}

func NewPassword(hasher shared_ports.Hasher, plainPassword string) (Password, error) {
	if len(plainPassword) < 8 {
		return Password{}, exceptions.ErrWeakPassword
	}

	if !hasUpperCase(plainPassword) || !hasLowerCase(plainPassword) || !hasDigit(plainPassword) {
		return Password{}, exceptions.ErrWeakPassword
	}

	hashed, err := hasher.Hash(plainPassword)
	if err != nil {
		return Password{}, shared_exceptions.NewBadRequestError("failed to hash password", err.Error())
	}

	return Password{hashedValue: hashed, hasher: hasher}, nil
}

func NewPasswordFromHash(hash string) Password {
	return Password{hashedValue: hash}
}

func (p Password) Hash() string {
	return p.hashedValue
}

func (p Password) Verify(plainPassword string) bool {
	return p.hasher.Verify(p.hashedValue, plainPassword)
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
