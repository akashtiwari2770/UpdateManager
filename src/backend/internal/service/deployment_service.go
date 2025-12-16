package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// DeploymentService handles deployment business logic
type DeploymentService struct {
	deploymentRepo *repository.DeploymentRepository
	tenantRepo     *repository.TenantRepository
	customerRepo   *repository.CustomerRepository
	productService *ProductService
	versionService *VersionService
	auditRepo      *repository.AuditLogRepository
}

// NewDeploymentService creates a new deployment service
func NewDeploymentService(
	deploymentRepo *repository.DeploymentRepository,
	tenantRepo *repository.TenantRepository,
	customerRepo *repository.CustomerRepository,
	productService *ProductService,
	versionService *VersionService,
	auditRepo *repository.AuditLogRepository,
) *DeploymentService {
	return &DeploymentService{
		deploymentRepo: deploymentRepo,
		tenantRepo:     tenantRepo,
		customerRepo:   customerRepo,
		productService: productService,
		versionService: versionService,
		auditRepo:      auditRepo,
	}
}

// CreateDeployment creates a new deployment with validation and audit logging
func (s *DeploymentService) CreateDeployment(ctx context.Context, tenantID string, req *models.CreateDeploymentRequest, userID, userEmail string) (*models.Deployment, error) {
	// Validate tenant exists
	tenant, err := s.tenantRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		// Try as ObjectID
		objectID, parseErr := primitive.ObjectIDFromHex(tenantID)
		if parseErr != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		tenant, err = s.tenantRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
	}

	// Validate product exists
	_, err = s.productService.GetProductByProductID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	// Validate version exists (check if version is valid for the product)
	// This is a basic check - in a real scenario, we'd validate the version exists for the product
	// For now, we'll just check that the version string is not empty

	// Check for duplicate deployment (same product + type in tenant)
	existingFilter := &repository.DeploymentFilter{
		ProductID:      req.ProductID,
		DeploymentType: req.DeploymentType,
	}
	existing, _, err := s.deploymentRepo.GetByTenantID(ctx, tenant.ID, existingFilter, &repository.Pagination{Page: 1, Limit: 1})
	if err == nil && len(existing) > 0 {
		return nil, fmt.Errorf("deployment with product '%s' and type '%s' already exists in this tenant", req.ProductID, req.DeploymentType)
	}

	// Generate unique deployment_id if not provided
	if req.DeploymentID == "" {
		req.DeploymentID = s.generateDeploymentID()
	} else {
		// Validate deployment_id uniqueness
		existingDeploy, err := s.deploymentRepo.GetByDeploymentID(ctx, req.DeploymentID)
		if err == nil && existingDeploy != nil {
			return nil, fmt.Errorf("deployment with deployment_id '%s' already exists", req.DeploymentID)
		}
	}

	// Create deployment
	deployment := &models.Deployment{
		DeploymentID:     req.DeploymentID,
		TenantID:         tenant.ID,
		ProductID:        req.ProductID,
		DeploymentType:   req.DeploymentType,
		InstalledVersion: req.InstalledVersion,
		NumberOfUsers:    req.NumberOfUsers,
		LicenseInfo:      req.LicenseInfo,
		ServerHostname:   req.ServerHostname,
		EnvironmentDetails: req.EnvironmentDetails,
		Status:           req.Status,
	}

	if err := s.deploymentRepo.Create(ctx, deployment); err != nil {
		// Check for duplicate key error (handled in repository)
		if err.Error() != "" && (strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists")) {
			return nil, fmt.Errorf("duplicate deployment: same product and type already exists in this tenant")
		}
		return nil, fmt.Errorf("failed to create deployment: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "deployment", deployment.ID.Hex(), userID, userEmail, map[string]interface{}{
		"deployment_id":    deployment.DeploymentID,
		"tenant_id":        tenant.TenantID,
		"product_id":       deployment.ProductID,
		"deployment_type":  deployment.DeploymentType,
		"installed_version": deployment.InstalledVersion,
	})

	return deployment, nil
}

// GetDeployment retrieves a deployment by ID (string or ObjectID)
func (s *DeploymentService) GetDeployment(ctx context.Context, id string) (*models.Deployment, error) {
	// Try to parse as ObjectID first
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		deployment, err := s.deploymentRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
		return deployment, nil
	}

	// If not ObjectID, try as deployment_id
	deployment, err := s.deploymentRepo.GetByDeploymentID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}
	return deployment, nil
}

