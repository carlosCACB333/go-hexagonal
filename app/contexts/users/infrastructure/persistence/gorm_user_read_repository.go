package persistence

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormUserReadRepository struct {
	db *gorm.DB
}

func NewGormUserReadRepository(db *gorm.DB) *GormUserReadRepository {
	return &GormUserReadRepository{db: db}
}

func (r *GormUserReadRepository) FindByID(ctx context.Context, tenantID string, id uuid.UUID) (*entities.UserRead, error) {
	var model UserReadModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, shared_exceptions.NewInternalServerError("failed to find user read model", err.Error())
	}

	return entities.NewUserRead(
		model.ID,
		model.TenantID,
		model.Name,
		model.Email,
		model.DisplayName,
		model.CreatedAt,
	), nil
}

func (r *GormUserReadRepository) Upsert(ctx context.Context, dto *entities.UserRead) error {
	model := &UserReadModel{
		ID:          dto.ID,
		TenantID:    dto.TenantID,
		Name:        dto.Name,
		Email:       dto.Email,
		DisplayName: dto.DisplayName,
		CreatedAt:   dto.CreatedAt,
	}

	// GORM upsert: Create or Update
	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", dto.ID, dto.TenantID).
		Assign(model).
		FirstOrCreate(&model).Error

	if err != nil {
		return shared_exceptions.NewInternalServerError("failed to upsert user read model", err.Error())
	}

	return nil
}
