package repository

import (
	"context"
	"project-POS-APP-golang-integer/internal/data/entity"
	"testing"

	"go.uber.org/zap"
)

func TestProfileRepository_CreateProfile_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)

	ctx := context.Background()

	user := entity.User{
		Email: "user@test.com",
		PasswordHash: "123",
		Role: "admin",
	}

	if err := db.Create(&user).Error; err != nil {
		t.Fatalf("failed to create user: %v", err)
	}

	profile := &entity.Profile{
		UserID:   user.ID,
		FullName: "Admin",
	}

	result, err := repo.CreateProfile(ctx, profile)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID == 0 {
		t.Fatalf("expected profile ID to be set")
	}
}

func TestProfileRepository_CreateProfile_Error(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)

	ctx := context.Background()

	// UserID does not exist
	profile := &entity.Profile{
		UserID: 9999,
		FullName: "Invalid",
	}

	_, err := repo.CreateProfile(ctx, profile)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestProfileRepository_GetProfileByID_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)
	ctx := context.Background()

	user := entity.User{
		Email: "user@test.com",
		PasswordHash: "123",
		Role: "admin",
	}
	db.Create(&user)

	profile := entity.Profile{
		UserID: user.ID,
		FullName: "Dandi",
	}
	db.Create(&profile)

	result, err := repo.GetProfileByID(ctx, user.ID)

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Profile.FullName != "Dandi" {
		t.Fatalf("expected profile name Dandi")
	}
}

func TestProfileRepository_GetProfileByID_NotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)

	_, err := repo.GetProfileByID(context.Background(), 999)

	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestProfileRepository_UpdateProfile_Success(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)

	user := entity.User{
		Email: "user@test.com",
		PasswordHash: "123",
		Role: "admin",
	}
	db.Create(&user)

	profile := entity.Profile{
		UserID: user.ID,
		FullName: "Old Name",
	}
	db.Create(&profile)

	profile.FullName = "New Name"
	err := repo.UpdateProfile(context.Background(), &profile)

	if err != nil {
		t.Fatalf("expected no error")
	}
}

func TestProfileRepository_UpdateProfile_ErrorNotFound(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)

	profile := &entity.Profile{
		UserID: 999,
		FullName: "Fail",
	}

	err := repo.UpdateProfile(context.Background(), profile)

	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestProfileRepository_UpdateProfile_SQL_Error(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	logger := zap.NewNop()
	repo := NewProfileRepo(db, logger)

	// close db to force error
	sqlDB, _ := db.DB()
	sqlDB.Close()

	profile := &entity.Profile{
		UserID: 1,
		FullName: "Fail",
	}

	err := repo.UpdateProfile(context.Background(), profile)

	if err == nil {
		t.Fatalf("expected sql error")
	}
}
