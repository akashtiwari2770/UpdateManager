package handlers

import (
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// UpdateDetectionHandler handles update detection-related HTTP requests
type UpdateDetectionHandler struct {
	updateDetectionService *service.UpdateDetectionService
}

// NewUpdateDetectionHandler creates a new update detection handler
func NewUpdateDetectionHandler(updateDetectionService *service.UpdateDetectionService) *UpdateDetectionHandler {
	return &UpdateDetectionHandler{
		updateDetectionService: updateDetectionService,
	}
}

// DetectUpdate handles POST /api/v1/update-detections
func (h *UpdateDetectionHandler) DetectUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var detection models.UpdateDetection
	if err := utils.ReadJSON(w, r, &detection); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	result, err := h.updateDetectionService.DetectUpdate(r.Context(), &detection)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "DETECTION_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, result)
}

// UpdateAvailableVersion handles PUT /api/v1/update-detections/:endpoint_id/:product_id/available-version
func (h *UpdateDetectionHandler) UpdateAvailableVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract endpoint_id and product_id from path
	pathParts := strings.Split(r.URL.Path, "/")
	var endpointID, productID string
	for i, part := range pathParts {
		if part == "update-detections" && i+1 < len(pathParts) {
			endpointID = pathParts[i+1]
			if i+2 < len(pathParts) {
				productID = pathParts[i+2]
			}
		}
	}

	if endpointID == "" || productID == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Endpoint ID and product ID are required")
		return
	}

	var req struct {
		AvailableVersion string `json:"available_version"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	if err := h.updateDetectionService.UpdateAvailableVersion(r.Context(), endpointID, productID, req.AvailableVersion); err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "DETECTION_NOT_FOUND", "Update detection not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	// Fetch the updated detection
	detection, err := h.updateDetectionService.GetDetection(r.Context(), endpointID, productID)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, detection)
}

// ListDetections handles GET /api/v1/update-detections
func (h *UpdateDetectionHandler) ListDetections(w http.ResponseWriter, r *http.Request) {
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
	if endpointID := r.URL.Query().Get("endpoint_id"); endpointID != "" {
		filter["endpoint_id"] = endpointID
	}
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter["product_id"] = productID
	}

	detections, total, err := h.updateDetectionService.ListDetections(r.Context(), filter, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, detections, page, limit, total)
}
