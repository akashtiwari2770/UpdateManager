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

// LicenseHandler handles license-related HTTP requests
type LicenseHandler struct {
	licenseService *service.LicenseService
}

// NewLicenseHandler creates a new license handler
func NewLicenseHandler(licenseService *service.LicenseService) *LicenseHandler {
	return &LicenseHandler{
		licenseService: licenseService,
	}
}

// extractPathParams extracts customer_id, subscription_id, and license_id from URL path
func (h *LicenseHandler) extractPathParams(path string) (customerID, subscriptionID, licenseID string) {
	// Path format: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/customers/"), "/")
	if len(parts) >= 4 && parts[1] == "subscriptions" && parts[3] == "licenses" {
		customerID = parts[0]
		subscriptionID = parts[2]
		if len(parts) >= 5 {
			licenseID = parts[4]
		}
	}
	return
}

// AssignLicense handles POST /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses
func (h *LicenseHandler) AssignLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, _ := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	var req models.CreateLicenseRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	license, err := h.licenseService.AssignLicense(r.Context(), customerID, subscriptionID, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_LICENSE", err.Error())
			return
		}
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "must have an end date") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "ASSIGN_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, license)
}

// GetLicense handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id
func (h *LicenseHandler) GetLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	license, err := h.licenseService.GetLicense(r.Context(), customerID, subscriptionID, licenseID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, license)
}

// ListLicenses handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses
func (h *LicenseHandler) ListLicenses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, _ := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID and Subscription ID are required")
		return
	}

	// Parse query parameters
	filter := &repository.LicenseFilter{}
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter.ProductID = productID
	}
	if licenseType := r.URL.Query().Get("license_type"); licenseType != "" {
		filter.LicenseType = models.LicenseType(licenseType)
	}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.LicenseStatus(status)
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

	licenses, paginationInfo, err := h.licenseService.ListLicenses(r.Context(), customerID, subscriptionID, filter, pagination)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response := map[string]interface{}{
		"licenses":   licenses,
		"pagination": paginationInfo,
	}
	utils.WriteSuccess(w, http.StatusOK, response)
}

// UpdateLicense handles PUT /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id
func (h *LicenseHandler) UpdateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	var req models.UpdateLicenseRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	license, err := h.licenseService.UpdateLicense(r.Context(), customerID, subscriptionID, licenseID, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "must have an end date") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, license)
}

// RevokeLicense handles DELETE /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id
func (h *LicenseHandler) RevokeLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	err := h.licenseService.RevokeLicense(r.Context(), customerID, subscriptionID, licenseID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "LICENSE_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "active allocations") {
			utils.WriteError(w, http.StatusConflict, "HAS_ALLOCATIONS", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "REVOKE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "License revoked successfully"})
}

// GetLicenseStatistics handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/statistics
func (h *LicenseHandler) GetLicenseStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	stats, err := h.licenseService.GetLicenseStatistics(r.Context(), customerID, subscriptionID, licenseID)
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

// RenewLicense handles POST /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/renew
func (h *LicenseHandler) RenewLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
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

	license, err := h.licenseService.RenewLicense(r.Context(), customerID, subscriptionID, licenseID, req.EndDate, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "only time-based") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "RENEW_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, license)
}

