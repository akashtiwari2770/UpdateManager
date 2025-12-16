package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var notificationRepo *NotificationRepository

func setupNotificationTestDB(t *testing.T) {
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
	notificationRepo = NewNotificationRepository(db.Collection("notifications"))
}

func teardownNotificationTestDB(t *testing.T) {
	if testDB != nil {
		_ = testDB.Collection("notifications").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

func TestNotificationCreate(t *testing.T) {
	setupNotificationTestDB(t)
	defer teardownNotificationTestDB(t)

	notification := &models.Notification{
		Type:        models.NotificationTypeNewVersion,
		RecipientID: "user-123",
		ProductID:   "test-product",
		Title:       "New Version Available",
		Message:     "Version 1.0.0 is now available",
		Priority:    models.NotificationPriorityNormal,
	}

	err := notificationRepo.Create(testCtx, notification)
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	if notification.ID.IsZero() {
		t.Error("Notification ID was not set after creation")
	}

	if notification.IsRead {
		t.Error("Notification should be unread by default")
	}

	t.Logf("Created notification with ID: %s", notification.ID.Hex())
}

func TestNotificationMarkAsRead(t *testing.T) {
	setupNotificationTestDB(t)
	defer teardownNotificationTestDB(t)

	notification := &models.Notification{
		Type:        models.NotificationTypeSecurityRelease,
		RecipientID: "user-456",
		Title:       "Security Update",
		Message:     "Important security update available",
		Priority:    models.NotificationPriorityHigh,
	}

	err := notificationRepo.Create(testCtx, notification)
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	err = notificationRepo.MarkAsRead(testCtx, notification.ID)
	if err != nil {
		t.Fatalf("Failed to mark notification as read: %v", err)
	}

	retrieved, err := notificationRepo.GetByID(testCtx, notification.ID)
	if err != nil {
		t.Fatalf("Failed to get notification: %v", err)
	}

	if !retrieved.IsRead {
		t.Error("Notification should be marked as read")
	}

	if retrieved.ReadAt == nil {
		t.Error("ReadAt should be set")
	}

	t.Logf("Marked notification as read: %+v", retrieved)
}

func TestNotificationGetUnread(t *testing.T) {
	setupNotificationTestDB(t)
	defer teardownNotificationTestDB(t)

	recipientID := "user-789"
	notifications := []*models.Notification{
		{Type: models.NotificationTypeNewVersion, RecipientID: recipientID, Title: "Notification 1", Message: "Message 1", Priority: models.NotificationPriorityLow},
		{Type: models.NotificationTypeEOLWarning, RecipientID: recipientID, Title: "Notification 2", Message: "Message 2", Priority: models.NotificationPriorityNormal},
	}

	for _, n := range notifications {
		err := notificationRepo.Create(testCtx, n)
		if err != nil {
			t.Fatalf("Failed to create notification: %v", err)
		}
	}

	// Mark one as read
	err := notificationRepo.MarkAsRead(testCtx, notifications[0].ID)
	if err != nil {
		t.Fatalf("Failed to mark notification as read: %v", err)
	}

	unread, err := notificationRepo.GetUnreadByRecipientID(testCtx, recipientID, nil)
	if err != nil {
		t.Fatalf("Failed to get unread notifications: %v", err)
	}

	if len(unread) < 1 {
		t.Errorf("Expected at least 1 unread notification, got %d", len(unread))
	}

	t.Logf("Found %d unread notifications", len(unread))
}
