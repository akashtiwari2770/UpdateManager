package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// UpgradePathService handles upgrade path business logic
type UpgradePathService struct {
	upgradePathRepo *repository.UpgradePathRepository
	versionRepo     *repository.VersionRepository
}

// NewUpgradePathService creates a new upgrade path service
func NewUpgradePathService(upgradePathRepo *repository.UpgradePathRepository, versionRepo *repository.VersionRepository) *UpgradePathService {
	return &UpgradePathService{
		upgradePathRepo: upgradePathRepo,
		versionRepo:     versionRepo,
	}
}

// CreateUpgradePath creates a new upgrade path with validation
func (s *UpgradePathService) CreateUpgradePath(ctx context.Context, path *models.UpgradePath) error {
	// Validate versions exist
	_, err := s.versionRepo.GetByProductIDAndVersion(ctx, path.ProductID, path.FromVersion)
	if err != nil {
		return fmt.Errorf("from_version '%s' not found: %w", path.FromVersion, err)
	}

	_, err = s.versionRepo.GetByProductIDAndVersion(ctx, path.ProductID, path.ToVersion)
	if err != nil {
		return fmt.Errorf("to_version '%s' not found: %w", path.ToVersion, err)
	}

	// Check if path already exists
	existing, err := s.upgradePathRepo.GetByProductIDAndVersions(ctx, path.ProductID, path.FromVersion, path.ToVersion)
	if err == nil && existing != nil {
		return fmt.Errorf("upgrade path from '%s' to '%s' already exists", path.FromVersion, path.ToVersion)
	}

	if err := s.upgradePathRepo.Create(ctx, path); err != nil {
		return fmt.Errorf("failed to create upgrade path: %w", err)
	}

	return nil
}

// GetUpgradePath retrieves an upgrade path
func (s *UpgradePathService) GetUpgradePath(ctx context.Context, productID, fromVersion, toVersion string) (*models.UpgradePath, error) {
	path, err := s.upgradePathRepo.GetByProductIDAndVersions(ctx, productID, fromVersion, toVersion)
	if err != nil {
		return nil, fmt.Errorf("upgrade path not found: %w", err)
	}
	return path, nil
}

// GetUpgradePathsByProduct retrieves all upgrade paths for a product
func (s *UpgradePathService) GetUpgradePathsByProduct(ctx context.Context, productID string, page, limit int) ([]*models.UpgradePath, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"created_at": -1})

	paths, err := s.upgradePathRepo.GetByProductID(ctx, productID, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get upgrade paths: %w", err)
	}

	total, err := s.upgradePathRepo.Count(ctx, bson.M{"product_id": productID})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count upgrade paths: %w", err)
	}

	return paths, total, nil
}

// BlockUpgradePath blocks an upgrade path
func (s *UpgradePathService) BlockUpgradePath(ctx context.Context, productID, fromVersion, toVersion, reason string) error {
	path, err := s.upgradePathRepo.GetByProductIDAndVersions(ctx, productID, fromVersion, toVersion)
	if err != nil {
		return fmt.Errorf("upgrade path not found: %w", err)
	}

	path.IsBlocked = true
	path.BlockReason = reason
	path.PathType = models.UpgradePathTypeBlocked

	if err := s.upgradePathRepo.Update(ctx, path); err != nil {
		return fmt.Errorf("failed to block upgrade path: %w", err)
	}

	return nil
}
