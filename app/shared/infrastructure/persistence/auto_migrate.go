package persistence

import (
	"fmt"

	user_persistence "github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/persistence"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB, logger *zap.Logger) error {
	err := db.AutoMigrate(
		&user_persistence.UserModel{},
		&IdempotencyKeyModel{},
		&user_persistence.UserReadModel{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	logger.Info("Auto-migration completed successfully")
	return nil
}

func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&user_persistence.UserModel{},
		&IdempotencyKeyModel{},
		&user_persistence.UserReadModel{},
	)
}
