package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// CustomerHandler handles customer-related HTTP requests
type CustomerHandler struct {
	customerService *service.CustomerService
}

// NewCustomerHandler creates a new customer handler
func NewCustomerHandler(customerService *service.CustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

// CreateCustomer handles POST /api/v1/customers
func (h *CustomerHandler) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var req models.CreateCustomerRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.Name) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Name is required")
		return
	}
	if strings.TrimSpace(req.Email) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Email is required")
		return
	}

	// Get user info from context
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	customer, err := h.customerService.CreateCustomer(r.Context(), &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_CUSTOMER", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, customer)
}

// GetCustomer handles GET /api/v1/customers/:id
func (h *CustomerHandler) GetCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Customer ID is required")
		return
	}

	customer, err := h.customerService.GetCustomer(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, customer)
}

// ListCustomers handles GET /api/v1/customers
func (h *CustomerHandler) ListCustomers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	query := &service.ListCustomersQuery{}

	// Parse query parameters
	if search := r.URL.Query().Get("search"); search != "" {
		query.Search = search
	}
	if status := r.URL.Query().Get("status"); status != "" {
		query.Status = models.CustomerStatus(status)
	}
	if email := r.URL.Query().Get("email"); email != "" {
		query.Email = email
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

	response, err := h.customerService.ListCustomers(r.Context(), query)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, response.Customers, response.Pagination.Page, response.Pagination.Limit, response.Pagination.Total)
}

// UpdateCustomer handles PUT /api/v1/customers/:id
func (h *CustomerHandler) UpdateCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Customer ID is required")
		return
	}

	var req models.UpdateCustomerRequest
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

	customer, err := h.customerService.UpdateCustomer(r.Context(), id, &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_CUSTOMER", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, customer)
}

// DeleteCustomer handles DELETE /api/v1/customers/:id
func (h *CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	if id == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Customer ID is required")
		return
	}

	// Get user info
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	err := h.customerService.DeleteCustomer(r.Context(), id, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "existing tenants") {
			utils.WriteError(w, http.StatusConflict, "CANNOT_DELETE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Customer deleted successfully"})
}

// GetCustomerTenants handles GET /api/v1/customers/:id/tenants
func (h *CustomerHandler) GetCustomerTenants(w http.ResponseWriter, r *http.Request) {
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

	response, err := h.customerService.GetCustomerTenants(r.Context(), customerID, query)
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

// GetCustomerStatistics handles GET /api/v1/customers/:id/statistics
func (h *CustomerHandler) GetCustomerStatistics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/customers/")
	id = strings.TrimSuffix(id, "/statistics")
	id = strings.TrimSuffix(id, "/")

	stats, err := h.customerService.GetCustomerStatistics(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, stats)
}

