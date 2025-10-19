package persistence

import "time"

type IdempotencyKeyModel struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	TenantID    string    `gorm:"type:varchar(100);not null;uniqueIndex:idx_unique_idem_key,composite:tenant_key"`
	Key         string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_unique_idem_key,composite:tenant_key"`
	ProcessedAt time.Time `gorm:"autoCreateTime;index:idx_idem_processed_at"`
}

func (IdempotencyKeyModel) TableName() string {
	return "idempotency_keys"
}
