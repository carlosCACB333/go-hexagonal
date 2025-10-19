package persistence

import (
	"context"

	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/entities"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/exceptions"
	"github.com/carloscacb333/go-hexagonal/app/contexts/users/domain/value_objects"
	shared_exceptions "github.com/carloscacb333/go-hexagonal/app/shared/domain/exceptions"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GormUserRepository struct {
	db *gorm.DB
}

func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

func (r *GormUserRepository) Save(ctx context.Context, user *entities.User) error {
	model := &UserModel{
		ID:          user.ID,
		TenantID:    user.TenantID,
		Name:        user.Name,
		Email:       user.Email.Value(),
		Password:    user.Password.Hash(),
		DisplayName: user.DisplayName,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
	}

	if err := r.db.WithContext(ctx).Create(model).Error; err != nil {
		if isDuplicateKeyError(err) {
			return exceptions.ErrDuplicateEmail
		}
		return shared_exceptions.NewInternalServerError("failed to save user", err.Error())
	}

	return nil
}

func (r *GormUserRepository) FindByID(ctx context.Context, tenantID string, id uuid.UUID) (*entities.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).
		Where("id = ? AND tenant_id = ?", id, tenantID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, shared_exceptions.NewInternalServerError("failed to find user", err.Error())
	}

	return r.toDomain(&model), nil
}

func (r *GormUserRepository) FindByEmail(ctx context.Context, tenantID string, email value_objects.Email) (*entities.User, error) {
	var model UserModel

	err := r.db.WithContext(ctx).
		Where("email = ? AND tenant_id = ?", email.Value(), tenantID).
		First(&model).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, exceptions.ErrUserNotFound
		}
		return nil, shared_exceptions.NewInternalServerError("failed to find user by email", err.Error())
	}

	return r.toDomain(&model), nil
}

func (r *GormUserRepository) ExistsByEmail(ctx context.Context, tenantID string, email value_objects.Email) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&UserModel{}).
		Where("email = ? AND tenant_id = ?", email.Value(), tenantID).
		Count(&count).Error

	if err != nil {
		return false, shared_exceptions.NewInternalServerError("failed to check if email exists", err.Error())
	}

	return count > 0, nil
}

func (r *GormUserRepository) toDomain(model *UserModel) *entities.User {
	email, _ := value_objects.NewEmail(model.Email)
	password := value_objects.NewPasswordFromHash(model.Password)

	return &entities.User{
		ID:          model.ID,
		TenantID:    model.TenantID,
		Name:        model.Name,
		Email:       email,
		Password:    password,
		DisplayName: model.DisplayName,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}
}

func isDuplicateKeyError(err error) bool {
	return err != nil && (gorm.ErrDuplicatedKey == err ||
		contains(err.Error(), "duplicate key") ||
		contains(err.Error(), "unique constraint"))
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findInString(s, substr)))
}

func findInString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
