package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
	"updatemanager/pkg/database"
)

func setupVersionHandlerTest(t *testing.T) (*VersionHandler, *service.ServiceFactory, func()) {
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

	// Clean up before test to ensure idempotency
	_ = db.Collection("products").Drop(ctx)
	_ = db.Collection("versions").Drop(ctx)
	_ = db.Collection("compatibility_matrices").Drop(ctx)
	_ = db.Collection("upgrade_paths").Drop(ctx)
	_ = db.Collection("notifications").Drop(ctx)
	_ = db.Collection("update_detections").Drop(ctx)
	_ = db.Collection("update_rollouts").Drop(ctx)
	_ = db.Collection("audit_logs").Drop(ctx)

	services := service.NewServiceFactory(db.Database)
	handler := NewVersionHandler(services.VersionService)

	cleanup := func() {
		// Drop all test collections to make tests idempotent
		_ = db.Collection("products").Drop(ctx)
		_ = db.Collection("versions").Drop(ctx)
		_ = db.Collection("compatibility_matrices").Drop(ctx)
		_ = db.Collection("upgrade_paths").Drop(ctx)
		_ = db.Collection("notifications").Drop(ctx)
		_ = db.Collection("update_detections").Drop(ctx)
		_ = db.Collection("update_rollouts").Drop(ctx)
		_ = db.Collection("audit_logs").Drop(ctx)
		db.Disconnect(ctx)
	}

	return handler, services, cleanup
}

func TestVersionHandler_CreateVersion(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a product first
	productReq := models.CreateProductRequest{
		ProductID: "version-test-product",
		Name:      "Version Test Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create version
	versionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
		ReleaseDate:   time.Now(),
	}

	body, _ := json.Marshal(versionReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products/"+product.ProductID+"/versions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "user1")

	w := httptest.NewRecorder()
	handler.CreateVersion(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}
}

