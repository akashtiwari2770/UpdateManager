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

func setupCompatibilityHandlerTest(t *testing.T) (*CompatibilityHandler, *service.ServiceFactory, func()) {
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
	handler := NewCompatibilityHandler(services.CompatibilityService)

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

func TestCompatibilityHandler_ValidateCompatibility(t *testing.T) {
	handler, services, cleanup := setupCompatibilityHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "compat-test-product",
		Name:      "Compatibility Test Product",
		Type:      models.ProductTypeClient,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	versionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create version: %v", err)
	}

	// Validate compatibility
	compatReq := models.ValidateCompatibilityRequest{
		MinServerVersion:         "1.0.0",
		MaxServerVersion:         "2.0.0",
		RecommendedServerVersion: "1.5.0",
		IncompatibleVersions:     []string{"0.9.0"},
	}

	body, _ := json.Marshal(compatReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products/"+product.ProductID+"/versions/"+version.VersionNumber+"/compatibility", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "validator")

	w := httptest.NewRecorder()
	handler.ValidateCompatibility(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}
}

func TestCompatibilityHandler_GetCompatibility(t *testing.T) {
	handler, services, cleanup := setupCompatibilityHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "get-compat-product",
		Name:      "Get Compatibility Product",
		Type:      models.ProductTypeClient,
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

	// Create compatibility matrix
	compatReq := models.ValidateCompatibilityRequest{
		MinServerVersion: "2.0.0",
		MaxServerVersion: "3.0.0",
	}
	_, err = services.CompatibilityService.ValidateCompatibility(ctx, product.ProductID, version.VersionNumber, &compatReq, "validator")
	if err != nil {
		t.Fatalf("Failed to create compatibility matrix: %v", err)
	}

	// Get compatibility
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+product.ProductID+"/versions/"+version.VersionNumber+"/compatibility", nil)
	w := httptest.NewRecorder()
	handler.GetCompatibility(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestCompatibilityHandler_ListCompatibility(t *testing.T) {
	handler, services, cleanup := setupCompatibilityHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and versions
	productReq := models.CreateProductRequest{
		ProductID: "list-compat-product",
		Name:      "List Compatibility Product",
		Type:      models.ProductTypeClient,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create multiple versions with compatibility
	for i := 0; i < 3; i++ {
		versionReq := models.CreateVersionRequest{
			VersionNumber: "3." + string(rune('0'+i)) + ".0",
			ReleaseType:   models.ReleaseTypeFeature,
		}
		version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
		if err != nil {
			t.Fatalf("Failed to create version: %v", err)
		}

		compatReq := models.ValidateCompatibilityRequest{
			MinServerVersion: "1.0.0",
			MaxServerVersion: "2.0.0",
		}
		_, err = services.CompatibilityService.ValidateCompatibility(ctx, product.ProductID, version.VersionNumber, &compatReq, "validator")
		if err != nil {
			t.Fatalf("Failed to create compatibility: %v", err)
		}
	}

	// List compatibility
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/compatibility?product_id="+product.ProductID+"&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.ListCompatibility(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 3 {
		t.Errorf("Expected at least 3 compatibility matrices, got %d", response.Meta.Total)
	}
}
