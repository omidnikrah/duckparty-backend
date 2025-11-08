package migration

import (
	"fmt"

	"github.com/omidnikrah/duckparty-backend/internal/model"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := Up(db)
	if err != nil {
		return fmt.Errorf("❌ Failed to up migration: %w", err)
	}

	return nil
}

func PerformMigration(db *gorm.DB, models ...interface{}) error {
	for _, model := range models {
		fmt.Printf("🔄 Migrating model: %T\n", model)
		err := db.AutoMigrate(model)
		if err != nil {
			return fmt.Errorf("❌ Failed to migrate model %T: %w", model, err)
		}
		fmt.Printf("✅ Successfully migrated model: %T\n", model)
	}

	return nil
}

func Up(db *gorm.DB) error {
	models := []interface{}{
		&model.User{},
		&model.Duck{},
		&model.DuckReactions{},
	}

	return PerformMigration(db, models...)
}

func Down(db *gorm.DB) error {
	models := []interface{}{
		&model.DuckReactions{},
		&model.Duck{},
		&model.User{},
	}

	for _, model := range models {
		fmt.Printf("🗑️ Dropping table for model: %T\n", model)
		err := db.Migrator().DropTable(model)
		if err != nil {
			return fmt.Errorf("❌ Failed to drop table for model %T: %w", model, err)
		}
		fmt.Printf("✅ Successfully dropped table for model: %T\n", model)
	}

	return nil
}
