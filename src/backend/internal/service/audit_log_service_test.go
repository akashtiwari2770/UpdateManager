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
	auditLogServiceTestDB  *database.MongoDB
	auditLogServiceTestCtx context.Context
	auditLogService        *AuditLogService
	auditLogRepo           *repository.AuditLogRepository
)

func setupAuditLogServiceTestDB(t *testing.T) {
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

	auditLogServiceTestDB = db
	auditLogServiceTestCtx = ctx
	auditLogRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	auditLogService = NewAuditLogService(auditLogRepo)
}

func teardownAuditLogServiceTestDB(t *testing.T) {
	if auditLogServiceTestDB != nil {
		_ = auditLogServiceTestDB.Collection("audit_logs").Drop(auditLogServiceTestCtx)
		_ = auditLogServiceTestDB.Disconnect(auditLogServiceTestCtx)
	}
}

func TestAuditLogService_GetAuditLogsByResource(t *testing.T) {
	setupAuditLogServiceTestDB(t)
	defer teardownAuditLogServiceTestDB(t)

	resourceType := "product"
	resourceID := "product-123"

	// Create audit logs
	logs := []*models.AuditLog{
		{Action: models.AuditActionCreate, ResourceType: resourceType, ResourceID: resourceID, UserID: "user-1", UserEmail: "user1@example.com", Timestamp: time.Now()},
		{Action: models.AuditActionUpdate, ResourceType: resourceType, ResourceID: resourceID, UserID: "user-2", UserEmail: "user2@example.com", Timestamp: time.Now()},
		{Action: models.AuditActionDelete, ResourceType: resourceType, ResourceID: resourceID, UserID: "user-3", UserEmail: "user3@example.com", Timestamp: time.Now()},
	}

	for _, log := range logs {
		auditLogRepo.Create(auditLogServiceTestCtx, log)
	}

	// Get audit logs by resource
	retrieved, total, err := auditLogService.GetAuditLogsByResource(auditLogServiceTestCtx, resourceType, resourceID, 1, 10)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(retrieved) < len(logs) {
		t.Errorf("Expected at least %d audit logs, got %d", len(logs), len(retrieved))
	}

	if total < int64(len(logs)) {
		t.Errorf("Expected total at least %d, got %d", len(logs), total)
	}

	t.Logf("Retrieved %d audit logs for resource (total: %d)", len(retrieved), total)
}

func TestAuditLogService_GetAuditLogsByUser(t *testing.T) {
	setupAuditLogServiceTestDB(t)
	defer teardownAuditLogServiceTestDB(t)

	userID := "user-456"

	// Create audit logs
	logs := []*models.AuditLog{
		{Action: models.AuditActionCreate, ResourceType: "product", ResourceID: "p1", UserID: userID, UserEmail: "user@example.com", Timestamp: time.Now()},
		{Action: models.AuditActionUpdate, ResourceType: "product", ResourceID: "p2", UserID: userID, UserEmail: "user@example.com", Timestamp: time.Now()},
	}

	for _, log := range logs {
		auditLogRepo.Create(auditLogServiceTestCtx, log)
	}

	// Get audit logs by user
	retrieved, total, err := auditLogService.GetAuditLogsByUser(auditLogServiceTestCtx, userID, 1, 10)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(retrieved) < len(logs) {
		t.Errorf("Expected at least %d audit logs, got %d", len(logs), len(retrieved))
	}

	if total < int64(len(logs)) {
		t.Errorf("Expected total at least %d, got %d", len(logs), total)
	}

	t.Logf("Retrieved %d audit logs for user (total: %d)", len(retrieved), total)
}

func TestAuditLogService_GetAuditLogsByAction(t *testing.T) {
	setupAuditLogServiceTestDB(t)
	defer teardownAuditLogServiceTestDB(t)

	action := models.AuditActionCreate

	// Create audit logs
	logs := []*models.AuditLog{
		{Action: action, ResourceType: "product", ResourceID: "p1", UserID: "user-1", UserEmail: "user1@example.com", Timestamp: time.Now()},
		{Action: action, ResourceType: "version", ResourceID: "v1", UserID: "user-2", UserEmail: "user2@example.com", Timestamp: time.Now()},
	}

	for _, log := range logs {
		auditLogRepo.Create(auditLogServiceTestCtx, log)
	}

	// Get audit logs by action
	retrieved, total, err := auditLogService.GetAuditLogsByAction(auditLogServiceTestCtx, action, 1, 10)
	if err != nil {
		t.Fatalf("Failed to get audit logs: %v", err)
	}

	if len(retrieved) < len(logs) {
		t.Errorf("Expected at least %d audit logs, got %d", len(logs), len(retrieved))
	}

	// Verify all are create actions
	for _, log := range retrieved {
		if log.Action != action {
			t.Errorf("Expected action %s, got %s", action, log.Action)
		}
	}

	t.Logf("Retrieved %d audit logs for action %s (total: %d)", len(retrieved), action, total)
}