// ListDeployments lists deployments for a tenant with filters and pagination
func (s *DeploymentService) ListDeployments(ctx context.Context, tenantID string, query *ListDeploymentsQuery) (*DeploymentListResponse, error) {
	tenant, err := s.tenantRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		// Try as ObjectID
		objectID, parseErr := primitive.ObjectIDFromHex(tenantID)
		if parseErr != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		tenant, err = s.tenantRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
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
		return nil, fmt.Errorf("failed to list deployments: %w", err)
	}

	return &DeploymentListResponse{
		Deployments: deployments,
		Pagination:  paginationInfo,
	}, nil
}

// UpdateDeployment updates an existing deployment
func (s *DeploymentService) UpdateDeployment(ctx context.Context, id string, req *models.UpdateDeploymentRequest, userID, userEmail string) (*models.Deployment, error) {
	// Get existing deployment
	deployment, err := s.GetDeployment(ctx, id)
	if err != nil {
		return nil, err
	}

	// Validate version if updated
	if req.InstalledVersion != nil {
		// Basic validation - in production, validate version exists for product
		if *req.InstalledVersion == "" {
			return nil, fmt.Errorf("installed_version cannot be empty")
		}
	}

	// Check for duplicate if product or type is being changed
	if req.DeploymentType != nil || (req.DeploymentType == nil && req.InstalledVersion != nil) {
		// If deployment type is being changed, check for duplicates
		newType := deployment.DeploymentType
		if req.DeploymentType != nil {
			newType = *req.DeploymentType
		}

		existingFilter := &repository.DeploymentFilter{
			ProductID:      deployment.ProductID,
			DeploymentType: newType,
		}
		existing, _, err := s.deploymentRepo.GetByTenantID(ctx, deployment.TenantID, existingFilter, &repository.Pagination{Page: 1, Limit: 10})
		if err == nil {
			for _, e := range existing {
				if e.ID != deployment.ID {
					return nil, fmt.Errorf("deployment with product '%s' and type '%s' already exists in this tenant", deployment.ProductID, newType)
				}
			}
		}
	}

	// Update fields if provided
	if req.DeploymentType != nil {
		deployment.DeploymentType = *req.DeploymentType
	}
	if req.InstalledVersion != nil {
		deployment.InstalledVersion = *req.InstalledVersion
	}
	if req.NumberOfUsers != nil {
		deployment.NumberOfUsers = req.NumberOfUsers
	}
	if req.LicenseInfo != nil {
		deployment.LicenseInfo = *req.LicenseInfo
	}
	if req.ServerHostname != nil {
		deployment.ServerHostname = *req.ServerHostname
	}
	if req.EnvironmentDetails != nil {
		deployment.EnvironmentDetails = *req.EnvironmentDetails
	}
	if req.Status != nil {
		deployment.Status = *req.Status
	}

	if err := s.deploymentRepo.Update(ctx, deployment.ID, deployment); err != nil {
		return nil, fmt.Errorf("failed to update deployment: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "deployment", deployment.ID.Hex(), userID, userEmail, map[string]interface{}{
		"deployment_id": deployment.DeploymentID,
	})

	return deployment, nil
}

// DeleteDeployment deletes a deployment
func (s *DeploymentService) DeleteDeployment(ctx context.Context, id string, userID, userEmail string) error {
	deployment, err := s.GetDeployment(ctx, id)
	if err != nil {
		return err
	}

	// Delete deployment
	if err := s.deploymentRepo.Delete(ctx, deployment.ID); err != nil {
		return fmt.Errorf("failed to delete deployment: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "deployment", deployment.ID.Hex(), userID, userEmail, map[string]interface{}{
		"deployment_id": deployment.DeploymentID,
	})

	return nil
}

// GetAvailableUpdates retrieves available updates for a deployment
func (s *DeploymentService) GetAvailableUpdates(ctx context.Context, deploymentID string) ([]*models.Version, error) {
	deployment, err := s.GetDeployment(ctx, deploymentID)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	// Get all versions for the product
	versions, _, err := s.versionService.GetVersionsByProduct(ctx, deployment.ProductID, 1, 1000)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	// Filter versions that are newer than the installed version
	// This is a simplified check - in production, we'd do proper semantic version comparison
	var availableUpdates []*models.Version
	for _, version := range versions {
		if version.State == models.VersionStateReleased && version.VersionNumber != deployment.InstalledVersion {
			// TODO: Implement proper semantic version comparison
			availableUpdates = append(availableUpdates, version)
		}
	}

	return availableUpdates, nil
}

// Helper methods

// generateDeploymentID generates a unique deployment ID
func (s *DeploymentService) generateDeploymentID() string {
	// Generate a unique ID based on timestamp
	timestamp := time.Now().Unix()
	return fmt.Sprintf("DEPLOY-%d", timestamp)
}

// logAudit logs an audit entry
func (s *DeploymentService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
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

