package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
)

// VersionHandler handles version-related HTTP requests
type VersionHandler struct {
	versionService      *service.VersionService
	pendingUpdatesService *service.PendingUpdatesService
}

// NewVersionHandler creates a new version handler
func NewVersionHandler(versionService *service.VersionService, pendingUpdatesService *service.PendingUpdatesService) *VersionHandler {
	return &VersionHandler{
		versionService:       versionService,
		pendingUpdatesService: pendingUpdatesService,
	}
}

// CreateVersion handles POST /api/v1/products/:product_id/versions
func (h *VersionHandler) CreateVersion(w http.ResponseWriter, r *http.Request) {
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

	var req models.CreateVersionRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	createdBy := r.Header.Get("X-User-ID")
	if createdBy == "" {
		createdBy = "anonymous"
	}

	version, err := h.versionService.CreateVersion(r.Context(), productID, &req, createdBy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
			return
		}
		if strings.Contains(err.Error(), "already exists") {
			utils.WriteError(w, http.StatusConflict, "DUPLICATE_VERSION", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "CREATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusCreated, version)
}

// ListVersions handles GET /api/v1/versions
func (h *VersionHandler) ListVersions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Build filter from query parameters
	filter := bson.M{}
	
	if productID := r.URL.Query().Get("product_id"); productID != "" {
		filter["product_id"] = productID
	}
	
	if state := r.URL.Query().Get("state"); state != "" {
		filter["state"] = state
	}
	
	if releaseType := r.URL.Query().Get("release_type"); releaseType != "" {
		filter["release_type"] = releaseType
	}

	page := utils.GetIntQueryParam(r, "page", 1)
	limit := utils.GetIntQueryParam(r, "limit", 25)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 25
	}

	versions, total, err := h.versionService.ListVersions(r.Context(), filter, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, versions, page, limit, total)
}

// GetVersion handles GET /api/v1/versions/:id
func (h *VersionHandler) GetVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	// Remove trailing slash if present
	idStr = strings.TrimSuffix(idStr, "/")
	
	// Check if ID is empty (should not happen, but handle gracefully)
	if idStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Version ID is required")
		return
	}
	
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	version, err := h.versionService.GetVersion(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, version)
}

// GetVersionsByProduct handles GET /api/v1/products/:product_id/versions
func (h *VersionHandler) GetVersionsByProduct(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
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

	page := utils.GetIntQueryParam(r, "page", 1)
	limit := utils.GetIntQueryParam(r, "limit", 10)
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	versions, total, err := h.versionService.GetVersionsByProduct(r.Context(), productID, page, limit)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "LIST_FAILED", err.Error())
		return
	}

	utils.WritePaginated(w, http.StatusOK, versions, page, limit, total)
}

// UpdateVersion handles PUT /api/v1/versions/:id
func (h *VersionHandler) UpdateVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	var req models.UpdateVersionRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	version, err := h.versionService.UpdateVersion(r.Context(), id, &req, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		if strings.Contains(err.Error(), "can only update draft") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPDATE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, version)
}

// SubmitForReview handles POST /api/v1/versions/:id/submit
func (h *VersionHandler) SubmitForReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	idStr = strings.TrimSuffix(idStr, "/submit")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	version, err := h.versionService.SubmitForReview(r.Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		if strings.Contains(err.Error(), "can only submit draft") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "SUBMIT_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, version)
}

// ApproveVersion handles POST /api/v1/versions/:id/approve
func (h *VersionHandler) ApproveVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	idStr = strings.TrimSuffix(idStr, "/approve")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	var req models.ApproveVersionRequest
	if err := utils.ReadJSON(w, r, &req); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid JSON: "+err.Error())
		return
	}

	version, err := h.versionService.ApproveVersion(r.Context(), id, &req)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		if strings.Contains(err.Error(), "can only approve") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "APPROVE_FAILED", err.Error())
		return
	}

	utils.WriteSuccess(w, http.StatusOK, version)
}

