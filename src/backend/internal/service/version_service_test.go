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
	versionServiceTestDB  *database.MongoDB
	versionServiceTestCtx context.Context
	versionService        *VersionService
	versionRepo           *repository.VersionRepository
	versionProductRepo    *repository.ProductRepository
	versionAuditRepo      *repository.AuditLogRepository
)

func setupVersionServiceTestDB(t *testing.T) {
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

	versionServiceTestDB = db
	versionServiceTestCtx = ctx
	versionRepo = repository.NewVersionRepository(db.Collection("versions"))
	versionProductRepo = repository.NewProductRepository(db.Collection("products"))
	versionAuditRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	versionService = NewVersionService(versionRepo, versionProductRepo, versionAuditRepo)
}

func teardownVersionServiceTestDB(t *testing.T) {
	if versionServiceTestDB != nil {
		_ = versionServiceTestDB.Collection("versions").Drop(versionServiceTestCtx)
		_ = versionServiceTestDB.Collection("products").Drop(versionServiceTestCtx)
		_ = versionServiceTestDB.Collection("audit_logs").Drop(versionServiceTestCtx)
		_ = versionServiceTestDB.Disconnect(versionServiceTestCtx)
	}
}

func TestVersionService_CreateVersion(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product first
	product := &models.Product{
		ProductID: "version-test-product",
		Name:      "Version Test Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	err := versionProductRepo.Create(versionServiceTestCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create version
	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}

	version, err := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Verify version
	if version.ID.IsZero() {
		t.Error("Version ID was not set")
	}
	if version.VersionNumber != req.VersionNumber {
		t.Errorf("VersionNumber mismatch: got %s, want %s", version.VersionNumber, req.VersionNumber)
	}
	if version.State != models.VersionStateDraft {
		t.Errorf("State mismatch: got %s, want %s", version.State, models.VersionStateDraft)
	}

	// Verify audit log
	auditLogs, _ := versionAuditRepo.GetByResource(versionServiceTestCtx, "version", version.ID.Hex(), nil)
	if len(auditLogs) == 0 {
		t.Error("Audit log should be created")
	}

	t.Logf("Created version: %+v", version)
}

func TestVersionService_CreateVersion_ProductNotFound(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}

	_, err := versionService.CreateVersion(versionServiceTestCtx, "non-existent-product", req, "user-123")
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}

	t.Logf("Correctly rejected version creation for non-existent product: %v", err)
}

