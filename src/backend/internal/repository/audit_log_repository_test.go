package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var auditLogRepo *AuditLogRepository

func setupAuditLogTestDB(t *testing.T) {
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
	auditLogRepo = NewAuditLogRepository(db.Collection("audit_logs"))
}

func teardownAuditLogTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("audit_logs").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestAuditLogCreate(t *testing.T) {
	setupAuditLogTestDB(t)
	defer teardownAuditLogTestDB(t)

	log := &models.AuditLog{
		Action:       models.AuditActionCreate,
		ResourceType: "product",
		ResourceID:   "product-123",
		UserID:       "user-123",
		UserEmail:    "user@example.com",
		Details: map[string]interface{}{
			"product_id": "test-product",
			"name":       "Test Product",
		},
		IPAddress: "127.0.0.1",
		UserAgent: "test-agent",
	}

	err := auditLogRepo.Create(testCtx, log)
	if err != nil {
		t.Fatalf("Failed to create audit log: %v", err)
	}

	if log.ID.IsZero() {
		t.Error("Audit log ID was not set after creation")
	}

	t.Logf("Created audit log with ID: %s", log.ID.Hex())
}

func TestAuditLogGetByResource(t *testing.T) {
	setupAuditLogTestDB(t)
	defer teardownAuditLogTestDB(t)

	resourceType := "version"
	resourceID := "version-123"
	logs := []*models.AuditLog{
		{Action: models.AuditActionCreate, ResourceType: resourceType, ResourceID: resourceID, UserID: "user-1", UserEmail: "user1@example.com"},
		{Action: models.AuditActionUpdate, ResourceType: resourceType, ResourceID: resourceID, UserID: "user-2", UserEmail: "user2@example.com"},
	}

	for _, l := range logs {
		err := auditLogRepo.Create(testCtx, l)
		if err != nil {
			t.Fatalf("Failed to create audit log: %v", err)
		}
	}

	retrieved, err := auditLogRepo.GetByResource(testCtx, resourceType, resourceID, nil)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(retrieved) < len(logs) {
		t.Errorf("Expected at least %d audit logs, got %d", len(logs), len(retrieved))
	}

	t.Logf("Retrieved %d audit logs for resource", len(retrieved))
}

func TestAuditLogGetByUserID(t *testing.T) {
	setupAuditLogTestDB(t)
	defer teardownAuditLogTestDB(t)

	userID := "user-456"
	logs := []*models.AuditLog{
		{Action: models.AuditActionCreate, ResourceType: "product", ResourceID: "product-1", UserID: userID, UserEmail: "user@example.com"},
		{Action: models.AuditActionUpdate, ResourceType: "product", ResourceID: "product-2", UserID: userID, UserEmail: "user@example.com"},
	}

	for _, l := range logs {
		err := auditLogRepo.Create(testCtx, l)
		if err != nil {
			t.Fatalf("Failed to create audit log: %v", err)
		}
	}

	retrieved, err := auditLogRepo.GetByUserID(testCtx, userID, nil)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(retrieved) < len(logs) {
		t.Errorf("Expected at least %d audit logs, got %d", len(logs), len(retrieved))
	}

	t.Logf("Retrieved %d audit logs for user", len(retrieved))
}
