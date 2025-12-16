package service

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// LicenseAllocationService handles license allocation business logic
type LicenseAllocationService struct {
	allocationRepo   *repository.LicenseAllocationRepository
	licenseRepo      *repository.LicenseRepository
	subscriptionRepo *repository.SubscriptionRepository
	customerRepo     *repository.CustomerRepository
	tenantRepo       *repository.TenantRepository
	deploymentRepo   *repository.DeploymentRepository
	auditRepo        *repository.AuditLogRepository
}

// NewLicenseAllocationService creates a new license allocation service
func NewLicenseAllocationService(
	allocationRepo *repository.LicenseAllocationRepository,
	licenseRepo *repository.LicenseRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	customerRepo *repository.CustomerRepository,
	tenantRepo *repository.TenantRepository,
	deploymentRepo *repository.DeploymentRepository,
	auditRepo *repository.AuditLogRepository,
) *LicenseAllocationService {
	return &LicenseAllocationService{
		allocationRepo:   allocationRepo,
		licenseRepo:      licenseRepo,
		subscriptionRepo: subscriptionRepo,
		customerRepo:     customerRepo,
		tenantRepo:       tenantRepo,
		deploymentRepo:   deploymentRepo,
		auditRepo:        auditRepo,
	}
}

// AllocateLicense allocates a license to a tenant or deployment
func (s *LicenseAllocationService) AllocateLicense(ctx context.Context, customerID, subscriptionID, licenseID string, req *models.AllocateLicenseRequest, userID string) (*models.LicenseAllocation, error) {
	// Validate license exists and belongs to subscription
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}

	// Validate subscription exists and belongs to customer
	subscription, err := s.subscriptionRepo.GetBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// Verify subscription belongs to customer
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	if subscription.CustomerID != customer.ID {
		return nil, fmt.Errorf("subscription does not belong to customer")
	}

	// Verify license belongs to subscription
	if license.SubscriptionID != subscription.ID {
		return nil, fmt.Errorf("license does not belong to subscription")
	}

	// Validate license is active
	if license.Status != models.LicenseStatusActive {
		return nil, fmt.Errorf("license is not active")
	}

	// Validate license expiration (for time-based licenses)
	if license.LicenseType == models.LicenseTypeTimeBased && license.EndDate != nil {
		if license.EndDate.Before(time.Now()) {
			return nil, fmt.Errorf("license has expired")
		}
	}

	// Validate tenant or deployment is provided
	if req.TenantID == nil && req.DeploymentID == nil {
		return nil, fmt.Errorf("either tenant_id or deployment_id must be provided")
	}

	// Validate tenant exists and belongs to customer (if provided)
	var tenantIDObj *primitive.ObjectID
	if req.TenantID != nil {
		tenant, err := s.tenantRepo.GetByTenantID(ctx, *req.TenantID)
		if err != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		if tenant.CustomerID != customer.ID {
			return nil, fmt.Errorf("tenant does not belong to customer")
		}
		tenantIDObj = &tenant.ID
	}

	// Validate deployment exists and belongs to tenant (if provided)
	var deploymentIDObj *primitive.ObjectID
	if req.DeploymentID != nil {
		if req.TenantID == nil {
			return nil, fmt.Errorf("tenant_id must be provided when deployment_id is provided")
		}
		deployment, err := s.deploymentRepo.GetByDeploymentID(ctx, *req.DeploymentID)
		if err != nil {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
		if deployment.TenantID != *tenantIDObj {
			return nil, fmt.Errorf("deployment does not belong to tenant")
		}

		// Validate product match
		if deployment.ProductID != license.ProductID {
			return nil, fmt.Errorf("deployment product does not match license product")
		}

		deploymentIDObj = &deployment.ID
	}

	// Check available seats
	availableSeats, err := s.GetAvailableSeats(ctx, licenseID)
	if err != nil {
		return nil, fmt.Errorf("failed to check available seats: %w", err)
	}

	if req.NumberOfSeatsAllocated > availableSeats {
		return nil, fmt.Errorf("insufficient available seats: requested %d, available %d", req.NumberOfSeatsAllocated, availableSeats)
	}

	// Generate allocation ID
	allocationID := fmt.Sprintf("ALLOC-%d", time.Now().Unix())

	// Create allocation
	allocation := &models.LicenseAllocation{
		AllocationID:          allocationID,
		LicenseID:             license.ID,
		TenantID:              tenantIDObj,
		DeploymentID:          deploymentIDObj,
		NumberOfSeatsAllocated: req.NumberOfSeatsAllocated,
		AllocationDate:        time.Now(),
		AllocatedBy:           userID,
		Status:                models.AllocationStatusActive,
		Notes:                 req.Notes,
	}

	if err := s.allocationRepo.Create(ctx, allocation); err != nil {
		return nil, fmt.Errorf("failed to allocate license: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "license_allocation", allocation.ID.Hex(), userID, "", map[string]interface{}{
		"allocation_id": allocation.AllocationID,
		"license_id":    licenseID,
		"seats":          req.NumberOfSeatsAllocated,
	})

	return allocation, nil
}

