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

// LicenseAllocationRepository handles license allocation database operations
type LicenseAllocationRepository struct {
	collection *mongo.Collection
}

// NewLicenseAllocationRepository creates a new license allocation repository
func NewLicenseAllocationRepository(collection *mongo.Collection) *LicenseAllocationRepository {
	return &LicenseAllocationRepository{
		collection: collection,
	}
}

// LicenseAllocationFilter represents filters for license allocation queries
type LicenseAllocationFilter struct {
	Status models.AllocationStatus
}

// Create creates a new license allocation in the database
func (r *LicenseAllocationRepository) Create(ctx context.Context, allocation *models.LicenseAllocation) error {
	// Set timestamps
	now := time.Now()
	allocation.CreatedAt = now
	allocation.UpdatedAt = now

	// Insert allocation
	result, err := r.collection.InsertOne(ctx, allocation)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("allocation with allocation_id '%s' already exists", allocation.AllocationID)
		}
		return fmt.Errorf("failed to create license allocation: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		allocation.ID = oid
	}

	return nil
}

// GetByID retrieves a license allocation by its ID
func (r *LicenseAllocationRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.LicenseAllocation, error) {
	var allocation models.LicenseAllocation
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&allocation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("license allocation not found")
		}
		return nil, fmt.Errorf("failed to get license allocation: %w", err)
	}
	return &allocation, nil
}

// GetByAllocationID retrieves a license allocation by its allocation_id
func (r *LicenseAllocationRepository) GetByAllocationID(ctx context.Context, allocationID string) (*models.LicenseAllocation, error) {
	var allocation models.LicenseAllocation
	err := r.collection.FindOne(ctx, bson.M{"allocation_id": allocationID}).Decode(&allocation)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("license allocation not found")
		}
		return nil, fmt.Errorf("failed to get license allocation: %w", err)
	}
	return &allocation, nil
}

// GetByLicenseID retrieves allocations for a license with filters and pagination
func (r *LicenseAllocationRepository) GetByLicenseID(ctx context.Context, licenseID primitive.ObjectID, filter *LicenseAllocationFilter, pagination *Pagination) ([]*models.LicenseAllocation, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"license_id": licenseID}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
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
		return nil, nil, fmt.Errorf("failed to count license allocations: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"allocation_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list license allocations: %w", err)
	}
	defer cursor.Close(ctx)

	var allocations []*models.LicenseAllocation
	if err := cursor.All(ctx, &allocations); err != nil {
		return nil, nil, fmt.Errorf("failed to decode license allocations: %w", err)
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

	return allocations, paginationInfo, nil
}

// GetByTenantID retrieves allocations for a tenant with filters and pagination
func (r *LicenseAllocationRepository) GetByTenantID(ctx context.Context, tenantID primitive.ObjectID, filter *LicenseAllocationFilter, pagination *Pagination) ([]*models.LicenseAllocation, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"tenant_id": tenantID}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
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
		return nil, nil, fmt.Errorf("failed to count license allocations: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"allocation_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list license allocations: %w", err)
	}
	defer cursor.Close(ctx)

	var allocations []*models.LicenseAllocation
	if err := cursor.All(ctx, &allocations); err != nil {
		return nil, nil, fmt.Errorf("failed to decode license allocations: %w", err)
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

	return allocations, paginationInfo, nil
}

// GetByDeploymentID retrieves allocations for a deployment with filters and pagination
func (r *LicenseAllocationRepository) GetByDeploymentID(ctx context.Context, deploymentID primitive.ObjectID, filter *LicenseAllocationFilter, pagination *Pagination) ([]*models.LicenseAllocation, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"deployment_id": deploymentID}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
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
		return nil, nil, fmt.Errorf("failed to count license allocations: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"allocation_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list license allocations: %w", err)
	}
	defer cursor.Close(ctx)

	var allocations []*models.LicenseAllocation
	if err := cursor.All(ctx, &allocations); err != nil {
		return nil, nil, fmt.Errorf("failed to decode license allocations: %w", err)
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

	return allocations, paginationInfo, nil
}

// GetActiveAllocationsByLicenseID retrieves all active allocations for a license
func (r *LicenseAllocationRepository) GetActiveAllocationsByLicenseID(ctx context.Context, licenseID primitive.ObjectID) ([]*models.LicenseAllocation, error) {
	filter := bson.M{
		"license_id": licenseID,
		"status":     models.AllocationStatusActive,
	}

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"allocation_date": -1})

	cursor, err := r.collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to find active allocations: %w", err)
	}
	defer cursor.Close(ctx)

	var allocations []*models.LicenseAllocation
	if err := cursor.All(ctx, &allocations); err != nil {
		return nil, fmt.Errorf("failed to decode allocations: %w", err)
	}

	return allocations, nil
}

// Update updates an existing license allocation
func (r *LicenseAllocationRepository) Update(ctx context.Context, id primitive.ObjectID, allocation *models.LicenseAllocation) error {
	allocation.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": allocation}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("allocation with allocation_id '%s' already exists", allocation.AllocationID)
		}
		return fmt.Errorf("failed to update license allocation: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("license allocation not found")
	}

	return nil
}

// Release releases a license allocation
func (r *LicenseAllocationRepository) Release(ctx context.Context, allocationID primitive.ObjectID, releasedBy string) error {
	now := time.Now()
	update := bson.M{
		"$set": bson.M{
			"status":       models.AllocationStatusReleased,
			"released_date": now,
			"released_by":   releasedBy,
			"updated_at":    now,
		},
	}

	filter := bson.M{"_id": allocationID}
	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to release license allocation: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("license allocation not found")
	}

	return nil
}

// Count counts license allocations matching the filter
func (r *LicenseAllocationRepository) Count(ctx context.Context, filter *LicenseAllocationFilter) (int64, error) {
	bsonFilter := bson.M{}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
	}

	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count license allocations: %w", err)
	}

	return count, nil
}

// GetTotalAllocatedSeats calculates the total number of seats allocated for a license
func (r *LicenseAllocationRepository) GetTotalAllocatedSeats(ctx context.Context, licenseID primitive.ObjectID) (int, error) {
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"license_id": licenseID,
				"status":     models.AllocationStatusActive,
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"total": bson.M{
					"$sum": "$number_of_seats_allocated",
				},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total allocated seats: %w", err)
	}
	defer cursor.Close(ctx)

	var result struct {
		Total int `bson:"total"`
	}

	if cursor.Next(ctx) {
		if err := cursor.Decode(&result); err != nil {
			return 0, fmt.Errorf("failed to decode aggregation result: %w", err)
		}
		return result.Total, nil
	}

	// No allocations found
	return 0, nil
}

