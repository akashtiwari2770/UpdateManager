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

// AuditLogRepository handles audit log database operations
type AuditLogRepository struct {
	collection *mongo.Collection
}

// NewAuditLogRepository creates a new audit log repository
func NewAuditLogRepository(collection *mongo.Collection) *AuditLogRepository {
	return &AuditLogRepository{
		collection: collection,
	}
}

// Create creates a new audit log entry in the database
func (r *AuditLogRepository) Create(ctx context.Context, log *models.AuditLog) error {
	// Set timestamp
	log.Timestamp = time.Now()

	// Insert audit log
	result, err := r.collection.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		log.ID = oid
	}

	return nil
}

// GetByID retrieves an audit log by its ID
func (r *AuditLogRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.AuditLog, error) {
	var log models.AuditLog
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&log)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("audit log not found")
		}
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}
	return &log, nil
}

// GetByResource retrieves audit logs for a specific resource
func (r *AuditLogRepository) GetByResource(ctx context.Context, resourceType, resourceID string, opts *options.FindOptions) ([]*models.AuditLog, error) {
	filter := bson.M{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode audit logs: %w", err)
	}

	return logs, nil
}

// GetByUserID retrieves audit logs for a specific user
func (r *AuditLogRepository) GetByUserID(ctx context.Context, userID string, opts *options.FindOptions) ([]*models.AuditLog, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode audit logs: %w", err)
	}

	return logs, nil
}

// GetByAction retrieves audit logs by action
func (r *AuditLogRepository) GetByAction(ctx context.Context, action models.AuditAction, opts *options.FindOptions) ([]*models.AuditLog, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"action": action}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode audit logs: %w", err)
	}

	return logs, nil
}

// List retrieves audit logs with optional filters
func (r *AuditLogRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.AuditLog, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer cursor.Close(ctx)

	var logs []*models.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("failed to decode audit logs: %w", err)
	}

	return logs, nil
}

// Count returns the count of audit logs matching the filter
func (r *AuditLogRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}
	return count, nil
}

// Delete deletes an audit log by ID (rarely used, but available)
func (r *AuditLogRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete audit log: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("audit log not found")
	}

	return nil
}
