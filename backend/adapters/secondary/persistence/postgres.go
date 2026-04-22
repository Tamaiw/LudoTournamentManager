package persistence

import (
	"context"
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"ludo-tournament/core/domain/models"
)

func NewPostgresDB(ctx context.Context) (*gorm.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=ludo password=changeme dbname=ludo_tournament port=5432 sslmode=disable"
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate schemas
	if err := db.WithContext(ctx).AutoMigrate(
		&models.User{},
		&models.Player{},
		&models.Tournament{},
		&models.League{},
		&models.Match{},
		&models.MatchAssignment{},
		&models.Invitation{},
		&models.UserInvite{},
	); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}

	return db, nil
}

func UpdateLastActive(db *gorm.DB, id string) error {
	return db.Model(&models.User{}).Where("id = ?", id).Update("last_active", time.Now()).Error
}
