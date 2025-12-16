package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var (
	subscriptionTestDB  *database.MongoDB
	subscriptionTestCtx context.Context
	subscriptionRepo    *SubscriptionRepository
	customerRepo        *CustomerRepository
)

func setupSubscriptionTestDB(t *testing.T) {
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

	subscriptionTestDB = db
	subscriptionTestCtx = ctx
	subscriptionRepo = NewSubscriptionRepository(db.Collection("subscriptions"))
	customerRepo = NewCustomerRepository(db.Collection("customers"))
}

func teardownSubscriptionTestDB(t *testing.T) {
	if subscriptionTestDB != nil {
		_ = subscriptionTestDB.Collection("subscriptions").Drop(subscriptionTestCtx)
		_ = subscriptionTestDB.Collection("customers").Drop(subscriptionTestCtx)
		_ = subscriptionTestDB.Disconnect(subscriptionTestCtx)
	}
}

func TestSubscriptionCreate(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	// Create a customer first
	customer := &models.Customer{
		CustomerID:   "sub-test-customer",
		Name:         "Subscription Test Customer",
		Email:        "subtest@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	err := customerRepo.Create(subscriptionTestCtx, customer)
	if err != nil {
		t.Fatalf("Failed to create customer: %v", err)
	}

	subscription := &models.Subscription{
		SubscriptionID: "SUB-001",
		CustomerID:     customer.ID,
		Name:           "Test Subscription",
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}

	err = subscriptionRepo.Create(subscriptionTestCtx, subscription)
	if err != nil {
		t.Fatalf("Failed to create subscription: %v", err)
	}

	if subscription.ID.IsZero() {
		t.Error("Subscription ID was not set after creation")
	}
	if subscription.CreatedAt.IsZero() {
		t.Error("CreatedAt was not set")
	}

	t.Logf("Created subscription with ID: %s", subscription.ID.Hex())
}

func TestSubscriptionGetByID(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-get-customer",
		Name:         "Get Customer",
		Email:        "get@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerRepo.Create(subscriptionTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-002",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionRepo.Create(subscriptionTestCtx, subscription)

	retrieved, err := subscriptionRepo.GetByID(subscriptionTestCtx, subscription.ID)
	if err != nil {
		t.Fatalf("Failed to get subscription: %v", err)
	}

	if retrieved.SubscriptionID != subscription.SubscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", retrieved.SubscriptionID, subscription.SubscriptionID)
	}

	t.Logf("Retrieved subscription: %+v", retrieved)
}

func TestSubscriptionGetBySubscriptionID(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-getbyid-customer",
		Name:         "GetByID Customer",
		Email:        "getbyid@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerRepo.Create(subscriptionTestCtx, customer)

	subscriptionID := "SUB-003"
	subscription := &models.Subscription{
		SubscriptionID: subscriptionID,
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionRepo.Create(subscriptionTestCtx, subscription)

	retrieved, err := subscriptionRepo.GetBySubscriptionID(subscriptionTestCtx, subscriptionID)
	if err != nil {
		t.Fatalf("Failed to get subscription by subscription_id: %v", err)
	}

	if retrieved.SubscriptionID != subscriptionID {
		t.Errorf("SubscriptionID mismatch: got %s, want %s", retrieved.SubscriptionID, subscriptionID)
	}

	t.Logf("Retrieved subscription by subscription_id: %+v", retrieved)
}

func TestSubscriptionGetByCustomerID(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-list-customer",
		Name:         "List Customer",
		Email:        "list@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerRepo.Create(subscriptionTestCtx, customer)

	// Create multiple subscriptions
	for i := 1; i <= 3; i++ {
		subscription := &models.Subscription{
			SubscriptionID: "SUB-LIST-" + string(rune('0'+i)),
			CustomerID:     customer.ID,
			StartDate:      time.Now(),
			Status:         models.SubscriptionStatusActive,
			CreatedBy:      "user-123",
		}
		subscriptionRepo.Create(subscriptionTestCtx, subscription)
	}

	subscriptions, pagination, err := subscriptionRepo.GetByCustomerID(subscriptionTestCtx, customer.ID, nil, &Pagination{Page: 1, Limit: 10})
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

func TestSubscriptionUpdate(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-update-customer",
		Name:         "Update Customer",
		Email:        "update@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerRepo.Create(subscriptionTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-004",
		CustomerID:     customer.ID,
		Name:           "Original Name",
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionRepo.Create(subscriptionTestCtx, subscription)

	subscription.Name = "Updated Name"
	subscription.Status = models.SubscriptionStatusInactive
	err := subscriptionRepo.Update(subscriptionTestCtx, subscription.ID, subscription)
	if err != nil {
		t.Fatalf("Failed to update subscription: %v", err)
	}

	retrieved, _ := subscriptionRepo.GetByID(subscriptionTestCtx, subscription.ID)
	if retrieved.Name != "Updated Name" {
		t.Errorf("Name not updated: got %s, want Updated Name", retrieved.Name)
	}
	if retrieved.Status != models.SubscriptionStatusInactive {
		t.Errorf("Status not updated: got %s, want inactive", retrieved.Status)
	}

	t.Logf("Updated subscription: %+v", retrieved)
}

func TestSubscriptionDelete(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-delete-customer",
		Name:         "Delete Customer",
		Email:        "delete@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerRepo.Create(subscriptionTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-005",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionRepo.Create(subscriptionTestCtx, subscription)

	err := subscriptionRepo.Delete(subscriptionTestCtx, subscription.ID)
	if err != nil {
		t.Fatalf("Failed to delete subscription: %v", err)
	}

	_, err = subscriptionRepo.GetByID(subscriptionTestCtx, subscription.ID)
	if err == nil {
		t.Error("Subscription should be deleted")
	}

	t.Logf("Subscription deleted successfully")
}

func TestSubscriptionGetExpiringSubscriptions(t *testing.T) {
	setupSubscriptionTestDB(t)
	defer teardownSubscriptionTestDB(t)

	customer := &models.Customer{
		CustomerID:   "sub-expire-customer",
		Name:         "Expire Customer",
		Email:        "expire@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerRepo.Create(subscriptionTestCtx, customer)

	// Create subscription expiring in 30 days
	endDate := time.Now().AddDate(0, 0, 30)
	subscription := &models.Subscription{
		SubscriptionID: "SUB-EXPIRE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		EndDate:        &endDate,
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionRepo.Create(subscriptionTestCtx, subscription)

	expiring, err := subscriptionRepo.GetExpiringSubscriptions(subscriptionTestCtx, 60)
	if err != nil {
		t.Fatalf("Failed to get expiring subscriptions: %v", err)
	}

	if len(expiring) == 0 {
		t.Error("Expected at least one expiring subscription")
	}

	t.Logf("Found %d expiring subscriptions", len(expiring))
}