func TestVersionService_CreateVersion_DuplicateVersion(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product
	product := &models.Product{
		ProductID: "duplicate-version-product",
		Name:      "Duplicate Version Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	err := versionProductRepo.Create(versionServiceTestCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create first version
	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	_, err = versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")
	if err != nil {
		t.Fatalf("Failed to create first version: %v", err)
	}

	// Try to create duplicate
	_, err = versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")
	if err == nil {
		t.Error("Expected error for duplicate version, got nil")
	}

	t.Logf("Correctly rejected duplicate version: %v", err)
}

func TestVersionService_SubmitForReview(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and version
	product := &models.Product{
		ProductID: "review-product",
		Name:      "Review Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Submit for review
	submitted, err := versionService.SubmitForReview(versionServiceTestCtx, version.ID, "user-123")
	if err != nil {
		t.Fatalf("Failed to submit for review: %v", err)
	}

	if submitted.State != models.VersionStatePendingReview {
		t.Errorf("State mismatch: got %s, want %s", submitted.State, models.VersionStatePendingReview)
	}

	t.Logf("Submitted version for review: %+v", submitted)
}

func TestVersionService_SubmitForReview_InvalidState(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and approved version
	product := &models.Product{
		ProductID: "invalid-state-product",
		Name:      "Invalid State Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Submit for review first, then approve
	versionService.SubmitForReview(versionServiceTestCtx, version.ID, "user-123")
	approveReq := &models.ApproveVersionRequest{ApprovedBy: "admin-123"}
	versionService.ApproveVersion(versionServiceTestCtx, version.ID, approveReq)

	// Verify version is now approved
	approvedVersion, _ := versionService.GetVersion(versionServiceTestCtx, version.ID)
	if approvedVersion.State != models.VersionStateApproved {
		t.Fatalf("Version should be approved, got state: %s", approvedVersion.State)
	}

	// Try to submit approved version for review (should fail)
	_, err := versionService.SubmitForReview(versionServiceTestCtx, version.ID, "user-123")
	if err == nil {
		t.Error("Expected error for submitting approved version, got nil")
	}

	t.Logf("Correctly rejected invalid state transition: %v", err)
}

func TestVersionService_ApproveVersion(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and version
	product := &models.Product{
		ProductID: "approve-product",
		Name:      "Approve Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Submit for review
	versionService.SubmitForReview(versionServiceTestCtx, version.ID, "user-123")

	// Approve version
	approveReq := &models.ApproveVersionRequest{ApprovedBy: "admin-123"}
	approved, err := versionService.ApproveVersion(versionServiceTestCtx, version.ID, approveReq)
	if err != nil {
		t.Fatalf("Failed to approve version: %v", err)
	}

	if approved.State != models.VersionStateApproved {
		t.Errorf("State mismatch: got %s, want %s", approved.State, models.VersionStateApproved)
	}
	if approved.ApprovedBy != "admin-123" {
		t.Errorf("ApprovedBy mismatch: got %s, want admin-123", approved.ApprovedBy)
	}
	if approved.ApprovedAt == nil {
		t.Error("ApprovedAt should be set")
	}

	// Verify audit log
	auditLogs, _ := versionAuditRepo.GetByResource(versionServiceTestCtx, "version", version.ID.Hex(), nil)
	approveFound := false
	for _, log := range auditLogs {
		if log.Action == models.AuditActionApprove {
			approveFound = true
			break
		}
	}
	if !approveFound {
		t.Error("Approve audit log should be created")
	}

	t.Logf("Approved version: %+v", approved)
}

func TestVersionService_ApproveVersion_InvalidState(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and draft version
	product := &models.Product{
		ProductID: "approve-invalid-product",
		Name:      "Approve Invalid Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Try to approve draft version (should fail)
	approveReq := &models.ApproveVersionRequest{ApprovedBy: "admin-123"}
	_, err := versionService.ApproveVersion(versionServiceTestCtx, version.ID, approveReq)
	if err == nil {
		t.Error("Expected error for approving draft version, got nil")
	}

	t.Logf("Correctly rejected invalid state for approval: %v", err)
}

func TestVersionService_ReleaseVersion(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and version
	product := &models.Product{
		ProductID: "release-product",
		Name:      "Release Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Submit, approve, then release
	versionService.SubmitForReview(versionServiceTestCtx, version.ID, "user-123")
	approveReq := &models.ApproveVersionRequest{ApprovedBy: "admin-123"}
	versionService.ApproveVersion(versionServiceTestCtx, version.ID, approveReq)

	// Release version
	released, err := versionService.ReleaseVersion(versionServiceTestCtx, version.ID, "admin-123")
	if err != nil {
		t.Fatalf("Failed to release version: %v", err)
	}

	if released.State != models.VersionStateReleased {
		t.Errorf("State mismatch: got %s, want %s", released.State, models.VersionStateReleased)
	}

	// Verify audit log
	auditLogs, _ := versionAuditRepo.GetByResource(versionServiceTestCtx, "version", version.ID.Hex(), nil)
	releaseFound := false
	for _, log := range auditLogs {
		if log.Action == models.AuditActionRelease {
			releaseFound = true
			break
		}
	}
	if !releaseFound {
		t.Error("Release audit log should be created")
	}

	t.Logf("Released version: %+v", released)
}

func TestVersionService_UpdateVersion(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and version
	product := &models.Product{
		ProductID: "update-version-product",
		Name:      "Update Version Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Update version
	newReleaseType := models.ReleaseTypeSecurity
	updateReq := &models.UpdateVersionRequest{
		ReleaseType: &newReleaseType,
	}

	updated, err := versionService.UpdateVersion(versionServiceTestCtx, version.ID, updateReq, "user-123")
	if err != nil {
		t.Fatalf("Failed to update version: %v", err)
	}

	if updated.ReleaseType != models.ReleaseTypeSecurity {
		t.Errorf("ReleaseType not updated: got %s, want %s", updated.ReleaseType, models.ReleaseTypeSecurity)
	}

	t.Logf("Updated version: %+v", updated)
}

func TestVersionService_UpdateVersion_NonDraft(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product and approved version
	product := &models.Product{
		ProductID: "update-nondraft-product",
		Name:      "Update Non-Draft Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	req := &models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version, _ := versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")

	// Submit and approve
	versionService.SubmitForReview(versionServiceTestCtx, version.ID, "user-123")
	approveReq := &models.ApproveVersionRequest{ApprovedBy: "admin-123"}
	versionService.ApproveVersion(versionServiceTestCtx, version.ID, approveReq)

	// Try to update approved version (should fail)
	newReleaseType := models.ReleaseTypeSecurity
	updateReq := &models.UpdateVersionRequest{
		ReleaseType: &newReleaseType,
	}
	_, err := versionService.UpdateVersion(versionServiceTestCtx, version.ID, updateReq, "user-123")
	if err == nil {
		t.Error("Expected error for updating non-draft version, got nil")
	}

	t.Logf("Correctly rejected update of non-draft version: %v", err)
}

func TestVersionService_GetVersionsByProduct(t *testing.T) {
	setupVersionServiceTestDB(t)
	defer teardownVersionServiceTestDB(t)

	// Create product
	product := &models.Product{
		ProductID: "list-versions-product",
		Name:      "List Versions Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	versionProductRepo.Create(versionServiceTestCtx, product)

	// Create multiple versions
	versions := []string{"1.0.0", "1.1.0", "2.0.0"}
	for _, v := range versions {
		req := &models.CreateVersionRequest{
			VersionNumber: v,
			ReleaseDate:   time.Now(),
			ReleaseType:   models.ReleaseTypeFeature,
		}
		versionService.CreateVersion(versionServiceTestCtx, product.ProductID, req, "user-123")
	}

	// List versions
	listVersions, total, err := versionService.GetVersionsByProduct(versionServiceTestCtx, product.ProductID, 1, 10)
	if err != nil {
		t.Fatalf("Failed to get versions: %v", err)
	}

	if len(listVersions) < len(versions) {
		t.Errorf("Expected at least %d versions, got %d", len(versions), len(listVersions))
	}

	if total < int64(len(versions)) {
		t.Errorf("Expected total at least %d, got %d", len(versions), total)
	}

	t.Logf("Listed %d versions (total: %d)", len(listVersions), total)
}
