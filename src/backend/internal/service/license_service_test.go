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
	licenseServiceTestDB          *database.MongoDB
	licenseServiceTestCtx         context.Context
	licenseService                *LicenseService
	licenseServiceRepo            *repository.LicenseRepository
	licenseServiceSubscriptionRepo *repository.SubscriptionRepository
	licenseServiceCustomerRepo    *repository.CustomerRepository
	licenseServiceAllocationRepo  *repository.LicenseAllocationRepository
	licenseServiceAuditRepo       *repository.AuditLogRepository
)

func setupLicenseServiceTestDB(t *testing.T) {
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

	licenseServiceTestDB = db
	licenseServiceTestCtx = ctx
	licenseServiceRepo = repository.NewLicenseRepository(db.Collection("licenses"))
	licenseServiceSubscriptionRepo = repository.NewSubscriptionRepository(db.Collection("subscriptions"))
	licenseServiceCustomerRepo = repository.NewCustomerRepository(db.Collection("customers"))
	licenseServiceAllocationRepo = repository.NewLicenseAllocationRepository(db.Collection("license_allocations"))
	licenseServiceAuditRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	licenseService = NewLicenseService(licenseServiceRepo, licenseServiceSubscriptionRepo, licenseServiceCustomerRepo, licenseServiceAllocationRepo, licenseServiceAuditRepo)
}

func teardownLicenseServiceTestDB(t *testing.T) {
	if licenseServiceTestDB != nil {
		_ = licenseServiceTestDB.Collection("licenses").Drop(licenseServiceTestCtx)
		_ = licenseServiceTestDB.Collection("subscriptions").Drop(licenseServiceTestCtx)
		_ = licenseServiceTestDB.Collection("customers").Drop(licenseServiceTestCtx)
		_ = licenseServiceTestDB.Collection("license_allocations").Drop(licenseServiceTestCtx)
		_ = licenseServiceTestDB.Collection("audit_logs").Drop(licenseServiceTestCtx)
		_ = licenseServiceTestDB.Disconnect(licenseServiceTestCtx)
	}
}

