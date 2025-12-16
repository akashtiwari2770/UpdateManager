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

// UpdateRolloutRepository handles update rollout database operations
type UpdateRolloutRepository struct {
	collection *mongo.Collection
}

// NewUpdateRolloutRepository creates a new update rollout repository
func NewUpdateRolloutRepository(collection *mongo.Collection) *UpdateRolloutRepository {
	return &UpdateRolloutRepository{
		collection: collection,
	}
}

// Create creates a new update rollout in the database
func (r *UpdateRolloutRepository) Create(ctx context.Context, rollout *models.UpdateRollout) error {
	// Set timestamp
	rollout.InitiatedAt = time.Now()

	// Insert update rollout
	result, err := r.collection.InsertOne(ctx, rollout)
	if err != nil {
		return fmt.Errorf("failed to create update rollout: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		rollout.ID = oid
	}

	return nil
}

// GetByID retrieves an update rollout by its ID
func (r *UpdateRolloutRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.UpdateRollout, error) {
	var rollout models.UpdateRollout
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&rollout)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("update rollout not found")
		}
		return nil, fmt.Errorf("failed to get update rollout: %w", err)
	}
	return &rollout, nil
}

// GetByEndpointID retrieves all rollouts for an endpoint
func (r *UpdateRolloutRepository) GetByEndpointID(ctx context.Context, endpointID string, opts *options.FindOptions) ([]*models.UpdateRollout, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"endpoint_id": endpointID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get update rollouts: %w", err)
	}
	defer cursor.Close(ctx)

	var rollouts []*models.UpdateRollout
	if err := cursor.All(ctx, &rollouts); err != nil {
		return nil, fmt.Errorf("failed to decode update rollouts: %w", err)
	}

	return rollouts, nil
}

// GetByStatus retrieves rollouts by status
func (r *UpdateRolloutRepository) GetByStatus(ctx context.Context, status models.RolloutStatus, opts *options.FindOptions) ([]*models.UpdateRollout, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"status": status}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get update rollouts by status: %w", err)
	}
	defer cursor.Close(ctx)

	var rollouts []*models.UpdateRollout
	if err := cursor.All(ctx, &rollouts); err != nil {
		return nil, fmt.Errorf("failed to decode update rollouts: %w", err)
	}

	return rollouts, nil
}

// UpdateStatus updates the status of a rollout
func (r *UpdateRolloutRepository) UpdateStatus(ctx context.Context, id primitive.ObjectID, status models.RolloutStatus, errorMessage string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status": status,
		},
	}

	switch status {
	case models.RolloutStatusInProgress:
		update["$set"].(bson.M)["started_at"] = now
	case models.RolloutStatusCompleted:
		update["$set"].(bson.M)["completed_at"] = now
	case models.RolloutStatusFailed:
		update["$set"].(bson.M)["failed_at"] = now
		if errorMessage != "" {
			update["$set"].(bson.M)["error_message"] = errorMessage
		}
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update rollout status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("update rollout not found")
	}

	return nil
}

// UpdateProgress updates the progress of a rollout
func (r *UpdateRolloutRepository) UpdateProgress(ctx context.Context, id primitive.ObjectID, progress int) error {
	if progress < 0 || progress > 100 {
		return fmt.Errorf("progress must be between 0 and 100")
	}

	update := bson.M{
		"$set": bson.M{
			"progress": progress,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update rollout progress: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("update rollout not found")
	}

	return nil
}

// Update updates an existing update rollout
func (r *UpdateRolloutRepository) Update(ctx context.Context, rollout *models.UpdateRollout) error {
	filter := bson.M{"_id": rollout.ID}
	update := bson.M{"$set": rollout}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update update rollout: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("update rollout not found")
	}

	return nil
}

// Delete deletes an update rollout by ID
func (r *UpdateRolloutRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete update rollout: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("update rollout not found")
	}

	return nil
}

// List retrieves update rollouts with optional filters
func (r *UpdateRolloutRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.UpdateRollout, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list update rollouts: %w", err)
	}
	defer cursor.Close(ctx)

	var rollouts []*models.UpdateRollout
	if err := cursor.All(ctx, &rollouts); err != nil {
		return nil, fmt.Errorf("failed to decode update rollouts: %w", err)
	}

	return rollouts, nil
}

// Count returns the count of update rollouts matching the filter
func (r *UpdateRolloutRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count update rollouts: %w", err)
	}
	return count, nil
}
