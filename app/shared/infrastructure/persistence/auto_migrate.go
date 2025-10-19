package persistence

import (
	"fmt"

	user_persistence "github.com/carloscacb333/go-hexagonal/app/contexts/users/infrastructure/persistence"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	// Migrar todos los modelos
	err := db.AutoMigrate(
		&user_persistence.UserModel{},
		&IdempotencyKeyModel{},
		&user_persistence.UserReadModel{},
	)

	if err != nil {
		return fmt.Errorf("failed to auto migrate: %w", err)
	}

	fmt.Println("✓ Auto-migration completed successfully")
	return nil
}

// DropAllTables elimina todas las tablas (útil para testing)
func DropAllTables(db *gorm.DB) error {
	return db.Migrator().DropTable(
		&user_persistence.UserModel{},
		&IdempotencyKeyModel{},
		&user_persistence.UserReadModel{},
	)
}
