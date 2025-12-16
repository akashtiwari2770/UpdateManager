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

func setupUpgradePathHandlerTest(t *testing.T) (*UpgradePathHandler, *service.ServiceFactory, func()) {
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
	handler := NewUpgradePathHandler(services.UpgradePathService)

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

func TestUpgradePathHandler_CreateUpgradePath(t *testing.T) {
	handler, services, cleanup := setupUpgradePathHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and versions
	productReq := models.CreateProductRequest{
		ProductID: "upgrade-path-product",
		Name:      "Upgrade Path Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create from version
	fromVersionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	fromVersion, err := services.VersionService.CreateVersion(ctx, product.ProductID, &fromVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create from version: %v", err)
	}

	// Create to version
	toVersionReq := models.CreateVersionRequest{
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	_, err = services.VersionService.CreateVersion(ctx, product.ProductID, &toVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create to version: %v", err)
	}

	// Create upgrade path
	path := models.UpgradePath{
		ProductID:   product.ProductID,
		FromVersion: fromVersion.VersionNumber,
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
		IsBlocked:   false,
	}

	body, _ := json.Marshal(path)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products/"+product.ProductID+"/upgrade-paths", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.CreateUpgradePath(w, httpReq)

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

func TestUpgradePathHandler_GetUpgradePath(t *testing.T) {
	handler, services, cleanup := setupUpgradePathHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and versions
	productReq := models.CreateProductRequest{
		ProductID: "get-upgrade-path-product",
		Name:      "Get Upgrade Path Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	fromVersionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	fromVersion, err := services.VersionService.CreateVersion(ctx, product.ProductID, &fromVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create from version: %v", err)
	}

	toVersionReq := models.CreateVersionRequest{
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	_, err = services.VersionService.CreateVersion(ctx, product.ProductID, &toVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create to version: %v", err)
	}

	// Create upgrade path
	path := models.UpgradePath{
		ProductID:   product.ProductID,
		FromVersion: fromVersion.VersionNumber,
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
	}
	if err := services.UpgradePathService.CreateUpgradePath(ctx, &path); err != nil {
		t.Fatalf("Failed to create upgrade path: %v", err)
	}

	// Get upgrade path
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+product.ProductID+"/upgrade-paths/"+fromVersion.VersionNumber+"/2.0.0", nil)
	w := httptest.NewRecorder()
	handler.GetUpgradePath(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestUpgradePathHandler_BlockUpgradePath(t *testing.T) {
	handler, services, cleanup := setupUpgradePathHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and versions
	productReq := models.CreateProductRequest{
		ProductID: "block-upgrade-path-product",
		Name:      "Block Upgrade Path Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	fromVersionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	fromVersion, err := services.VersionService.CreateVersion(ctx, product.ProductID, &fromVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create from version: %v", err)
	}

	toVersionReq := models.CreateVersionRequest{
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	_, err = services.VersionService.CreateVersion(ctx, product.ProductID, &toVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create to version: %v", err)
	}

	// Create upgrade path
	path := models.UpgradePath{
		ProductID:   product.ProductID,
		FromVersion: fromVersion.VersionNumber,
		ToVersion:   "2.0.0",
		PathType:    models.UpgradePathTypeDirect,
	}
	if err := services.UpgradePathService.CreateUpgradePath(ctx, &path); err != nil {
		t.Fatalf("Failed to create upgrade path: %v", err)
	}

	// Block upgrade path
	blockReq := map[string]string{
		"block_reason": "Incompatible changes",
	}
	body, _ := json.Marshal(blockReq)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products/"+product.ProductID+"/upgrade-paths/"+fromVersion.VersionNumber+"/2.0.0/block", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.BlockUpgradePath(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify path is blocked
	updatedPath, err := services.UpgradePathService.GetUpgradePath(ctx, product.ProductID, fromVersion.VersionNumber, "2.0.0")
	if err != nil {
		t.Fatalf("Failed to get upgrade path: %v", err)
	}

	if !updatedPath.IsBlocked {
		t.Error("Expected upgrade path to be blocked")
	}
}
