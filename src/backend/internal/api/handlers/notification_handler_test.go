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

func setupNotificationHandlerTest(t *testing.T) (*NotificationHandler, func()) {
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
	handler := NewNotificationHandler(services.NotificationService)

	cleanup := func() {
		db.Disconnect(ctx)
	}

	return handler, cleanup
}

func TestNotificationHandler_CreateNotification(t *testing.T) {
	handler, cleanup := setupNotificationHandlerTest(t)
	defer cleanup()

	notification := models.Notification{
		Type:        models.NotificationTypeNewVersion,
		RecipientID: "user1",
		ProductID:   "test-product",
		Title:       "New Version Available",
		Message:     "Version 1.0.0 is now available",
		Priority:    models.NotificationPriorityNormal,
		IsRead:      false,
		CreatedAt:   time.Now(),
	}

	body, _ := json.Marshal(notification)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/notifications", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.CreateNotification(w, httpReq)

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

func TestNotificationHandler_GetNotifications(t *testing.T) {
	handler, cleanup := setupNotificationHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple notifications
	recipientID := "user2"
	for i := 0; i < 5; i++ {
		notification := models.Notification{
			Type:        models.NotificationTypeNewVersion,
			RecipientID: recipientID,
			ProductID:   "test-product",
			Title:       "Notification " + string(rune('A'+i)),
			Message:     "Test message",
			Priority:    models.NotificationPriorityNormal,
			IsRead:      i%2 == 0, // Alternate read/unread
		}
		err := handler.notificationService.CreateNotification(ctx, &notification)
		if err != nil {
			t.Fatalf("Failed to create notification: %v", err)
		}
	}

	// Get notifications
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?recipient_id="+recipientID+"&page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.GetNotifications(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta.Total < 5 {
		t.Errorf("Expected at least 5 notifications, got %d", response.Meta.Total)
	}
}

func TestNotificationHandler_GetUnreadCount(t *testing.T) {
	handler, cleanup := setupNotificationHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	recipientID := "user3"
	// Create unread notifications
	for i := 0; i < 3; i++ {
		notification := models.Notification{
			Type:        models.NotificationTypeNewVersion,
			RecipientID: recipientID,
			Title:       "Unread Notification",
			Message:     "Test",
			IsRead:      false,
		}
		err := handler.notificationService.CreateNotification(ctx, &notification)
		if err != nil {
			t.Fatalf("Failed to create notification: %v", err)
		}
	}

	// Get unread count
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread-count?recipient_id="+recipientID, nil)
	w := httptest.NewRecorder()
	handler.GetUnreadCount(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	data := response.Data.(map[string]interface{})
	count := int(data["unread_count"].(float64))
	if count < 3 {
		t.Errorf("Expected at least 3 unread notifications, got %d", count)
	}
}

func TestNotificationHandler_MarkAllAsRead(t *testing.T) {
	handler, cleanup := setupNotificationHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	recipientID := "user4"
	// Create unread notifications
	for i := 0; i < 3; i++ {
		notification := models.Notification{
			Type:        models.NotificationTypeNewVersion,
			RecipientID: recipientID,
			Title:       "Unread Notification",
			Message:     "Test",
			IsRead:      false,
		}
		err := handler.notificationService.CreateNotification(ctx, &notification)
		if err != nil {
			t.Fatalf("Failed to create notification: %v", err)
		}
	}

	// Mark all as read
	req := map[string]string{
		"recipient_id": recipientID,
	}
	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/notifications/mark-all-read", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	handler.MarkAllAsRead(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify all are read
	count, err := handler.notificationService.GetUnreadCount(ctx, recipientID)
	if err != nil {
		t.Fatalf("Failed to get unread count: %v", err)
	}

	if count != 0 {
		t.Errorf("Expected 0 unread notifications, got %d", count)
	}
}
