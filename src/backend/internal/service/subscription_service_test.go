package service

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/pkg/database"
)

var (
	subscriptionServiceTestDB  *database.MongoDB
	subscriptionServiceTestCtx context.Context
	subscriptionService        *SubscriptionService
	subscriptionServiceRepo    *repository.SubscriptionRepository
	subscriptionServiceCustomerRepo *repository.CustomerRepository
	subscriptionServiceLicenseRepo  *repository.LicenseRepository
	subscriptionServiceAuditRepo    *repository.AuditLogRepository
)

func setupSubscriptionServiceTestDB(t *testing.T) {
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

	subscriptionServiceTestDB = db
	subscriptionServiceTestCtx = ctx
	subscriptionServiceRepo = repository.NewSubscriptionRepository(db.Collection("subscriptions"))
	subscriptionServiceCustomerRepo = repository.NewCustomerRepository(db.Collection("customers"))
	subscriptionServiceLicenseRepo = repository.NewLicenseRepository(db.Collection("licenses"))
	subscriptionServiceAuditRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	subscriptionService = NewSubscriptionService(subscriptionServiceRepo, subscriptionServiceCustomerRepo, subscriptionServiceLicenseRepo, subscriptionServiceAuditRepo)
}

func teardownSubscriptionServiceTestDB(t *testing.T) {
	if subscriptionServiceTestDB != nil {
		_ = subscriptionServiceTestDB.Collection("subscriptions").Drop(subscriptionServiceTestCtx)
		_ = subscriptionServiceTestDB.Collection("customers").Drop(subscriptionServiceTestCtx)
		_ = subscriptionServiceTestDB.Collection("licenses").Drop(subscriptionServiceTestCtx)
		_ = subscriptionServiceTestDB.Collection("audit_logs").Drop(subscriptionServiceTestCtx)
		_ = subscriptionServiceTestDB.Disconnect(subscriptionServiceTestCtx)
	}
}

func TestSubscriptionService_CreateSubscription(t *testing.T) {
	setupSubscriptionServiceTestDB(t)
	defer teardownSubscriptionServiceTestDB(t)

	// Create customer
	customer := &models.Customer{
		CustomerID:   "sub-svc-customer",
		Name:         "Subscription Service Customer",
		Email:        "subsvc@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	subscriptionServiceCustomerRepo.Create(subscriptionServiceTestCtx, customer)

	req := &models.CreateSubscriptionRequest{
		SubscriptionID: "SUB-SVC-001",
		Name:           "Test Subscription Service",
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
	}

	subscription, err := subscriptionService.CreateSubscription(subscriptionServiceTestCtx, customer.CustomerID, req, "user-123")
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	if subscription.ID.IsZero() {
		t.Error("Subscription ID was not set")
	}
	if subscription.SubscriptionID != req.SubscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", subscription.SubscriptionID, req.SubscriptionID)
	}

	t.Logf("Created subscription: %+v", subscription)
}

func TestSubscriptionService_GetSubscription(t *testing.T) {
	setupSubscriptionServiceTestDB(t)
	defer teardownSubscriptionServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-svc-get-customer",
		Name:         "Get Customer",
		Email:        "getsub@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	subscriptionServiceCustomerRepo.Create(subscriptionServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-SVC-GET",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionServiceRepo.Create(subscriptionServiceTestCtx, subscription)

	retrieved, err := subscriptionService.GetSubscription(subscriptionServiceTestCtx, customer.CustomerID, subscription.SubscriptionID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}

	if retrieved.SubscriptionID != subscription.SubscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", retrieved.SubscriptionID, subscription.SubscriptionID)
	}

	t.Logf("Retrieved subscription: %+v", retrieved)
}

func TestSubscriptionService_ListSubscriptions(t *testing.T) {
	setupSubscriptionServiceTestDB(t)
	defer teardownSubscriptionServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-svc-list-customer",
		Name:         "List Customer",
		Email:        "listsub@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	subscriptionServiceCustomerRepo.Create(subscriptionServiceTestCtx, customer)

	// Create multiple subscriptions
	for i := 1; i <= 3; i++ {
		subscription := &models.Subscription{
			SubscriptionID: "SUB-SVC-LIST-" + string(rune('0'+i)),
			CustomerID:     customer.ID,
			StartDate:      time.Now(),
			Status:         models.SubscriptionStatusActive,
			CreatedBy:      "user-123",
		}
		subscriptionServiceRepo.Create(subscriptionServiceTestCtx, subscription)
	}

	subscriptions, pagination, err := subscriptionService.ListSubscriptions(subscriptionServiceTestCtx, customer.CustomerID, nil, &repository.Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("Failed to list subscriptions: %v", err)
	}

	if len(subscriptions) != 3 {
		t.Errorf("Expected 3 subscriptions, got %d", len(subscriptions))
	}
	if pagination.Total != 3 {
		t.Errorf("Expected total 3, got %d", pagination.Total)
	}

	t.Logf("Listed %d subscriptions", len(subscriptions))
}

func TestSubscriptionService_DeleteSubscription_WithLicenses(t *testing.T) {
	setupSubscriptionServiceTestDB(t)
	defer teardownSubscriptionServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-svc-delete-customer",
		Name:         "Delete Customer",
		Email:        "deletesub@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	subscriptionServiceCustomerRepo.Create(subscriptionServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-SVC-DELETE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionServiceRepo.Create(subscriptionServiceTestCtx, subscription)

	// Create a license for the subscription
	license := &models.License{
		LicenseID:      "LIC-DELETE",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-DELETE",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	subscriptionServiceLicenseRepo.Create(subscriptionServiceTestCtx, license)

	// Try to delete subscription with licenses - should fail
	err := subscriptionService.DeleteSubscription(subscriptionServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, "user-123")
	if err == nil {
		t.Error("Expected error when deleting subscription with licenses")
	}

	t.Logf("Correctly prevented deletion of subscription with licenses: %v", err)
}