// ReleaseAllocation releases a license allocation
func (s *LicenseAllocationService) ReleaseAllocation(ctx context.Context, customerID, subscriptionID, licenseID, allocationID string, userID string) error {
	// Validate license exists and belongs to subscription
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("license not found: %w", err)
	}

	// Validate subscription exists and belongs to customer
	subscription, err := s.subscriptionRepo.GetBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	// Verify subscription belongs to customer
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("customer not found: %w", err)
	}
	if subscription.CustomerID != customer.ID {
		return fmt.Errorf("subscription does not belong to customer")
	}

	// Verify license belongs to subscription
	if license.SubscriptionID != subscription.ID {
		return fmt.Errorf("license does not belong to subscription")
	}

	// Get allocation
	allocation, err := s.allocationRepo.GetByAllocationID(ctx, allocationID)
	if err != nil {
		return fmt.Errorf("allocation not found: %w", err)
	}

	// Verify allocation belongs to license
	if allocation.LicenseID != license.ID {
		return fmt.Errorf("allocation does not belong to license")
	}

	// Check if already released
	if allocation.Status == models.AllocationStatusReleased {
		return fmt.Errorf("allocation is already released")
	}

	// Release allocation
	if err := s.allocationRepo.Release(ctx, allocation.ID, userID); err != nil {
		return fmt.Errorf("failed to release allocation: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "license_allocation", allocation.ID.Hex(), userID, "", map[string]interface{}{
		"allocation_id": allocationID,
		"action":        "release",
	})

	return nil
}

// GetAllocations lists allocations for a license with filters and pagination
func (s *LicenseAllocationService) GetAllocations(ctx context.Context, customerID, subscriptionID, licenseID string, filter *repository.LicenseAllocationFilter, pagination *repository.Pagination) ([]*models.LicenseAllocation, *repository.PaginationInfo, error) {
	// Validate license exists and belongs to subscription
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return nil, nil, fmt.Errorf("license not found: %w", err)
	}

	// Validate subscription exists and belongs to customer
	subscription, err := s.subscriptionRepo.GetBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, nil, fmt.Errorf("subscription not found: %w", err)
	}

	// Verify subscription belongs to customer
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, nil, fmt.Errorf("customer not found: %w", err)
	}
	if subscription.CustomerID != customer.ID {
		return nil, nil, fmt.Errorf("subscription does not belong to customer")
	}

	// Verify license belongs to subscription
	if license.SubscriptionID != subscription.ID {
		return nil, nil, fmt.Errorf("license does not belong to subscription")
	}

	// Get allocations
	allocations, paginationInfo, err := s.allocationRepo.GetByLicenseID(ctx, license.ID, filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list allocations: %w", err)
	}

	return allocations, paginationInfo, nil
}

// GetAllocationsByTenant lists allocations for a tenant
func (s *LicenseAllocationService) GetAllocationsByTenant(ctx context.Context, customerID, tenantID string, filter *repository.LicenseAllocationFilter, pagination *repository.Pagination) ([]*models.LicenseAllocation, *repository.PaginationInfo, error) {
	// Validate tenant exists and belongs to customer
	tenant, err := s.tenantRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Verify tenant belongs to customer
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, nil, fmt.Errorf("customer not found: %w", err)
	}
	if tenant.CustomerID != customer.ID {
		return nil, nil, fmt.Errorf("tenant does not belong to customer")
	}

	// Get allocations
	allocations, paginationInfo, err := s.allocationRepo.GetByTenantID(ctx, tenant.ID, filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list allocations: %w", err)
	}

	return allocations, paginationInfo, nil
}

