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

// UpdateDetectionRepository handles update detection database operations
type UpdateDetectionRepository struct {
	collection *mongo.Collection
}

// NewUpdateDetectionRepository creates a new update detection repository
func NewUpdateDetectionRepository(collection *mongo.Collection) *UpdateDetectionRepository {
	return &UpdateDetectionRepository{
		collection: collection,
	}
}

// Create creates a new update detection in the database
func (r *UpdateDetectionRepository) Create(ctx context.Context, detection *models.UpdateDetection) error {
	// Set timestamps
	now := time.Now()
	detection.DetectedAt = now
	detection.LastCheckedAt = now

	// Insert update detection
	result, err := r.collection.InsertOne(ctx, detection)
	if err != nil {
		return fmt.Errorf("failed to create update detection: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		detection.ID = oid
	}

	return nil
}

// GetByID retrieves an update detection by its ID
func (r *UpdateDetectionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.UpdateDetection, error) {
	var detection models.UpdateDetection
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&detection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("update detection not found")
		}
		return nil, fmt.Errorf("failed to get update detection: %w", err)
	}
	return &detection, nil
}

// GetByEndpointIDAndProductID retrieves an update detection by endpoint_id and product_id
func (r *UpdateDetectionRepository) GetByEndpointIDAndProductID(ctx context.Context, endpointID, productID string) (*models.UpdateDetection, error) {
	var detection models.UpdateDetection
	err := r.collection.FindOne(ctx, bson.M{
		"endpoint_id": endpointID,
		"product_id":  productID,
	}).Decode(&detection)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("update detection not found")
		}
		return nil, fmt.Errorf("failed to get update detection: %w", err)
	}
	return &detection, nil
}

// UpdateLastChecked updates the last_checked_at timestamp
func (r *UpdateDetectionRepository) UpdateLastChecked(ctx context.Context, id primitive.ObjectID) error {
	update := bson.M{
		"$set": bson.M{
			"last_checked_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update last checked: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("update detection not found")
	}

	return nil
}

// UpdateAvailableVersion updates the available version for a detection
func (r *UpdateDetectionRepository) UpdateAvailableVersion(ctx context.Context, id primitive.ObjectID, availableVersion string) error {
	update := bson.M{
		"$set": bson.M{
			"available_version": availableVersion,
			"detected_at":       time.Now(),
			"last_checked_at":   time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update available version: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("update detection not found")
	}

	return nil
}

// Update updates an existing update detection
func (r *UpdateDetectionRepository) Update(ctx context.Context, detection *models.UpdateDetection) error {
	detection.LastCheckedAt = time.Now()

	filter := bson.M{"_id": detection.ID}
	update := bson.M{"$set": detection}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update update detection: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("update detection not found")
	}

	return nil
}

// Delete deletes an update detection by ID
func (r *UpdateDetectionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete update detection: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("update detection not found")
	}

	return nil
}

// List retrieves update detections with optional filters
func (r *UpdateDetectionRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.UpdateDetection, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list update detections: %w", err)
	}
	defer cursor.Close(ctx)

	var detections []*models.UpdateDetection
	if err := cursor.All(ctx, &detections); err != nil {
		return nil, fmt.Errorf("failed to decode update detections: %w", err)
	}

	return detections, nil
}

// Count returns the count of update detections matching the filter
func (r *UpdateDetectionRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count update detections: %w", err)
	}
	return count, nil
}
