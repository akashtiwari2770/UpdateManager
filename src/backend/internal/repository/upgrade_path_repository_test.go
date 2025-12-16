package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var upgradePathRepo *UpgradePathRepository

func setupUpgradePathTestDB(t *testing.T) {
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
	upgradePathRepo = NewUpgradePathRepository(db.Collection("upgrade_paths"))
}

func teardownUpgradePathTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("upgrade_paths").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestUpgradePathCreate(t *testing.T) {
	setupUpgradePathTestDB(t)
	defer teardownUpgradePathTestDB(t)

	path := &models.UpgradePath{
		ProductID:   "test-product",
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
		IsBlocked:   false,
	}

	err := upgradePathRepo.Create(testCtx, path)
	if err != nil {
		t.Fatalf("Failed to create upgrade path: %v", err)
	}

	if path.ID.IsZero() {
		t.Error("Upgrade path ID was not set after creation")
	}

	t.Logf("Created upgrade path with ID: %s", path.ID.Hex())
}

func TestUpgradePathGetByProductIDAndVersions(t *testing.T) {
	setupUpgradePathTestDB(t)
	defer teardownUpgradePathTestDB(t)

	productID := "test-product"
	fromVersion := "1.0.0"
	toVersion := "2.0.0"
	path := &models.UpgradePath{
		ProductID:   productID,
		FromVersion: fromVersion,
		ToVersion:   toVersion,
		PathType:    models.UpgradePathTypeDirect,
	}

	err := upgradePathRepo.Create(testCtx, path)
	if err != nil {
		t.Fatalf("Failed to create upgrade path: %v", err)
	}

	retrieved, err := upgradePathRepo.GetByProductIDAndVersions(testCtx, productID, fromVersion, toVersion)
	if err != nil {
		t.Fatalf("Failed to get upgrade path: %v", err)
	}

	if retrieved.FromVersion != fromVersion {
		t.Errorf("FromVersion mismatch: got %s, want %s", retrieved.FromVersion, fromVersion)
	}

	t.Logf("Retrieved upgrade path: %+v", retrieved)
}
