package persistence

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type GormIdempotencyRepository struct {
	db *gorm.DB
}

func NewGormIdempotencyRepository(db *gorm.DB) *GormIdempotencyRepository {
	return &GormIdempotencyRepository{db: db}
}

func (s *GormIdempotencyRepository) IsProcessed(ctx context.Context, tenantID, key string) (bool, error) {
	var count int64

	err := s.db.WithContext(ctx).
		Model(&IdempotencyKeyModel{}).
		Where("tenant_id = ? AND key = ?", tenantID, key).
		Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check idempotency key: %w", err)
	}

	return count > 0, nil
}

func (s *GormIdempotencyRepository) MarkAsProcessed(ctx context.Context, tenantID, key string) error {
	model := &IdempotencyKeyModel{
		TenantID:    tenantID,
		Key:         key,
		ProcessedAt: time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(model).Error; err != nil {
		return fmt.Errorf("failed to mark as processed: %w", err)
	}

	return nil
}
