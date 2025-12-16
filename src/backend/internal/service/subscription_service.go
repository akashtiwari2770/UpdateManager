package service

import (
	"context"
	"fmt"
	"time"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// SubscriptionService handles subscription business logic
type SubscriptionService struct {
	subscriptionRepo *repository.SubscriptionRepository
	customerRepo     *repository.CustomerRepository
	licenseRepo      *repository.LicenseRepository
	auditRepo        *repository.AuditLogRepository
}

// NewSubscriptionService creates a new subscription service
func NewSubscriptionService(
	subscriptionRepo *repository.SubscriptionRepository,
	customerRepo *repository.CustomerRepository,
	licenseRepo *repository.LicenseRepository,
	auditRepo *repository.AuditLogRepository,
) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
		customerRepo:     customerRepo,
		licenseRepo:      licenseRepo,
		auditRepo:        auditRepo,
	}
}

// CreateSubscription creates a new subscription with validation
func (s *SubscriptionService) CreateSubscription(ctx context.Context, customerID string, req *models.CreateSubscriptionRequest, userID string) (*models.Subscription, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Validate subscription_id uniqueness
	existing, err := s.subscriptionRepo.GetBySubscriptionID(ctx, req.SubscriptionID)
	if err == nil && existing != nil {
		return nil, fmt.Errorf("subscription with subscription_id '%s' already exists", req.SubscriptionID)
	}

	// Validate dates
	if req.EndDate != nil && req.EndDate.Before(req.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Create subscription
	subscription := &models.Subscription{
		SubscriptionID: req.SubscriptionID,
		CustomerID:     customer.ID,
		Name:           req.Name,
		Description:    req.Description,
		StartDate:      req.StartDate,
		EndDate:        req.EndDate,
		Status:         req.Status,
		CreatedBy:      userID,
		Notes:          req.Notes,
	}

	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, fmt.Errorf("failed to create subscription: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "subscription", subscription.ID.Hex(), userID, "", map[string]interface{}{
		"subscription_id": subscription.SubscriptionID,
		"customer_id":     customerID,
		"status":          subscription.Status,
	})

	return subscription, nil
}

// GetSubscription retrieves a subscription by customer ID and subscription ID
func (s *SubscriptionService) GetSubscription(ctx context.Context, customerID, subscriptionID string) (*models.Subscription, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Get subscription
	subscription, err := s.subscriptionRepo.GetBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, fmt.Errorf("subscription not found: %w", err)
	}

	// Verify subscription belongs to customer
	if subscription.CustomerID != customer.ID {
		return nil, fmt.Errorf("subscription does not belong to customer")
	}

	return subscription, nil
}

// ListSubscriptions lists subscriptions for a customer with filters and pagination
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, customerID string, filter *repository.SubscriptionFilter, pagination *repository.Pagination) ([]*models.Subscription, *repository.PaginationInfo, error) {
	// Validate customer exists
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return nil, nil, fmt.Errorf("customer not found: %w", err)
	}

	// Get subscriptions
	subscriptions, paginationInfo, err := s.subscriptionRepo.GetByCustomerID(ctx, customer.ID, filter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to list subscriptions: %w", err)
	}

	return subscriptions, paginationInfo, nil
}

// UpdateSubscription updates an existing subscription
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, customerID, subscriptionID string, req *models.UpdateSubscriptionRequest, userID string) (*models.Subscription, error) {
	// Get subscription
	subscription, err := s.GetSubscription(ctx, customerID, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Update fields
	if req.Name != nil {
		subscription.Name = *req.Name
	}
	if req.Description != nil {
		subscription.Description = *req.Description
	}
	if req.StartDate != nil {
		subscription.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		subscription.EndDate = req.EndDate
	}
	if req.Status != nil {
		subscription.Status = *req.Status
	}
	if req.Notes != nil {
		subscription.Notes = *req.Notes
	}

	// Validate dates
	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		return nil, fmt.Errorf("end date must be after start date")
	}

	// Update subscription
	if err := s.subscriptionRepo.Update(ctx, subscription.ID, subscription); err != nil {
		return nil, fmt.Errorf("failed to update subscription: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "subscription", subscription.ID.Hex(), userID, "", map[string]interface{}{
		"subscription_id": subscription.SubscriptionID,
		"customer_id":     customerID,
	})

	return subscription, nil
}

