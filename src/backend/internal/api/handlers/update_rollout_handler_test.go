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

func setupUpdateRolloutHandlerTest(t *testing.T) (*UpdateRolloutHandler, *service.ServiceFactory, func()) {
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
	handler := NewUpdateRolloutHandler(services.UpdateRolloutService)

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

func TestUpdateRolloutHandler_InitiateRollout(t *testing.T) {
	handler, services, cleanup := setupUpdateRolloutHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and version
	productReq := models.CreateProductRequest{
		ProductID: "rollout-test-product",
		Name:      "Rollout Test Product",
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
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version2, err := services.VersionService.CreateVersion(ctx, product.ProductID, &version2Req, "user1")
	if err != nil {
		t.Fatalf("Failed to create version 2.0.0: %v", err)
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
		EndpointID:       "endpoint-rollout",
		ProductID:        product.ProductID,
		CurrentVersion:   "1.0.0",
		AvailableVersion: "2.0.0",
	}
	_, err = services.UpdateDetectionService.DetectUpdate(ctx, &detection)
	if err != nil {
		t.Fatalf("Failed to create detection: %v", err)
	}

	// Initiate rollout
	rollout := models.UpdateRollout{
		EndpointID:  "endpoint-rollout",
		ProductID:   product.ProductID,
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusPending,
		Progress:    0,
	}

	body, _ := json.Marshal(rollout)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/update-rollouts", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "user1")

	w := httptest.NewRecorder()
	handler.InitiateRollout(w, httpReq)

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

func TestUpdateRolloutHandler_UpdateRolloutStatus(t *testing.T) {
	handler, services, cleanup := setupUpdateRolloutHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and detection
	productReq := models.CreateProductRequest{
		ProductID: "rollout-status-product",
		Name:      "Rollout Status Product",
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
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version2, err := services.VersionService.CreateVersion(ctx, product.ProductID, &version2Req, "user1")
	if err != nil {
		t.Fatalf("Failed to create version 2.0.0: %v", err)
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

	detection := models.UpdateDetection{
		EndpointID:       "endpoint-status",
		ProductID:        product.ProductID,
		CurrentVersion:   "1.0.0",
		AvailableVersion: "2.0.0",
	}
	_, err = services.UpdateDetectionService.DetectUpdate(ctx, &detection)
	if err != nil {
		t.Fatalf("Failed to create detection: %v", err)
	}

	// Create rollout
	rollout := models.UpdateRollout{
		EndpointID:  "endpoint-status",
		ProductID:   product.ProductID,
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusPending,
		InitiatedBy: "user1",
	}
	result, err := services.UpdateRolloutService.InitiateRollout(ctx, &rollout)
	if err != nil {
		t.Fatalf("Failed to create rollout: %v", err)
	}

	// Update status
	statusReq := map[string]interface{}{
		"status": models.RolloutStatusInProgress,
	}
	body, _ := json.Marshal(statusReq)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/update-rollouts/"+result.ID.Hex()+"/status", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.UpdateRolloutStatus(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify status updated
	updated, err := services.UpdateRolloutService.GetRollout(ctx, result.ID)
	if err != nil {
		t.Fatalf("Failed to get rollout: %v", err)
	}

	if updated.Status != models.RolloutStatusInProgress {
		t.Errorf("Expected status=%s, got %s", models.RolloutStatusInProgress, updated.Status)
	}
}

func TestUpdateRolloutHandler_UpdateRolloutProgress(t *testing.T) {
	handler, services, cleanup := setupUpdateRolloutHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and detection
	productReq := models.CreateProductRequest{
		ProductID: "rollout-progress-product",
		Name:      "Rollout Progress Product",
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
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version2, err := services.VersionService.CreateVersion(ctx, product.ProductID, &version2Req, "user1")
	if err != nil {
		t.Fatalf("Failed to create version 2.0.0: %v", err)
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

	detection := models.UpdateDetection{
		EndpointID:       "endpoint-progress",
		ProductID:        product.ProductID,
		CurrentVersion:   "1.0.0",
		AvailableVersion: "2.0.0",
	}
	_, err = services.UpdateDetectionService.DetectUpdate(ctx, &detection)
	if err != nil {
		t.Fatalf("Failed to create detection: %v", err)
	}

	// Create rollout
	rollout := models.UpdateRollout{
		EndpointID:  "endpoint-progress",
		ProductID:   product.ProductID,
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusInProgress,
		InitiatedBy: "user1",
	}
	result, err := services.UpdateRolloutService.InitiateRollout(ctx, &rollout)
	if err != nil {
		t.Fatalf("Failed to create rollout: %v", err)
	}

	// Update progress
	progressReq := map[string]int{
		"progress": 50,
	}
	body, _ := json.Marshal(progressReq)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/update-rollouts/"+result.ID.Hex()+"/progress", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.UpdateRolloutProgress(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify progress updated
	updated, err := services.UpdateRolloutService.GetRollout(ctx, result.ID)
	if err != nil {
		t.Fatalf("Failed to get rollout: %v", err)
	}

	if updated.Progress != 50 {
		t.Errorf("Expected progress=50, got %d", updated.Progress)
	}
}

func TestUpdateRolloutHandler_GetRollout(t *testing.T) {
	handler, services, cleanup := setupUpdateRolloutHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product and detection
	productReq := models.CreateProductRequest{
		ProductID: "get-rollout-product",
		Name:      "Get Rollout Product",
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
		VersionNumber: "2.0.0",
		ReleaseType:   models.ReleaseTypeMajor,
	}
	version2, err := services.VersionService.CreateVersion(ctx, product.ProductID, &version2Req, "user1")
	if err != nil {
		t.Fatalf("Failed to create version 2.0.0: %v", err)
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

	detection := models.UpdateDetection{
		EndpointID:       "endpoint-get",
		ProductID:        product.ProductID,
		CurrentVersion:   "1.0.0",
		AvailableVersion: "2.0.0",
	}
	_, err = services.UpdateDetectionService.DetectUpdate(ctx, &detection)
	if err != nil {
		t.Fatalf("Failed to create detection: %v", err)
	}

	// Create rollout
	rollout := models.UpdateRollout{
		EndpointID:  "endpoint-get",
		ProductID:   product.ProductID,
		FromVersion: "1.0.0",
		ToVersion:   "2.0.0",
		Status:      models.RolloutStatusPending,
		InitiatedBy: "user1",
	}
	result, err := services.UpdateRolloutService.InitiateRollout(ctx, &rollout)
	if err != nil {
		t.Fatalf("Failed to create rollout: %v", err)
	}

	// Get rollout
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts/"+result.ID.Hex(), nil)
	w := httptest.NewRecorder()
	handler.GetRollout(w, httpReq)

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

	// Test invalid ID
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts/invalid-id", nil)
	w = httptest.NewRecorder()
	handler.GetRollout(w, httpReq)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid ID, got %d", http.StatusBadRequest, w.Code)
	}

	// Test non-existent rollout
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts/507f1f77bcf86cd799439011", nil)
	w = httptest.NewRecorder()
	handler.GetRollout(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent rollout, got %d", http.StatusNotFound, w.Code)
	}
}

func TestUpdateRolloutHandler_ListRollouts(t *testing.T) {
	handler, services, cleanup := setupUpdateRolloutHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create product
	productReq := models.CreateProductRequest{
		ProductID: "list-rollouts-product",
		Name:      "List Rollouts Product",
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

	// Create detections
	detections := []models.UpdateDetection{
		{
			EndpointID:       "endpoint-list-rollout-1",
			ProductID:        product.ProductID,
			CurrentVersion:   "1.0.0",
			AvailableVersion: "1.1.0",
		},
		{
			EndpointID:       "endpoint-list-rollout-2",
			ProductID:        product.ProductID,
			CurrentVersion:   "1.1.0",
			AvailableVersion: "2.0.0",
		},
		{
			EndpointID:       "endpoint-list-rollout-3",
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

	// Create multiple rollouts
	rollouts := []models.UpdateRollout{
		{
			EndpointID:  "endpoint-list-rollout-1",
			ProductID:   product.ProductID,
			FromVersion: "1.0.0",
			ToVersion:   "1.1.0",
			Status:      models.RolloutStatusPending,
			InitiatedBy: "user1",
		},
		{
			EndpointID:  "endpoint-list-rollout-2",
			ProductID:   product.ProductID,
			FromVersion: "1.1.0",
			ToVersion:   "2.0.0",
			Status:      models.RolloutStatusInProgress,
			InitiatedBy: "user1",
		},
		{
			EndpointID:  "endpoint-list-rollout-3",
			ProductID:   "other-product",
			FromVersion: "1.0.0",
			ToVersion:   "1.5.0",
			Status:      models.RolloutStatusCompleted,
			InitiatedBy: "user2",
		},
	}

	for _, rollout := range rollouts {
		_, err = services.UpdateRolloutService.InitiateRollout(ctx, &rollout)
		if err != nil {
			t.Fatalf("Failed to create rollout: %v", err)
		}
	}

	// Test list all rollouts
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts", nil)
	w := httptest.NewRecorder()
	handler.ListRollouts(w, httpReq)

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
		t.Errorf("Expected at least 3 rollouts, got %d", len(data))
	}

	// Test filter by product_id
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts?product_id="+product.ProductID, nil)
	w = httptest.NewRecorder()
	handler.ListRollouts(w, httpReq)

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
		t.Errorf("Expected 2 rollouts for product, got %d", len(data))
	}

	// Test filter by status
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts?status="+string(models.RolloutStatusPending), nil)
	w = httptest.NewRecorder()
	handler.ListRollouts(w, httpReq)

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

	// Should have at least 1 pending rollout
	if len(data) < 1 {
		t.Errorf("Expected at least 1 pending rollout, got %d", len(data))
	}

	// Test filter by endpoint_id
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts?endpoint_id=endpoint-list-rollout-1", nil)
	w = httptest.NewRecorder()
	handler.ListRollouts(w, httpReq)

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
		t.Errorf("Expected 1 rollout for endpoint, got %d", len(data))
	}

	// Test pagination
	httpReq = httptest.NewRequest(http.MethodGet, "/api/v1/update-rollouts?page=1&limit=2", nil)
	w = httptest.NewRecorder()
	handler.ListRollouts(w, httpReq)

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
		t.Errorf("Expected at most 2 rollouts with limit=2, got %d", len(data))
	}
}
