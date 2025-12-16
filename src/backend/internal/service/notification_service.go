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

// NotificationService handles notification business logic
type NotificationService struct {
	notificationRepo *repository.NotificationRepository
}

// NewNotificationService creates a new notification service
func NewNotificationService(notificationRepo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

// CreateNotification creates a new notification
func (s *NotificationService) CreateNotification(ctx context.Context, notification *models.Notification) error {
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return fmt.Errorf("failed to create notification: %w", err)
	}
	return nil
}

// GetNotifications retrieves notifications for a recipient
func (s *NotificationService) GetNotifications(ctx context.Context, recipientID string, page, limit int, unreadOnly bool) ([]*models.Notification, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"created_at": -1})

	var notifications []*models.Notification
	var err error
	var total int64

	if unreadOnly {
		notifications, err = s.notificationRepo.GetUnreadByRecipientID(ctx, recipientID, opts)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get unread notifications: %w", err)
		}
		filter := bson.M{"recipient_id": recipientID, "is_read": false}
		total, err = s.notificationRepo.Count(ctx, filter)
	} else {
		notifications, err = s.notificationRepo.GetByRecipientID(ctx, recipientID, opts)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to get notifications: %w", err)
		}
		filter := bson.M{"recipient_id": recipientID}
		total, err = s.notificationRepo.Count(ctx, filter)
	}

	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	return notifications, total, nil
}

// MarkAsRead marks a notification as read
func (s *NotificationService) MarkAsRead(ctx context.Context, notificationID string) error {
	// This would need to convert string to ObjectID in real implementation
	// For now, assuming the repository handles it
	return fmt.Errorf("not implemented: need ObjectID conversion")
}

// MarkAllAsRead marks all notifications for a recipient as read
func (s *NotificationService) MarkAllAsRead(ctx context.Context, recipientID string) error {
	if err := s.notificationRepo.MarkAllAsRead(ctx, recipientID); err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}
	return nil
}

// GetUnreadCount returns the count of unread notifications
func (s *NotificationService) GetUnreadCount(ctx context.Context, recipientID string) (int64, error) {
	filter := bson.M{
		"recipient_id": recipientID,
		"is_read":      false,
	}
	count, err := s.notificationRepo.Count(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread notifications: %w", err)
	}
	return count, nil
}

// NotifyCustomersOnVersionRelease generates notifications for customers when a new version is released
func (s *NotificationService) NotifyCustomersOnVersionRelease(ctx context.Context, productID string, versionID string, deploymentRepo *repository.DeploymentRepository, tenantRepo *repository.TenantRepository, customerRepo *repository.CustomerRepository) error {
	// Get all active deployments for the product
	deployments, err := deploymentRepo.GetDeploymentsForNotification(ctx, productID)
	if err != nil {
		return fmt.Errorf("failed to get deployments for notification: %w", err)
	}

	// Group deployments by customer
	customerDeployments := make(map[primitive.ObjectID][]*models.Deployment)
	for _, deployment := range deployments {
		// Get tenant to find customer
		tenant, err := tenantRepo.GetByID(ctx, deployment.TenantID)
		if err != nil {
			continue
		}
		customerDeployments[tenant.CustomerID] = append(customerDeployments[tenant.CustomerID], deployment)
	}

	// Create notifications per customer
	for customerID, deployments := range customerDeployments {
		// Get customer to get customer_id string
		customer, err := customerRepo.GetByID(ctx, customerID)
		if err != nil {
			continue
		}

		// Build deployment list for notification message
		deploymentList := ""
		for i, dep := range deployments {
			if i > 0 {
				deploymentList += ", "
			}
			deploymentList += fmt.Sprintf("%s (%s)", dep.ProductID, dep.DeploymentType)
		}

		// Determine priority based on deployment types
		priority := models.NotificationPriorityNormal
		for _, dep := range deployments {
			if dep.DeploymentType == models.DeploymentTypeProduction {
				priority = models.NotificationPriorityHigh
				break
			}
		}

		notification := &models.Notification{
			Type:         models.NotificationTypeNewVersion,
			RecipientID:  customer.CustomerID,
			CustomerID:   customer.CustomerID,
			ProductID:    productID,
			VersionID:    versionID,
			Title:        "New Version Available",
			Message:      fmt.Sprintf("A new version is available for product %s. Affected deployments: %s", productID, deploymentList),
			Priority:     priority,
			IsRead:       false,
			CreatedAt:    time.Now(),
		}

		if err := s.CreateNotification(ctx, notification); err != nil {
			// Log error but continue with other customers
			continue
		}
	}

	return nil
}

// getNotificationPriority determines notification priority based on deployment type
func (s *NotificationService) getNotificationPriority(deploymentType models.DeploymentType) models.NotificationPriority {
	switch deploymentType {
	case models.DeploymentTypeProduction:
		return models.NotificationPriorityHigh
	case models.DeploymentTypeUAT, models.DeploymentTypeTesting:
		return models.NotificationPriorityNormal
	default:
		return models.NotificationPriorityNormal
	}
}
