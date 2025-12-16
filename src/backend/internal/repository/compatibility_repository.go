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

// CompatibilityRepository handles compatibility matrix database operations
type CompatibilityRepository struct {
	collection *mongo.Collection
}

// NewCompatibilityRepository creates a new compatibility repository
func NewCompatibilityRepository(collection *mongo.Collection) *CompatibilityRepository {
	return &CompatibilityRepository{
		collection: collection,
	}
}

// Create creates a new compatibility matrix in the database
func (r *CompatibilityRepository) Create(ctx context.Context, matrix *models.CompatibilityMatrix) error {
	// Set timestamp
	matrix.ValidatedAt = time.Now()

	// Insert compatibility matrix
	result, err := r.collection.InsertOne(ctx, matrix)
	if err != nil {
		return fmt.Errorf("failed to create compatibility matrix: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		matrix.ID = oid
	}

	return nil
}

// GetByID retrieves a compatibility matrix by its ID
func (r *CompatibilityRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.CompatibilityMatrix, error) {
	var matrix models.CompatibilityMatrix
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&matrix)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("compatibility matrix not found")
		}
		return nil, fmt.Errorf("failed to get compatibility matrix: %w", err)
	}
	return &matrix, nil
}

// GetByProductIDAndVersion retrieves a compatibility matrix by product_id and version_number
func (r *CompatibilityRepository) GetByProductIDAndVersion(ctx context.Context, productID, versionNumber string) (*models.CompatibilityMatrix, error) {
	var matrix models.CompatibilityMatrix
	err := r.collection.FindOne(ctx, bson.M{
		"product_id":     productID,
		"version_number": versionNumber,
	}).Decode(&matrix)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("compatibility matrix not found")
		}
		return nil, fmt.Errorf("failed to get compatibility matrix: %w", err)
	}
	return &matrix, nil
}

// Update updates an existing compatibility matrix
func (r *CompatibilityRepository) Update(ctx context.Context, matrix *models.CompatibilityMatrix) error {
	matrix.ValidatedAt = time.Now()

	filter := bson.M{"_id": matrix.ID}
	update := bson.M{"$set": matrix}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update compatibility matrix: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("compatibility matrix not found")
	}

	return nil
}

// Delete deletes a compatibility matrix by ID
func (r *CompatibilityRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete compatibility matrix: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("compatibility matrix not found")
	}

	return nil
}

// List retrieves compatibility matrices with optional filters
func (r *CompatibilityRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.CompatibilityMatrix, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list compatibility matrices: %w", err)
	}
	defer cursor.Close(ctx)

	var matrices []*models.CompatibilityMatrix
	if err := cursor.All(ctx, &matrices); err != nil {
		return nil, fmt.Errorf("failed to decode compatibility matrices: %w", err)
	}

	return matrices, nil
}

// Count returns the count of compatibility matrices matching the filter
func (r *CompatibilityRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count compatibility matrices: %w", err)
	}
	return count, nil
}
