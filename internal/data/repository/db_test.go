package repository

import (
	"testing"

	"project-POS-APP-golang-integer/internal/data/entity"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) (*gorm.DB, func()) {
	dsn := "host=localhost user=postgres password=postgres dbname=pos_system_test port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	err = db.AutoMigrate(&entity.User{}, &entity.Profile{})
	if err != nil {
		t.Fatalf("failed migrate: %v", err)
	}

	cleanup := func() {
		db.Exec("TRUNCATE users, profiles RESTART IDENTITY CASCADE")
	}

	return db, cleanup
}
