package repository

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var (
	testDB      *database.MongoDB
	testCtx     context.Context
	productRepo *ProductRepository
)

// setupTestDB initializes the test database connection
func setupTestDB(t *testing.T) {
	ctx := context.Background()

	cfg := &database.Config{
		URI:      "mongodb://admin:admin123@localhost:27017/updatemanager_test?authSource=admin",
		Database: "updatemanager_test",
		Timeout:  10 * time.Second,
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	testDB = db
	testCtx = ctx
	productRepo = NewProductRepository(db.Collection("products"))
}

// teardownTestDB cleans up test database
func teardownTestDB(t *testing.T) {
	if testDB != nil {
		// Drop the test collection
		_ = testDB.Collection("products").Drop(testCtx)
		_ = testDB.Disconnect(testCtx)
	}
}

// TestProductCreate tests creating a product
func TestProductCreate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	product := &models.Product{
		ProductID:   "test-product-1",
		Name:        "Test Product",
		Type:        models.ProductTypeServer,
		Description: "A test product",
		Vendor:      "Test Vendor",
		IsActive:    true,
	}

	err := productRepo.Create(testCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Verify ID was set
	if product.ID.IsZero() {
		t.Error("Product ID was not set after creation")
	}

	// Verify timestamps were set
	if product.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}
	if product.UpdatedAt.IsZero() {
		t.Error("UpdatedAt was not set")
	}

	t.Logf("Created product with ID: %s", product.ID.Hex())
}

// TestProductGetByID tests retrieving a product by ID
func TestProductGetByID(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create a product first
	product := &models.Product{
		ProductID: "test-product-2",
		Name:      "Test Product 2",
		Type:      models.ProductTypeClient,
		IsActive:  true,
	}

	err := productRepo.Create(testCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Retrieve it
	retrieved, err := productRepo.GetByID(testCtx, product.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	// Verify fields
	if retrieved.ProductID != product.ProductID {
		t.Errorf("ProductID mismatch: got %s, want %s", retrieved.ProductID, product.ProductID)
	}
	if retrieved.Name != product.Name {
		t.Errorf("Name mismatch: got %s, want %s", retrieved.Name, product.Name)
	}
	if retrieved.Type != product.Type {
		t.Errorf("Type mismatch: got %s, want %s", retrieved.Type, product.Type)
	}

	t.Logf("Retrieved product: %+v", retrieved)
}

// TestProductGetByProductID tests retrieving a product by product_id
func TestProductGetByProductID(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	productID := "test-product-3"
	product := &models.Product{
		ProductID: productID,
		Name:      "Test Product 3",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}

	err := productRepo.Create(testCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Retrieve by product_id
	retrieved, err := productRepo.GetByProductID(testCtx, productID)
	if err != nil {
		t.Fatalf("Failed to get product by product_id: %v", err)
	}

	if retrieved.ProductID != productID {
		t.Errorf("ProductID mismatch: got %s, want %s", retrieved.ProductID, productID)
	}

	t.Logf("Retrieved product by product_id: %+v", retrieved)
}

// TestProductUpdate tests updating a product
func TestProductUpdate(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create a product
	product := &models.Product{
		ProductID: "test-product-4",
		Name:      "Original Name",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}

	err := productRepo.Create(testCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	originalUpdatedAt := product.UpdatedAt

	// Update the product
	time.Sleep(100 * time.Millisecond) // Ensure UpdatedAt changes
	product.Name = "Updated Name"
	product.Description = "Updated Description"

	err = productRepo.Update(testCtx, product)
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	// Verify update
	retrieved, err := productRepo.GetByID(testCtx, product.ID)
	if err != nil {
		t.Fatalf("Failed to get updated product: %v", err)
	}

	if retrieved.Name != "Updated Name" {
		t.Errorf("Name not updated: got %s, want Updated Name", retrieved.Name)
	}
	if retrieved.Description != "Updated Description" {
		t.Errorf("Description not updated: got %s, want Updated Description", retrieved.Description)
	}
	if !retrieved.UpdatedAt.After(originalUpdatedAt) {
		t.Error("UpdatedAt was not updated")
	}

	t.Logf("Updated product: %+v", retrieved)
}

// TestProductDelete tests deleting a product
func TestProductDelete(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create a product
	product := &models.Product{
		ProductID: "test-product-5",
		Name:      "Product to Delete",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}

	err := productRepo.Create(testCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Delete it
	err = productRepo.Delete(testCtx, product.ID)
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	// Verify it's deleted
	_, err = productRepo.GetByID(testCtx, product.ID)
	if err == nil {
		t.Error("Product should be deleted but was found")
	}

	t.Log("Product deleted successfully")
}

// TestProductList tests listing products
func TestProductList(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create multiple products
	products := []*models.Product{
		{ProductID: "list-product-1", Name: "List Product 1", Type: models.ProductTypeServer, IsActive: true},
		{ProductID: "list-product-2", Name: "List Product 2", Type: models.ProductTypeClient, IsActive: true},
		{ProductID: "list-product-3", Name: "List Product 3", Type: models.ProductTypeServer, IsActive: false},
	}

	for _, p := range products {
		err := productRepo.Create(testCtx, p)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	// List all products
	allProducts, err := productRepo.List(testCtx, bson.M{}, nil)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if len(allProducts) < len(products) {
		t.Errorf("Expected at least %d products, got %d", len(products), len(allProducts))
	}

	// List only active products
	activeProducts, err := productRepo.List(testCtx, bson.M{"is_active": true}, nil)
	if err != nil {
		t.Fatalf("Failed to list active products: %v", err)
	}

	activeCount := 0
	for _, p := range products {
		if p.IsActive {
			activeCount++
		}
	}

	if len(activeProducts) < activeCount {
		t.Errorf("Expected at least %d active products, got %d", activeCount, len(activeProducts))
	}

	t.Logf("Listed %d total products, %d active", len(allProducts), len(activeProducts))
}

// TestProductCount tests counting products
func TestProductCount(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	// Create some products
	products := []*models.Product{
		{ProductID: "count-product-1", Name: "Count Product 1", Type: models.ProductTypeServer, IsActive: true},
		{ProductID: "count-product-2", Name: "Count Product 2", Type: models.ProductTypeServer, IsActive: true},
	}

	for _, p := range products {
		err := productRepo.Create(testCtx, p)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	// Count all products
	totalCount, err := productRepo.Count(testCtx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to count products: %v", err)
	}

	if totalCount < int64(len(products)) {
		t.Errorf("Expected at least %d products, got %d", len(products), totalCount)
	}

	// Count by type
	serverCount, err := productRepo.Count(testCtx, bson.M{"type": models.ProductTypeServer})
	if err != nil {
		t.Fatalf("Failed to count server products: %v", err)
	}

	if serverCount < int64(len(products)) {
		t.Errorf("Expected at least %d server products, got %d", len(products), serverCount)
	}

	t.Logf("Total products: %d, Server products: %d", totalCount, serverCount)
}

// TestProductNotFound tests error handling for non-existent products
func TestProductNotFound(t *testing.T) {
	setupTestDB(t)
	defer teardownTestDB(t)

	fakeID := primitive.NewObjectID()

	// Try to get non-existent product
	_, err := productRepo.GetByID(testCtx, fakeID)
	if err == nil {
		t.Error("Expected error for non-existent product")
	}

	// Try to get by non-existent product_id
	_, err = productRepo.GetByProductID(testCtx, "non-existent-id")
	if err == nil {
		t.Error("Expected error for non-existent product_id")
	}

	t.Log("Error handling works correctly for non-existent products")
}
