package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// VersionService handles version business logic
type VersionService struct {
	versionRepo *repository.VersionRepository
	productRepo *repository.ProductRepository
	auditRepo   *repository.AuditLogRepository
}

// NewVersionService creates a new version service
func NewVersionService(versionRepo *repository.VersionRepository, productRepo *repository.ProductRepository, auditRepo *repository.AuditLogRepository) *VersionService {
	return &VersionService{
		versionRepo: versionRepo,
		productRepo: productRepo,
		auditRepo:   auditRepo,
	}
}

// CreateVersion creates a new version with validation
func (s *VersionService) CreateVersion(ctx context.Context, productID string, req *models.CreateVersionRequest, createdBy string) (*models.Version, error) {
	// Validate product exists
	_, err := s.productRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check version uniqueness
	existing, err := s.versionRepo.GetByProductIDAndVersion(ctx, productID, req.VersionNumber)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("version '%s' already exists for product '%s'", req.VersionNumber, productID)
	}

	// Create version
	version := &models.Version{
		ProductID:                productID,
		VersionNumber:            req.VersionNumber,
		ReleaseDate:              req.ReleaseDate,
		ReleaseType:              req.ReleaseType,
		State:                    models.VersionStateDraft,
		EOLDate:                  req.EOLDate,
		MinServerVersion:         req.MinServerVersion,
		MaxServerVersion:         req.MaxServerVersion,
		RecommendedServerVersion: req.RecommendedServerVersion,
		ReleaseNotes:             req.ReleaseNotes,
		CreatedBy:                createdBy,
	}

	if err := s.versionRepo.Create(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to create version: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "version", version.ID.Hex(), createdBy, "", map[string]interface{}{
		"product_id":     productID,
		"version_number": req.VersionNumber,
		"release_type":   req.ReleaseType,
	})

	return version, nil
}

// GetVersion retrieves a version by ID
func (s *VersionService) GetVersion(ctx context.Context, id primitive.ObjectID) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}
	return version, nil
}

// GetVersionByProductAndVersion retrieves a version by product_id and version_number
func (s *VersionService) GetVersionByProductAndVersion(ctx context.Context, productID, versionNumber string) (*models.Version, error) {
	version, err := s.versionRepo.GetByProductIDAndVersion(ctx, productID, versionNumber)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}
	return version, nil
}

// GetVersionsByProduct retrieves all versions for a product
func (s *VersionService) GetVersionsByProduct(ctx context.Context, productID string, page, limit int) ([]*models.Version, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"release_date": -1})

	versions, err := s.versionRepo.GetByProductID(ctx, productID, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get versions: %w", err)
	}

	total, err := s.versionRepo.Count(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count versions: %w", err)
	}

	return versions, total, nil
}

// UpdateVersion updates an existing version
func (s *VersionService) UpdateVersion(ctx context.Context, id primitive.ObjectID, req *models.UpdateVersionRequest, userID string) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Only allow updates to draft versions
	if version.State != models.VersionStateDraft {
		return nil, fmt.Errorf("can only update draft versions, current state: %s", version.State)
	}

	// Update fields
	if req.ReleaseDate != nil {
		version.ReleaseDate = *req.ReleaseDate
	}
	if req.ReleaseType != nil {
		version.ReleaseType = *req.ReleaseType
	}
	if req.EOLDate != nil {
		version.EOLDate = req.EOLDate
	}
	if req.MinServerVersion != nil {
		version.MinServerVersion = *req.MinServerVersion
	}
	if req.MaxServerVersion != nil {
		version.MaxServerVersion = *req.MaxServerVersion
	}
	if req.RecommendedServerVersion != nil {
		version.RecommendedServerVersion = *req.RecommendedServerVersion
	}
	if req.ReleaseNotes != nil {
		version.ReleaseNotes = req.ReleaseNotes
	}

	if err := s.versionRepo.Update(ctx, version); err != nil {
		return nil, fmt.Errorf("failed to update version: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "version", version.ID.Hex(), userID, "", map[string]interface{}{
		"product_id":     version.ProductID,
		"version_number": version.VersionNumber,
	})

	return version, nil
}