func TestLicenseService_AssignLicense(t *testing.T) {
	setupLicenseServiceTestDB(t)
	defer teardownLicenseServiceTestDB(t)

	// Create customer and subscription
	customer := &models.Customer{
		CustomerID:   "lic-svc-customer",
		Name:         "License Service Customer",
		Email:        "licsvc@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	licenseServiceCustomerRepo.Create(licenseServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-SVC",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	licenseServiceSubscriptionRepo.Create(licenseServiceTestCtx, subscription)

	req := &models.CreateLicenseRequest{
		LicenseID:     "LIC-SVC-001",
		ProductID:     "PROD-001",
		LicenseType:   models.LicenseTypePerpetual,
		NumberOfSeats: 100,
		StartDate:     time.Now(),
		Status:        models.LicenseStatusActive,
	}

	license, err := licenseService.AssignLicense(licenseServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, req, "sales-user")
	if err != nil {
		t.Fatalf("Failed to assign license: %v", err)
	}

	if license.ID.IsZero() {
		t.Error("License ID was not set")
	}
	if license.LicenseID != req.LicenseID {
		t.Errorf("LicenseID mismatch: got %s, want %s", license.LicenseID, req.LicenseID)
	}
	if license.NumberOfSeats != req.NumberOfSeats {
		t.Errorf("NumberOfSeats mismatch: got %d, want %d", license.NumberOfSeats, req.NumberOfSeats)
	}

	t.Logf("Assigned license: %+v", license)
}

func TestLicenseService_AssignLicense_TimeBasedRequiresEndDate(t *testing.T) {
	setupLicenseServiceTestDB(t)
	defer teardownLicenseServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-svc-time-customer",
		Name:         "Time Customer",
		Email:        "time@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	licenseServiceCustomerRepo.Create(licenseServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-TIME",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	licenseServiceSubscriptionRepo.Create(licenseServiceTestCtx, subscription)

	req := &models.CreateLicenseRequest{
		LicenseID:     "LIC-TIME-001",
		ProductID:     "PROD-001",
		LicenseType:   models.LicenseTypeTimeBased,
		NumberOfSeats: 100,
		StartDate:     time.Now(),
		Status:        models.LicenseStatusActive,
		// EndDate is nil - should fail
	}

	_, err := licenseService.AssignLicense(licenseServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, req, "sales-user")
	if err == nil {
		t.Error("Expected error for time-based license without end date")
	}

	t.Logf("Correctly rejected time-based license without end date: %v", err)
}

func TestLicenseService_GetLicense(t *testing.T) {
	setupLicenseServiceTestDB(t)
	defer teardownLicenseServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-svc-get-customer",
		Name:         "Get Customer",
		Email:        "getlic@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	licenseServiceCustomerRepo.Create(licenseServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-GET",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	licenseServiceSubscriptionRepo.Create(licenseServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-GET",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-GET",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  50,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	licenseServiceRepo.Create(licenseServiceTestCtx, license)

	retrieved, err := licenseService.GetLicense(licenseServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID)
	if err != nil {
		t.Fatalf("Failed to get license: %v", err)
	}

	if retrieved.LicenseID != license.LicenseID {
		t.Errorf("LicenseID mismatch: got %s, want %s", retrieved.LicenseID, license.LicenseID)
	}

	t.Logf("Retrieved license: %+v", retrieved)
}

func TestLicenseService_GetAvailableSeats(t *testing.T) {
	setupLicenseServiceTestDB(t)
	defer teardownLicenseServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-svc-seats-customer",
		Name:         "Seats Customer",
		Email:        "seats@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	licenseServiceCustomerRepo.Create(licenseServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-SEATS",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	licenseServiceSubscriptionRepo.Create(licenseServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-SEATS",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-SEATS",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	licenseServiceRepo.Create(licenseServiceTestCtx, license)

	// Allocate some seats
	allocation := &models.LicenseAllocation{
		AllocationID:          "ALLOC-SEATS",
		LicenseID:             license.ID,
		NumberOfSeatsAllocated: 30,
		AllocationDate:        time.Now(),
		AllocatedBy:           "customer-user",
		Status:                models.AllocationStatusActive,
	}
	licenseServiceAllocationRepo.Create(licenseServiceTestCtx, allocation)

	available, err := licenseService.GetAvailableSeats(licenseServiceTestCtx, license.LicenseID)
	if err != nil {
		t.Fatalf("Failed to get available seats: %v", err)
	}

	expectedAvailable := 100 - 30
	if available != expectedAvailable {
		t.Errorf("Available seats mismatch: got %d, want %d", available, expectedAvailable)
	}

	t.Logf("Available seats: %d", available)
}

func TestLicenseService_RevokeLicense_WithAllocations(t *testing.T) {
	setupLicenseServiceTestDB(t)
	defer teardownLicenseServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-svc-revoke-customer",
		Name:         "Revoke Customer",
		Email:        "revoke@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	licenseServiceCustomerRepo.Create(licenseServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-REVOKE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	licenseServiceSubscriptionRepo.Create(licenseServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-REVOKE",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-REVOKE",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	licenseServiceRepo.Create(licenseServiceTestCtx, license)

	// Create active allocation
	allocation := &models.LicenseAllocation{
		AllocationID:          "ALLOC-REVOKE",
		LicenseID:             license.ID,
		NumberOfSeatsAllocated: 50,
		AllocationDate:        time.Now(),
		AllocatedBy:           "customer-user",
		Status:                models.AllocationStatusActive,
	}
	licenseServiceAllocationRepo.Create(licenseServiceTestCtx, allocation)

	// Try to revoke license with active allocations - should fail
	err := licenseService.RevokeLicense(licenseServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID, "admin-user")
	if err == nil {
		t.Error("Expected error when revoking license with active allocations")
	}

	t.Logf("Correctly prevented revocation of license with active allocations: %v", err)
}

func TestLicenseService_GetLicenseStatistics(t *testing.T) {
	setupLicenseServiceTestDB(t)
	defer teardownLicenseServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "lic-svc-stats-customer",
		Name:         "Stats Customer",
		Email:        "stats@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	licenseServiceCustomerRepo.Create(licenseServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-LIC-STATS",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	licenseServiceSubscriptionRepo.Create(licenseServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-STATS",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-STATS",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	licenseServiceRepo.Create(licenseServiceTestCtx, license)

	// Create allocations
	for i := 1; i <= 2; i++ {
		allocation := &models.LicenseAllocation{
			AllocationID:          "ALLOC-STATS-" + string(rune('0'+i)),
			LicenseID:             license.ID,
			NumberOfSeatsAllocated: 20 * i,
			AllocationDate:        time.Now(),
			AllocatedBy:           "customer-user",
			Status:                models.AllocationStatusActive,
		}
		licenseServiceAllocationRepo.Create(licenseServiceTestCtx, allocation)
	}

	stats, err := licenseService.GetLicenseStatistics(licenseServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID)
	if err != nil {
		t.Fatalf("Failed to get license statistics: %v", err)
	}

	if stats["total_seats"] != 100 {
		t.Errorf("Total seats mismatch: got %v, want 100", stats["total_seats"])
	}
	if stats["allocated_seats"] != 60 {
		t.Errorf("Allocated seats mismatch: got %v, want 60", stats["allocated_seats"])
	}
	if stats["available_seats"] != 40 {
		t.Errorf("Available seats mismatch: got %v, want 40", stats["available_seats"])
	}

	t.Logf("License statistics: %+v", stats)
}

