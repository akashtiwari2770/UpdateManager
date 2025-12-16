package service

import (
	"context"
	"fmt"
	"time"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// LicenseService handles license business logic
type LicenseService struct {
	licenseRepo      *repository.LicenseRepository
	subscriptionRepo *repository.SubscriptionRepository
	customerRepo     *repository.CustomerRepository
	allocationRepo   *repository.LicenseAllocationRepository
	auditRepo        *repository.AuditLogRepository
}

// NewLicenseService creates a new license service
func NewLicenseService(
	licenseRepo *repository.LicenseRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	customerRepo *repository.CustomerRepository,
	allocationRepo *repository.LicenseAllocationRepository,
	auditRepo *repository.AuditLogRepository,
) *LicenseService {
	return &LicenseService{
		licenseRepo:      licenseRepo,
		subscriptionRepo: subscriptionRepo,
		customerRepo:     customerRepo,
		allocationRepo:   allocationRepo,
		auditRepo:        auditRepo,
	}
}

// AssignLicense assigns a license to a subscription
func (s *LicenseService) AssignLicense(ctx context.Context, customerID, subscriptionID string, req *models.CreateLicenseRequest, userID string) (*models.License, error) {
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

	// Validate license_id uniqueness
	existing, err := s.licenseRepo.GetByLicenseID(ctx, req.LicenseID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("license with license_id '%s' already exists", req.LicenseID)
	}

	// Validate time-based license has end date
	if req.LicenseType == models.LicenseTypeTimeBased && req.EndDate == nil {
		return nil, fmt.Errorf("time-based license must have an end date")
	}

	// Validate dates
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Create license
	license := &models.License{
		LicenseID:      req.LicenseID,
		SubscriptionID: subscription.ID,
		ProductID:      req.ProductID,
		LicenseType:    req.LicenseType,
		NumberOfSeats:  req.NumberOfSeats,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         req.Status,
		AssignedBy:     userID,
		AssignmentDate: time.Now(),
		Notes:          req.Notes,
	}

	if err := s.licenseRepo.Create(ctx, license); err != nil {
		return nil, fmt.Errorf("failed to assign license: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "license", license.ID.Hex(), userID, "", map[string]interface{}{
		"license_id":     license.LicenseID,
		"subscription_id": subscriptionID,
		"product_id":     license.ProductID,
		"license_type":   license.LicenseType,
		"seats":          license.NumberOfSeats,
	})

	return license, nil
}

// GetLicense retrieves a license by customer ID, subscription ID, and license ID
func (s *LicenseService) GetLicense(ctx context.Context, customerID, subscriptionID, licenseID string) (*models.License, error) {
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

	// Get license
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return nil, fmt.Errorf("license not found: %w", err)
	}

	// Verify license belongs to subscription
	if license.SubscriptionID != subscription.ID {
		return nil, fmt.Errorf("license does not belong to subscription")
	}

	return license, nil
}

// ListLicenses lists licenses for a subscription with filters and pagination
func (s *LicenseService) ListLicenses(ctx context.Context, customerID, subscriptionID string, filter *repository.LicenseFilter, pagination *repository.Pagination) ([]*models.License, *repository.PaginationInfo, error) {
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

	// Get licenses
	licenses, paginationInfo, err := s.licenseRepo.GetBySubscriptionID(ctx, subscription.ID, filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list licenses: %w", err)
	}

	return licenses, paginationInfo, nil
}