// ReleaseVersion handles POST /api/v1/versions/:id/release
func (h *VersionHandler) ReleaseVersion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	idStr = strings.TrimSuffix(idStr, "/release")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		userID = "anonymous"
	}

	version, err := h.versionService.ReleaseVersion(r.Context(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		if strings.Contains(err.Error(), "can only release") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "RELEASE_FAILED", err.Error())
		return
	}

	// Invalidate pending updates cache for this product when a new version is released
	// This ensures pending updates are recalculated with the new version
	if h.pendingUpdatesService != nil {
		h.pendingUpdatesService.InvalidateCacheForProduct(r.Context(), version.ProductID)
	}

	utils.WriteSuccess(w, http.StatusOK, version)
}

// ListPackages handles GET /api/v1/versions/:id/packages
func (h *VersionHandler) ListPackages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract version ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	idStr = strings.TrimSuffix(idStr, "/packages")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	// Get version to retrieve packages
	version, err := h.versionService.GetVersion(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	// Return packages array (empty array if nil)
	packages := version.Packages
	if packages == nil {
		packages = []models.PackageInfo{}
	}

	response := map[string]interface{}{
		"packages": packages,
	}
	utils.WriteSuccess(w, http.StatusOK, response)
}

// UploadPackage handles POST /api/v1/versions/:id/packages
func (h *VersionHandler) UploadPackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Parse multipart form (max 10GB)
	if err := r.ParseMultipartForm(10 << 30); err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Failed to parse multipart form: "+err.Error())
		return
	}

	// Extract version ID from path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/versions/")
	idStr = strings.TrimSuffix(idStr, "/packages")
	idStr = strings.TrimSuffix(idStr, "/")
	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	// Get version to check if it exists and is in draft state
	version, err := h.versionService.GetVersion(r.Context(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	// Only allow package uploads for draft versions
	if version.State != models.VersionStateDraft {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", "Packages can only be uploaded for draft versions")
		return
	}

	// Get file from form
	file, header, err := r.FormFile("file")
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "File is required: "+err.Error())
		return
	}
	defer file.Close()

	// Get package metadata
	packageTypeStr := r.FormValue("package_type")
	if packageTypeStr == "" {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "package_type is required")
		return
	}

	packageType := models.PackageType(packageTypeStr)
	if packageType != models.PackageTypeFullInstaller &&
		packageType != models.PackageTypeUpdate &&
		packageType != models.PackageTypeDelta &&
		packageType != models.PackageTypeRollback {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_REQUEST", "Invalid package_type")
		return
	}

	// Calculate file size and checksum while reading
	fileSize := header.Size
	hasher := sha256.New()
	teeReader := io.TeeReader(file, hasher)

	// Create storage directory structure: storage/packages/{version_id}/
	storageBaseDir := "storage/packages"
	versionDir := filepath.Join(storageBaseDir, id.Hex())
	if err := os.MkdirAll(versionDir, 0755); err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_FAILED", "Failed to create storage directory: "+err.Error())
		return
	}

	// Generate unique filename to avoid conflicts
	packageID := primitive.NewObjectID()
	fileExt := filepath.Ext(header.Filename)
	fileName := packageID.Hex() + fileExt
	filePath := filepath.Join(versionDir, fileName)

	// Create and write file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_FAILED", "Failed to create file: "+err.Error())
		return
	}
	defer dst.Close()

	// Copy file content while calculating checksum
	written, err := io.Copy(dst, teeReader)
	if err != nil {
		os.Remove(filePath) // Clean up on error
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_FAILED", "Failed to save file: "+err.Error())
		return
	}

	// Verify file size matches
	if written != fileSize {
		os.Remove(filePath) // Clean up on error
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_FAILED", "File size mismatch")
		return
	}

	// Get checksum
	checksum := hex.EncodeToString(hasher.Sum(nil))

	// Create package info
	uploadedBy := r.Header.Get("X-User-ID")
	if uploadedBy == "" {
		uploadedBy = "anonymous"
	}

	packageInfo := models.PackageInfo{
		ID:             packageID,
		PackageType:    packageType,
		FileName:       header.Filename,
		FileSize:       fileSize,
		ChecksumSHA256: checksum,
		OS:             r.FormValue("os"),
		Architecture:   r.FormValue("architecture"),
		UploadedAt:     time.Now(),
		UploadedBy:     uploadedBy,
		DownloadURL:    fmt.Sprintf("/api/v1/versions/%s/packages/%s/download", id.Hex(), packageID.Hex()),
	}

	// Add package to version
	_, err = h.versionService.AddPackageToVersion(r.Context(), id, &packageInfo)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		if strings.Contains(err.Error(), "can only be added to draft") {
			utils.WriteError(w, http.StatusBadRequest, "INVALID_STATE", err.Error())
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "UPLOAD_FAILED", err.Error())
		return
	}

	// Return success response
	response := map[string]interface{}{
		"message": "Package uploaded successfully",
		"package": packageInfo,
	}
	utils.WriteSuccess(w, http.StatusCreated, response)
}

