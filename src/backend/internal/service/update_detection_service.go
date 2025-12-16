package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// UpdateDetectionService handles update detection business logic
type UpdateDetectionService struct {
	detectionRepo *repository.UpdateDetectionRepository
	versionRepo   *repository.VersionRepository
	productRepo   *repository.ProductRepository
}

// NewUpdateDetectionService creates a new update detection service
func NewUpdateDetectionService(detectionRepo *repository.UpdateDetectionRepository, versionRepo *repository.VersionRepository, productRepo *repository.ProductRepository) *UpdateDetectionService {
	return &UpdateDetectionService{
		detectionRepo: detectionRepo,
		versionRepo:   versionRepo,
		productRepo:   productRepo,
	}
}

// DetectUpdate detects or updates detection for an endpoint
func (s *UpdateDetectionService) DetectUpdate(ctx context.Context, detection *models.UpdateDetection) (*models.UpdateDetection, error) {
	// Validate product exists
	_, err := s.productRepo.GetByProductID(ctx, detection.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product %s not found: %w", detection.ProductID, err)
	}

	// Validate that versions exist for the product
	_, err = s.versionRepo.GetByProductIDAndVersion(ctx, detection.ProductID, detection.CurrentVersion)
	if err != nil {
		return nil, fmt.Errorf("current version %s not found for product %s: %w", detection.CurrentVersion, detection.ProductID, err)
	}

	availableVersion, err := s.versionRepo.GetByProductIDAndVersion(ctx, detection.ProductID, detection.AvailableVersion)
	if err != nil {
		return nil, fmt.Errorf("available version %s not found for product %s: %w", detection.AvailableVersion, detection.ProductID, err)
	}

	// Validate that available version is in Released state
	if availableVersion.State != models.VersionStateReleased {
		return nil, fmt.Errorf("available version %s must be in Released state, current state: %s", detection.AvailableVersion, availableVersion.State)
	}

	// Check if detection already exists
	existing, err := s.detectionRepo.GetByEndpointIDAndProductID(ctx, detection.EndpointID, detection.ProductID)
	if err == nil && existing != nil {
		// Update existing detection
		existing.CurrentVersion = detection.CurrentVersion
		existing.AvailableVersion = detection.AvailableVersion
		if err := s.detectionRepo.Update(ctx, existing); err != nil {
			return nil, fmt.Errorf("failed to update detection: %w", err)
		}
		return existing, nil
	}

	// Create new detection
	if err := s.detectionRepo.Create(ctx, detection); err != nil {
		return nil, fmt.Errorf("failed to create detection: %w", err)
	}

	return detection, nil
}

// GetDetection retrieves detection for an endpoint and product
func (s *UpdateDetectionService) GetDetection(ctx context.Context, endpointID, productID string) (*models.UpdateDetection, error) {
	detection, err := s.detectionRepo.GetByEndpointIDAndProductID(ctx, endpointID, productID)
	if err != nil {
		return nil, fmt.Errorf("detection not found: %w", err)
	}
	return detection, nil
}

// UpdateAvailableVersion updates the available version for a detection
func (s *UpdateDetectionService) UpdateAvailableVersion(ctx context.Context, endpointID, productID, availableVersion string) error {
	detection, err := s.detectionRepo.GetByEndpointIDAndProductID(ctx, endpointID, productID)
	if err != nil {
		return fmt.Errorf("detection not found: %w", err)
	}

	// Validate that the new available version exists for the product
	version, err := s.versionRepo.GetByProductIDAndVersion(ctx, productID, availableVersion)
	if err != nil {
		return fmt.Errorf("available version %s not found for product %s: %w", availableVersion, productID, err)
	}

	// Validate that available version is in Released state
	if version.State != models.VersionStateReleased {
		return fmt.Errorf("available version %s must be in Released state, current state: %s", availableVersion, version.State)
	}

	if err := s.detectionRepo.UpdateAvailableVersion(ctx, detection.ID, availableVersion); err != nil {
		return fmt.Errorf("failed to update available version: %w", err)
	}

	return nil
}

// ListDetections lists detections with filters
func (s *UpdateDetectionService) ListDetections(ctx context.Context, filter bson.M, page, limit int) ([]*models.UpdateDetection, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"detected_at": -1})

	detections, err := s.detectionRepo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list detections: %w", err)
	}

	total, err := s.detectionRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count detections: %w", err)
	}

	return detections, total, nil
}
