package handlers

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// CompatibilityHandler handles compatibility-related HTTP requests
type CompatibilityHandler struct {
	compatibilityService *service.CompatibilityService
}

// NewCompatibilityHandler creates a new compatibility handler
func NewCompatibilityHandler(compatibilityService *service.CompatibilityService) *CompatibilityHandler {
	return &CompatibilityHandler{
		compatibilityService: compatibilityService,
	}
}

// ValidateCompatibility handles POST /api/v1/products/:product_id/versions/:version_number/compatibility
func (h *CompatibilityHandler) ValidateCompatibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract product_id and version_number from path
	pathParts := strings.Split(r.URL.Path, "/")
	var productID, versionNumber string
	for i, part := range pathParts {
		if part == "products" && i+1 < len(pathParts) {
			productID = pathParts[i+1]
		}
		if part == "versions" && i+1 < len(pathParts) {
			versionNumber = pathParts[i+1]
		}
	}

	if productID == "" || versionNumber == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Product ID and version number are required")
		return
	}

	var req models.ValidateCompatibilityRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	validatedBy := r.Header.Get("X-User-ID")
	if validatedBy == "" {
		validatedBy = "anonymous"
	}

	matrix, err := h.compatibilityService.ValidateCompatibility(r.Context(), productID, versionNumber, &req, validatedBy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "VALIDATION_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, matrix)
}

// GetCompatibility handles GET /api/v1/products/:product_id/versions/:version_number/compatibility
func (h *CompatibilityHandler) GetCompatibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract product_id and version_number from path
	pathParts := strings.Split(r.URL.Path, "/")
	var productID, versionNumber string
	for i, part := range pathParts {
		if part == "products" && i+1 < len(pathParts) {
			productID = pathParts[i+1]
		}
		if part == "versions" && i+1 < len(pathParts) {
			versionNumber = pathParts[i+1]
		}
	}

	if productID == "" || versionNumber == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Product ID and version number are required")
		return
	}

	matrix, err := h.compatibilityService.GetCompatibility(r.Context(), productID, versionNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "COMPATIBILITY_NOT_FOUND", "Compatibility matrix not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, matrix)
}

// ListCompatibility handles GET /api/v1/compatibility
func (h *CompatibilityHandler) ListCompatibility(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	page := utils.GetIntQueryParam(r, "page", 1)
	limit := utils.GetIntQueryParam(r, "limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	filter := bson.M{}
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter["product_id"] = productID
	}
	if status := r.URL.Query().Get("validation_status"); status != "" {
		filter["validation_status"] = status
	}

	matrices, total, err := h.compatibilityService.ListCompatibility(r.Context(), filter, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, matrices, page, limit, total)
}
