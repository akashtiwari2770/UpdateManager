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

// TenantRepository handles tenant database operations
type TenantRepository struct {
	collection *mongo.Collection
}

// NewTenantRepository creates a new tenant repository
func NewTenantRepository(collection *mongo.Collection) *TenantRepository {
	return &TenantRepository{
		collection: collection,
	}
}

// TenantFilter represents filters for tenant queries
type TenantFilter struct {
	Status models.TenantStatus
}

// Create creates a new tenant in the database
func (r *TenantRepository) Create(ctx context.Context, tenant *models.CustomerTenant) error {
	// Set timestamps
	now := time.Now()
	tenant.CreatedAt = now
	tenant.UpdatedAt = now

	// Insert tenant
	result, err := r.collection.InsertOne(ctx, tenant)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("tenant with tenant_id '%s' already exists", tenant.TenantID)
		}
		return fmt.Errorf("failed to create tenant: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		tenant.ID = oid
	}

	return nil
}

// GetByID retrieves a tenant by its ID
func (r *TenantRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.CustomerTenant, error) {
	var tenant models.CustomerTenant
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return &tenant, nil
}

// GetByTenantID retrieves a tenant by its tenant_id
func (r *TenantRepository) GetByTenantID(ctx context.Context, tenantID string) (*models.CustomerTenant, error) {
	var tenant models.CustomerTenant
	err := r.collection.FindOne(ctx, bson.M{"tenant_id": tenantID}).Decode(&tenant)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("tenant not found")
		}
		return nil, fmt.Errorf("failed to get tenant: %w", err)
	}
	return &tenant, nil
}

// GetByCustomerID retrieves tenants for a customer with filters and pagination
func (r *TenantRepository) GetByCustomerID(ctx context.Context, customerID primitive.ObjectID, filter *TenantFilter, pagination *Pagination) ([]*models.CustomerTenant, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"customer_id": customerID}

	if filter != nil && filter.Status != "" {
		bsonFilter["status"] = filter.Status
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
		return nil, nil, fmt.Errorf("failed to count tenants: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"created_at": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list tenants: %w", err)
	}
	defer cursor.Close(ctx)

	var tenants []*models.CustomerTenant
	if err := cursor.All(ctx, &tenants); err != nil {
		return nil, nil, fmt.Errorf("failed to decode tenants: %w", err)
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

	return tenants, paginationInfo, nil
}

// Update updates an existing tenant
func (r *TenantRepository) Update(ctx context.Context, id primitive.ObjectID, tenant *models.CustomerTenant) error {
	tenant.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": tenant}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("tenant with tenant_id '%s' already exists", tenant.TenantID)
		}
		return fmt.Errorf("failed to update tenant: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

// Delete deletes a tenant by ID
func (r *TenantRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("tenant not found")
	}

	return nil
}

// CountByCustomerID returns the count of tenants for a customer matching the filter
func (r *TenantRepository) CountByCustomerID(ctx context.Context, customerID primitive.ObjectID, filter *TenantFilter) (int64, error) {
	bsonFilter := bson.M{"customer_id": customerID}

	if filter != nil && filter.Status != "" {
		bsonFilter["status"] = filter.Status
	}

	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count tenants: %w", err)
	}
	return count, nil
}
