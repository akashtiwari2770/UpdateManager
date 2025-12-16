package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/internal/service"
)

// PendingUpdatesHandler handles pending updates API requests
type PendingUpdatesHandler struct {
	pendingUpdatesService *service.PendingUpdatesService
}

// NewPendingUpdatesHandler creates a new pending updates handler
func NewPendingUpdatesHandler(pendingUpdatesService *service.PendingUpdatesService) *PendingUpdatesHandler {
	return &PendingUpdatesHandler{
		pendingUpdatesService: pendingUpdatesService,
	}
}

// GetDeploymentPendingUpdates handles GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/updates
func (h *PendingUpdatesHandler) GetDeploymentPendingUpdates(w http.ResponseWriter, r *http.Request) {
	// Extract IDs from URL path
	// Path: /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/updates
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/customers/"), "/")
	if len(pathParts) < 6 || pathParts[1] != "tenants" || pathParts[3] != "deployments" || pathParts[5] != "updates" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	deploymentID := pathParts[4]

	// Get pending updates for deployment
	pendingUpdates, err := h.pendingUpdatesService.GetPendingUpdatesForDeployment(r.Context(), deploymentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    pendingUpdates,
	})
}

// GetTenantPendingUpdates handles GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/pending-updates
func (h *PendingUpdatesHandler) GetTenantPendingUpdates(w http.ResponseWriter, r *http.Request) {
	// Extract IDs from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/customers/"), "/")
	if len(pathParts) < 4 || pathParts[1] != "tenants" || pathParts[3] != "deployments" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	customerID := pathParts[0]
	tenantID := pathParts[2]

	// Parse query parameters
	filter := &models.PendingUpdatesFilter{}
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter.ProductID = productID
	}
	if deploymentType := r.URL.Query().Get("deployment_type"); deploymentType != "" {
		filter.DeploymentType = models.DeploymentType(deploymentType)
	}
	if priority := r.URL.Query().Get("priority"); priority != "" {
		filter.Priority = priority
	}

	// Get pending updates for tenant
	summary, err := h.pendingUpdatesService.GetPendingUpdatesForTenant(r.Context(), customerID, tenantID, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    summary,
	})
}

// GetCustomerPendingUpdates handles GET /api/v1/customers/{customer_id}/deployments/pending-updates
func (h *PendingUpdatesHandler) GetCustomerPendingUpdates(w http.ResponseWriter, r *http.Request) {
	// Extract customer ID from URL path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/customers/"), "/")
	if len(pathParts) < 2 || pathParts[1] != "deployments" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	customerID := pathParts[0]

	// Parse query parameters
	filter := &models.PendingUpdatesFilter{}
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter.ProductID = productID
	}
	if deploymentType := r.URL.Query().Get("deployment_type"); deploymentType != "" {
		filter.DeploymentType = models.DeploymentType(deploymentType)
	}
	if priority := r.URL.Query().Get("priority"); priority != "" {
		filter.Priority = priority
	}

	// Get pending updates for customer
	summary, err := h.pendingUpdatesService.GetPendingUpdatesForCustomer(r.Context(), customerID, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    summary,
	})
}

// GetAllPendingUpdates handles GET /api/v1/updates/pending
func (h *PendingUpdatesHandler) GetAllPendingUpdates(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	filter := &models.PendingUpdatesFilter{}
	if customerID := r.URL.Query().Get("customer_id"); customerID != "" {
		filter.CustomerID = customerID
	}
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter.ProductID = productID
	}
	if tenantID := r.URL.Query().Get("tenant_id"); tenantID != "" {
		filter.TenantID = tenantID
	}
	if deploymentType := r.URL.Query().Get("deployment_type"); deploymentType != "" {
		filter.DeploymentType = models.DeploymentType(deploymentType)
	}
	if priority := r.URL.Query().Get("priority"); priority != "" {
		filter.Priority = priority
	}

	// Parse pagination
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

	// Get all pending updates
	results, paginationInfo, err := h.pendingUpdatesService.GetAllPendingUpdates(r.Context(), filter, pagination)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    results,
		"meta": map[string]interface{}{
			"page":        paginationInfo.Page,
			"limit":       paginationInfo.Limit,
			"total":       paginationInfo.Total,
			"total_pages": paginationInfo.TotalPages,
		},
	})
}

