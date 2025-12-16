package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// ProductService handles product business logic
type ProductService struct {
	productRepo *repository.ProductRepository
	auditRepo   *repository.AuditLogRepository
}

// NewProductService creates a new product service
func NewProductService(productRepo *repository.ProductRepository, auditRepo *repository.AuditLogRepository) *ProductService {
	return &ProductService{
		productRepo: productRepo,
		auditRepo:   auditRepo,
	}
}

// CreateProduct creates a new product with validation and audit logging
func (s *ProductService) CreateProduct(ctx context.Context, req *models.CreateProductRequest, userID, userEmail string) (*models.Product, error) {
	// Validate product_id uniqueness
	existing, err := s.productRepo.GetByProductID(ctx, req.ProductID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("product with product_id '%s' already exists", req.ProductID)
	}

	// Create product
	product := &models.Product{
		ProductID:   req.ProductID,
		Name:        req.Name,
		Type:        req.Type,
		Description: req.Description,
		Vendor:      req.Vendor,
		IsActive:    true,
	}

	if err := s.productRepo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to create product: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "product", product.ID.Hex(), userID, userEmail, map[string]interface{}{
		"product_id": product.ProductID,
		"name":       product.Name,
		"type":       product.Type,
	})

	return product, nil
}

// GetProduct retrieves a product by ID
func (s *ProductService) GetProduct(ctx context.Context, id primitive.ObjectID) (*models.Product, error) {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// GetProductByProductID retrieves a product by product_id
func (s *ProductService) GetProductByProductID(ctx context.Context, productID string) (*models.Product, error) {
	product, err := s.productRepo.GetByProductID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}
	return product, nil
}

// UpdateProduct updates an existing product
func (s *ProductService) UpdateProduct(ctx context.Context, id primitive.ObjectID, req *models.CreateProductRequest, userID, userEmail string) (*models.Product, error) {
	// Get existing product
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Check product_id uniqueness if changed
	if req.ProductID != product.ProductID {
		existing, err := s.productRepo.GetByProductID(ctx, req.ProductID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("product with product_id '%s' already exists", req.ProductID)
		}
	}

	// Update fields
	product.ProductID = req.ProductID
	product.Name = req.Name
	product.Type = req.Type
	product.Description = req.Description
	product.Vendor = req.Vendor

	if err := s.productRepo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("failed to update product: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "product", product.ID.Hex(), userID, userEmail, map[string]interface{}{
		"product_id": product.ProductID,
		"name":       product.Name,
	})

	return product, nil
}

// DeleteProduct deletes a product (soft delete by setting IsActive to false)
func (s *ProductService) DeleteProduct(ctx context.Context, id primitive.ObjectID, userID, userEmail string) error {
	product, err := s.productRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("product not found: %w", err)
	}

	// Soft delete
	product.IsActive = false
	if err := s.productRepo.Update(ctx, product); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "product", product.ID.Hex(), userID, userEmail, map[string]interface{}{
		"product_id": product.ProductID,
	})

	return nil
}

// ListProducts lists products with optional filters
func (s *ProductService) ListProducts(ctx context.Context, filter bson.M, page, limit int) ([]*models.Product, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"created_at": -1})

	products, err := s.productRepo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list products: %w", err)
	}

	total, err := s.productRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count products: %w", err)
	}

	return products, total, nil
}

// GetActiveProducts returns only active products
func (s *ProductService) GetActiveProducts(ctx context.Context) ([]*models.Product, error) {
	filter := bson.M{"is_active": true}
	products, err := s.productRepo.List(ctx, filter, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get active products: %w", err)
	}
	return products, nil
}

// logAudit logs an audit entry
func (s *ProductService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
	if s.auditRepo == nil {
		return
	}

	auditLog := &models.AuditLog{
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserID:       userID,
		UserEmail:    userEmail,
		Details:      details,
		Timestamp:    time.Now(),
	}

	_ = s.auditRepo.Create(ctx, auditLog) // Ignore errors for audit logging
}
