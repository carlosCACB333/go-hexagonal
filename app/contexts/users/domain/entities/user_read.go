package entities

import "github.com/google/uuid"

type UserRead struct {
	ID          uuid.UUID `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	DisplayName *string   `json:"display_name,omitempty"`
	CreatedAt   string    `json:"created_at"`
}

func NewUserRead(id uuid.UUID, tenantID, name, email string, displayName *string, createdAt string) *UserRead {
	return &UserRead{
		ID:          id,
		TenantID:    tenantID,
		Name:        name,
		Email:       email,
		DisplayName: displayName,
		CreatedAt:   createdAt,
	}
}
