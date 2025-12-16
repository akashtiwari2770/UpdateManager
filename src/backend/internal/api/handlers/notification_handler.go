package handlers

import (
	"net/http"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// NotificationHandler handles notification-related HTTP requests
type NotificationHandler struct {
	notificationService *service.NotificationService
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(notificationService *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		notificationService: notificationService,
	}
}

// CreateNotification handles POST /api/v1/notifications
func (h *NotificationHandler) CreateNotification(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var notification models.Notification
	if err := utils.ReadJSON(w, r, &notification); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	if err := h.notificationService.CreateNotification(r.Context(), &notification); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, notification)
}

// GetNotifications handles GET /api/v1/notifications
func (h *NotificationHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	recipientID := r.URL.Query().Get("recipient_id")
	if recipientID == "" {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_RECIPIENT_ID", "recipient_id query parameter is required")
		return
	}

	page := utils.GetIntQueryParam(r, "page", 1)
	limit := utils.GetIntQueryParam(r, "limit", 10)
	unreadOnly := r.URL.Query().Get("unread_only") == "true"

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	notifications, total, err := h.notificationService.GetNotifications(r.Context(), recipientID, page, limit, unreadOnly)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, notifications, page, limit, total)
}

// GetUnreadCount handles GET /api/v1/notifications/unread-count
func (h *NotificationHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	recipientID := r.URL.Query().Get("recipient_id")
	if recipientID == "" {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_RECIPIENT_ID", "recipient_id query parameter is required")
		return
	}

	count, err := h.notificationService.GetUnreadCount(r.Context(), recipientID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]interface{}{
		"recipient_id": recipientID,
		"unread_count": count,
	})
}

// MarkAllAsRead handles POST /api/v1/notifications/mark-all-read
func (h *NotificationHandler) MarkAllAsRead(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var req struct {
		RecipientID string `json:"recipient_id"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	if req.RecipientID == "" {
		utils.WriteError(w, http.StatusBadRequest, "MISSING_RECIPIENT_ID", "recipient_id is required")
		return
	}

	if err := h.notificationService.MarkAllAsRead(r.Context(), req.RecipientID); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "MARK_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "All notifications marked as read"})
}
