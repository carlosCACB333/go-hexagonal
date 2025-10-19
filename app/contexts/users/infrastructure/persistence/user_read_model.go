package persistence

import (
	"github.com/google/uuid"
)

type UserReadModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID    string    `gorm:"type:varchar(100);primaryKey;index:idx_read_tenant"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);not null"`
	DisplayName *string   `gorm:"type:varchar(255)"`
	CreatedAt   string    `gorm:"type:varchar(50);not null"`
}

func (UserReadModel) TableName() string {
	return "users_read"
}
