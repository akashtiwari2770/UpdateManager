package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var versionRepo *VersionRepository

func setupVersionTestDB(t *testing.T) {
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
	versionRepo = NewVersionRepository(db.Collection("versions"))
}

func teardownVersionTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("versions").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestVersionCreate(t *testing.T) {
	setupVersionTestDB(t)
	defer teardownVersionTestDB(t)

	version := &models.Version{
		ProductID:     "test-product",
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateDraft,
		CreatedBy:     "test-user",
	}

	err := versionRepo.Create(testCtx, version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	if version.ID.IsZero() {
		t.Error("Version ID was not set after creation")
	}

	t.Logf("Created version with ID: %s", version.ID.Hex())
}

func TestVersionGetByID(t *testing.T) {
	setupVersionTestDB(t)
	defer teardownVersionTestDB(t)

	version := &models.Version{
		ProductID:     "test-product",
		VersionNumber: "1.0.1",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeSecurity,
		State:         models.VersionStateDraft,
		CreatedBy:     "test-user",
	}

	err := versionRepo.Create(testCtx, version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	retrieved, err := versionRepo.GetByID(testCtx, version.ID)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if retrieved.VersionNumber != version.VersionNumber {
		t.Errorf("VersionNumber mismatch: got %s, want %s", retrieved.VersionNumber, version.VersionNumber)
	}

	t.Logf("Retrieved version: %+v", retrieved)
}

func TestVersionGetByProductIDAndVersion(t *testing.T) {
	setupVersionTestDB(t)
	defer teardownVersionTestDB(t)

	productID := "test-product-2"
	versionNumber := "2.0.0"
	version := &models.Version{
		ProductID:     productID,
		VersionNumber: versionNumber,
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeMajor,
		State:         models.VersionStateDraft,
		CreatedBy:     "test-user",
	}

	err := versionRepo.Create(testCtx, version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	retrieved, err := versionRepo.GetByProductIDAndVersion(testCtx, productID, versionNumber)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if retrieved.VersionNumber != versionNumber {
		t.Errorf("VersionNumber mismatch: got %s, want %s", retrieved.VersionNumber, versionNumber)
	}

	t.Logf("Retrieved version by product_id and version: %+v", retrieved)
}

func TestVersionUpdateState(t *testing.T) {
	setupVersionTestDB(t)
	defer teardownVersionTestDB(t)

	version := &models.Version{
		ProductID:     "test-product",
		VersionNumber: "1.0.2",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateDraft,
		CreatedBy:     "test-user",
	}

	err := versionRepo.Create(testCtx, version)
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	err = versionRepo.UpdateState(testCtx, version.ID, models.VersionStateApproved, "admin-user")
	if err != nil {
		t.Fatalf("Failed to update version state: %v", err)
	}

	retrieved, err := versionRepo.GetByID(testCtx, version.ID)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if retrieved.State != models.VersionStateApproved {
		t.Errorf("State not updated: got %s, want %s", retrieved.State, models.VersionStateApproved)
	}

	if retrieved.ApprovedBy != "admin-user" {
		t.Errorf("ApprovedBy not set: got %s, want admin-user", retrieved.ApprovedBy)
	}

	if retrieved.ApprovedAt == nil {
		t.Error("ApprovedAt was not set")
	}

	t.Logf("Updated version state: %+v", retrieved)
}

func TestVersionList(t *testing.T) {
	setupVersionTestDB(t)
	defer teardownVersionTestDB(t)

	productID := "test-product-3"
	versions := []*models.Version{
		{ProductID: productID, VersionNumber: "1.0.0", ReleaseType: models.ReleaseTypeFeature, State: models.VersionStateDraft, CreatedBy: "user1"},
		{ProductID: productID, VersionNumber: "1.1.0", ReleaseType: models.ReleaseTypeFeature, State: models.VersionStateApproved, CreatedBy: "user1"},
		{ProductID: productID, VersionNumber: "2.0.0", ReleaseType: models.ReleaseTypeMajor, State: models.VersionStateReleased, CreatedBy: "user1"},
	}

	for _, v := range versions {
		v.ReleaseDate = time.Now()
		err := versionRepo.Create(testCtx, v)
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}
	}

	allVersions, err := versionRepo.GetByProductID(testCtx, productID, nil)
	if err != nil {
		t.Fatalf("Failed to list versions: %v", err)
	}

	if len(allVersions) < len(versions) {
		t.Errorf("Expected at least %d versions, got %d", len(versions), len(allVersions))
	}

	t.Logf("Listed %d versions for product %s", len(allVersions), productID)
}
