package entities

import (
	"time"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID              `json:"id"`
	TenantID    string                 `json:"tenant_id"`
	Name        string                 `json:"name"`
	Email       value_objects.Email    `json:"email"`
	Password    value_objects.Password `json:"-"`
	DisplayName *string                `json:"display_name,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

func NewUser(tenantID, name string, email value_objects.Email, password value_objects.Password, displayName *string) (*User, error) {
	user := &User{
		ID:          uuid.New(),
		TenantID:    tenantID,
		Name:        name,
		Email:       email,
		Password:    password,
		DisplayName: displayName,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return user, nil
}
