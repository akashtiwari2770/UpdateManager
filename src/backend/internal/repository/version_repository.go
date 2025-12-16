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

// VersionRepository handles version database operations
type VersionRepository struct {
	collection *mongo.Collection
}

// NewVersionRepository creates a new version repository
func NewVersionRepository(collection *mongo.Collection) *VersionRepository {
	return &VersionRepository{
		collection: collection,
	}
}

// Create creates a new version in the database
func (r *VersionRepository) Create(ctx context.Context, version *models.Version) error {
	// Set timestamps
	now := time.Now()
	version.CreatedAt = now
	version.UpdatedAt = now

	// Insert version
	result, err := r.collection.InsertOne(ctx, version)
	if err != nil {
		return fmt.Errorf("failed to create version: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		version.ID = oid
	}

	return nil
}

// GetByID retrieves a version by its ID
func (r *VersionRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Version, error) {
	var version models.Version
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("version not found")
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	return &version, nil
}

// GetByProductIDAndVersion retrieves a version by product_id and version_number
func (r *VersionRepository) GetByProductIDAndVersion(ctx context.Context, productID, versionNumber string) (*models.Version, error) {
	var version models.Version
	err := r.collection.FindOne(ctx, bson.M{
		"product_id":     productID,
		"version_number": versionNumber,
	}).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("version not found")
		}
		return nil, fmt.Errorf("failed to get version: %w", err)
	}
	return &version, nil
}

// GetByProductID retrieves all versions for a product
func (r *VersionRepository) GetByProductID(ctx context.Context, productID string, opts *options.FindOptions) ([]*models.Version, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"product_id": productID}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}
	defer cursor.Close(ctx)

	var versions []*models.Version
	if err := cursor.All(ctx, &versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	return versions, nil
}

// GetByState retrieves versions by state
func (r *VersionRepository) GetByState(ctx context.Context, state models.VersionState, opts *options.FindOptions) ([]*models.Version, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"state": state}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions by state: %w", err)
	}
	defer cursor.Close(ctx)

	var versions []*models.Version
	if err := cursor.All(ctx, &versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	return versions, nil
}

// Update updates an existing version
func (r *VersionRepository) Update(ctx context.Context, version *models.Version) error {
	version.UpdatedAt = time.Now()

	filter := bson.M{"_id": version.ID}
	update := bson.M{"$set": version}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to update version: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("version not found")
	}

	return nil
}

// UpdateState updates the state of a version
func (r *VersionRepository) UpdateState(ctx context.Context, id primitive.ObjectID, state models.VersionState, approvedBy string) error {
	update := bson.M{
		"$set": bson.M{
			"state":      state,
			"updated_at": time.Now(),
		},
	}

	if state == models.VersionStateApproved {
		now := time.Now()
		update["$set"].(bson.M)["approved_by"] = approvedBy
		update["$set"].(bson.M)["approved_at"] = now
	}

	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": id}, update)
	if err != nil {
		return fmt.Errorf("failed to update version state: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("version not found")
	}

	return nil
}

// Delete deletes a version by ID
func (r *VersionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete version: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("version not found")
	}

	return nil
}

// List retrieves versions with optional filters
func (r *VersionRepository) List(ctx context.Context, filter bson.M, opts *options.FindOptions) ([]*models.Version, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list versions: %w", err)
	}
	defer cursor.Close(ctx)

	var versions []*models.Version
	if err := cursor.All(ctx, &versions); err != nil {
		return nil, fmt.Errorf("failed to decode versions: %w", err)
	}

	return versions, nil
}

// Count returns the count of versions matching the filter
func (r *VersionRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count versions: %w", err)
	}
	return count, nil
}

// AddPackage adds a package to a version's packages array
func (r *VersionRepository) AddPackage(ctx context.Context, versionID primitive.ObjectID, packageInfo models.PackageInfo) error {
	// First, ensure packages field exists and is an array
	// If packages is null or doesn't exist, initialize it as an empty array
	filter := bson.M{"_id": versionID}
	
	// Check if packages exists and is null, then initialize it
	var version models.Version
	err := r.collection.FindOne(ctx, filter).Decode(&version)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return fmt.Errorf("version not found")
		}
		return fmt.Errorf("failed to get version: %w", err)
	}

	// If packages is nil, initialize it first
	if version.Packages == nil {
		initResult, err := r.collection.UpdateOne(ctx, filter, bson.M{
			"$set": bson.M{
				"packages": []models.PackageInfo{},
				"updated_at": time.Now(),
			},
		})
		if err != nil {
			return fmt.Errorf("failed to initialize packages: %w", err)
		}
		if initResult.MatchedCount == 0 {
			return fmt.Errorf("version not found")
		}
	}

	// Now push the package
	update := bson.M{
		"$push": bson.M{
			"packages": packageInfo,
		},
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("failed to add package: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("version not found")
	}

	return nil
}
