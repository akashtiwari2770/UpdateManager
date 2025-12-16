package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/internal/service"
	"updatemanager/pkg/database"
)

func setupAuditLogHandlerTest(t *testing.T) (*AuditLogHandler, *service.ServiceFactory, *database.MongoDB, func()) {
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

	services := service.NewServiceFactory(db.Database)
	handler := NewAuditLogHandler(services.AuditLogService)

	cleanup := func() {
		db.Disconnect(ctx)
	}

	return handler, services, db, cleanup
}

func TestAuditLogHandler_GetAuditLogs(t *testing.T) {
	handler, _, db, cleanup := setupAuditLogHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create some audit logs
	auditRepo := repository.NewAuditLogRepository(db.Collection("audit_logs"))
	for i := 0; i < 5; i++ {
		auditLog := models.AuditLog{
			Action:       models.AuditActionCreate,
			ResourceType: "product",
			ResourceID:   "product-" + string(rune('a'+i)),
			UserID:       "user1",
			UserEmail:    "user1@example.com",
			Details:      map[string]interface{}{"test": "data"},
			Timestamp:    time.Now(),
		}
		err := auditRepo.Create(ctx, &auditLog)
		if err != nil {
			t.Fatalf("Failed to create audit log: %v", err)
		}
	}

	// Get audit logs
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.GetAuditLogs(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 5 {
		t.Errorf("Expected at least 5 audit logs, got %d", response.Meta.Total)
	}
}

func TestAuditLogHandler_GetAuditLogs_WithFilters(t *testing.T) {
	handler, _, db, cleanup := setupAuditLogHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	auditRepo := repository.NewAuditLogRepository(db.Collection("audit_logs"))

	// Create audit logs with different filters
	auditLog1 := models.AuditLog{
		Action:       models.AuditActionCreate,
		ResourceType: "product",
		ResourceID:   "filter-product",
		UserID:       "user-filter",
		Timestamp:    time.Now(),
	}
	err := auditRepo.Create(ctx, &auditLog1)
	if err != nil {
		t.Fatalf("Failed to create audit log: %v", err)
	}

	auditLog2 := models.AuditLog{
		Action:       models.AuditActionUpdate,
		ResourceType: "version",
		ResourceID:   "filter-version",
		UserID:       "user-filter",
		Timestamp:    time.Now(),
	}
	err = auditRepo.Create(ctx, &auditLog2)
	if err != nil {
		t.Fatalf("Failed to create audit log: %v", err)
	}

	// Get audit logs by resource type
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs?resource_type=product&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.GetAuditLogs(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Should have at least one product audit log
	if response.Meta.Total < 1 {
		t.Errorf("Expected at least 1 audit log, got %d", response.Meta.Total)
	}

	// Get audit logs by user
	httpReq2 := httptest.NewRequest(http.MethodGet, "/api/v1/audit-logs?user_id=user-filter&page=1&limit=10", nil)
	w2 := httptest.NewRecorder()
	handler.GetAuditLogs(w2, httpReq2)

	if w2.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w2.Code)
	}

	var response2 utils.JSONResponse
	if err := json.Unmarshal(w2.Body.Bytes(), &response2); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response2.Meta.Total < 2 {
		t.Errorf("Expected at least 2 audit logs for user, got %d", response2.Meta.Total)
	}
}