// UpdateLicense updates an existing license
func (s *LicenseService) UpdateLicense(ctx context.Context, customerID, subscriptionID, licenseID string, req *models.UpdateLicenseRequest, userID string) (*models.License, error) {
	// Get license
	license, err := s.GetLicense(ctx, customerID, subscriptionID, licenseID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.LicenseType != nil {
		license.LicenseType = *req.LicenseType
	}
	if req.NumberOfSeats != nil {
		license.NumberOfSeats = *req.NumberOfSeats
	}
	if req.StartDate != nil {
		license.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		license.EndDate = req.EndDate
	}
	if req.Status != nil {
		license.Status = *req.Status
	}
	if req.Notes != nil {
		license.Notes = *req.Notes
	}

	// Validate time-based license has end date
	if license.LicenseType == models.LicenseTypeTimeBased && license.EndDate == nil {
		return nil, fmt.Errorf("time-based license must have an end date")
	}

	// Validate dates
	if license.EndDate != nil && license.EndDate.Before(license.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Update license
	if err := s.licenseRepo.Update(ctx, license.ID, license); err != nil {
		return nil, fmt.Errorf("failed to update license: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "license", license.ID.Hex(), userID, "", map[string]interface{}{
		"license_id": license.LicenseID,
	})

	return license, nil
}

// RevokeLicense revokes a license
func (s *LicenseService) RevokeLicense(ctx context.Context, customerID, subscriptionID, licenseID string, userID string) error {
	// Get license
	license, err := s.GetLicense(ctx, customerID, subscriptionID, licenseID)
	if err != nil {
		return err
	}

	// Check if license has active allocations
	allocations, err := s.allocationRepo.GetActiveAllocationsByLicenseID(ctx, license.ID)
	if err == nil && len(allocations) > 0 {
		return fmt.Errorf("cannot revoke license with active allocations")
	}

	// Update license status
	license.Status = models.LicenseStatusRevoked
	if err := s.licenseRepo.Update(ctx, license.ID, license); err != nil {
		return fmt.Errorf("failed to revoke license: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "license", license.ID.Hex(), userID, "", map[string]interface{}{
		"license_id": license.LicenseID,
		"action":     "revoke",
	})

	return nil
}

// GetLicenseStatistics retrieves statistics for a license
func (s *LicenseService) GetLicenseStatistics(ctx context.Context, customerID, subscriptionID, licenseID string) (map[string]interface{}, error) {
	// Get license
	license, err := s.GetLicense(ctx, customerID, subscriptionID, licenseID)
	if err != nil {
		return nil, err
	}

	// Get allocations
	allocations, err := s.allocationRepo.GetActiveAllocationsByLicenseID(ctx, license.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get allocations: %w", err)
	}

	// Calculate total allocated seats
	totalAllocated, err := s.allocationRepo.GetTotalAllocatedSeats(ctx, license.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate allocated seats: %w", err)
	}

	availableSeats := license.NumberOfSeats - totalAllocated
	utilizationPercent := float64(0)
	if license.NumberOfSeats > 0 {
		utilizationPercent = (float64(totalAllocated) / float64(license.NumberOfSeats)) * 100
	}

	stats := map[string]interface{}{
		"total_seats":         license.NumberOfSeats,
		"allocated_seats":     totalAllocated,
		"available_seats":     availableSeats,
		"utilization_percent": utilizationPercent,
		"active_allocations":  len(allocations),
		"license_type":        license.LicenseType,
		"status":              license.Status,
	}

	return stats, nil
}

// ValidateLicenseStatus validates and updates license status based on dates
func (s *LicenseService) ValidateLicenseStatus(ctx context.Context, licenseID string) (*models.License, error) {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	needsUpdate := false

	// Check if time-based license should be expired
	if license.LicenseType == models.LicenseTypeTimeBased && license.EndDate != nil && license.EndDate.Before(now) {
		if license.Status != models.LicenseStatusExpired {
			license.Status = models.LicenseStatusExpired
			needsUpdate = true
		}
	}

	// Update if needed
	if needsUpdate {
		if err := s.licenseRepo.Update(ctx, license.ID, license); err != nil {
			return nil, fmt.Errorf("failed to update license status: %w", err)
		}
	}

	return license, nil
}

// GetAvailableSeats calculates available seats for a license
func (s *LicenseService) GetAvailableSeats(ctx context.Context, licenseID string) (int, error) {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return 0, err
	}

	totalAllocated, err := s.allocationRepo.GetTotalAllocatedSeats(ctx, license.ID)
	if err != nil {
		return 0, err
	}

	return license.NumberOfSeats - totalAllocated, nil
}

// CheckLicenseExpiration checks if a license is expiring soon
func (s *LicenseService) CheckLicenseExpiration(ctx context.Context, licenseID string, days int) (bool, error) {
	license, err := s.licenseRepo.GetByLicenseID(ctx, licenseID)
	if err != nil {
		return false, err
	}

	if license.LicenseType != models.LicenseTypeTimeBased || license.EndDate == nil {
		return false, nil
	}

	now := time.Now()
	expirationDate := license.EndDate
	thresholdDate := now.AddDate(0, 0, days)

	return expirationDate.Before(thresholdDate) && expirationDate.After(now), nil
}

// RenewLicense renews a license by updating the end date
func (s *LicenseService) RenewLicense(ctx context.Context, customerID, subscriptionID, licenseID string, newEndDate time.Time, userID string) (*models.License, error) {
	// Get license
	license, err := s.GetLicense(ctx, customerID, subscriptionID, licenseID)
	if err != nil {
		return nil, err
	}

	// Validate license is time-based
	if license.LicenseType != models.LicenseTypeTimeBased {
		return nil, fmt.Errorf("only time-based licenses can be renewed")
	}

	// Validate new end date
	if newEndDate.Before(license.StartDate) {
		return nil, fmt.Errorf("new end date must be after start date")
	}

	// Update end date
	license.EndDate = &newEndDate
	if license.Status == models.LicenseStatusExpired {
		license.Status = models.LicenseStatusActive
	}

	// Update license
	if err := s.licenseRepo.Update(ctx, license.ID, license); err != nil {
		return nil, fmt.Errorf("failed to renew license: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "license", license.ID.Hex(), userID, "", map[string]interface{}{
		"license_id":  license.LicenseID,
		"action":      "renew",
		"new_end_date": newEndDate,
	})

	return license, nil
}

// logAudit logs an audit entry
func (s *LicenseService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
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