// GetAllocationsByDeployment lists allocations for a deployment
func (s *LicenseAllocationService) GetAllocationsByDeployment(ctx context.Context, customerID, tenantID, deploymentID string, filter *repository.LicenseAllocationFilter, pagination *repository.Pagination) ([]*models.LicenseAllocation, *repository.PaginationInfo, error) {
	// Validate deployment exists and belongs to tenant
	deployment, err := s.deploymentRepo.GetByDeploymentID(ctx, deploymentID)
	if err != nil {
		return nil, nil, fmt.Errorf("deployment not found: %w", err)
	}

	// Validate tenant exists and belongs to customer
	tenant, err := s.tenantRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		return nil, nil, fmt.Errorf("tenant not found: %w", err)
	}

	// Verify tenant belongs to customer
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, nil, fmt.Errorf("customer not found: %w", err)
	}
	if tenant.CustomerID != customer.ID {
		return nil, nil, fmt.Errorf("tenant does not belong to customer")
	}

	// Verify deployment belongs to tenant
	if deployment.TenantID != tenant.ID {
		return nil, nil, fmt.Errorf("deployment does not belong to tenant")
	}

	// Get allocations
	allocations, paginationInfo, err := s.allocationRepo.GetByDeploymentID(ctx, deployment.ID, filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list allocations: %w", err)
	}

	return allocations, paginationInfo, nil
}

// GetLicenseUtilization retrieves utilization metrics for a license
func (s *LicenseAllocationService) GetLicenseUtilization(ctx context.Context, licenseID string) (map[string]interface{}, error) {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}

	// Get total allocated seats
	totalAllocated, err := s.allocationRepo.GetTotalAllocatedSeats(ctx, license.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate allocated seats: %w", err)
	}

	availableSeats := license.NumberOfSeats - totalAllocated
	utilizationPercent := float64(0)
	if license.NumberOfSeats > 0 {
		utilizationPercent = (float64(totalAllocated) / float64(license.NumberOfSeats)) * 100
	}

	// Get active allocations
	allocations, err := s.allocationRepo.GetActiveAllocationsByLicenseID(ctx, license.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	utilization := map[string]interface{}{
		"total_seats":         license.NumberOfSeats,
		"allocated_seats":     totalAllocated,
		"available_seats":     availableSeats,
		"utilization_percent": utilizationPercent,
		"active_allocations":  len(allocations),
	}

	return utilization, nil
}

// ValidateAllocation validates if an allocation can be made
func (s *LicenseAllocationService) ValidateAllocation(ctx context.Context, licenseID string, seats int, productID string) error {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("license not found: %w", err)
	}

	// Validate product match
	if license.ProductID != productID {
		return fmt.Errorf("license product does not match requested product")
	}

	// Validate license is active
	if license.Status != models.LicenseStatusActive {
		return fmt.Errorf("license is not active")
	}

	// Validate license expiration (for time-based licenses)
	if license.LicenseType == models.LicenseTypeTimeBased && license.EndDate != nil {
		if license.EndDate.Before(time.Now()) {
			return fmt.Errorf("license has expired")
		}
	}

	// Check available seats
	availableSeats, err := s.GetAvailableSeats(ctx, licenseID)
	if err != nil {
		return fmt.Errorf("failed to check available seats: %w", err)
	}

	if seats > availableSeats {
		return fmt.Errorf("insufficient available seats: requested %d, available %d", seats, availableSeats)
	}

	return nil
}

// GetTotalAllocatedSeats gets the total allocated seats for a license
func (s *LicenseAllocationService) GetTotalAllocatedSeats(ctx context.Context, licenseID string) (int, error) {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return 0, fmt.Errorf("license not found: %w", err)
	}

	totalAllocated, err := s.allocationRepo.GetTotalAllocatedSeats(ctx, license.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate allocated seats: %w", err)
	}

	return totalAllocated, nil
}

// GetAvailableSeats calculates available seats for a license
func (s *LicenseAllocationService) GetAvailableSeats(ctx context.Context, licenseID string) (int, error) {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return 0, fmt.Errorf("license not found: %w", err)
	}

	totalAllocated, err := s.allocationRepo.GetTotalAllocatedSeats(ctx, license.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate allocated seats: %w", err)
	}

	return license.NumberOfSeats - totalAllocated, nil
}

// logAudit logs an audit entry
func (s *LicenseAllocationService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
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

	_ = s.auditRepo.Create(ctx, auditLog)
}

