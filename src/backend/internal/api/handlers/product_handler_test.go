package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/api/utils"
	"updatemanager/internal/models"
	"updatemanager/internal/service"
	"updatemanager/pkg/database"
)

func setupProductHandlerTest(t *testing.T) (*ProductHandler, func()) {
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

	// Clean up before test to ensure idempotency
	_ = db.Collection("products").Drop(ctx)
	_ = db.Collection("versions").Drop(ctx)
	_ = db.Collection("compatibility_matrices").Drop(ctx)
	_ = db.Collection("upgrade_paths").Drop(ctx)
	_ = db.Collection("notifications").Drop(ctx)
	_ = db.Collection("update_detections").Drop(ctx)
	_ = db.Collection("update_rollouts").Drop(ctx)
	_ = db.Collection("audit_logs").Drop(ctx)

	services := service.NewServiceFactory(db.Database)
	handler := NewProductHandler(services.ProductService)

	cleanup := func() {
		// Drop all test collections to make tests idempotent
		_ = db.Collection("products").Drop(ctx)
		_ = db.Collection("versions").Drop(ctx)
		_ = db.Collection("compatibility_matrices").Drop(ctx)
		_ = db.Collection("upgrade_paths").Drop(ctx)
		_ = db.Collection("notifications").Drop(ctx)
		_ = db.Collection("update_detections").Drop(ctx)
		_ = db.Collection("update_rollouts").Drop(ctx)
		_ = db.Collection("audit_logs").Drop(ctx)
		db.Disconnect(ctx)
	}

	return handler, cleanup
}

func TestProductHandler_CreateProduct(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	req := models.CreateProductRequest{
		ProductID:   "test-product-1",
		Name:        "Test Product",
		Type:        models.ProductTypeServer,
		Description: "A test product",
		Vendor:      "Test Vendor",
	}

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "test-user")
	httpReq.Header.Set("X-User-Email", "test@example.com")

	w := httptest.NewRecorder()
	handler.CreateProduct(w, httpReq)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}

	productData, ok := response.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected product data, got %T", response.Data)
	}

	if productData["product_id"] != "test-product-1" {
		t.Errorf("Expected product_id=test-product-1, got %v", productData["product_id"])
	}
}

func TestProductHandler_CreateProduct_Duplicate(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create first product
	req1 := models.CreateProductRequest{
		ProductID: "duplicate-product",
		Name:      "Duplicate Product",
		Type:      models.ProductTypeServer,
	}
	_, err := handler.productService.CreateProduct(ctx, &req1, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create first product: %v", err)
	}

	// Try to create duplicate
	req2 := models.CreateProductRequest{
		ProductID: "duplicate-product",
		Name:      "Another Product",
		Type:      models.ProductTypeClient,
	}

	body, _ := json.Marshal(req2)
	httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "test-user")

	w := httptest.NewRecorder()
	handler.CreateProduct(w, httpReq)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status %d, got %d", http.StatusConflict, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Success {
		t.Errorf("Expected success=false, got true")
	}
}

func TestProductHandler_GetProduct(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a product first
	req := models.CreateProductRequest{
		ProductID: "get-product-test",
		Name:      "Get Test Product",
		Type:      models.ProductTypeServer,
	}
	product, err := handler.productService.CreateProduct(ctx, &req, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Get the product
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+product.ID.Hex(), nil)
	w := httptest.NewRecorder()
	handler.GetProduct(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if !response.Success {
		t.Errorf("Expected success=true, got false")
	}
}

func TestProductHandler_GetProduct_NotFound(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	fakeID := primitive.NewObjectID()
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/"+fakeID.Hex(), nil)
	w := httptest.NewRecorder()
	handler.GetProduct(w, httpReq)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status %d, got %d", http.StatusNotFound, w.Code)
	}
}