// DownloadPackage handles GET /api/v1/versions/:id/packages/:package_id/download
func (h *VersionHandler) DownloadPackage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.WriteError(w, http.StatusMethodNotAllowed, "METHOD_NOT_ALLOWED", "Method not allowed")
		return
	}

	// Extract version ID and package ID from path
	pathParts := strings.Split(strings.TrimPrefix(r.URL.Path, "/api/v1/versions/"), "/")
	if len(pathParts) < 3 || pathParts[1] != "packages" || !strings.HasSuffix(pathParts[2], "/download") {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_PATH", "Invalid path format")
		return
	}

	versionIDStr := pathParts[0]
	packageIDStr := strings.TrimSuffix(pathParts[2], "/download")

	versionID, err := primitive.ObjectIDFromHex(versionIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid version ID format")
		return
	}

	packageID, err := primitive.ObjectIDFromHex(packageIDStr)
	if err != nil {
		utils.WriteError(w, http.StatusBadRequest, "INVALID_ID", "Invalid package ID format")
		return
	}

	// Get version to find the package
	version, err := h.versionService.GetVersion(r.Context(), versionID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			utils.WriteError(w, http.StatusNotFound, "VERSION_NOT_FOUND", "Version not found")
			return
		}
		utils.WriteError(w, http.StatusInternalServerError, "GET_FAILED", err.Error())
		return
	}

	// Find the package
	var packageInfo *models.PackageInfo
	for i := range version.Packages {
		if version.Packages[i].ID == packageID {
			packageInfo = &version.Packages[i]
			break
		}
	}

	if packageInfo == nil {
		utils.WriteError(w, http.StatusNotFound, "PACKAGE_NOT_FOUND", "Package not found")
		return
	}

	// Construct file path
	storageBaseDir := "storage/packages"
	versionDir := filepath.Join(storageBaseDir, versionID.Hex())
	fileExt := filepath.Ext(packageInfo.FileName)
	fileName := packageID.Hex() + fileExt
	filePath := filepath.Join(versionDir, fileName)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		utils.WriteError(w, http.StatusNotFound, "FILE_NOT_FOUND", "Package file not found on disk")
		return
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "DOWNLOAD_FAILED", "Failed to open file: "+err.Error())
		return
	}
	defer file.Close()

	// Get file info for Content-Length
	fileInfo, err := file.Stat()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError, "DOWNLOAD_FAILED", "Failed to get file info: "+err.Error())
		return
	}

	// Set headers
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", packageInfo.FileName))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	w.Header().Set("X-Checksum-SHA256", packageInfo.ChecksumSHA256)

	// Stream file to response
	http.ServeContent(w, r, packageInfo.FileName, fileInfo.ModTime(), file)
}
