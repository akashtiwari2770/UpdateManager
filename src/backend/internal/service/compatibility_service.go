package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// CompatibilityService handles compatibility matrix business logic
type CompatibilityService struct {
	compatibilityRepo *repository.CompatibilityRepository
	versionRepo       *repository.VersionRepository
	auditRepo         *repository.AuditLogRepository
}

// NewCompatibilityService creates a new compatibility service
func NewCompatibilityService(compatibilityRepo *repository.CompatibilityRepository, versionRepo *repository.VersionRepository, auditRepo *repository.AuditLogRepository) *CompatibilityService {
	return &CompatibilityService{
		compatibilityRepo: compatibilityRepo,
		versionRepo:       versionRepo,
		auditRepo:         auditRepo,
	}
}

// ValidateCompatibility validates compatibility for a version
func (s *CompatibilityService) ValidateCompatibility(ctx context.Context, productID, versionNumber string, req *models.ValidateCompatibilityRequest, validatedBy string) (*models.CompatibilityMatrix, error) {
	// Verify version exists
	_, err := s.versionRepo.GetByProductIDAndVersion(ctx, productID, versionNumber)
	if err != nil {
		return nil, fmt.Errorf("version not found: %w", err)
	}

	// Check if compatibility matrix already exists
	existing, err := s.compatibilityRepo.GetByProductIDAndVersion(ctx, productID, versionNumber)
	if err == nil && existing != nil {
		// Update existing
		existing.MinServerVersion = req.MinServerVersion
		existing.MaxServerVersion = req.MaxServerVersion
		existing.RecommendedServerVersion = req.RecommendedServerVersion
		existing.IncompatibleVersions = req.IncompatibleVersions
		existing.ValidatedBy = validatedBy
		existing.ValidationStatus = models.ValidationStatusPassed
		existing.ValidationErrors = []string{}

		if err := s.compatibilityRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update compatibility matrix: %w", err)
		}

		return existing, nil
	}

	// Create new compatibility matrix
	matrix := &models.CompatibilityMatrix{
		ProductID:                productID,
		VersionNumber:            versionNumber,
		MinServerVersion:         req.MinServerVersion,
		MaxServerVersion:         req.MaxServerVersion,
		RecommendedServerVersion: req.RecommendedServerVersion,
		IncompatibleVersions:     req.IncompatibleVersions,
		ValidatedBy:              validatedBy,
		ValidationStatus:         models.ValidationStatusPassed,
		ValidationErrors:         []string{},
	}

	if err := s.compatibilityRepo.Create(ctx, matrix); err != nil {
		return nil, fmt.Errorf("failed to create compatibility matrix: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "compatibility_matrix", matrix.ID.Hex(), validatedBy, "", map[string]interface{}{
		"product_id":     productID,
		"version_number": versionNumber,
	})

	return matrix, nil
}

// GetCompatibility retrieves a compatibility matrix
func (s *CompatibilityService) GetCompatibility(ctx context.Context, productID, versionNumber string) (*models.CompatibilityMatrix, error) {
	matrix, err := s.compatibilityRepo.GetByProductIDAndVersion(ctx, productID, versionNumber)
	if err != nil {
		return nil, fmt.Errorf("compatibility matrix not found: %w", err)
	}
	return matrix, nil
}

// ListCompatibility lists compatibility matrices
func (s *CompatibilityService) ListCompatibility(ctx context.Context, filter bson.M, page, limit int) ([]*models.CompatibilityMatrix, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"validated_at": -1})

	matrices, err := s.compatibilityRepo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list compatibility matrices: %w", err)
	}

	total, err := s.compatibilityRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count compatibility matrices: %w", err)
	}

	return matrices, total, nil
}

// logAudit logs an audit entry
func (s *CompatibilityService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
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
