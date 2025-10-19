package persistence

import (
	"time"

	"github.com/google/uuid"
)

type UserModel struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	TenantID    string    `gorm:"type:varchar(100);not null;index:idx_users_tenant"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Email       string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_unique_email_per_tenant"`
	Password    string    `gorm:"type:varchar(255);not null"`
	DisplayName *string   `gorm:"type:varchar(255)"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}