func TestVersionHandler_CreateVersion_ProductNotFound(t *testing.T) {
	handler, _, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	versionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}

	body, _ := json.Marshal(versionReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products/nonexistent/versions", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "user1")

	w := httptest.NewRecorder()
	handler.CreateVersion(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestVersionHandler_GetVersion(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "get-version-product",
		Name:      "Get Version Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	versionReq := models.CreateVersionRequest{
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Get version
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/versions/"+version.ID.Hex(), nil)
	w := httptest.NewRecorder()
	handler.GetVersion(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestVersionHandler_GetVersionsByProduct(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product
	productReq := models.CreateProductRequest{
		ProductID: "list-versions-product",
		Name:      "List Versions Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create multiple versions
	for i := 0; i < 3; i++ {
		versionReq := models.CreateVersionRequest{
			VersionNumber: "1." + string(rune('0'+i)) + ".0",
			ReleaseType:   models.ReleaseTypeFeature,
		}
		_, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}
	}

	// Get versions by product
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+product.ProductID+"/versions?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.GetVersionsByProduct(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 3 {
		t.Errorf("Expected at least 3 versions, got %d", response.Meta.Total)
	}
}

func TestVersionHandler_ListVersions(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple products and versions
	product1Req := models.CreateProductRequest{
		ProductID: "list-versions-product1",
		Name:      "List Versions Product 1",
		Type:      models.ProductTypeServer,
	}
	product1, err := services.ProductService.CreateProduct(ctx, &product1Req, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	product2Req := models.CreateProductRequest{
		ProductID: "list-versions-product2",
		Name:      "List Versions Product 2",
		Type:      models.ProductTypeClient,
	}
	product2, err := services.ProductService.CreateProduct(ctx, &product2Req, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create versions for product1
	for i := 0; i < 2; i++ {
		versionReq := models.CreateVersionRequest{
			VersionNumber: "1." + string(rune('0'+i)) + ".0",
			ReleaseType:   models.ReleaseTypeFeature,
		}
		_, err := services.VersionService.CreateVersion(ctx, product1.ProductID, &versionReq, "user1")
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}
	}

	// Create versions for product2
	for i := 0; i < 2; i++ {
		versionReq := models.CreateVersionRequest{
			VersionNumber: "2." + string(rune('0'+i)) + ".0",
			ReleaseType:   models.ReleaseTypeMaintenance,
		}
		_, err := services.VersionService.CreateVersion(ctx, product2.ProductID, &versionReq, "user1")
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}
	}

	// Test list all versions
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/versions?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.ListVersions(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 4 {
		t.Errorf("Expected at least 4 versions, got %d", response.Meta.Total)
	}

	// Test filter by product_id
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/versions?product_id="+product1.ProductID+"&page=1&limit=10", nil)
	w = httptest.NewRecorder()
	handler.ListVersions(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 2 {
		t.Errorf("Expected at least 2 versions for product1, got %d", response.Meta.Total)
	}

	// Test filter by state
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/versions?state=draft&page=1&limit=10", nil)
	w = httptest.NewRecorder()
	handler.ListVersions(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 4 {
		t.Errorf("Expected at least 4 draft versions, got %d", response.Meta.Total)
	}

	// Test filter by release_type
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/versions?release_type=feature&page=1&limit=10", nil)
	w = httptest.NewRecorder()
	handler.ListVersions(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 2 {
		t.Errorf("Expected at least 2 feature versions, got %d", response.Meta.Total)
	}
}

func TestVersionHandler_SubmitForReview(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "submit-review-product",
		Name:      "Submit Review Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	versionReq := models.CreateVersionRequest{
		VersionNumber: "3.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Submit for review
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/versions/"+version.ID.Hex()+"/submit", nil)
	httpReq.Header.Set("X-User-ID", "user1")

	w := httptest.NewRecorder()
	handler.SubmitForReview(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify state changed
	updated, err := services.VersionService.GetVersion(ctx, version.ID)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if updated.State != models.VersionStatePendingReview {
		t.Errorf("Expected state=%s, got %s", models.VersionStatePendingReview, updated.State)
	}
}

func TestVersionHandler_ApproveVersion(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "approve-version-product",
		Name:      "Approve Version Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	versionReq := models.CreateVersionRequest{
		VersionNumber: "4.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Submit for review first
	_, err = services.VersionService.SubmitForReview(ctx, version.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to submit for review: %v", err)
	}

	// Approve version
	approveReq := models.ApproveVersionRequest{
		ApprovedBy: "admin",
	}

	body, _ := json.Marshal(approveReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/versions/"+version.ID.Hex()+"/approve", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.ApproveVersion(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify state changed
	updated, err := services.VersionService.GetVersion(ctx, version.ID)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if updated.State != models.VersionStateApproved {
		t.Errorf("Expected state=%s, got %s", models.VersionStateApproved, updated.State)
	}
}

func TestVersionHandler_ReleaseVersion(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "release-version-product",
		Name:      "Release Version Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	versionReq := models.CreateVersionRequest{
		VersionNumber: "5.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Submit and approve first
	_, err = services.VersionService.SubmitForReview(ctx, version.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to submit for review: %v", err)
	}

	approveReq := models.ApproveVersionRequest{
		ApprovedBy: "admin",
	}
	_, err = services.VersionService.ApproveVersion(ctx, version.ID, &approveReq)
	if err != nil {
		t.Fatalf("Failed to approve version: %v", err)
	}

	// Release version
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/versions/"+version.ID.Hex()+"/release", nil)
	httpReq.Header.Set("X-User-ID", "user1")

	w := httptest.NewRecorder()
	handler.ReleaseVersion(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify state changed
	updated, err := services.VersionService.GetVersion(ctx, version.ID)
	if err != nil {
		t.Fatalf("Failed to get version: %v", err)
	}

	if updated.State != models.VersionStateReleased {
		t.Errorf("Expected state=%s, got %s", models.VersionStateReleased, updated.State)
	}
}

func TestVersionHandler_UpdateVersion(t *testing.T) {
	handler, services, cleanup := setupVersionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "update-version-product",
		Name:      "Update Version Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	versionReq := models.CreateVersionRequest{
		VersionNumber: "6.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Update version (only draft can be updated)
	releaseType := models.ReleaseTypeFeature
	updateReq := models.UpdateVersionRequest{
		ReleaseType: &releaseType,
	}

	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/versions/"+version.ID.Hex(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "user1")

	w := httptest.NewRecorder()
	handler.UpdateVersion(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}
