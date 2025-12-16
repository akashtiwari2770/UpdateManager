package service

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/pkg/database"
)

var (
	detectionServiceTestDB  *database.MongoDB
	detectionServiceTestCtx context.Context
	detectionService        *UpdateDetectionService
	detectionRepo           *repository.UpdateDetectionRepository
	detectionVersionRepo    *repository.VersionRepository
	detectionProductRepo    *repository.ProductRepository
)

func setupDetectionServiceTestDB(t *testing.T) {
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

	detectionServiceTestDB = db
	detectionServiceTestCtx = ctx
	detectionRepo = repository.NewUpdateDetectionRepository(db.Collection("update_detections"))
	detectionVersionRepo = repository.NewVersionRepository(db.Collection("versions"))
	detectionProductRepo = repository.NewProductRepository(db.Collection("products"))
	detectionService = NewUpdateDetectionService(detectionRepo, detectionVersionRepo, detectionProductRepo)
}

func teardownDetectionServiceTestDB(t *testing.T) {
	if detectionServiceTestDB != nil {
		_ = detectionServiceTestDB.Collection("update_detections").Drop(detectionServiceTestCtx)
		_ = detectionServiceTestDB.Collection("versions").Drop(detectionServiceTestCtx)
		_ = detectionServiceTestDB.Disconnect(detectionServiceTestCtx)
	}
}

func TestUpdateDetectionService_DetectUpdate(t *testing.T) {
	setupDetectionServiceTestDB(t)
	defer teardownDetectionServiceTestDB(t)

	detection := &models.UpdateDetection{
		EndpointID:       "endpoint-123",
		ProductID:        "test-product",
		CurrentVersion:   "1.0.0",
		AvailableVersion: "1.1.0",
	}

	detected, err := detectionService.DetectUpdate(detectionServiceTestCtx, detection)
	if err != nil {
		t.Fatalf("Failed to detect update: %v", err)
	}

	if detected.ID.IsZero() {
		t.Error("Detection ID was not set")
	}
	if detected.CurrentVersion != detection.CurrentVersion {
		t.Errorf("CurrentVersion mismatch: got %s, want %s", detected.CurrentVersion, detection.CurrentVersion)
	}

	t.Logf("Detected update: %+v", detected)
}

func TestUpdateDetectionService_DetectUpdate_UpdateExisting(t *testing.T) {
	setupDetectionServiceTestDB(t)
	defer teardownDetectionServiceTestDB(t)

	// Create initial detection
	detection := &models.UpdateDetection{
		EndpointID:       "endpoint-456",
		ProductID:        "test-product",
		CurrentVersion:   "1.0.0",
		AvailableVersion: "1.1.0",
	}
	detectionService.DetectUpdate(detectionServiceTestCtx, detection)

	// Update detection
	updatedDetection := &models.UpdateDetection{
		EndpointID:       "endpoint-456",
		ProductID:        "test-product",
		CurrentVersion:   "1.1.0",
		AvailableVersion: "1.2.0",
	}

	updated, err := detectionService.DetectUpdate(detectionServiceTestCtx, updatedDetection)
	if err != nil {
		t.Fatalf("Failed to update detection: %v", err)
	}

	if updated.AvailableVersion != "1.2.0" {
		t.Errorf("AvailableVersion not updated: got %s, want 1.2.0", updated.AvailableVersion)
	}

	t.Logf("Updated detection: %+v", updated)
}

func TestUpdateDetectionService_UpdateAvailableVersion(t *testing.T) {
	setupDetectionServiceTestDB(t)
	defer teardownDetectionServiceTestDB(t)

	// Create detection
	detection := &models.UpdateDetection{
		EndpointID:       "endpoint-789",
		ProductID:        "test-product",
		CurrentVersion:   "1.0.0",
		AvailableVersion: "1.1.0",
	}
	detected, _ := detectionService.DetectUpdate(detectionServiceTestCtx, detection)

	// Update available version
	err := detectionService.UpdateAvailableVersion(detectionServiceTestCtx, detected.EndpointID, detected.ProductID, "1.2.0")
	if err != nil {
		t.Fatalf("Failed to update available version: %v", err)
	}

	// Verify update
	retrieved, _ := detectionRepo.GetByEndpointIDAndProductID(detectionServiceTestCtx, detected.EndpointID, detected.ProductID)
	if retrieved.AvailableVersion != "1.2.0" {
		t.Errorf("AvailableVersion not updated: got %s, want 1.2.0", retrieved.AvailableVersion)
	}

	t.Logf("Updated available version: %+v", retrieved)
}
