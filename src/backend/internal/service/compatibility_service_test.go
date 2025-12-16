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
	compatibilityServiceTestDB  *database.MongoDB
	compatibilityServiceTestCtx context.Context
	compatibilityService        *CompatibilityService
	compatibilityRepo           *repository.CompatibilityRepository
	compatibilityVersionRepo    *repository.VersionRepository
	compatibilityAuditRepo      *repository.AuditLogRepository
	compatibilityProductRepo    *repository.ProductRepository
)

func setupCompatibilityServiceTestDB(t *testing.T) {
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

	compatibilityServiceTestDB = db
	compatibilityServiceTestCtx = ctx
	compatibilityRepo = repository.NewCompatibilityRepository(db.Collection("compatibility_matrices"))
	compatibilityVersionRepo = repository.NewVersionRepository(db.Collection("versions"))
	compatibilityAuditRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	compatibilityProductRepo = repository.NewProductRepository(db.Collection("products"))
	compatibilityService = NewCompatibilityService(compatibilityRepo, compatibilityVersionRepo, compatibilityAuditRepo)
}

func teardownCompatibilityServiceTestDB(t *testing.T) {
	if compatibilityServiceTestDB != nil {
		_ = compatibilityServiceTestDB.Collection("compatibility_matrices").Drop(compatibilityServiceTestCtx)
		_ = compatibilityServiceTestDB.Collection("versions").Drop(compatibilityServiceTestCtx)
		_ = compatibilityServiceTestDB.Collection("products").Drop(compatibilityServiceTestCtx)
		_ = compatibilityServiceTestDB.Collection("audit_logs").Drop(compatibilityServiceTestCtx)
		_ = compatibilityServiceTestDB.Disconnect(compatibilityServiceTestCtx)
	}
}

func TestCompatibilityService_ValidateCompatibility(t *testing.T) {
	setupCompatibilityServiceTestDB(t)
	defer teardownCompatibilityServiceTestDB(t)

	// Create product and version
	product := &models.Product{
		ProductID: "compat-product",
		Name:      "Compatibility Product",
		Type:      models.ProductTypeClient,
		IsActive:  true,
	}
	compatibilityProductRepo.Create(compatibilityServiceTestCtx, product)

	version := &models.Version{
		ProductID:     product.ProductID,
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateDraft,
		CreatedBy:     "user-123",
	}
	compatibilityVersionRepo.Create(compatibilityServiceTestCtx, version)

	// Validate compatibility
	req := &models.ValidateCompatibilityRequest{
		MinServerVersion:         "1.0.0",
		MaxServerVersion:         "2.0.0",
		RecommendedServerVersion: "1.5.0",
		IncompatibleVersions:     []string{"0.9.0"},
	}

	matrix, err := compatibilityService.ValidateCompatibility(compatibilityServiceTestCtx, product.ProductID, version.VersionNumber, req, "validator-123")
	if err != nil {
		t.Fatalf("Failed to validate compatibility: %v", err)
	}

	if matrix.MinServerVersion != req.MinServerVersion {
		t.Errorf("MinServerVersion mismatch: got %s, want %s", matrix.MinServerVersion, req.MinServerVersion)
	}
	if matrix.ValidationStatus != models.ValidationStatusPassed {
		t.Errorf("ValidationStatus mismatch: got %s, want %s", matrix.ValidationStatus, models.ValidationStatusPassed)
	}

	// Verify audit log
	auditLogs, _ := compatibilityAuditRepo.GetByResource(compatibilityServiceTestCtx, "compatibility_matrix", matrix.ID.Hex(), nil)
	if len(auditLogs) == 0 {
		t.Error("Audit log should be created")
	}

	t.Logf("Validated compatibility: %+v", matrix)
}

func TestCompatibilityService_ValidateCompatibility_VersionNotFound(t *testing.T) {
	setupCompatibilityServiceTestDB(t)
	defer teardownCompatibilityServiceTestDB(t)

	req := &models.ValidateCompatibilityRequest{
		MinServerVersion: "1.0.0",
	}

	_, err := compatibilityService.ValidateCompatibility(compatibilityServiceTestCtx, "non-existent-product", "1.0.0", req, "validator-123")
	if err == nil {
		t.Error("Expected error for non-existent version, got nil")
	}

	t.Logf("Correctly rejected validation for non-existent version: %v", err)
}

func TestCompatibilityService_GetCompatibility(t *testing.T) {
	setupCompatibilityServiceTestDB(t)
	defer teardownCompatibilityServiceTestDB(t)

	// Create product, version, and compatibility matrix
	product := &models.Product{
		ProductID: "get-compat-product",
		Name:      "Get Compatibility Product",
		Type:      models.ProductTypeClient,
		IsActive:  true,
	}
	compatibilityProductRepo.Create(compatibilityServiceTestCtx, product)

	version := &models.Version{
		ProductID:     product.ProductID,
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateDraft,
		CreatedBy:     "user-123",
	}
	compatibilityVersionRepo.Create(compatibilityServiceTestCtx, version)

	matrix := &models.CompatibilityMatrix{
		ProductID:        product.ProductID,
		VersionNumber:    version.VersionNumber,
		MinServerVersion: "1.0.0",
		ValidatedBy:      "validator-123",
		ValidationStatus: models.ValidationStatusPassed,
	}
	compatibilityRepo.Create(compatibilityServiceTestCtx, matrix)

	// Get compatibility
	retrieved, err := compatibilityService.GetCompatibility(compatibilityServiceTestCtx, product.ProductID, version.VersionNumber)
	if err != nil {
		t.Fatalf("Failed to get compatibility: %v", err)
	}

	if retrieved.MinServerVersion != matrix.MinServerVersion {
		t.Errorf("MinServerVersion mismatch: got %s, want %s", retrieved.MinServerVersion, matrix.MinServerVersion)
	}

	t.Logf("Retrieved compatibility: %+v", retrieved)
}
