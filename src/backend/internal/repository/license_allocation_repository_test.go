package repository

import (
	"context"
	"testing"
	"time"

	"updatemanager/internal/models"
	"updatemanager/pkg/database"
)

var (
	allocationTestDB        *database.MongoDB
	allocationTestCtx       context.Context
	allocationRepo          *LicenseAllocationRepository
	allocationLicenseRepo   *LicenseRepository
	allocationSubscriptionRepo *SubscriptionRepository
	allocationCustomerRepo  *CustomerRepository
	allocationTenantRepo    *TenantRepository
	allocationDeploymentRepo *DeploymentRepository
)

func setupAllocationTestDB(t *testing.T) {
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

	allocationTestDB = db
	allocationTestCtx = ctx
	allocationRepo = NewLicenseAllocationRepository(db.Collection("license_allocations"))
	allocationLicenseRepo = NewLicenseRepository(db.Collection("licenses"))
	allocationSubscriptionRepo = NewSubscriptionRepository(db.Collection("subscriptions"))
	allocationCustomerRepo = NewCustomerRepository(db.Collection("customers"))
	allocationTenantRepo = NewTenantRepository(db.Collection("customer_tenants"))
	allocationDeploymentRepo = NewDeploymentRepository(db.Collection("deployments"))
}

func teardownAllocationTestDB(t *testing.T) {
	if allocationTestDB != nil {
		_ = allocationTestDB.Collection("license_allocations").Drop(allocationTestCtx)
		_ = allocationTestDB.Collection("licenses").Drop(allocationTestCtx)
		_ = allocationTestDB.Collection("subscriptions").Drop(allocationTestCtx)
		_ = allocationTestDB.Collection("customers").Drop(allocationTestCtx)
		_ = allocationTestDB.Collection("customer_tenants").Drop(allocationTestCtx)
		_ = allocationTestDB.Collection("deployments").Drop(allocationTestCtx)
		_ = allocationTestDB.Disconnect(allocationTestCtx)
	}
}

func TestLicenseAllocationCreate(t *testing.T) {
	setupAllocationTestDB(t)
	defer teardownAllocationTestDB(t)

	// Create customer, subscription, license, tenant
	customer := &models.Customer{
		CustomerID:   "alloc-test-customer",
		Name:         "Allocation Test Customer",
		Email:        "alloctest@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationCustomerRepo.Create(allocationTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-001",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationSubscriptionRepo.Create(allocationTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-001",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-ALLOC",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationLicenseRepo.Create(allocationTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-ALLOC",
		CustomerID: customer.ID,
		Name:       "Allocation Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationTenantRepo.Create(allocationTestCtx, tenant)

	allocation := &models.LicenseAllocation{
		AllocationID:          "ALLOC-001",
		LicenseID:             license.ID,
		TenantID:              &tenant.ID,
		NumberOfSeatsAllocated: 50,
		AllocationDate:        time.Now(),
		AllocatedBy:           "customer-user",
		Status:                models.AllocationStatusActive,
	}

	err := allocationRepo.Create(allocationTestCtx, allocation)
	if err != nil {
		t.Fatalf("Failed to create allocation: %v", err)
	}

	if allocation.ID.IsZero() {
		t.Error("Allocation ID was not set after creation")
	}

	t.Logf("Created allocation with ID: %s", allocation.ID.Hex())
}

func TestLicenseAllocationGetTotalAllocatedSeats(t *testing.T) {
	setupAllocationTestDB(t)
	defer teardownAllocationTestDB(t)

	customer := &models.Customer{
		CustomerID:   "alloc-seats-customer",
		Name:         "Seats Customer",
		Email:        "seats@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationCustomerRepo.Create(allocationTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-SEATS",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationSubscriptionRepo.Create(allocationTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-SEATS",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-SEATS",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationLicenseRepo.Create(allocationTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-SEATS",
		CustomerID: customer.ID,
		Name:       "Seats Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationTenantRepo.Create(allocationTestCtx, tenant)

	// Create multiple allocations
	totalAllocated := 0
	for i := 1; i <= 3; i++ {
		allocation := &models.LicenseAllocation{
			AllocationID:          "ALLOC-SEATS-" + string(rune('0'+i)),
			LicenseID:             license.ID,
			TenantID:              &tenant.ID,
			NumberOfSeatsAllocated: 10 * i,
			AllocationDate:        time.Now(),
			AllocatedBy:           "customer-user",
			Status:                models.AllocationStatusActive,
		}
		allocationRepo.Create(allocationTestCtx, allocation)
		totalAllocated += 10 * i
	}

	calculatedTotal, err := allocationRepo.GetTotalAllocatedSeats(allocationTestCtx, license.ID)
	if err != nil {
		t.Fatalf("Failed to calculate total allocated seats: %v", err)
	}

	if calculatedTotal != totalAllocated {
		t.Errorf("Total allocated seats mismatch: got %d, want %d", calculatedTotal, totalAllocated)
	}

	t.Logf("Total allocated seats: %d", calculatedTotal)
}

func TestLicenseAllocationRelease(t *testing.T) {
	setupAllocationTestDB(t)
	defer teardownAllocationTestDB(t)

	customer := &models.Customer{
		CustomerID:   "alloc-release-customer",
		Name:         "Release Customer",
		Email:        "release@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationCustomerRepo.Create(allocationTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-RELEASE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationSubscriptionRepo.Create(allocationTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-RELEASE",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-RELEASE",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationLicenseRepo.Create(allocationTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-RELEASE",
		CustomerID: customer.ID,
		Name:       "Release Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationTenantRepo.Create(allocationTestCtx, tenant)

	allocation := &models.LicenseAllocation{
		AllocationID:          "ALLOC-RELEASE",
		LicenseID:             license.ID,
		TenantID:              &tenant.ID,
		NumberOfSeatsAllocated: 50,
		AllocationDate:        time.Now(),
		AllocatedBy:           "customer-user",
		Status:                models.AllocationStatusActive,
	}
	allocationRepo.Create(allocationTestCtx, allocation)

	err := allocationRepo.Release(allocationTestCtx, allocation.ID, "admin-user")
	if err != nil {
		t.Fatalf("Failed to release allocation: %v", err)
	}

	retrieved, _ := allocationRepo.GetByID(allocationTestCtx, allocation.ID)
	if retrieved.Status != models.AllocationStatusReleased {
		t.Errorf("Status not updated: got %s, want released", retrieved.Status)
	}
	if retrieved.ReleasedBy == nil || *retrieved.ReleasedBy != "admin-user" {
		t.Error("ReleasedBy not set correctly")
	}

	t.Logf("Allocation released successfully")
}

