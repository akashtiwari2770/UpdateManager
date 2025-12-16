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

// UpdateRolloutHandler handles update rollout-related HTTP requests
type UpdateRolloutHandler struct {
	updateRolloutService *service.UpdateRolloutService
}

// NewUpdateRolloutHandler creates a new update rollout handler
func NewUpdateRolloutHandler(updateRolloutService *service.UpdateRolloutService) *UpdateRolloutHandler {
	return &UpdateRolloutHandler{
		updateRolloutService: updateRolloutService,
	}
}

// InitiateRollout handles POST /api/v1/update-rollouts
func (h *UpdateRolloutHandler) InitiateRollout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	var rollout models.UpdateRollout
	if err := utils.ReadJSON(w, r, &rollout); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	initiatedBy := r.Header.Get("X-User-ID")
	if initiatedBy == "" {
		initiatedBy = "anonymous"
	}
	rollout.InitiatedBy = initiatedBy

	result, err := h.updateRolloutService.InitiateRollout(r.Context(), &rollout)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "ROLLOUT_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, result)
}

// UpdateRolloutStatus handles PUT /api/v1/update-rollouts/:id/status
func (h *UpdateRolloutHandler) UpdateRolloutStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/update-rollouts/")
	idStr = strings.TrimSuffix(idStr, "/status")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rollout ID format")
		return
	}

	var req struct {
		Status       models.RolloutStatus `json:"status"`
		ErrorMessage string               `json:"error_message,omitempty"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	rollout, err := h.updateRolloutService.UpdateRolloutStatus(r.Context(), id, req.Status, req.ErrorMessage)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "ROLLOUT_NOT_FOUND", "Update rollout not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, rollout)
}

// UpdateRolloutProgress handles PUT /api/v1/update-rollouts/:id/progress
func (h *UpdateRolloutHandler) UpdateRolloutProgress(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/update-rollouts/")
	idStr = strings.TrimSuffix(idStr, "/progress")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rollout ID format")
		return
	}

	var req struct {
		Progress int `json:"progress"`
	}
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	if req.Progress < 0 || req.Progress > 100 {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PROGRESS", "Progress must be between 0 and 100")
		return
	}

	rollout, err := h.updateRolloutService.UpdateRolloutProgress(r.Context(), id, req.Progress)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "ROLLOUT_NOT_FOUND", "Update rollout not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, rollout)
}

// GetRollout handles GET /api/v1/update-rollouts/:id
func (h *UpdateRolloutHandler) GetRollout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/update-rollouts/")
	idStr = strings.TrimSuffix(idStr, "/status")
	idStr = strings.TrimSuffix(idStr, "/progress")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid rollout ID format")
		return
	}

	rollout, err := h.updateRolloutService.GetRollout(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "ROLLOUT_NOT_FOUND", "Update rollout not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, rollout)
}

// ListRollouts handles GET /api/v1/update-rollouts
func (h *UpdateRolloutHandler) ListRollouts(w http.ResponseWriter, r *http.Request) {
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
	if status := r.URL.Query().Get("status"); status != "" {
		filter["status"] = status
	}

	rollouts, total, err := h.updateRolloutService.ListRollouts(r.Context(), filter, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, rollouts, page, limit, total)
}
