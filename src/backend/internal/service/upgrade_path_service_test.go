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
	upgradePathServiceTestDB  *database.MongoDB
	upgradePathServiceTestCtx context.Context
	upgradePathService        *UpgradePathService
	upgradePathRepo           *repository.UpgradePathRepository
	upgradePathVersionRepo    *repository.VersionRepository
	upgradePathProductRepo    *repository.ProductRepository
)

func setupUpgradePathServiceTestDB(t *testing.T) {
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

	upgradePathServiceTestDB = db
	upgradePathServiceTestCtx = ctx
	upgradePathRepo = repository.NewUpgradePathRepository(db.Collection("upgrade_paths"))
	upgradePathVersionRepo = repository.NewVersionRepository(db.Collection("versions"))
	upgradePathProductRepo = repository.NewProductRepository(db.Collection("products"))
	upgradePathService = NewUpgradePathService(upgradePathRepo, upgradePathVersionRepo)
}

func teardownUpgradePathServiceTestDB(t *testing.T) {
	if upgradePathServiceTestDB != nil {
		_ = upgradePathServiceTestDB.Collection("upgrade_paths").Drop(upgradePathServiceTestCtx)
		_ = upgradePathServiceTestDB.Collection("versions").Drop(upgradePathServiceTestCtx)
		_ = upgradePathServiceTestDB.Collection("products").Drop(upgradePathServiceTestCtx)
		_ = upgradePathServiceTestDB.Disconnect(upgradePathServiceTestCtx)
	}
}

func TestUpgradePathService_CreateUpgradePath(t *testing.T) {
	setupUpgradePathServiceTestDB(t)
	defer teardownUpgradePathServiceTestDB(t)

	// Create product and versions
	product := &models.Product{
		ProductID: "upgrade-path-product",
		Name:      "Upgrade Path Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	upgradePathProductRepo.Create(upgradePathServiceTestCtx, product)

	fromVersion := &models.Version{
		ProductID:     product.ProductID,
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateReleased,
		CreatedBy:     "user-123",
	}
	upgradePathVersionRepo.Create(upgradePathServiceTestCtx, fromVersion)

	toVersion := &models.Version{
		ProductID:     product.ProductID,
		VersionNumber: "2.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeMajor,
		State:         models.VersionStateReleased,
		CreatedBy:     "user-123",
	}
	upgradePathVersionRepo.Create(upgradePathServiceTestCtx, toVersion)

	// Create upgrade path
	path := &models.UpgradePath{
		ProductID:   product.ProductID,
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
		IsBlocked:   false,
	}

	err := upgradePathService.CreateUpgradePath(upgradePathServiceTestCtx, path)
	if err != nil {
		t.Fatalf("Failed to create upgrade path: %v", err)
	}

	if path.ID.IsZero() {
		t.Error("Upgrade path ID was not set")
	}

	t.Logf("Created upgrade path: %+v", path)
}

func TestUpgradePathService_CreateUpgradePath_VersionNotFound(t *testing.T) {
	setupUpgradePathServiceTestDB(t)
	defer teardownUpgradePathServiceTestDB(t)

	// Create product
	product := &models.Product{
		ProductID: "upgrade-path-invalid-product",
		Name:      "Upgrade Path Invalid Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	upgradePathProductRepo.Create(upgradePathServiceTestCtx, product)

	path := &models.UpgradePath{
		ProductID:   product.ProductID,
		FromVersion: "non-existent-version",
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
	}

	err := upgradePathService.CreateUpgradePath(upgradePathServiceTestCtx, path)
	if err == nil {
		t.Error("Expected error for non-existent version, got nil")
	}

	t.Logf("Correctly rejected upgrade path with non-existent version: %v", err)
}

func TestUpgradePathService_BlockUpgradePath(t *testing.T) {
	setupUpgradePathServiceTestDB(t)
	defer teardownUpgradePathServiceTestDB(t)

	// Create product, versions, and upgrade path
	product := &models.Product{
		ProductID: "block-path-product",
		Name:      "Block Path Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	upgradePathProductRepo.Create(upgradePathServiceTestCtx, product)

	fromVersion := &models.Version{
		ProductID:     product.ProductID,
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateReleased,
		CreatedBy:     "user-123",
	}
	upgradePathVersionRepo.Create(upgradePathServiceTestCtx, fromVersion)

	toVersion := &models.Version{
		ProductID:     product.ProductID,
		VersionNumber: "2.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeMajor,
		State:         models.VersionStateReleased,
		CreatedBy:     "user-123",
	}
	upgradePathVersionRepo.Create(upgradePathServiceTestCtx, toVersion)

	path := &models.UpgradePath{
		ProductID:   product.ProductID,
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
		IsBlocked:   false,
	}
	upgradePathService.CreateUpgradePath(upgradePathServiceTestCtx, path)

	// Block the path
	reason := "Breaking changes detected"
	err := upgradePathService.BlockUpgradePath(upgradePathServiceTestCtx, product.ProductID, "1.0.0", "2.0.0", reason)
	if err != nil {
		t.Fatalf("Failed to block upgrade path: %v", err)
	}

	// Verify it's blocked
	retrieved, _ := upgradePathRepo.GetByProductIDAndVersions(upgradePathServiceTestCtx, product.ProductID, "1.0.0", "2.0.0")
	if !retrieved.IsBlocked {
		t.Error("Upgrade path should be blocked")
	}
	if retrieved.BlockReason != reason {
		t.Errorf("BlockReason mismatch: got %s, want %s", retrieved.BlockReason, reason)
	}
	if retrieved.PathType != models.UpgradePathTypeBlocked {
		t.Errorf("PathType mismatch: got %s, want %s", retrieved.PathType, models.UpgradePathTypeBlocked)
	}

	t.Logf("Blocked upgrade path: %+v", retrieved)
}
