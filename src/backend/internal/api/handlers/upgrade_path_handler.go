package handlers

import (
	"net/http"
	"strings"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// UpgradePathHandler handles upgrade path-related HTTP requests
type UpgradePathHandler struct {
	upgradePathService *service.UpgradePathService
}

// NewUpgradePathHandler creates a new upgrade path handler
func NewUpgradePathHandler(upgradePathService *service.UpgradePathService) *UpgradePathHandler {
	return &UpgradePathHandler{
		upgradePathService: upgradePathService,
	}
}

// CreateUpgradePath handles POST /api/v1/products/:product_id/upgrade-paths
func (h *UpgradePathHandler) CreateUpgradePath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract product_id from path
	pathParts := strings.Split(r.URL.Path, "/")
	var productID string
	for i, part := range pathParts {
		if part == "products" && i+1 < len(pathParts) {
			productID = pathParts[i+1]
			break
		}
	}

	if productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PRODUCT_ID", "Product ID is required")
		return
	}

	var path models.UpgradePath
	if err := utils.ReadJSON(w, r, &path); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	path.ProductID = productID
	if err := h.upgradePathService.CreateUpgradePath(r.Context(), &path); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_PATH", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, path)
}

// GetUpgradePath handles GET /api/v1/products/:product_id/upgrade-paths/:from_version/:to_version
func (h *UpgradePathHandler) GetUpgradePath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract product_id, from_version, to_version from path
	pathParts := strings.Split(r.URL.Path, "/")
	var productID, fromVersion, toVersion string
	for i, part := range pathParts {
		if part == "products" && i+1 < len(pathParts) {
			productID = pathParts[i+1]
		}
		if part == "upgrade-paths" && i+1 < len(pathParts) {
			fromVersion = pathParts[i+1]
			if i+2 < len(pathParts) {
				toVersion = pathParts[i+2]
			}
		}
	}

	if productID == "" || fromVersion == "" || toVersion == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Product ID, from version, and to version are required")
		return
	}

	path, err := h.upgradePathService.GetUpgradePath(r.Context(), productID, fromVersion, toVersion)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "UPGRADE_PATH_NOT_FOUND", "Upgrade path not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, path)
}

// BlockUpgradePath handles POST /api/v1/products/:product_id/upgrade-paths/:from_version/:to_version/block
func (h *UpgradePathHandler) BlockUpgradePath(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract product_id, from_version, to_version from path
	pathParts := strings.Split(r.URL.Path, "/")
	var productID, fromVersion, toVersion string
	for i, part := range pathParts {
		if part == "products" && i+1 < len(pathParts) {
			productID = pathParts[i+1]
		}
		if part == "upgrade-paths" && i+1 < len(pathParts) {
			fromVersion = pathParts[i+1]
			if i+2 < len(pathParts) {
				toVersion = pathParts[i+2]
			}
		}
	}

	if productID == "" || fromVersion == "" || toVersion == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Product ID, from version, and to version are required")
		return
	}

	var req struct {
		BlockReason string `json:"block_reason"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	if err := h.upgradePathService.BlockUpgradePath(r.Context(), productID, fromVersion, toVersion, req.BlockReason); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "UPGRADE_PATH_NOT_FOUND", "Upgrade path not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "BLOCK_FAILED", err.Error())
		return
	}

	// Fetch the updated path
	path, err := h.upgradePathService.GetUpgradePath(r.Context(), productID, fromVersion, toVersion)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, path)
}