func TestProductHandler_GetProductByProductID(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a product
	req := models.CreateProductRequest{
		ProductID: "by-product-id-test",
		Name:      "By Product ID Test",
		Type:      models.ProductTypeServer,
	}
	_, err := handler.productService.CreateProduct(ctx, &req, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Get by product ID
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/by-product-id/by-product-id-test", nil)
	w := httptest.NewRecorder()
	handler.GetProductByProductID(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}
}

func TestProductHandler_UpdateProduct(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a product
	req := models.CreateProductRequest{
		ProductID: "update-product-test",
		Name:      "Original Name",
		Type:      models.ProductTypeServer,
	}
	product, err := handler.productService.CreateProduct(ctx, &req, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Update the product
	updateReq := models.CreateProductRequest{
		ProductID:   "update-product-test",
		Name:        "Updated Name",
		Type:        models.ProductTypeClient,
		Description: "Updated description",
	}

	body, _ := json.Marshal(updateReq)
	httpReq := httptest.NewRequest(http.MethodPut, "/api/v1/products/"+product.ID.Hex(), bytes.NewBuffer(body))
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-User-ID", "user2")

	w := httptest.NewRecorder()
	handler.UpdateProduct(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	productData := response.Data.(map[string]interface{})
	if productData["name"] != "Updated Name" {
		t.Errorf("Expected name=Updated Name, got %v", productData["name"])
	}
}

func TestProductHandler_DeleteProduct(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create a product
	req := models.CreateProductRequest{
		ProductID: "delete-product-test",
		Name:      "Delete Test Product",
		Type:      models.ProductTypeServer,
	}
	product, err := handler.productService.CreateProduct(ctx, &req, "user1", "user1@example.com")
	if err != nil {
		t.Fatalf("Failed to create product: %v", err)
	}

	// Delete the product
	httpReq := httptest.NewRequest(http.MethodDelete, "/api/v1/products/"+product.ID.Hex(), nil)
	httpReq.Header.Set("X-User-ID", "user1")
	httpReq.Header.Set("X-User-Email", "user1@example.com")

	w := httptest.NewRecorder()
	handler.DeleteProduct(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Verify product is soft deleted (IsActive = false)
	deletedProduct, err := handler.productService.GetProduct(ctx, product.ID)
	if err != nil {
		t.Fatalf("Failed to get product after deletion: %v", err)
	}
	if deletedProduct.IsActive {
		t.Error("Expected product to be inactive after soft delete, but IsActive is still true")
	}
}

func TestProductHandler_ListProducts(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create multiple products
	for i := 0; i < 5; i++ {
		req := models.CreateProductRequest{
			ProductID: "list-product-" + string(rune('a'+i)),
			Name:      "List Product " + string(rune('A'+i)),
			Type:      models.ProductTypeServer,
		}
		_, err := handler.productService.CreateProduct(ctx, &req, "user1", "user1@example.com")
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	// List products
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products?page=1&limit=10", nil)
	w := httptest.NewRecorder()
	handler.ListProducts(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Meta == nil {
		t.Fatal("Expected pagination metadata")
	}

	if response.Meta.Total < 5 {
		t.Errorf("Expected at least 5 products, got %d", response.Meta.Total)
	}
}

func TestProductHandler_GetActiveProducts(t *testing.T) {
	handler, cleanup := setupProductHandlerTest(t)
	defer cleanup()

	ctx := context.Background()

	// Create active products
	for i := 0; i < 3; i++ {
		req := models.CreateProductRequest{
			ProductID: "active-product-" + string(rune('a'+i)),
			Name:      "Active Product " + string(rune('A'+i)),
			Type:      models.ProductTypeServer,
		}
		_, err := handler.productService.CreateProduct(ctx, &req, "user1", "user1@example.com")
		if err != nil {
			t.Fatalf("Failed to create product: %v", err)
		}
	}

	// Get active products
	httpReq := httptest.NewRequest(http.MethodGet, "/api/v1/products/active", nil)
	w := httptest.NewRecorder()
	handler.GetActiveProducts(w, httpReq)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response utils.JSONResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	products, ok := response.Data.([]interface{})
	if !ok {
		t.Fatalf("Expected array of products, got %T", response.Data)
	}

	if len(products) < 3 {
		t.Errorf("Expected at least 3 active products, got %d", len(products))
	}
}
