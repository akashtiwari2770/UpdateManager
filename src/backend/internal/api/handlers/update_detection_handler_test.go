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

func setupUpdateDetectionHandlerTest(t *testing.T) (*UpdateDetectionHandler, *service.ServiceFactory, func()) {
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
	handler := NewUpdateDetectionHandler(services.UpdateDetectionService)

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

func TestUpdateDetectionHandler_DetectUpdate(t *testing.T) {
	handler, services, cleanup := setupUpdateDetectionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "detection-test-product",
		Name:      "Detection Test Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create versions
	currentVersionReq := models.CreateVersionRequest{
		VersionNumber: "0.9.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	currentVersion, err := services.VersionService.CreateVersion(ctx, product.ProductID, &currentVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create current version: %v", err)
	}
	// Release the current version
	_, err = services.VersionService.ReleaseVersion(ctx, currentVersion.ID, "user1")
	if err != nil {
		// If release fails, try to approve first
		_, _ = services.VersionService.ApproveVersion(ctx, currentVersion.ID, &models.ApproveVersionRequest{ApprovedBy: "user1"})
		_, _ = services.VersionService.ReleaseVersion(ctx, currentVersion.ID, "user1")
	}

	availableVersionReq := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	availableVersion, err := services.VersionService.CreateVersion(ctx, product.ProductID, &availableVersionReq, "user1")
	if err != nil {
		t.Fatalf("Failed to create available version: %v", err)
	}
	// Release the available version
	_, err = services.VersionService.SubmitForReview(ctx, availableVersion.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to submit version for review: %v", err)
	}
	_, err = services.VersionService.ApproveVersion(ctx, availableVersion.ID, &models.ApproveVersionRequest{ApprovedBy: "user1"})
	if err != nil {
		t.Fatalf("Failed to approve version: %v", err)
	}
	_, err = services.VersionService.ReleaseVersion(ctx, availableVersion.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to release version: %v", err)
	}

	// Detect update
	detection := models.UpdateDetection{
		EndpointID:       "endpoint-1",
		ProductID:        product.ProductID,
		CurrentVersion:   "0.9.0",
		AvailableVersion: "1.0.0",
		DetectedAt:       time.Now(),
		LastCheckedAt:    time.Now(),
	}

	body, _ := json.Marshal(detection)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/update-detections", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.DetectUpdate(w, httpReq)

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

func TestUpdateDetectionHandler_UpdateAvailableVersion(t *testing.T) {
	handler, services, cleanup := setupUpdateDetectionHandlerTest(t)
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

	// Create and release versions
	version1Req := models.CreateVersionRequest{
		VersionNumber: "1.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version1, err := services.VersionService.CreateVersion(ctx, product.ProductID, &version1Req, "user1")
	if err != nil {
		t.Fatalf("Failed to create version 1.0.0: %v", err)
	}
	_, err = services.VersionService.SubmitForReview(ctx, version1.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to submit version for review: %v", err)
	}
	_, err = services.VersionService.ApproveVersion(ctx, version1.ID, &models.ApproveVersionRequest{ApprovedBy: "user1"})
	if err != nil {
		t.Fatalf("Failed to approve version: %v", err)
	}
	_, err = services.VersionService.ReleaseVersion(ctx, version1.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to release version: %v", err)
	}

	version2Req := models.CreateVersionRequest{
		VersionNumber: "1.1.0",
		ReleaseType:   models.ReleaseTypeFeature,
	}
	version2, err := services.VersionService.CreateVersion(ctx, product.ProductID, &version2Req, "user1")
	if err != nil {
		t.Fatalf("Failed to create version 1.1.0: %v", err)
	}
	_, err = services.VersionService.SubmitForReview(ctx, version2.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to submit version for review: %v", err)
	}
	_, err = services.VersionService.ApproveVersion(ctx, version2.ID, &models.ApproveVersionRequest{ApprovedBy: "user1"})
	if err != nil {
		t.Fatalf("Failed to approve version: %v", err)
	}
	_, err = services.VersionService.ReleaseVersion(ctx, version2.ID, "user1")
	if err != nil {
		t.Fatalf("Failed to release version: %v", err)
	}

	// Create detection
	detection := models.UpdateDetection{
		EndpointID:       "endpoint-2",
		ProductID:        product.ProductID,
		CurrentVersion:   "1.0.0",
		AvailableVersion: "1.1.0",
	}
	_, err = services.UpdateDetectionService.DetectUpdate(ctx, &detection)
	if err != nil {
		t.Fatalf("Failed to create detection: %v", err)
	}

	// Update available version
	updateReq := map[string]string{
		"available_version": "1.1.0",
	}
	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/update-detections/endpoint-2/"+product.ProductID+"/available-version", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.UpdateAvailableVersion(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify update
	updated, err := services.UpdateDetectionService.GetDetection(ctx, "endpoint-2", product.ProductID)
	if err != nil {
		t.Fatalf("Failed to get detection: %v", err)
	}

	if updated.AvailableVersion != "1.1.0" {
		t.Errorf("Expected available_version=1.1.0, got %s", updated.AvailableVersion)
	}
}

func TestUpdateDetectionHandler_ListDetections(t *testing.T) {
	handler, services, cleanup := setupUpdateDetectionHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product
	productReq := models.CreateProductRequest{
		ProductID: "list-detections-product",
		Name:      "List Detections Product",
		Type:      models.ProductTypeServer,
	}
	product, err := services.ProductService.CreateProduct(ctx, &productReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Create and release versions for the product
	versions := []string{"1.0.0", "1.1.0", "2.0.0"}
	for _, versionNum := range versions {
		versionReq := models.CreateVersionRequest{
			VersionNumber: versionNum,
			ReleaseType:   models.ReleaseTypeMajor,
		}
		version, err := services.VersionService.CreateVersion(ctx, product.ProductID, &versionReq, "user1")
		if err != nil {
			t.Fatalf("Failed to create version %s: %v", versionNum, err)
		}
		_, err = services.VersionService.SubmitForReview(ctx, version.ID, "user1")
		if err != nil {
			t.Fatalf("Failed to submit version %s for review: %v", versionNum, err)
		}
		_, err = services.VersionService.ApproveVersion(ctx, version.ID, &models.ApproveVersionRequest{ApprovedBy: "user1"})
		if err != nil {
			t.Fatalf("Failed to approve version %s: %v", versionNum, err)
		}
		_, err = services.VersionService.ReleaseVersion(ctx, version.ID, "user1")
		if err != nil {
			t.Fatalf("Failed to release version %s: %v", versionNum, err)
		}
	}

	// Create other product and versions
	otherProductReq := models.CreateProductRequest{
		ProductID: "other-product",
		Name:      "Other Product",
		Type:      models.ProductTypeServer,
	}
	otherProduct, err := services.ProductService.CreateProduct(ctx, &otherProductReq, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create other product: %v", err)
	}

	otherVersions := []string{"1.0.0", "1.5.0"}
	for _, versionNum := range otherVersions {
		versionReq := models.CreateVersionRequest{
			VersionNumber: versionNum,
			ReleaseType:   models.ReleaseTypeMajor,
		}
		version, err := services.VersionService.CreateVersion(ctx, otherProduct.ProductID, &versionReq, "user1")
		if err != nil {
			t.Fatalf("Failed to create version %s: %v", versionNum, err)
		}
		_, err = services.VersionService.SubmitForReview(ctx, version.ID, "user1")
		if err != nil {
			t.Fatalf("Failed to submit version %s for review: %v", versionNum, err)
		}
		_, err = services.VersionService.ApproveVersion(ctx, version.ID, &models.ApproveVersionRequest{ApprovedBy: "user1"})
		if err != nil {
			t.Fatalf("Failed to approve version %s: %v", versionNum, err)
		}
		_, err = services.VersionService.ReleaseVersion(ctx, version.ID, "user1")
		if err != nil {
			t.Fatalf("Failed to release version %s: %v", versionNum, err)
		}
	}

	// Create multiple detections
	detections := []models.UpdateDetection{
		{
			EndpointID:       "endpoint-list-1",
			ProductID:        product.ProductID,
			CurrentVersion:   "1.0.0",
			AvailableVersion: "1.1.0",
		},
		{
			EndpointID:       "endpoint-list-2",
			ProductID:        product.ProductID,
			CurrentVersion:   "1.1.0",
			AvailableVersion: "2.0.0",
		},
		{
			EndpointID:       "endpoint-list-3",
			ProductID:        "other-product",
			CurrentVersion:   "1.0.0",
			AvailableVersion: "1.5.0",
		},
	}

	for _, detection := range detections {
		_, err = services.UpdateDetectionService.DetectUpdate(ctx, &detection)
		if err != nil {
			t.Fatalf("Failed to create detection: %v", err)
		}
	}

	// Test list all detections
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/update-detections", nil)
	w := httptest.NewRecorder()
	handler.ListDetections(w, httpReq)

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

	data, ok := response.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be array, got %T", response.Data)
	}

	if len(data) < 3 {
		t.Errorf("Expected at least 3 detections, got %d", len(data))
	}

	// Test filter by product_id
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-detections?product_id="+product.ProductID, nil)
	w = httptest.NewRecorder()
	handler.ListDetections(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok = response.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be array, got %T", response.Data)
	}

	if len(data) != 2 {
		t.Errorf("Expected 2 detections for product, got %d", len(data))
	}

	// Test filter by endpoint_id
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-detections?endpoint_id=endpoint-list-1", nil)
	w = httptest.NewRecorder()
	handler.ListDetections(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok = response.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be array, got %T", response.Data)
	}

	if len(data) != 1 {
		t.Errorf("Expected 1 detection for endpoint, got %d", len(data))
	}

	// Test pagination
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-detections?page=1&limit=2", nil)
	w = httptest.NewRecorder()
	handler.ListDetections(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data, ok = response.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected data to be array, got %T", response.Data)
	}

	if len(data) > 2 {
		t.Errorf("Expected at most 2 detections with limit=2, got %d", len(data))
	}
}
