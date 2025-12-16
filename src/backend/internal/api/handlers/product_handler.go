package handlers

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productService *service.ProductService
}

// NewProductHandler creates a new product handler
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{
		productService: productService,
	}
}

// CreateProduct handles POST /api/v1/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var req models.CreateProductRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.ProductID) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product ID is required and cannot be empty")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product name is required and cannot be empty")
		return
	}
	if req.Type == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product type is required")
		return
	}

	// Validate ProductID format (alphanumeric, underscore, hyphen only)
	if len(req.ProductID) > 100 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product ID must be 100 characters or less")
		return
	}
	if len(req.Name) > 200 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product name must be 200 characters or less")
		return
	}

	// Get user info from context (would come from auth middleware)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	product, err := h.productService.CreateProduct(r.Context(), &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_PRODUCT", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, product)
}

// GetProduct handles GET /api/v1/products/:id
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID format")
		return
	}

	product, err := h.productService.GetProduct(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, product)
}

// GetProductByProductID handles GET /api/v1/products/by-product-id/:product_id
func (h *ProductHandler) GetProductByProductID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	productID := strings.TrimPrefix(r.URL.Path, "/api/v1/products/by-product-id/")
	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PRODUCT_ID", "Product ID is required")
		return
	}

	product, err := h.productService.GetProductByProductID(r.Context(), productID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, product)
}

// UpdateProduct handles PUT /api/v1/products/:id
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID format")
		return
	}

	var req models.CreateProductRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	// Validate required fields
	if strings.TrimSpace(req.ProductID) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product ID is required and cannot be empty")
		return
	}
	if strings.TrimSpace(req.Name) == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product name is required and cannot be empty")
		return
	}
	if req.Type == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product type is required")
		return
	}

	// Validate field lengths
	if len(req.ProductID) > 100 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product ID must be 100 characters or less")
		return
	}
	if len(req.Name) > 200 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Product name must be 200 characters or less")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	product, err := h.productService.UpdateProduct(r.Context(), id, &req, userID, userEmail)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_PRODUCT", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, product)
}

// DeleteProduct handles DELETE /api/v1/products/:id
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/products/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid product ID format")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}
	userEmail := r.Header.Get("X-User-Email")

	if err := h.productService.DeleteProduct(r.Context(), id, userID, userEmail); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", "Product not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "DELETE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}

// ListProducts handles GET /api/v1/products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Parse query parameters
	page := utils.GetIntQueryParam(r, "page", 1)
	limit := utils.GetIntQueryParam(r, "limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Build filter from query params
	filter := bson.M{}
	if productType := r.URL.Query().Get("type"); productType != "" {
		filter["type"] = productType
	}
	if isActive := r.URL.Query().Get("is_active"); isActive != "" {
		if isActive == "true" {
			filter["is_active"] = true
		} else if isActive == "false" {
			filter["is_active"] = false
		}
	}

	products, total, err := h.productService.ListProducts(r.Context(), filter, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, products, page, limit, total)
}

// GetActiveProducts handles GET /api/v1/products/active
func (h *ProductHandler) GetActiveProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	products, err := h.productService.GetActiveProducts(r.Context())
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, products)
}
