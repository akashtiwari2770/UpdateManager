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

// NotificationRepository handles notification database operations
type NotificationRepository struct {
	collection *mongo.Collection
}

// NewNotificationRepository creates a new notification repository
func NewNotificationRepository(collection *mongo.Collection) *NotificationRepository {
	return &NotificationRepository{
		collection: collection,
	}
}

// Create creates a new notification in the database
func (r *NotificationRepository) Create(ctx context.Context, notification *models.Notification) error {
	// Set timestamp
	notification.CreatedAt = time.Now()
	notification.IsRead = false

	// Insert notification
	result, err := r.collection.InsertOne(ctx, notification)
	if err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		notification.ID = oid
	}

	return nil
}

// GetByID retrieves a notification by its ID
func (r *NotificationRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Notification, error) {
	var notification models.Notification
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&notification)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("notification not found")
		}
		return nil, fmt.Errorf("failed to get notification: %w", err)
	}
	return &notification, nil
}

// GetByRecipientID retrieves notifications for a recipient
func (r *NotificationRepository) GetByRecipientID(ctx context.Context, recipientID string, opts *options.FindOptions) ([]*models.Notification, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"recipient_id": recipientID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, nil
}

// GetUnreadByRecipientID retrieves unread notifications for a recipient
func (r *NotificationRepository) GetUnreadByRecipientID(ctx context.Context, recipientID string, opts *options.FindOptions) ([]*models.Notification, error) {
	filter := bson.M{
		"recipient_id": recipientID,
		"is_read":      false,
	}
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread notifications: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, nil
}

// MarkAsRead marks a notification as read
func (r *NotificationRepository) MarkAsRead(ctx context.Context, id primitive.ObjectID) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
			"read_at": now,
		},
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// MarkAllAsRead marks all notifications for a recipient as read
func (r *NotificationRepository) MarkAllAsRead(ctx context.Context, recipientID string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"is_read": true,
			"read_at": now,
		},
	}

	_, err := r.collection.UpdateMany(ctx, bson.M{
		"recipient_id": recipientID,
		"is_read":      false,
	}, update)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

// Update updates an existing notification
func (r *NotificationRepository) Update(ctx context.Context, notification *models.Notification) error {
	filter := bson.M{"_id": notification.ID}
	update := bson.M{"$set": notification}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update notification: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// Delete deletes a notification by ID
func (r *NotificationRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("notification not found")
	}

	return nil
}

// List retrieves notifications with optional filters
func (r *NotificationRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.Notification, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer cursor.Close(ctx)

	var notifications []*models.Notification
	if err := cursor.All(ctx, &notifications); err != nil {
		return nil, fmt.Errorf("failed to decode notifications: %w", err)
	}

	return notifications, nil
}

// Count returns the count of notifications matching the filter
func (r *NotificationRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count notifications: %w", err)
	}
	return count, nil
}
