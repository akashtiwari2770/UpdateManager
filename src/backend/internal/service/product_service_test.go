package service

import (
	"context"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/pkg/database"
)

var (
	productServiceTestDB  *database.MongoDB
	productServiceTestCtx context.Context
	productService        *ProductService
	productRepo           *repository.ProductRepository
	auditRepo             *repository.AuditLogRepository
)

func setupProductServiceTestDB(t *testing.T) {
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

	productServiceTestDB = db
	productServiceTestCtx = ctx
	productRepo = repository.NewProductRepository(db.Collection("products"))
	auditRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	productService = NewProductService(productRepo, auditRepo)
}

func teardownProductServiceTestDB(t *testing.T) {
	if productServiceTestDB != nil {
		_ = productServiceTestDB.Collection("products").Drop(productServiceTestCtx)
		_ = productServiceTestDB.Collection("audit_logs").Drop(productServiceTestCtx)
		_ = productServiceTestDB.Disconnect(productServiceTestCtx)
	}
}

func TestProductService_CreateProduct(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	req := &models.CreateProductRequest{
		ProductID:   "test-product-svc-1",
		Name:        "Test Product Service",
		Type:        models.ProductTypeServer,
		Description: "A test product for service layer",
		Vendor:      "Test Vendor",
	}

	product, err := productService.CreateProduct(productServiceTestCtx, req, "user-123", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Verify product was created
	if product.ID.IsZero() {
		t.Error("Product ID was not set")
	}
	if product.ProductID != req.ProductID {
		t.Errorf("ProductID mismatch: got %s, want %s", product.ProductID, req.ProductID)
	}
	if product.Name != req.Name {
		t.Errorf("Name mismatch: got %s, want %s", product.Name, req.Name)
	}
	if !product.IsActive {
		t.Error("Product should be active by default")
	}

	// Verify audit log was created
	auditLogs, _ := auditRepo.GetByResource(productServiceTestCtx, "product", product.ID.Hex(), nil)
	if len(auditLogs) == 0 {
		t.Error("Audit log should be created")
	} else {
		if auditLogs[0].Action != models.AuditActionCreate {
			t.Errorf("Audit action mismatch: got %s, want %s", auditLogs[0].Action, models.AuditActionCreate)
		}
	}

	t.Logf("Created product: %+v", product)
}

func TestProductService_CreateProduct_DuplicateProductID(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	req := &models.CreateProductRequest{
		ProductID: "duplicate-product",
		Name:      "First Product",
		Type:      models.ProductTypeServer,
	}

	// Create first product
	_, err := productService.CreateProduct(productServiceTestCtx, req, "user-123", "user@example.com")
	if err != nil {
		t.Fatalf("Failed to create first product: %v", err)
	}

	// Try to create duplicate
	req2 := &models.CreateProductRequest{
		ProductID: "duplicate-product",
		Name:      "Second Product",
		Type:      models.ProductTypeClient,
	}
	_, err = productService.CreateProduct(productServiceTestCtx, req2, "user-123", "user@example.com")
	if err == nil {
		t.Error("Expected error for duplicate product_id, got nil")
	}

	t.Logf("Correctly rejected duplicate product_id: %v", err)
}

func TestProductService_GetProduct(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	// Create product directly via repository
	product := &models.Product{
		ProductID: "get-product-test",
		Name:      "Get Product Test",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	err := productRepo.Create(productServiceTestCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Get via service
	retrieved, err := productService.GetProduct(productServiceTestCtx, product.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if retrieved.ProductID != product.ProductID {
		t.Errorf("ProductID mismatch: got %s, want %s", retrieved.ProductID, product.ProductID)
	}

	t.Logf("Retrieved product: %+v", retrieved)
}

func TestProductService_GetProduct_NotFound(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	fakeID := primitive.NewObjectID()
	_, err := productService.GetProduct(productServiceTestCtx, fakeID)
	if err == nil {
		t.Error("Expected error for non-existent product, got nil")
	}

	t.Logf("Correctly returned error for non-existent product: %v", err)
}

func TestProductService_GetProductByProductID(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	productID := "get-by-product-id"
	product := &models.Product{
		ProductID: productID,
		Name:      "Get By Product ID",
		Type:      models.ProductTypeClient,
		IsActive:  true,
	}
	err := productRepo.Create(productServiceTestCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	retrieved, err := productService.GetProductByProductID(productServiceTestCtx, productID)
	if err != nil {
		t.Fatalf("Failed to get product by product_id: %v", err)
	}

	if retrieved.ProductID != productID {
		t.Errorf("ProductID mismatch: got %s, want %s", retrieved.ProductID, productID)
	}

	t.Logf("Retrieved product by product_id: %+v", retrieved)
}

func TestProductService_UpdateProduct(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	// Create product
	product := &models.Product{
		ProductID: "update-product-test",
		Name:      "Original Name",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	err := productRepo.Create(productServiceTestCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Update via service
	req := &models.CreateProductRequest{
		ProductID:   "update-product-test",
		Name:        "Updated Name",
		Type:        models.ProductTypeClient,
		Description: "Updated Description",
		Vendor:      "Updated Vendor",
	}

	updated, err := productService.UpdateProduct(productServiceTestCtx, product.ID, req, "user-456", "user2@example.com")
	if err != nil {
		t.Fatalf("Failed to update product: %v", err)
	}

	if updated.Name != "Updated Name" {
		t.Errorf("Name not updated: got %s, want Updated Name", updated.Name)
	}
	if updated.Type != models.ProductTypeClient {
		t.Errorf("Type not updated: got %s, want client", updated.Type)
	}

	// Verify audit log
	auditLogs, _ := auditRepo.GetByResource(productServiceTestCtx, "product", product.ID.Hex(), nil)
	updateFound := false
	for _, log := range auditLogs {
		if log.Action == models.AuditActionUpdate {
			updateFound = true
			break
		}
	}
	if !updateFound {
		t.Error("Update audit log should be created")
	}

	t.Logf("Updated product: %+v", updated)
}

func TestProductService_DeleteProduct(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	// Create product
	product := &models.Product{
		ProductID: "delete-product-test",
		Name:      "Product to Delete",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	err := productRepo.Create(productServiceTestCtx, product)
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Delete via service (soft delete)
	err = productService.DeleteProduct(productServiceTestCtx, product.ID, "user-789", "user3@example.com")
	if err != nil {
		t.Fatalf("Failed to delete product: %v", err)
	}

	// Verify it's soft deleted
	retrieved, err := productRepo.GetByID(productServiceTestCtx, product.ID)
	if err != nil {
		t.Fatalf("Failed to get product: %v", err)
	}

	if retrieved.IsActive {
		t.Error("Product should be inactive after deletion")
	}

	// Verify audit log
	auditLogs, _ := auditRepo.GetByResource(productServiceTestCtx, "product", product.ID.Hex(), nil)
	deleteFound := false
	for _, log := range auditLogs {
		if log.Action == models.AuditActionDelete {
			deleteFound = true
			break
		}
	}
	if !deleteFound {
		t.Error("Delete audit log should be created")
	}

	t.Logf("Soft deleted product: %+v", retrieved)
}

func TestProductService_ListProducts(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	// Create multiple products
	products := []*models.Product{
		{ProductID: "list-product-1", Name: "List Product 1", Type: models.ProductTypeServer, IsActive: true},
		{ProductID: "list-product-2", Name: "List Product 2", Type: models.ProductTypeClient, IsActive: true},
		{ProductID: "list-product-3", Name: "List Product 3", Type: models.ProductTypeServer, IsActive: false},
	}

	for _, p := range products {
		err := productRepo.Create(productServiceTestCtx, p)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	// List all products
	allProducts, total, err := productService.ListProducts(productServiceTestCtx, nil, 1, 10)
	if err != nil {
		t.Fatalf("Failed to list products: %v", err)
	}

	if len(allProducts) < len(products) {
		t.Errorf("Expected at least %d products, got %d", len(products), len(allProducts))
	}

	if total < int64(len(products)) {
		t.Errorf("Expected total at least %d, got %d", len(products), total)
	}

	// List with filter
	filter := map[string]interface{}{"type": models.ProductTypeServer}
	serverProducts, _, err := productService.ListProducts(productServiceTestCtx, filter, 1, 10)
	if err != nil {
		t.Fatalf("Failed to list server products: %v", err)
	}

	serverCount := 0
	for _, p := range products {
		if p.Type == models.ProductTypeServer {
			serverCount++
		}
	}

	if len(serverProducts) < serverCount {
		t.Errorf("Expected at least %d server products, got %d", serverCount, len(serverProducts))
	}

	t.Logf("Listed %d total products, %d server products", len(allProducts), len(serverProducts))
}

func TestProductService_GetActiveProducts(t *testing.T) {
	setupProductServiceTestDB(t)
	defer teardownProductServiceTestDB(t)

	// Create active and inactive products
	products := []*models.Product{
		{ProductID: "active-1", Name: "Active 1", Type: models.ProductTypeServer, IsActive: true},
		{ProductID: "active-2", Name: "Active 2", Type: models.ProductTypeClient, IsActive: true},
		{ProductID: "inactive-1", Name: "Inactive 1", Type: models.ProductTypeServer, IsActive: false},
	}

	for _, p := range products {
		err := productRepo.Create(productServiceTestCtx, p)
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	activeProducts, err := productService.GetActiveProducts(productServiceTestCtx)
	if err != nil {
		t.Fatalf("Failed to get active products: %v", err)
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

	// Verify all returned are active
	for _, p := range activeProducts {
		if !p.IsActive {
			t.Errorf("Product %s should be active", p.ProductID)
		}
	}

	t.Logf("Found %d active products", len(activeProducts))
}
