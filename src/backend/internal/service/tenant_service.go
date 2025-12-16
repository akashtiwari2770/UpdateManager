package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// TenantService handles tenant business logic
type TenantService struct {
	tenantRepo     *repository.TenantRepository
	customerRepo   *repository.CustomerRepository
	deploymentRepo *repository.DeploymentRepository
	auditRepo      *repository.AuditLogRepository
}

// NewTenantService creates a new tenant service
func NewTenantService(
	tenantRepo *repository.TenantRepository,
	customerRepo *repository.CustomerRepository,
	deploymentRepo *repository.DeploymentRepository,
	auditRepo *repository.AuditLogRepository,
) *TenantService {
	return &TenantService{
		tenantRepo:     tenantRepo,
		customerRepo:   customerRepo,
		deploymentRepo: deploymentRepo,
		auditRepo:      auditRepo,
	}
}

// ListTenantsQuery represents query parameters for listing tenants
type ListTenantsQuery struct {
	Status models.TenantStatus
	Page   int
	Limit  int
}

// TenantListResponse represents a paginated list of tenants
type TenantListResponse struct {
	Tenants    []*models.CustomerTenant
	Pagination *repository.PaginationInfo
}

// ListDeploymentsQuery represents query parameters for listing deployments
type ListDeploymentsQuery struct {
	ProductID      string
	DeploymentType models.DeploymentType
	Status         models.DeploymentStatus
	Version        string
	Page           int
	Limit          int
}

// DeploymentListResponse represents a paginated list of deployments
type DeploymentListResponse struct {
	Deployments []*models.Deployment
	Pagination  *repository.PaginationInfo
}

// CreateTenant creates a new tenant with validation and audit logging
func (s *TenantService) CreateTenant(ctx context.Context, customerID string, req *models.CreateTenantRequest, userID, userEmail string) (*models.CustomerTenant, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		// Try as ObjectID
		objectID, parseErr := primitive.ObjectIDFromHex(customerID)
		if parseErr != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
		customer, err = s.customerRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
	}

	// Validate tenant_id uniqueness
	if req.TenantID != "" {
		existing, err := s.tenantRepo.GetByTenantID(ctx, req.TenantID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("tenant with tenant_id '%s' already exists", req.TenantID)
		}
	} else {
		// Generate unique tenant_id if not provided
		req.TenantID = s.generateTenantID()
	}

	// Create tenant
	tenant := &models.CustomerTenant{
		TenantID:    req.TenantID,
		CustomerID:  customer.ID,
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
	}

	if err := s.tenantRepo.Create(ctx, tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "tenant", tenant.ID.Hex(), userID, userEmail, map[string]interface{}{
		"tenant_id":  tenant.TenantID,
		"customer_id": customer.CustomerID,
		"name":        tenant.Name,
	})

	return tenant, nil
}

// GetTenant retrieves a tenant by ID (string or ObjectID)
func (s *TenantService) GetTenant(ctx context.Context, id string) (*models.CustomerTenant, error) {
	// Try to parse as ObjectID first
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		tenant, err := s.tenantRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		return tenant, nil
	}

	// If not ObjectID, try as tenant_id
	tenant, err := s.tenantRepo.GetByTenantID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}
	return tenant, nil
}

// ListTenants lists tenants for a customer with filters and pagination
func (s *TenantService) ListTenants(ctx context.Context, customerID string, query *ListTenantsQuery) (*TenantListResponse, error) {
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		// Try as ObjectID
		objectID, parseErr := primitive.ObjectIDFromHex(customerID)
		if parseErr != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
		customer, err = s.customerRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
	}

	filter := &repository.TenantFilter{}
	if query != nil {
		filter.Status = query.Status
	}

	pagination := &repository.Pagination{}
	if query != nil {
		pagination.Page = query.Page
		pagination.Limit = query.Limit
	}

	tenants, paginationInfo, err := s.tenantRepo.GetByCustomerID(ctx, customer.ID, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list tenants: %w", err)
	}

	return &TenantListResponse{
		Tenants:    tenants,
		Pagination: paginationInfo,
	}, nil
}

