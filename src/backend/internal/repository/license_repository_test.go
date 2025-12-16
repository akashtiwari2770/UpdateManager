package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var (
	licenseTestDB      *database.MongoDB
	licenseTestCtx     context.Context
	licenseRepo        *LicenseRepository
	subscriptionTestRepo *SubscriptionRepository
	customerTestRepo   *CustomerRepository
)

func setupLicenseTestDB(t *testing.T) {
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

	licenseTestDB = db
	licenseTestCtx = ctx
	licenseRepo = NewLicenseRepository(db.Collection("licenses"))
	subscriptionTestRepo = NewSubscriptionRepository(db.Collection("subscriptions"))
	customerTestRepo = NewCustomerRepository(db.Collection("customers"))
}

func teardownLicenseTestDB(t *testing.T) {
	if licenseTestDB != nil {
		_ = licenseTestDB.Collection("licenses").Drop(licenseTestCtx)
		_ = licenseTestDB.Collection("subscriptions").Drop(licenseTestCtx)
		_ = licenseTestDB.Collection("customers").Drop(licenseTestCtx)
		_ = licenseTestDB.Disconnect(licenseTestCtx)
	}
}

func TestLicenseCreate(t *testing.T) {
	setupLicenseTestDB(t)
	defer teardownLicenseTestDB(t)

	// Create customer and subscription
	customer := &models.Customer{
		CustomerID:   "lic-test-customer",
		Name:         "License Test Customer",
		Email:        "lictest@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerTestRepo.Create(licenseTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-001",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionTestRepo.Create(licenseTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-001",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-001",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}

	err := licenseRepo.Create(licenseTestCtx, license)
	if err != nil {
		t.Fatalf("Failed to create license: %v", err)
	}

	if license.ID.IsZero() {
		t.Error("License ID was not set after creation")
	}

	t.Logf("Created license with ID: %s", license.ID.Hex())
}

func TestLicenseGetByID(t *testing.T) {
	setupLicenseTestDB(t)
	defer teardownLicenseTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-get-customer",
		Name:         "Get Customer",
		Email:        "getlic@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerTestRepo.Create(licenseTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-002",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionTestRepo.Create(licenseTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-002",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-002",
		LicenseType:    models.LicenseTypeTimeBased,
		NumberOfSeats:  50,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	licenseRepo.Create(licenseTestCtx, license)

	retrieved, err := licenseRepo.GetByID(licenseTestCtx, license.ID)
	if err != nil {
		t.Fatalf("Failed to get license: %v", err)
	}

	if retrieved.LicenseID != license.LicenseID {
		t.Errorf("LicenseID mismatch: got %s, want %s", retrieved.LicenseID, license.LicenseID)
	}

	t.Logf("Retrieved license: %+v", retrieved)
}

func TestLicenseGetBySubscriptionID(t *testing.T) {
	setupLicenseTestDB(t)
	defer teardownLicenseTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-list-customer",
		Name:         "List Customer",
		Email:        "listlic@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerTestRepo.Create(licenseTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-003",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionTestRepo.Create(licenseTestCtx, subscription)

	// Create multiple licenses
	for i := 1; i <= 3; i++ {
		license := &models.License{
			LicenseID:      "LIC-LIST-" + string(rune('0'+i)),
			SubscriptionID: subscription.ID,
			ProductID:      "PROD-003",
			LicenseType:    models.LicenseTypePerpetual,
			NumberOfSeats:  10 * i,
			StartDate:      time.Now(),
			Status:         models.LicenseStatusActive,
			AssignedBy:     "sales-user",
			AssignmentDate: time.Now(),
		}
		licenseRepo.Create(licenseTestCtx, license)
	}

	licenses, pagination, err := licenseRepo.GetBySubscriptionID(licenseTestCtx, subscription.ID, nil, &Pagination{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("Failed to list licenses: %v", err)
	}

	if len(licenses) != 3 {
		t.Errorf("Expected 3 licenses, got %d", len(licenses))
	}
	if pagination.Total != 3 {
		t.Errorf("Expected total 3, got %d", pagination.Total)
	}

	t.Logf("Listed %d licenses", len(licenses))
}

func TestLicenseGetExpiringLicenses(t *testing.T) {
	setupLicenseTestDB(t)
	defer teardownLicenseTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-expire-customer",
		Name:         "Expire Customer",
		Email:        "expirelic@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	customerTestRepo.Create(licenseTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-EXPIRE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	subscriptionTestRepo.Create(licenseTestCtx, subscription)

	// Create time-based license expiring in 30 days
	endDate := time.Now().AddDate(0, 0, 30)
	license := &models.License{
		LicenseID:      "LIC-EXPIRE",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-EXPIRE",
		LicenseType:    models.LicenseTypeTimeBased,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		EndDate:        &endDate,
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	licenseRepo.Create(licenseTestCtx, license)

	expiring, err := licenseRepo.GetExpiringLicenses(licenseTestCtx, 60)
	if err != nil {
		t.Fatalf("Failed to get expiring licenses: %v", err)
	}

	if len(expiring) == 0 {
		t.Error("Expected at least one expiring license")
	}

	t.Logf("Found %d expiring licenses", len(expiring))
}

