package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// UpdateRolloutService handles update rollout business logic
type UpdateRolloutService struct {
	rolloutRepo   *repository.UpdateRolloutRepository
	detectionRepo *repository.UpdateDetectionRepository
	versionRepo   *repository.VersionRepository
	productRepo   *repository.ProductRepository
}

// NewUpdateRolloutService creates a new update rollout service
func NewUpdateRolloutService(rolloutRepo *repository.UpdateRolloutRepository, detectionRepo *repository.UpdateDetectionRepository, versionRepo *repository.VersionRepository, productRepo *repository.ProductRepository) *UpdateRolloutService {
	return &UpdateRolloutService{
		rolloutRepo:   rolloutRepo,
		detectionRepo: detectionRepo,
		versionRepo:   versionRepo,
		productRepo:   productRepo,
	}
}

// InitiateRollout initiates a new update rollout
func (s *UpdateRolloutService) InitiateRollout(ctx context.Context, rollout *models.UpdateRollout) (*models.UpdateRollout, error) {
	// Validate product exists
	_, err := s.productRepo.GetByProductID(ctx, rollout.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product %s not found: %w", rollout.ProductID, err)
	}

	// Validate that versions exist for the product
	_, err = s.versionRepo.GetByProductIDAndVersion(ctx, rollout.ProductID, rollout.FromVersion)
	if err != nil {
		return nil, fmt.Errorf("from version %s not found for product %s: %w", rollout.FromVersion, rollout.ProductID, err)
	}

	toVersion, err := s.versionRepo.GetByProductIDAndVersion(ctx, rollout.ProductID, rollout.ToVersion)
	if err != nil {
		return nil, fmt.Errorf("to version %s not found for product %s: %w", rollout.ToVersion, rollout.ProductID, err)
	}

	// Validate that to version is in Released state
	if toVersion.State != models.VersionStateReleased {
		return nil, fmt.Errorf("to version %s must be in Released state, current state: %s", rollout.ToVersion, toVersion.State)
	}

	// Verify detection exists
	_, err = s.detectionRepo.GetByEndpointIDAndProductID(ctx, rollout.EndpointID, rollout.ProductID)
	if err != nil {
		return nil, fmt.Errorf("update detection not found: %w", err)
	}

	rollout.Status = models.RolloutStatusPending
	rollout.Progress = 0

	if err := s.rolloutRepo.Create(ctx, rollout); err != nil {
		return nil, fmt.Errorf("failed to initiate rollout: %w", err)
	}

	return rollout, nil
}

// GetRollout retrieves a rollout by ID
func (s *UpdateRolloutService) GetRollout(ctx context.Context, id primitive.ObjectID) (*models.UpdateRollout, error) {
	rollout, err := s.rolloutRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("rollout not found: %w", err)
	}
	return rollout, nil
}

// UpdateRolloutStatus updates the status of a rollout
func (s *UpdateRolloutService) UpdateRolloutStatus(ctx context.Context, id primitive.ObjectID, status models.RolloutStatus, errorMessage string) (*models.UpdateRollout, error) {
	if err := s.rolloutRepo.UpdateStatus(ctx, id, status, errorMessage); err != nil {
		return nil, fmt.Errorf("failed to update rollout status: %w", err)
	}
	return s.GetRollout(ctx, id)
}

// UpdateRolloutProgress updates the progress of a rollout
func (s *UpdateRolloutService) UpdateRolloutProgress(ctx context.Context, id primitive.ObjectID, progress int) (*models.UpdateRollout, error) {
	if err := s.rolloutRepo.UpdateProgress(ctx, id, progress); err != nil {
		return nil, fmt.Errorf("failed to update rollout progress: %w", err)
	}
	return s.GetRollout(ctx, id)
}

// ListRollouts lists rollouts with filters
func (s *UpdateRolloutService) ListRollouts(ctx context.Context, filter bson.M, page, limit int) ([]*models.UpdateRollout, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"initiated_at": -1})

	rollouts, err := s.rolloutRepo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list rollouts: %w", err)
	}

	total, err := s.rolloutRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count rollouts: %w", err)
	}

	return rollouts, total, nil
}