// UpdateTenant updates an existing tenant
func (s *TenantService) UpdateTenant(ctx context.Context, id string, req *models.UpdateTenantRequest, userID, userEmail string) (*models.CustomerTenant, error) {
	// Get existing tenant
	tenant, err := s.GetTenant(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		tenant.Name = *req.Name
	}
	if req.Description != nil {
		tenant.Description = *req.Description
	}
	if req.Status != nil {
		tenant.Status = *req.Status
	}

	if err := s.tenantRepo.Update(ctx, tenant.ID, tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "tenant", tenant.ID.Hex(), userID, userEmail, map[string]interface{}{
		"tenant_id": tenant.TenantID,
	})

	return tenant, nil
}

// DeleteTenant deletes a tenant
func (s *TenantService) DeleteTenant(ctx context.Context, id string, userID, userEmail string) error {
	tenant, err := s.GetTenant(ctx, id)
	if err != nil {
		return err
	}

	// Check for existing deployments
	deployments, _, err := s.deploymentRepo.GetByTenantID(ctx, tenant.ID, nil, &repository.Pagination{Page: 1, Limit: 1})
	if err == nil && len(deployments) > 0 {
		return fmt.Errorf("cannot delete tenant: tenant has existing deployments")
	}

	// Delete tenant
	if err := s.tenantRepo.Delete(ctx, tenant.ID); err != nil {
		return fmt.Errorf("failed to delete tenant: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "tenant", tenant.ID.Hex(), userID, userEmail, map[string]interface{}{
		"tenant_id": tenant.TenantID,
	})

	return nil
}

// GetTenantDeployments retrieves deployments for a tenant
func (s *TenantService) GetTenantDeployments(ctx context.Context, tenantID string, query *ListDeploymentsQuery) (*DeploymentListResponse, error) {
	tenant, err := s.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	filter := &repository.DeploymentFilter{}
	if query != nil {
		filter.ProductID = query.ProductID
		filter.DeploymentType = query.DeploymentType
		filter.Status = query.Status
		filter.Version = query.Version
	}

	pagination := &repository.Pagination{}
	if query != nil {
		pagination.Page = query.Page
		pagination.Limit = query.Limit
	}

	deployments, paginationInfo, err := s.deploymentRepo.GetByTenantID(ctx, tenant.ID, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get tenant deployments: %w", err)
	}

	return &DeploymentListResponse{
		Deployments: deployments,
		Pagination:  paginationInfo,
	}, nil
}

// TenantStatistics represents statistics for a tenant
type TenantStatistics struct {
	TotalDeployments int64
	TotalUsers       int64
	DeploymentsByProduct map[string]int64
	DeploymentsByType    map[string]int64
}

// GetTenantStatistics retrieves statistics for a tenant
func (s *TenantService) GetTenantStatistics(ctx context.Context, tenantID string) (*TenantStatistics, error) {
	tenant, err := s.GetTenant(ctx, tenantID)
	if err != nil {
		return nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Get all deployments for tenant
	deployments, _, err := s.deploymentRepo.GetByTenantID(ctx, tenant.ID, nil, &repository.Pagination{Page: 1, Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	stats := &TenantStatistics{
		TotalDeployments:     int64(len(deployments)),
		DeploymentsByProduct: make(map[string]int64),
		DeploymentsByType:    make(map[string]int64),
	}

	var totalUsers int64
	for _, deployment := range deployments {
		stats.DeploymentsByProduct[deployment.ProductID]++
		stats.DeploymentsByType[string(deployment.DeploymentType)]++
		if deployment.NumberOfUsers != nil {
			totalUsers += int64(*deployment.NumberOfUsers)
		}
	}

	stats.TotalUsers = totalUsers

	return stats, nil
}

// Helper methods

// generateTenantID generates a unique tenant ID
func (s *TenantService) generateTenantID() string {
	// Generate a unique ID based on timestamp
	timestamp := time.Now().Unix()
	return fmt.Sprintf("TENANT-%d", timestamp)
}

// logAudit logs an audit entry
func (s *TenantService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
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