// ApproveVersion approves a version
func (s *VersionService) ApproveVersion(ctx context.Context, id primitive.ObjectID, req *models.ApproveVersionRequest) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Validate state transition
	if version.State != models.VersionStatePendingReview {
		return nil, fmt.Errorf("can only approve versions in pending_review state, current state: %s", version.State)
	}

	// Update state
	if err := s.versionRepo.UpdateState(ctx, id, models.VersionStateApproved, req.ApprovedBy); err != nil {
		return nil, fmt.Errorf("failed to approve version: %w", err)
	}

	// Get updated version
	version, err = s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated version: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionApprove, "version", version.ID.Hex(), req.ApprovedBy, "", map[string]interface{}{
		"product_id":     version.ProductID,
		"version_number": version.VersionNumber,
	})

	return version, nil
}

// SubmitForReview submits a draft version for review
func (s *VersionService) SubmitForReview(ctx context.Context, id primitive.ObjectID, userID string) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	if version.State != models.VersionStateDraft {
		return nil, fmt.Errorf("can only submit draft versions for review, current state: %s", version.State)
	}

	if err := s.versionRepo.UpdateState(ctx, id, models.VersionStatePendingReview, ""); err != nil {
		return nil, fmt.Errorf("failed to submit version for review: %w", err)
	}

	version, err = s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated version: %w", err)
	}

	s.logAudit(ctx, models.AuditActionUpdate, "version", version.ID.Hex(), userID, "", map[string]interface{}{
		"action":         "submit_for_review",
		"product_id":     version.ProductID,
		"version_number": version.VersionNumber,
	})

	return version, nil
}

// ReleaseVersion releases an approved version
func (s *VersionService) ReleaseVersion(ctx context.Context, id primitive.ObjectID, userID string) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	if version.State != models.VersionStateApproved {
		return nil, fmt.Errorf("can only release approved versions, current state: %s", version.State)
	}

	if err := s.versionRepo.UpdateState(ctx, id, models.VersionStateReleased, ""); err != nil {
		return nil, fmt.Errorf("failed to release version: %w", err)
	}

	version, err = s.versionRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated version: %w", err)
	}

	s.logAudit(ctx, models.AuditActionRelease, "version", version.ID.Hex(), userID, "", map[string]interface{}{
		"product_id":     version.ProductID,
		"version_number": version.VersionNumber,
	})

	return version, nil
}

// GetVersionsByState retrieves versions by state
func (s *VersionService) GetVersionsByState(ctx context.Context, state models.VersionState, page, limit int) ([]*models.Version, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"created_at": -1})

	versions, err := s.versionRepo.GetByState(ctx, state, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get versions by state: %w", err)
	}

	total, err := s.versionRepo.Count(ctx, bson.M{"state": state})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count versions: %w", err)
	}

	return versions, total, nil
}

// ListVersions lists all versions with optional filters
func (s *VersionService) ListVersions(ctx context.Context, filter bson.M, page, limit int) ([]*models.Version, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"release_date": -1})

	versions, err := s.versionRepo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list versions: %w", err)
	}

	// Ensure we return an empty slice instead of nil
	if versions == nil {
		versions = []*models.Version{}
	}

	total, err := s.versionRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count versions: %w", err)
	}

	return versions, total, nil
}

// logAudit logs an audit entry
func (s *VersionService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
	if s.auditRepo == nil {
		return
	}

	auditLog := &models.AuditLog{
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserID:       userID,
		UserEmail:    userEmail,
		Details:      details,
		Timestamp:    time.Now(),
	}

	_ = s.auditRepo.Create(ctx, auditLog)
}

// AddPackageToVersion adds a package to a version
func (s *VersionService) AddPackageToVersion(ctx context.Context, versionID primitive.ObjectID, packageInfo *models.PackageInfo) (*models.Version, error) {
	version, err := s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Only allow adding packages to draft versions
	if version.State != models.VersionStateDraft {
		return nil, fmt.Errorf("packages can only be added to draft versions")
	}

	// Add package to version
	if version.Packages == nil {
		version.Packages = []models.PackageInfo{}
	}
	version.Packages = append(version.Packages, *packageInfo)
	version.UpdatedAt = time.Now()

	// Update version in database using a specific update for packages
	if err := s.versionRepo.AddPackage(ctx, versionID, *packageInfo); err != nil {
		return nil, fmt.Errorf("failed to update version: %w", err)
	}

	// Reload version to get the updated data
	version, err = s.versionRepo.GetByID(ctx, versionID)
	if err != nil {
		return nil, fmt.Errorf("failed to reload version: %w", err)
	}

	return version, nil
}
