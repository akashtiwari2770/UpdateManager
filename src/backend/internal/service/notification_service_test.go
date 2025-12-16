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
	notificationServiceTestDB  *database.MongoDB
	notificationServiceTestCtx context.Context
	notificationService        *NotificationService
	notificationRepo           *repository.NotificationRepository
)

func setupNotificationServiceTestDB(t *testing.T) {
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

	notificationServiceTestDB = db
	notificationServiceTestCtx = ctx
	notificationRepo = repository.NewNotificationRepository(db.Collection("notifications"))
	notificationService = NewNotificationService(notificationRepo)
}

func teardownNotificationServiceTestDB(t *testing.T) {
	if notificationServiceTestDB != nil {
		_ = notificationServiceTestDB.Collection("notifications").Drop(notificationServiceTestCtx)
		_ = notificationServiceTestDB.Disconnect(notificationServiceTestCtx)
	}
}

func TestNotificationService_CreateNotification(t *testing.T) {
	setupNotificationServiceTestDB(t)
	defer teardownNotificationServiceTestDB(t)

	notification := &models.Notification{
		Type:        models.NotificationTypeNewVersion,
		RecipientID: "user-123",
		ProductID:   "test-product",
		Title:       "New Version Available",
		Message:     "Version 1.0.0 is now available",
		Priority:    models.NotificationPriorityNormal,
	}

	err := notificationService.CreateNotification(notificationServiceTestCtx, notification)
	if err != nil {
		t.Fatalf("Failed to create notification: %v", err)
	}

	if notification.ID.IsZero() {
		t.Error("Notification ID was not set")
	}

	t.Logf("Created notification: %+v", notification)
}

func TestNotificationService_GetNotifications(t *testing.T) {
	setupNotificationServiceTestDB(t)
	defer teardownNotificationServiceTestDB(t)

	recipientID := "user-456"
	notifications := []*models.Notification{
		{Type: models.NotificationTypeNewVersion, RecipientID: recipientID, Title: "Notification 1", Message: "Message 1", Priority: models.NotificationPriorityLow},
		{Type: models.NotificationTypeSecurityRelease, RecipientID: recipientID, Title: "Notification 2", Message: "Message 2", Priority: models.NotificationPriorityHigh},
		{Type: models.NotificationTypeEOLWarning, RecipientID: recipientID, Title: "Notification 3", Message: "Message 3", Priority: models.NotificationPriorityNormal},
	}

	for _, n := range notifications {
		notificationService.CreateNotification(notificationServiceTestCtx, n)
	}

	// Mark one as read
	notificationRepo.MarkAsRead(notificationServiceTestCtx, notifications[0].ID)

	// Get all notifications
	allNotifications, total, err := notificationService.GetNotifications(notificationServiceTestCtx, recipientID, 1, 10, false)
	if err != nil {
		t.Fatalf("Failed to get notifications: %v", err)
	}

	if len(allNotifications) < len(notifications) {
		t.Errorf("Expected at least %d notifications, got %d", len(notifications), len(allNotifications))
	}

	if total < int64(len(notifications)) {
		t.Errorf("Expected total at least %d, got %d", len(notifications), total)
	}

	// Get unread only
	unreadNotifications, _, err := notificationService.GetNotifications(notificationServiceTestCtx, recipientID, 1, 10, true)
	if err != nil {
		t.Fatalf("Failed to get unread notifications: %v", err)
	}

	if len(unreadNotifications) >= len(notifications) {
		t.Errorf("Expected fewer unread notifications, got %d", len(unreadNotifications))
	}

	t.Logf("Found %d total notifications, %d unread", len(allNotifications), len(unreadNotifications))
}

func TestNotificationService_GetUnreadCount(t *testing.T) {
	setupNotificationServiceTestDB(t)
	defer teardownNotificationServiceTestDB(t)

	recipientID := "user-789"
	notifications := []*models.Notification{
		{Type: models.NotificationTypeNewVersion, RecipientID: recipientID, Title: "N1", Message: "M1", Priority: models.NotificationPriorityLow},
		{Type: models.NotificationTypeNewVersion, RecipientID: recipientID, Title: "N2", Message: "M2", Priority: models.NotificationPriorityNormal},
	}

	for _, n := range notifications {
		notificationService.CreateNotification(notificationServiceTestCtx, n)
	}

	// Mark one as read
	notificationRepo.MarkAsRead(notificationServiceTestCtx, notifications[0].ID)

	// Get unread count
	count, err := notificationService.GetUnreadCount(notificationServiceTestCtx, recipientID)
	if err != nil {
		t.Fatalf("Failed to get unread count: %v", err)
	}

	if count < 1 {
		t.Errorf("Expected at least 1 unread notification, got %d", count)
	}

	t.Logf("Unread count: %d", count)
}

func TestNotificationService_MarkAllAsRead(t *testing.T) {
	setupNotificationServiceTestDB(t)
	defer teardownNotificationServiceTestDB(t)

	recipientID := "user-all-read"
	notifications := []*models.Notification{
		{Type: models.NotificationTypeNewVersion, RecipientID: recipientID, Title: "N1", Message: "M1", Priority: models.NotificationPriorityLow},
		{Type: models.NotificationTypeNewVersion, RecipientID: recipientID, Title: "N2", Message: "M2", Priority: models.NotificationPriorityNormal},
	}

	for _, n := range notifications {
		notificationService.CreateNotification(notificationServiceTestCtx, n)
	}

	// Mark all as read
	err := notificationService.MarkAllAsRead(notificationServiceTestCtx, recipientID)
	if err != nil {
		t.Fatalf("Failed to mark all as read: %v", err)
	}

	// Verify all are read
	unreadCount, _ := notificationService.GetUnreadCount(notificationServiceTestCtx, recipientID)
	if unreadCount != 0 {
		t.Errorf("Expected 0 unread notifications, got %d", unreadCount)
	}

	t.Logf("All notifications marked as read")
}
