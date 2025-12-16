package repository

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
)

// UpgradePathRepository handles upgrade path database operations
type UpgradePathRepository struct {
	collection *mongo.Collection
}

// NewUpgradePathRepository creates a new upgrade path repository
func NewUpgradePathRepository(collection *mongo.Collection) *UpgradePathRepository {
	return &UpgradePathRepository{
		collection: collection,
	}
}

// Create creates a new upgrade path in the database
func (r *UpgradePathRepository) Create(ctx context.Context, path *models.UpgradePath) error {
	// Set timestamp
	path.CreatedAt = time.Now()

	// Insert upgrade path
	result, err := r.collection.InsertOne(ctx, path)
	if err != nil {
		return fmt.Errorf("failed to create upgrade path: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		path.ID = oid
	}

	return nil
}

// GetByID retrieves an upgrade path by its ID
func (r *UpgradePathRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.UpgradePath, error) {
	var path models.UpgradePath
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&path)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("upgrade path not found")
		}
		return nil, fmt.Errorf("failed to get upgrade path: %w", err)
	}
	return &path, nil
}

// GetByProductIDAndVersions retrieves an upgrade path by product_id, from_version, and to_version
func (r *UpgradePathRepository) GetByProductIDAndVersions(ctx context.Context, productID, fromVersion, toVersion string) (*models.UpgradePath, error) {
	var path models.UpgradePath
	err := r.collection.FindOne(ctx, bson.M{
		"product_id":   productID,
		"from_version": fromVersion,
		"to_version":   toVersion,
	}).Decode(&path)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("upgrade path not found")
		}
		return nil, fmt.Errorf("failed to get upgrade path: %w", err)
	}
	return &path, nil
}

// GetByProductID retrieves all upgrade paths for a product
func (r *UpgradePathRepository) GetByProductID(ctx context.Context, productID string, opts *options.FindOptions) ([]*models.UpgradePath, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get upgrade paths: %w", err)
	}
	defer cursor.Close(ctx)

	var paths []*models.UpgradePath
	if err := cursor.All(ctx, &paths); err != nil {
		return nil, fmt.Errorf("failed to decode upgrade paths: %w", err)
	}

	return paths, nil
}

// Update updates an existing upgrade path
func (r *UpgradePathRepository) Update(ctx context.Context, path *models.UpgradePath) error {
	filter := bson.M{"_id": path.ID}
	update := bson.M{"$set": path}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update upgrade path: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("upgrade path not found")
	}

	return nil
}

// Delete deletes an upgrade path by ID
func (r *UpgradePathRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete upgrade path: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("upgrade path not found")
	}

	return nil
}

// List retrieves upgrade paths with optional filters
func (r *UpgradePathRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.UpgradePath, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list upgrade paths: %w", err)
	}
	defer cursor.Close(ctx)

	var paths []*models.UpgradePath
	if err := cursor.All(ctx, &paths); err != nil {
		return nil, fmt.Errorf("failed to decode upgrade paths: %w", err)
	}

	return paths, nil
}

// Count returns the count of upgrade paths matching the filter
func (r *UpgradePathRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count upgrade paths: %w", err)
	}
	return count, nil
}
