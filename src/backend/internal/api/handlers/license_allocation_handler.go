package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/internal/service"
)

// LicenseAllocationHandler handles license allocation-related HTTP requests
type LicenseAllocationHandler struct {
	allocationService *service.LicenseAllocationService
}

// NewLicenseAllocationHandler creates a new license allocation handler
func NewLicenseAllocationHandler(allocationService *service.LicenseAllocationService) *LicenseAllocationHandler {
	return &LicenseAllocationHandler{
		allocationService: allocationService,
	}
}

// extractPathParams extracts customer_id, subscription_id, license_id, and allocation_id from URL path
func (h *LicenseAllocationHandler) extractPathParams(path string) (customerID, subscriptionID, licenseID, allocationID string) {
	// Path format: /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations/:allocation_id
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/customers/"), "/")
	if len(parts) >= 6 && parts[1] == "subscriptions" && parts[3] == "licenses" && parts[5] == "allocations" {
		customerID = parts[0]
		subscriptionID = parts[2]
		licenseID = parts[4]
		if len(parts) >= 7 {
			allocationID = parts[6]
		}
	}
	return
}

// AllocateLicense handles POST /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocate
func (h *LicenseAllocationHandler) AllocateLicense(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID, _ := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	var req models.AllocateLicenseRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	allocation, err := h.allocationService.AllocateLicense(r.Context(), customerID, subscriptionID, licenseID, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "insufficient available seats") {
			utils.WriteError(w, http.StatusBadRequest, "INSUFFICIENT_SEATS", err.Error())
			return
		}
		if strings.Contains(err.Error(), "expired") {
			utils.WriteError(w, http.StatusBadRequest, "LICENSE_EXPIRED", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "ALLOCATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, allocation)
}

// ReleaseAllocation handles POST /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations/:allocation_id/release
func (h *LicenseAllocationHandler) ReleaseAllocation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID, allocationID := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" || allocationID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, License ID, and Allocation ID are required")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	err := h.allocationService.ReleaseAllocation(r.Context(), customerID, subscriptionID, licenseID, allocationID, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already released") {
			utils.WriteError(w, http.StatusBadRequest, "ALREADY_RELEASED", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "RELEASE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Allocation released successfully"})
}

// GetAllocations handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/allocations
func (h *LicenseAllocationHandler) GetAllocations(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID, _ := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	// Parse query parameters
	filter := &repository.LicenseAllocationFilter{}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.AllocationStatus(status)
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

	allocations, paginationInfo, err := h.allocationService.GetAllocations(r.Context(), customerID, subscriptionID, licenseID, filter, pagination)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response := map[string]interface{}{
		"allocations": allocations,
		"pagination": paginationInfo,
	}
	utils.WriteSuccess(w, http.StatusOK, response)
}

// GetAllocationsByTenant handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/licenses
func (h *LicenseAllocationHandler) GetAllocationsByTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract customer_id and tenant_id from path
	// Path format: /api/v1/customers/:customer_id/tenants/:tenant_id/licenses
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/customers/"), "/")
	if len(parts) < 4 || parts[2] != "tenants" || parts[4] != "licenses" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path format")
		return
	}

	customerID := parts[0]
	tenantID := parts[3]

	// Parse query parameters
	filter := &repository.LicenseAllocationFilter{}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.AllocationStatus(status)
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

	allocations, paginationInfo, err := h.allocationService.GetAllocationsByTenant(r.Context(), customerID, tenantID, filter, pagination)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response := map[string]interface{}{
		"allocations": allocations,
		"pagination": paginationInfo,
	}
	utils.WriteSuccess(w, http.StatusOK, response)
}

// GetAllocationsByDeployment handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:deployment_id/licenses
func (h *LicenseAllocationHandler) GetAllocationsByDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract customer_id, tenant_id, and deployment_id from path
	// Path format: /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:deployment_id/licenses
	path := r.URL.Path
	parts := strings.Split(strings.TrimPrefix(path, "/api/v1/customers/"), "/")
	if len(parts) < 6 || parts[2] != "tenants" || parts[4] != "deployments" || parts[6] != "licenses" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path format")
		return
	}

	customerID := parts[0]
	tenantID := parts[3]
	deploymentID := parts[5]

	// Parse query parameters
	filter := &repository.LicenseAllocationFilter{}
	if status := r.URL.Query().Get("status"); status != "" {
		filter.Status = models.AllocationStatus(status)
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

	allocations, paginationInfo, err := h.allocationService.GetAllocationsByDeployment(r.Context(), customerID, tenantID, deploymentID, filter, pagination)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	response := map[string]interface{}{
		"allocations": allocations,
		"pagination": paginationInfo,
	}
	utils.WriteSuccess(w, http.StatusOK, response)
}

// GetLicenseUtilization handles GET /api/v1/customers/:customer_id/subscriptions/:subscription_id/licenses/:license_id/utilization
func (h *LicenseAllocationHandler) GetLicenseUtilization(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	path := r.URL.Path
	customerID, subscriptionID, licenseID, _ := h.extractPathParams(path)
	if customerID == "" || subscriptionID == "" || licenseID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Customer ID, Subscription ID, and License ID are required")
		return
	}

	utilization, err := h.allocationService.GetLicenseUtilization(r.Context(), licenseID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, utilization)
}

