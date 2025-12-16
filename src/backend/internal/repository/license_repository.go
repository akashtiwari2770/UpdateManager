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

// LicenseRepository handles license database operations
type LicenseRepository struct {
	collection *mongo.Collection
}

// NewLicenseRepository creates a new license repository
func NewLicenseRepository(collection *mongo.Collection) *LicenseRepository {
	return &LicenseRepository{
		collection: collection,
	}
}

// LicenseFilter represents filters for license queries
type LicenseFilter struct {
	ProductID   string
	LicenseType models.LicenseType
	Status      models.LicenseStatus
	StartDate   *time.Time
	EndDate     *time.Time
}

// Create creates a new license in the database
func (r *LicenseRepository) Create(ctx context.Context, license *models.License) error {
	// Set timestamps
	now := time.Now()
	license.CreatedAt = now
	license.UpdatedAt = now

	// Insert license
	result, err := r.collection.InsertOne(ctx, license)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("license with license_id '%s' already exists", license.LicenseID)
		}
		return fmt.Errorf("failed to create license: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		license.ID = oid
	}

	return nil
}

// GetByID retrieves a license by its ID
func (r *LicenseRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.License, error) {
	var license models.License
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&license)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("license not found")
		}
		return nil, fmt.Errorf("failed to get license: %w", err)
	}
	return &license, nil
}

// GetByLicenseID retrieves a license by its license_id
func (r *LicenseRepository) GetByLicenseID(ctx context.Context, licenseID string) (*models.License, error) {
	var license models.License
	err := r.collection.FindOne(ctx, bson.M{"license_id": licenseID}).Decode(&license)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("license not found")
		}
		return nil, fmt.Errorf("failed to get license: %w", err)
	}
	return &license, nil
}

// GetBySubscriptionID retrieves licenses for a subscription with filters and pagination
func (r *LicenseRepository) GetBySubscriptionID(ctx context.Context, subscriptionID primitive.ObjectID, filter *LicenseFilter, pagination *Pagination) ([]*models.License, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"subscription_id": subscriptionID}

	if filter != nil {
		if filter.ProductID != "" {
			bsonFilter["product_id"] = filter.ProductID
		}
		if filter.LicenseType != "" {
			bsonFilter["license_type"] = filter.LicenseType
		}
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
		return nil, nil, fmt.Errorf("failed to count licenses: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"assignment_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	defer cursor.Close(ctx)

	var licenses []*models.License
	if err := cursor.All(ctx, &licenses); err != nil {
		return nil, nil, fmt.Errorf("failed to decode licenses: %w", err)
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

	return licenses, paginationInfo, nil
}

// GetByProductID retrieves licenses for a product with filters and pagination
func (r *LicenseRepository) GetByProductID(ctx context.Context, productID string, filter *LicenseFilter, pagination *Pagination) ([]*models.License, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"product_id": productID}

	if filter != nil {
		if filter.LicenseType != "" {
			bsonFilter["license_type"] = filter.LicenseType
		}
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
		return nil, nil, fmt.Errorf("failed to count licenses: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"assignment_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list licenses: %w", err)
	}
	defer cursor.Close(ctx)

	var licenses []*models.License
	if err := cursor.All(ctx, &licenses); err != nil {
		return nil, nil, fmt.Errorf("failed to decode licenses: %w", err)
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

	return licenses, paginationInfo, nil
}

// Update updates an existing license
func (r *LicenseRepository) Update(ctx context.Context, id primitive.ObjectID, license *models.License) error {
	license.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": license}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("license with license_id '%s' already exists", license.LicenseID)
		}
		return fmt.Errorf("failed to update license: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("license not found")
	}

	return nil
}

// Delete deletes a license by ID
func (r *LicenseRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to delete license: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("license not found")
	}

	return nil
}

// Count counts licenses matching the filter
func (r *LicenseRepository) Count(ctx context.Context, filter *LicenseFilter) (int64, error) {
	bsonFilter := bson.M{}

	if filter != nil {
		if filter.ProductID != "" {
			bsonFilter["product_id"] = filter.ProductID
		}
		if filter.LicenseType != "" {
			bsonFilter["license_type"] = filter.LicenseType
		}
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
		return 0, fmt.Errorf("failed to count licenses: %w", err)
	}

	return count, nil
}

// GetExpiringLicenses retrieves licenses expiring within the specified number of days
func (r *LicenseRepository) GetExpiringLicenses(ctx context.Context, days int) ([]*models.License, error) {
	now := time.Now()
	expirationDate := now.AddDate(0, 0, days)

	filter := bson.M{
		"end_date": bson.M{
			"$gte": now,
			"$lte": expirationDate,
		},
		"status":       models.LicenseStatusActive,
		"license_type": models.LicenseTypeTimeBased,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"end_date": 1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find expiring licenses: %w", err)
	}
	defer cursor.Close(ctx)

	var licenses []*models.License
	if err := cursor.All(ctx, &licenses); err != nil {
		return nil, fmt.Errorf("failed to decode licenses: %w", err)
	}

	return licenses, nil
}

