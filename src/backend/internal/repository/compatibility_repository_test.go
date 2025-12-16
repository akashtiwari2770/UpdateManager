package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var compatibilityRepo *CompatibilityRepository

func setupCompatibilityTestDB(t *testing.T) {
	ctx := context.Background()

	cfg := &database.Config{
		URI:      "mongodb://admin:admin123@localhost:27017/updatemanager_test?authSource=admin",
		Database: "updatemanager_test",
		Timeout:  10 * time.Second,
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	testDB = db
	testCtx = ctx
	compatibilityRepo = NewCompatibilityRepository(db.Collection("compatibility_matrices"))
}

func teardownCompatibilityTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("compatibility_matrices").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestCompatibilityCreate(t *testing.T) {
	setupCompatibilityTestDB(t)
	defer teardownCompatibilityTestDB(t)

	matrix := &models.CompatibilityMatrix{
		ProductID:        "test-product",
		VersionNumber:    "1.0.0",
		MinServerVersion: "1.0.0",
		MaxServerVersion: "2.0.0",
		ValidatedBy:      "test-user",
		ValidationStatus: models.ValidationStatusPassed,
	}

	err := compatibilityRepo.Create(testCtx, matrix)
	if err != nil {
		t.Fatalf("Failed to create compatibility matrix: %v", err)
	}

	if matrix.ID.IsZero() {
		t.Error("Compatibility matrix ID was not set after creation")
	}

	t.Logf("Created compatibility matrix with ID: %s", matrix.ID.Hex())
}

func TestCompatibilityGetByProductIDAndVersion(t *testing.T) {
	setupCompatibilityTestDB(t)
	defer teardownCompatibilityTestDB(t)

	productID := "test-product"
	versionNumber := "1.0.0"
	matrix := &models.CompatibilityMatrix{
		ProductID:        productID,
		VersionNumber:    versionNumber,
		MinServerVersion: "1.0.0",
		ValidatedBy:      "test-user",
		ValidationStatus: models.ValidationStatusPassed,
	}

	err := compatibilityRepo.Create(testCtx, matrix)
	if err != nil {
		t.Fatalf("Failed to create compatibility matrix: %v", err)
	}

	retrieved, err := compatibilityRepo.GetByProductIDAndVersion(testCtx, productID, versionNumber)
	if err != nil {
		t.Fatalf("Failed to get compatibility matrix: %v", err)
	}

	if retrieved.VersionNumber != versionNumber {
		t.Errorf("VersionNumber mismatch: got %s, want %s", retrieved.VersionNumber, versionNumber)
	}

	t.Logf("Retrieved compatibility matrix: %+v", retrieved)
}
