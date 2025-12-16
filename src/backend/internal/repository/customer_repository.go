package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
)

// CustomerRepository handles customer database operations
type CustomerRepository struct {
	collection *mongo.Collection
}

// NewCustomerRepository creates a new customer repository
func NewCustomerRepository(collection *mongo.Collection) *CustomerRepository {
	return &CustomerRepository{
		collection: collection,
	}
}

// CustomerFilter represents filters for customer queries
type CustomerFilter struct {
	Search string
	Status models.CustomerStatus
	Email  string
}

// Pagination represents pagination options
type Pagination struct {
	Page  int
	Limit int
}

// PaginationInfo represents pagination metadata
type PaginationInfo struct {
	Page       int
	Limit      int
	Total      int64
	TotalPages int64
}

// Create creates a new customer in the database
func (r *CustomerRepository) Create(ctx context.Context, customer *models.Customer) error {
	// Set timestamps
	now := time.Now()
	customer.CreatedAt = now
	customer.UpdatedAt = now

	// Insert customer
	result, err := r.collection.InsertOne(ctx, customer)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("customer with customer_id '%s' already exists", customer.CustomerID)
		}
		return fmt.Errorf("failed to create customer: %w", err)
	}

	// Set the generated ID
	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		customer.ID = oid
	}

	return nil
}

// GetByID retrieves a customer by its ID
func (r *CustomerRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Customer, error) {
	var customer models.Customer
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

// GetByCustomerID retrieves a customer by its customer_id
func (r *CustomerRepository) GetByCustomerID(ctx context.Context, customerID string) (*models.Customer, error) {
	var customer models.Customer
	err := r.collection.FindOne(ctx, bson.M{"customer_id": customerID}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

// GetByEmail retrieves a customer by email
func (r *CustomerRepository) GetByEmail(ctx context.Context, email string) (*models.Customer, error) {
	var customer models.Customer
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&customer)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("customer not found")
		}
		return nil, fmt.Errorf("failed to get customer: %w", err)
	}
	return &customer, nil
}

// List retrieves customers with filters and pagination
func (r *CustomerRepository) List(ctx context.Context, filter *CustomerFilter, pagination *Pagination) ([]*models.Customer, *PaginationInfo, error) {
	// Build MongoDB filter
	bsonFilter := bson.M{}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["account_status"] = filter.Status
		}
		if filter.Email != "" {
			bsonFilter["email"] = filter.Email
		}
		if filter.Search != "" {
			searchRegex := bson.M{"$regex": filter.Search, "$options": "i"}
			bsonFilter["$or"] = []bson.M{
				{"name": searchRegex},
				{"organization_name": searchRegex},
				{"email": searchRegex},
			}
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
		return nil, nil, fmt.Errorf("failed to count customers: %w", err)
	}

	// Find options
	findOptions := options.Find()
	findOptions.SetSkip(int64(skip))
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.M{"created_at": -1})

	// Execute query
	cursor, err := r.collection.Find(ctx, bsonFilter, findOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list customers: %w", err)
	}
	defer cursor.Close(ctx)

	var customers []*models.Customer
	if err := cursor.All(ctx, &customers); err != nil {
		return nil, nil, fmt.Errorf("failed to decode customers: %w", err)
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

	return customers, paginationInfo, nil
}

// Update updates an existing customer
func (r *CustomerRepository) Update(ctx context.Context, id primitive.ObjectID, customer *models.Customer) error {
	customer.UpdatedAt = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": customer}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("customer with customer_id '%s' already exists", customer.CustomerID)
		}
		return fmt.Errorf("failed to update customer: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Delete deletes a customer by ID
func (r *CustomerRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("customer not found")
	}

	return nil
}

// Count returns the count of customers matching the filter
func (r *CustomerRepository) Count(ctx context.Context, filter *CustomerFilter) (int64, error) {
	bsonFilter := bson.M{}

	if filter != nil {
		if filter.Status != "" {
			bsonFilter["account_status"] = filter.Status
		}
		if filter.Email != "" {
			bsonFilter["email"] = filter.Email
		}
		if filter.Search != "" {
			searchRegex := bson.M{"$regex": filter.Search, "$options": "i"}
			bsonFilter["$or"] = []bson.M{
				{"name": searchRegex},
				{"organization_name": searchRegex},
				{"email": searchRegex},
			}
		}
	}

	count, err := r.collection.CountDocuments(ctx, bsonFilter)
	if err != nil {
		return 0, fmt.Errorf("failed to count customers: %w", err)
	}
	return count, nil
}

// Helper function to check if string is empty (used for validation)
func isEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}
