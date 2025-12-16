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

// SubscriptionRepository handles subscription database operations
type SubscriptionRepository struct {
	collection *mongo.Collection
}

// NewSubscriptionRepository creates a new subscription repository
func NewSubscriptionRepository(collection *mongo.Collection) *SubscriptionRepository {
	return &SubscriptionRepository{
		collection: collection,
	}
}

// SubscriptionFilter represents filters for subscription queries
type SubscriptionFilter struct {
	Status    models.SubscriptionStatus
	StartDate *time.Time
	EndDate   *time.Time
}

// Create creates a new subscription in the database
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *models.Subscription) error {
	// Set timestamps
	now := time.Now()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	// Insert subscription
	result, err := r.collection.InsertOne(ctx, subscription)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("subscription with subscription_id '%s' already exists", subscription.SubscriptionID)
		}
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		subscription.ID = oid
	}

	return nil
}

// GetByID retrieves a subscription by its ID
func (r *SubscriptionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return &subscription, nil
}

// GetBySubscriptionID retrieves a subscription by its subscription_id
func (r *SubscriptionRepository) GetBySubscriptionID(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	var subscription models.Subscription
	err := r.collection.FindOne(ctx, bson.M{"subscription_id": subscriptionID}).Decode(&subscription)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("subscription not found")
		}
		return nil, fmt.Errorf("failed to get subscription: %w", err)
	}
	return &subscription, nil
}

// GetByCustomerID retrieves subscriptions for a customer with filters and pagination
func (r *SubscriptionRepository) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID, filter *SubscriptionFilter, pagination *Pagination) ([]*models.Subscription, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"customer_id": customerID}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
		if filter.StartDate != nil {
			bsonFilter["start_date"] = bson.M{"$gte": *filter.StartDate}
		}
		if filter.EndDate != nil {
			bsonFilter["end_date"] = bson.M{"$lte": *filter.EndDate}
		}
	}

	// Set up pagination
	page := 1
	limit := 20
	if pagination != nil {
		if pagination.Page > 0 {
			page = pagination.Page
		}
		if pagination.Limit > 0 && pagination.Limit <= 100 {
			limit = pagination.Limit
		}
	}

	skip := (page - 1) * limit

	// Count total documents
	total, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"start_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*models.Subscription
	if err := cursor.All(ctx, &subscriptions); err != nil {
		return nil, nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}

	// Calculate total pages
	totalPages := total / int64(limit)
	if total%int64(limit) > 0 {
		totalPages++
	}

	paginationInfo := &PaginationInfo{
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}

	return subscriptions, paginationInfo, nil
}

// Update updates an existing subscription
func (r *SubscriptionRepository) Update(ctx context.Context, id primitive.ObjectID, subscription *models.Subscription) error {
	subscription.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": subscription}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("subscription with subscription_id '%s' already exists", subscription.SubscriptionID)
		}
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// Delete deletes a subscription by ID
func (r *SubscriptionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

// Count counts subscriptions matching the filter
func (r *SubscriptionRepository) Count(ctx context.Context, filter *SubscriptionFilter) (int64, error) {
	bsonFilter := bson.M{}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
		if filter.StartDate != nil {
			bsonFilter["start_date"] = bson.M{"$gte": *filter.StartDate}
		}
		if filter.EndDate != nil {
			bsonFilter["end_date"] = bson.M{"$lte": *filter.EndDate}
		}
	}

	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count subscriptions: %w", err)
	}

	return count, nil
}

// GetExpiringSubscriptions retrieves subscriptions expiring within the specified number of days
func (r *SubscriptionRepository) GetExpiringSubscriptions(ctx context.Context, days int) ([]*models.Subscription, error) {
	now := time.Now()
	expirationDate := now.AddDate(0, 0, days)

	filter := bson.M{
		"end_date": bson.M{
			"$gte": now,
			"$lte": expirationDate,
		},
		"status": models.SubscriptionStatusActive,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"end_date": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find expiring subscriptions: %w", err)
	}
	defer cursor.Close(ctx)

	var subscriptions []*models.Subscription
	if err := cursor.All(ctx, &subscriptions); err != nil {
		return nil, fmt.Errorf("failed to decode subscriptions: %w", err)
	}

	return subscriptions, nil
}