// DeleteSubscription deletes a subscription
func (s *SubscriptionService) DeleteSubscription(ctx context.Context, customerID, subscriptionID string, userID string) error {
	// Get subscription
	subscription, err := s.GetSubscription(ctx, customerID, subscriptionID)
	if err != nil {
		return err
	}

	// Check if subscription has licenses
	licenseFilter := &repository.LicenseFilter{}
	licenses, _, err := s.licenseRepo.GetBySubscriptionID(ctx, subscription.ID, licenseFilter, &repository.Pagination{Page: 1, Limit: 1})
	if err == nil && len(licenses) > 0 {
		return fmt.Errorf("cannot delete subscription with existing licenses")
	}

	// Delete subscription
	if err := s.subscriptionRepo.Delete(ctx, subscription.ID); err != nil {
		return fmt.Errorf("failed to delete subscription: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "subscription", subscription.ID.Hex(), userID, "", map[string]interface{}{
		"subscription_id": subscription.SubscriptionID,
		"customer_id":     customerID,
	})

	return nil
}

// GetSubscriptionStatistics retrieves statistics for a subscription
func (s *SubscriptionService) GetSubscriptionStatistics(ctx context.Context, customerID, subscriptionID string) (map[string]interface{}, error) {
	// Get subscription
	subscription, err := s.GetSubscription(ctx, customerID, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Get licenses count
	licenseFilter := &repository.LicenseFilter{}
	licenses, _, err := s.licenseRepo.GetBySubscriptionID(ctx, subscription.ID, licenseFilter, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get licenses: %w", err)
	}

	// Calculate statistics
	stats := map[string]interface{}{
		"total_licenses":        len(licenses),
		"active_licenses":       0,
		"expired_licenses":       0,
		"total_seats":           0,
		"perpetual_licenses":     0,
		"time_based_licenses":    0,
	}

	for _, license := range licenses {
		if license.Status == models.LicenseStatusActive {
			stats["active_licenses"] = stats["active_licenses"].(int) + 1
		}
		if license.Status == models.LicenseStatusExpired {
			stats["expired_licenses"] = stats["expired_licenses"].(int) + 1
		}
		stats["total_seats"] = stats["total_seats"].(int) + license.NumberOfSeats
		if license.LicenseType == models.LicenseTypePerpetual {
			stats["perpetual_licenses"] = stats["perpetual_licenses"].(int) + 1
		}
		if license.LicenseType == models.LicenseTypeTimeBased {
			stats["time_based_licenses"] = stats["time_based_licenses"].(int) + 1
		}
	}

	return stats, nil
}

// RenewSubscription renews a subscription by updating the end date
func (s *SubscriptionService) RenewSubscription(ctx context.Context, customerID, subscriptionID string, newEndDate time.Time, userID string) (*models.Subscription, error) {
	// Get subscription
	subscription, err := s.GetSubscription(ctx, customerID, subscriptionID)
	if err != nil {
		return nil, err
	}

	// Validate new end date
	if newEndDate.Before(subscription.StartDate) {
		return nil, fmt.Errorf("new end date must be after start date")
	}

	// Update end date
	subscription.EndDate = &newEndDate
	if subscription.Status == models.SubscriptionStatusExpired {
		subscription.Status = models.SubscriptionStatusActive
	}

	// Update subscription
	if err := s.subscriptionRepo.Update(ctx, subscription.ID, subscription); err != nil {
		return nil, fmt.Errorf("failed to renew subscription: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "subscription", subscription.ID.Hex(), userID, "", map[string]interface{}{
		"subscription_id": subscription.SubscriptionID,
		"action":          "renew",
		"new_end_date":    newEndDate,
	})

	return subscription, nil
}

// ValidateSubscriptionStatus validates and updates subscription status based on dates
func (s *SubscriptionService) ValidateSubscriptionStatus(ctx context.Context, subscriptionID string) (*models.Subscription, error) {
	subscription, err := s.subscriptionRepo.GetBySubscriptionID(ctx, subscriptionID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	needsUpdate := false

	// Check if subscription should be expired
	if subscription.EndDate != nil && subscription.EndDate.Before(now) {
		if subscription.Status != models.SubscriptionStatusExpired {
			subscription.Status = models.SubscriptionStatusExpired
			needsUpdate = true
		}
	}

	// Update if needed
	if needsUpdate {
		if err := s.subscriptionRepo.Update(ctx, subscription.ID, subscription); err != nil {
			return nil, fmt.Errorf("failed to update subscription status: %w", err)
		}
	}

	return subscription, nil
}

// logAudit logs an audit entry
func (s *SubscriptionService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
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

