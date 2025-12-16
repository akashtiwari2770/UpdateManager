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

// DeploymentRepository handles deployment database operations
type DeploymentRepository struct {
	collection *mongo.Collection
}

// NewDeploymentRepository creates a new deployment repository
func NewDeploymentRepository(collection *mongo.Collection) *DeploymentRepository {
	return &DeploymentRepository{
		collection: collection,
	}
}

// DeploymentFilter represents filters for deployment queries
type DeploymentFilter struct {
	ProductID      string
	DeploymentType models.DeploymentType
	Status         models.DeploymentStatus
	Version        string
}

// Create creates a new deployment in the database
func (r *DeploymentRepository) Create(ctx context.Context, deployment *models.Deployment) error {
	// Set timestamps
	now := time.Now()
	deployment.DeploymentDate = now
	deployment.LastUpdatedDate = now

	// Insert deployment
	result, err := r.collection.InsertOne(ctx, deployment)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("deployment with deployment_id '%s' already exists, or duplicate product+type in tenant", deployment.DeploymentID)
		}
		return fmt.Errorf("failed to create deployment: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		deployment.ID = oid
	}

	return nil
}

// GetByID retrieves a deployment by its ID
func (r *DeploymentRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Deployment, error) {
	var deployment models.Deployment
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&deployment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("deployment not found")
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return &deployment, nil
}

// GetByDeploymentID retrieves a deployment by its deployment_id
func (r *DeploymentRepository) GetByDeploymentID(ctx context.Context, deploymentID string) (*models.Deployment, error) {
	var deployment models.Deployment
	err := r.collection.FindOne(ctx, bson.M{"deployment_id": deploymentID}).Decode(&deployment)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("deployment not found")
		}
		return nil, fmt.Errorf("failed to get deployment: %w", err)
	}
	return &deployment, nil
}

// GetByTenantID retrieves deployments for a tenant with filters and pagination
func (r *DeploymentRepository) GetByTenantID(ctx context.Context, tenantID primitive.ObjectID, filter *DeploymentFilter, pagination *Pagination) ([]*models.Deployment, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"tenant_id": tenantID}

	if filter != nil {
		if filter.ProductID != "" {
			bsonFilter["product_id"] = filter.ProductID
		}
		if filter.DeploymentType != "" {
			bsonFilter["deployment_type"] = filter.DeploymentType
		}
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
		if filter.Version != "" {
			bsonFilter["installed_version"] = filter.Version
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
		return nil, nil, fmt.Errorf("failed to count deployments: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"deployment_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer cursor.Close(ctx)

	var deployments []*models.Deployment
	if err := cursor.All(ctx, &deployments); err != nil {
		return nil, nil, fmt.Errorf("failed to decode deployments: %w", err)
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

	return deployments, paginationInfo, nil
}

// GetByProductID retrieves deployments for a product with filters and pagination
func (r *DeploymentRepository) GetByProductID(ctx context.Context, productID string, filter *DeploymentFilter, pagination *Pagination) ([]*models.Deployment, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{"product_id": productID}

	if filter != nil {
		if filter.DeploymentType != "" {
			bsonFilter["deployment_type"] = filter.DeploymentType
		}
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
		if filter.Version != "" {
			bsonFilter["installed_version"] = filter.Version
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
		return nil, nil, fmt.Errorf("failed to count deployments: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"deployment_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer cursor.Close(ctx)

	var deployments []*models.Deployment
	if err := cursor.All(ctx, &deployments); err != nil {
		return nil, nil, fmt.Errorf("failed to decode deployments: %w", err)
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

	return deployments, paginationInfo, nil
}

// Update updates an existing deployment
func (r *DeploymentRepository) Update(ctx context.Context, id primitive.ObjectID, deployment *models.Deployment) error {
	deployment.LastUpdatedDate = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": deployment}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("duplicate deployment (product+type) in tenant")
		}
		return fmt.Errorf("failed to update deployment: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("deployment not found")
	}

	return nil
}

// Delete deletes a deployment by ID
func (r *DeploymentRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("deployment not found")
	}

	return nil
}

// CountByTenantID returns the count of deployments for a tenant matching the filter
func (r *DeploymentRepository) CountByTenantID(ctx context.Context, tenantID primitive.ObjectID, filter *DeploymentFilter) (int64, error) {
	bsonFilter := bson.M{"tenant_id": tenantID}

	if filter != nil {
		if filter.ProductID != "" {
			bsonFilter["product_id"] = filter.ProductID
		}
		if filter.DeploymentType != "" {
			bsonFilter["deployment_type"] = filter.DeploymentType
		}
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
		if filter.Version != "" {
			bsonFilter["installed_version"] = filter.Version
		}
	}

	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count deployments: %w", err)
	}
	return count, nil
}

// GetAll retrieves all deployments with filters and pagination
func (r *DeploymentRepository) GetAll(ctx context.Context, filter *DeploymentFilter, pagination *Pagination) ([]*models.Deployment, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{}

	if filter != nil {
		if filter.ProductID != "" {
			bsonFilter["product_id"] = filter.ProductID
		}
		if filter.DeploymentType != "" {
			bsonFilter["deployment_type"] = filter.DeploymentType
		}
		if filter.Status != "" {
			bsonFilter["status"] = filter.Status
		}
		if filter.Version != "" {
			bsonFilter["installed_version"] = filter.Version
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
		return nil, nil, fmt.Errorf("failed to count deployments: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"deployment_date": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list deployments: %w", err)
	}
	defer cursor.Close(ctx)

	var deployments []*models.Deployment
	if err := cursor.All(ctx, &deployments); err != nil {
		return nil, nil, fmt.Errorf("failed to decode deployments: %w", err)
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

	return deployments, paginationInfo, nil
}

// GetDeploymentsForNotification retrieves all active deployments for a product that need notifications
func (r *DeploymentRepository) GetDeploymentsForNotification(ctx context.Context, productID string) ([]*models.Deployment, error) {
	filter := bson.M{
		"product_id": productID,
		"status":     models.DeploymentStatusActive,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("failed to find deployments: %w", err)
	}
	defer cursor.Close(ctx)

	var deployments []*models.Deployment
	if err := cursor.All(ctx, &deployments); err != nil {
		return nil, fmt.Errorf("failed to decode deployments: %w", err)
	}

	return deployments, nil
}
