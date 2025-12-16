package handlers

import (
	"net/http"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/service"
)

// AuditLogHandler handles audit log-related HTTP requests
type AuditLogHandler struct {
	auditLogService *service.AuditLogService
}

// NewAuditLogHandler creates a new audit log handler
func NewAuditLogHandler(auditLogService *service.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{
		auditLogService: auditLogService,
	}
}

// GetAuditLogs handles GET /api/v1/audit-logs
func (h *AuditLogHandler) GetAuditLogs(w http.ResponseWriter, r *http.Request) {
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

	// Build filter from query params
	filter := make(map[string]interface{})
	if resourceType := r.URL.Query().Get("resource_type"); resourceType != "" {
		filter["resource_type"] = resourceType
	}
	if resourceID := r.URL.Query().Get("resource_id"); resourceID != "" {
		filter["resource_id"] = resourceID
	}
	if userID := r.URL.Query().Get("user_id"); userID != "" {
		filter["user_id"] = userID
	}
	if action := r.URL.Query().Get("action"); action != "" {
		filter["action"] = action
	}

	logs, total, err := h.auditLogService.ListAuditLogs(r.Context(), filter, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, logs, page, limit, total)
}
