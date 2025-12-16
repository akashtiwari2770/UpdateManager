package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// DeploymentHandler handles deployment-related HTTP requests
type DeploymentHandler struct {
	deploymentService *service.DeploymentService
}

// NewDeploymentHandler creates a new deployment handler
func NewDeploymentHandler(deploymentService *service.DeploymentService) *DeploymentHandler {
	return &DeploymentHandler{
		deploymentService: deploymentService,
	}
}

// CreateDeployment handles POST /api/v1/customers/:customer_id/tenants/:tenant_id/deployments
func (h *DeploymentHandler) CreateDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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

	var req models.CreateDeploymentRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.ProductID) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product ID is required")
		return
	}
	if req.DeploymentType == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Deployment type is required")
		return
	}
	if strings.TrimSpace(req.InstalledVersion) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Installed version is required")
		return
	}

	// Get user info
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	deployment, err := h.deploymentService.CreateDeployment(r.Context(), tenantID, &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "TENANT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_DEPLOYMENT", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, deployment)
}

// GetDeployment handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id
func (h *DeploymentHandler) GetDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 5 || parts[1] != "tenants" || parts[3] != "deployments" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	deploymentID := parts[4]

	deployment, err := h.deploymentService.GetDeployment(r.Context(), deploymentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "DEPLOYMENT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, deployment)
}

// ListDeployments handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments
func (h *DeploymentHandler) ListDeployments(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.deploymentService.ListDeployments(r.Context(), tenantID, query)
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

// UpdateDeployment handles PUT /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id
func (h *DeploymentHandler) UpdateDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 5 || parts[1] != "tenants" || parts[3] != "deployments" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	deploymentID := parts[4]

	var req models.UpdateDeploymentRequest
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

	deployment, err := h.deploymentService.UpdateDeployment(r.Context(), deploymentID, &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "DEPLOYMENT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_DEPLOYMENT", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, deployment)
}

// DeleteDeployment handles DELETE /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id
func (h *DeploymentHandler) DeleteDeployment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 5 || parts[1] != "tenants" || parts[3] != "deployments" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	deploymentID := parts[4]

	// Get user info
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	err := h.deploymentService.DeleteDeployment(r.Context(), deploymentID, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "DEPLOYMENT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Deployment deleted successfully"})
}

// GetAvailableUpdates handles GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id/updates
func (h *DeploymentHandler) GetAvailableUpdates(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract IDs from path
	path := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	parts := strings.Split(path, "/")
	if len(parts) < 6 || parts[1] != "tenants" || parts[3] != "deployments" || parts[5] != "updates" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path")
		return
	}
	deploymentID := parts[4]

	updates, err := h.deploymentService.GetAvailableUpdates(r.Context(), deploymentID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "DEPLOYMENT_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, updates)
}

