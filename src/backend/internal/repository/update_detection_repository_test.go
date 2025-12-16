package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var updateDetectionRepo *UpdateDetectionRepository

func setupUpdateDetectionTestDB(t *testing.T) {
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
	updateDetectionRepo = NewUpdateDetectionRepository(db.Collection("update_detections"))
}

func teardownUpdateDetectionTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("update_detections").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestUpdateDetectionCreate(t *testing.T) {
	setupUpdateDetectionTestDB(t)
	defer teardownUpdateDetectionTestDB(t)

	detection := &models.UpdateDetection{
		EndpointID:       "endpoint-123",
		ProductID:        "test-product",
		CurrentVersion:   "1.0.0",
		AvailableVersion: "1.1.0",
	}

	err := updateDetectionRepo.Create(testCtx, detection)
	if err != nil {
		t.Fatalf("Failed to create update detection: %v", err)
	}

	if detection.ID.IsZero() {
		t.Error("Update detection ID was not set after creation")
	}

	t.Logf("Created update detection with ID: %s", detection.ID.Hex())
}

func TestUpdateDetectionUpdateAvailableVersion(t *testing.T) {
	setupUpdateDetectionTestDB(t)
	defer teardownUpdateDetectionTestDB(t)

	detection := &models.UpdateDetection{
		EndpointID:       "endpoint-456",
		ProductID:        "test-product",
		CurrentVersion:   "1.0.0",
		AvailableVersion: "1.1.0",
	}

	err := updateDetectionRepo.Create(testCtx, detection)
	if err != nil {
		t.Fatalf("Failed to create update detection: %v", err)
	}

	newVersion := "1.2.0"
	err = updateDetectionRepo.UpdateAvailableVersion(testCtx, detection.ID, newVersion)
	if err != nil {
		t.Fatalf("Failed to update available version: %v", err)
	}

	retrieved, err := updateDetectionRepo.GetByID(testCtx, detection.ID)
	if err != nil {
		t.Fatalf("Failed to get update detection: %v", err)
	}

	if retrieved.AvailableVersion != newVersion {
		t.Errorf("AvailableVersion not updated: got %s, want %s", retrieved.AvailableVersion, newVersion)
	}

	t.Logf("Updated available version: %+v", retrieved)
}
