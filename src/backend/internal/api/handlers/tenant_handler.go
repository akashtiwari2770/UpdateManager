package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// TenantHandler handles tenant-related HTTP requests
type TenantHandler struct {
	tenantService *service.TenantService
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(tenantService *service.TenantService) *TenantHandler {
	return &TenantHandler{
		tenantService: tenantService,
	}
}

// CreateTenant handles POST /api/v1/customers/:customer_id/tenants
func (h *TenantHandler) CreateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract customer ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "tenants" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	customerID := parts[0]

	var req models.CreateTenantRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Name is required")
		return
	}

	// Get user info
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	tenant, err := h.tenantService.CreateTenant(r.Context(), customerID, &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_TENANT", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, tenant)
}

// GetTenant handles GET /api/v1/customers/:customer_id/tenants/:id
func (h *TenantHandler) GetTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[1] != "tenants" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	tenantID := parts[2]

	tenant, err := h.tenantService.GetTenant(r.Context(), tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "TENANT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, tenant)
}

// ListTenants handles GET /api/v1/customers/:customer_id/tenants
func (h *TenantHandler) ListTenants(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract customer ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 2 || parts[1] != "tenants" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	customerID := parts[0]

	query := &service.ListTenantsQuery{}

	// Parse query parameters
	if status := r.URL.Query().Get("status"); status != "" {
		query.Status = models.TenantStatus(status)
	}
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			query.Limit = limit
		}
	}

	// Default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}

	response, err := h.tenantService.ListTenants(r.Context(), customerID, query)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, response.Tenants, response.Pagination.Page, response.Pagination.Limit, response.Pagination.Total)
}

// UpdateTenant handles PUT /api/v1/customers/:customer_id/tenants/:id
func (h *TenantHandler) UpdateTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[1] != "tenants" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	tenantID := parts[2]

	var req models.UpdateTenantRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	// Get user info
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	tenant, err := h.tenantService.UpdateTenant(r.Context(), tenantID, &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "TENANT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, tenant)
}

// DeleteTenant handles DELETE /api/v1/customers/:customer_id/tenants/:id
func (h *TenantHandler) DeleteTenant(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[1] != "tenants" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	tenantID := parts[2]

	// Get user info
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	err := h.tenantService.DeleteTenant(r.Context(), tenantID, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "TENANT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "existing deployments") {
			utils.WriteError(w, http.StatusConflict, "CANNOT_DELETE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Tenant deleted successfully"})
}

// GetTenantDeployments handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments
func (h *TenantHandler) GetTenantDeployments(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 4 || parts[1] != "tenants" || parts[3] != "deployments" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	tenantID := parts[2]

	query := &service.ListDeploymentsQuery{}

	// Parse query parameters
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		query.ProductID = productID
	}
	if deploymentType := r.URL.Query().Get("deployment_type"); deploymentType != "" {
		query.DeploymentType = models.DeploymentType(deploymentType)
	}
	if status := r.URL.Query().Get("status"); status != "" {
		query.Status = models.DeploymentStatus(status)
	}
	if version := r.URL.Query().Get("version"); version != "" {
		query.Version = version
	}
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			query.Page = page
		}
	}
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			query.Limit = limit
		}
	}

	// Default values
	if query.Page == 0 {
		query.Page = 1
	}
	if query.Limit == 0 {
		query.Limit = 20
	}

	response, err := h.tenantService.GetTenantDeployments(r.Context(), tenantID, query)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "TENANT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, response.Deployments, response.Pagination.Page, response.Pagination.Limit, response.Pagination.Total)
}

// GetTenantStatistics handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/statistics
func (h *TenantHandler) GetTenantStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract tenant ID from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 3 || parts[1] != "tenants" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	tenantID := parts[2]

	stats, err := h.tenantService.GetTenantStatistics(r.Context(), tenantID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "TENANT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, stats)
}

