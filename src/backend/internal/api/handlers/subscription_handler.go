package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/internal/service"
)

// SubscriptionHandler handles subscription-related HTTP requests
type SubscriptionHandler struct {
	subscriptionService *service.SubscriptionService
}

// NewSubscriptionHandler creates a new subscription handler
func NewSubscriptionHandler(subscriptionService *service.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// extractPathParams extracts customer_id and subscription_id from URL path
func (h *SubscriptionHandler) extractPathParams(path string) (customerID, subscriptionID string) {
	// Path format: /api/v1/customers/:customer_id/subscriptions/:subscription_id
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/customers/"), "/")
	if len(parts) >= 2 && parts[1] == "subscriptions" {
		customerID = parts[0]
		if len(parts) >= 3 {
			subscriptionID = parts[2]
		}
	}
	return
}

// CreateSubscription handles POST /api/v1/customers/:customer_id/subscriptions
func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, _ := h.extractPathParams(path)
	if customerID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID is required")
		return
	}

	var req models.CreateSubscriptionRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	subscription, err := h.subscriptionService.CreateSubscription(r.Context(), customerID, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_SUBSCRIPTION", err.Error())
			return
		}
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, subscription)
}

// GetSubscription handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id
func (h *SubscriptionHandler) GetSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	subscription, err := h.subscriptionService.GetSubscription(r.Context(), customerID, subscriptionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "SUBSCRIPTION_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, subscription)
}

// ListSubscriptions handles GET /api/v1/customers/:customer_id/subscriptions
func (h *SubscriptionHandler) ListSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, _ := h.extractPathParams(path)
	if customerID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID is required")
		return
	}

	// Parse query parameters
	filter := &repository.SubscriptionFilter{}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.SubscriptionStatus(status)
	}

	pagination := &repository.Pagination{
		Page:  1,
		Limit: 20,
	}
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			pagination.Page = page
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			pagination.Limit = limit
		}
	}

	subscriptions, paginationInfo, err := h.subscriptionService.ListSubscriptions(r.Context(), customerID, filter, pagination)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response := map[string]interface{}{
		"subscriptions": subscriptions,
		"pagination":    paginationInfo,
	}
	utils.WriteSuccess(w, http.StatusOK, response)
}

// UpdateSubscription handles PUT /api/v1/customers/:customer_id/subscriptions/:subscription_id
func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	var req models.UpdateSubscriptionRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	subscription, err := h.subscriptionService.UpdateSubscription(r.Context(), customerID, subscriptionID, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "SUBSCRIPTION_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, subscription)
}

// DeleteSubscription handles DELETE /api/v1/customers/:customer_id/subscriptions/:subscription_id
func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	err := h.subscriptionService.DeleteSubscription(r.Context(), customerID, subscriptionID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "SUBSCRIPTION_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "existing licenses") {
			utils.WriteError(w, http.StatusConflict, "HAS_LICENSES", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Subscription deleted successfully"})
}

// GetSubscriptionStatistics handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/statistics
func (h *SubscriptionHandler) GetSubscriptionStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	stats, err := h.subscriptionService.GetSubscriptionStatistics(r.Context(), customerID, subscriptionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, stats)
}

// RenewSubscription handles POST /api/v1/customers/:customer_id/subscriptions/:subscription_id/renew
func (h *SubscriptionHandler) RenewSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	var req struct {
		EndDate time.Time `json:"end_date"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	subscription, err := h.subscriptionService.RenewSubscription(r.Context(), customerID, subscriptionID, req.EndDate, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "RENEW_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, subscription)
}

