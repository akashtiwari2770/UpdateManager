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
	allocationServiceTestDB            *database.MongoDB
	allocationServiceTestCtx           context.Context
	allocationService                  *LicenseAllocationService
	allocationServiceRepo              *repository.LicenseAllocationRepository
	allocationServiceLicenseRepo       *repository.LicenseRepository
	allocationServiceSubscriptionRepo  *repository.SubscriptionRepository
	allocationServiceCustomerRepo      *repository.CustomerRepository
	allocationServiceTenantRepo        *repository.TenantRepository
	allocationServiceDeploymentRepo    *repository.DeploymentRepository
	allocationServiceAuditRepo         *repository.AuditLogRepository
)

func setupAllocationServiceTestDB(t *testing.T) {
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

	allocationServiceTestDB = db
	allocationServiceTestCtx = ctx
	allocationServiceRepo = repository.NewLicenseAllocationRepository(db.Collection("license_allocations"))
	allocationServiceLicenseRepo = repository.NewLicenseRepository(db.Collection("licenses"))
	allocationServiceSubscriptionRepo = repository.NewSubscriptionRepository(db.Collection("subscriptions"))
	allocationServiceCustomerRepo = repository.NewCustomerRepository(db.Collection("customers"))
	allocationServiceTenantRepo = repository.NewTenantRepository(db.Collection("customer_tenants"))
	allocationServiceDeploymentRepo = repository.NewDeploymentRepository(db.Collection("deployments"))
	allocationServiceAuditRepo = repository.NewAuditLogRepository(db.Collection("audit_logs"))
	allocationService = NewLicenseAllocationService(
		allocationServiceRepo,
		allocationServiceLicenseRepo,
		allocationServiceSubscriptionRepo,
		allocationServiceCustomerRepo,
		allocationServiceTenantRepo,
		allocationServiceDeploymentRepo,
		allocationServiceAuditRepo,
	)
}

func teardownAllocationServiceTestDB(t *testing.T) {
	if allocationServiceTestDB != nil {
		_ = allocationServiceTestDB.Collection("license_allocations").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Collection("licenses").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Collection("subscriptions").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Collection("customers").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Collection("customer_tenants").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Collection("deployments").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Collection("audit_logs").Drop(allocationServiceTestCtx)
		_ = allocationServiceTestDB.Disconnect(allocationServiceTestCtx)
	}
}

func TestLicenseAllocationService_AllocateLicense(t *testing.T) {
	setupAllocationServiceTestDB(t)
	defer teardownAllocationServiceTestDB(t)

	// Create customer, subscription, license, tenant
	customer := &models.Customer{
		CustomerID:   "alloc-svc-customer",
		Name:         "Allocation Service Customer",
		Email:        "allocsvc@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationServiceCustomerRepo.Create(allocationServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-SVC",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationServiceSubscriptionRepo.Create(allocationServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-SVC",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-ALLOC",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationServiceLicenseRepo.Create(allocationServiceTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-ALLOC-SVC",
		CustomerID: customer.ID,
		Name:       "Allocation Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationServiceTenantRepo.Create(allocationServiceTestCtx, tenant)

	req := &models.AllocateLicenseRequest{
		TenantID:              &tenant.TenantID,
		NumberOfSeatsAllocated: 50,
	}

	allocation, err := allocationService.AllocateLicense(allocationServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID, req, "customer-user")
	if err != nil {
		t.Fatalf("Failed to allocate license: %v", err)
	}

	if allocation.ID.IsZero() {
		t.Error("Allocation ID was not set")
	}
	if allocation.NumberOfSeatsAllocated != req.NumberOfSeatsAllocated {
		t.Errorf("NumberOfSeatsAllocated mismatch: got %d, want %d", allocation.NumberOfSeatsAllocated, req.NumberOfSeatsAllocated)
	}
	if allocation.Status != models.AllocationStatusActive {
		t.Errorf("Status mismatch: got %s, want active", allocation.Status)
	}

	t.Logf("Allocated license: %+v", allocation)
}

func TestLicenseAllocationService_AllocateLicense_InsufficientSeats(t *testing.T) {
	setupAllocationServiceTestDB(t)
	defer teardownAllocationServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "alloc-svc-insufficient-customer",
		Name:         "Insufficient Customer",
		Email:        "insufficient@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationServiceCustomerRepo.Create(allocationServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-INSUFFICIENT",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationServiceSubscriptionRepo.Create(allocationServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-INSUFFICIENT",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-INSUFFICIENT",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationServiceLicenseRepo.Create(allocationServiceTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-INSUFFICIENT",
		CustomerID: customer.ID,
		Name:       "Insufficient Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationServiceTenantRepo.Create(allocationServiceTestCtx, tenant)

	// Allocate 60 seats first
	req1 := &models.AllocateLicenseRequest{
		TenantID:              &tenant.TenantID,
		NumberOfSeatsAllocated: 60,
	}
	allocationService.AllocateLicense(allocationServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID, req1, "customer-user")

	// Try to allocate 50 more seats (only 40 available) - should fail
	req2 := &models.AllocateLicenseRequest{
		TenantID:              &tenant.TenantID,
		NumberOfSeatsAllocated: 50,
	}
	_, err := allocationService.AllocateLicense(allocationServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID, req2, "customer-user")
	if err == nil {
		t.Error("Expected error when allocating more seats than available")
	}

	t.Logf("Correctly rejected allocation with insufficient seats: %v", err)
}

func TestLicenseAllocationService_ReleaseAllocation(t *testing.T) {
	setupAllocationServiceTestDB(t)
	defer teardownAllocationServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "alloc-svc-release-customer",
		Name:         "Release Customer",
		Email:        "release@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationServiceCustomerRepo.Create(allocationServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-RELEASE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationServiceSubscriptionRepo.Create(allocationServiceTestCtx, subscription)

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
	allocationServiceLicenseRepo.Create(allocationServiceTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-RELEASE",
		CustomerID: customer.ID,
		Name:       "Release Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationServiceTenantRepo.Create(allocationServiceTestCtx, tenant)

	// Create allocation
	allocation := &models.LicenseAllocation{
		AllocationID:          "ALLOC-RELEASE",
		LicenseID:             license.ID,
		TenantID:              &tenant.ID,
		NumberOfSeatsAllocated: 50,
		AllocationDate:        time.Now(),
		AllocatedBy:           "customer-user",
		Status:                models.AllocationStatusActive,
	}
	allocationServiceRepo.Create(allocationServiceTestCtx, allocation)

	err := allocationService.ReleaseAllocation(allocationServiceTestCtx, customer.CustomerID, subscription.SubscriptionID, license.LicenseID, allocation.AllocationID, "admin-user")
	if err != nil {
		t.Fatalf("Failed to release allocation: %v", err)
	}

	retrieved, _ := allocationServiceRepo.GetByID(allocationServiceTestCtx, allocation.ID)
	if retrieved.Status != models.AllocationStatusReleased {
		t.Errorf("Status not updated: got %s, want released", retrieved.Status)
	}

	t.Logf("Allocation released successfully")
}

func TestLicenseAllocationService_GetLicenseUtilization(t *testing.T) {
	setupAllocationServiceTestDB(t)
	defer teardownAllocationServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "alloc-svc-util-customer",
		Name:         "Utilization Customer",
		Email:        "util@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationServiceCustomerRepo.Create(allocationServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-UTIL",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationServiceSubscriptionRepo.Create(allocationServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-UTIL",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-UTIL",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationServiceLicenseRepo.Create(allocationServiceTestCtx, license)

	tenant := &models.CustomerTenant{
		TenantID:   "TENANT-UTIL",
		CustomerID: customer.ID,
		Name:       "Utilization Tenant",
		Status:     models.TenantStatusActive,
	}
	allocationServiceTenantRepo.Create(allocationServiceTestCtx, tenant)

	// Create allocations
	for i := 1; i <= 2; i++ {
		allocation := &models.LicenseAllocation{
			AllocationID:          "ALLOC-UTIL-" + string(rune('0'+i)),
			LicenseID:             license.ID,
			TenantID:              &tenant.ID,
			NumberOfSeatsAllocated: 25 * i,
			AllocationDate:        time.Now(),
			AllocatedBy:           "customer-user",
			Status:                models.AllocationStatusActive,
		}
		allocationServiceRepo.Create(allocationServiceTestCtx, allocation)
	}

	utilization, err := allocationService.GetLicenseUtilization(allocationServiceTestCtx, license.LicenseID)
	if err != nil {
		t.Fatalf("Failed to get license utilization: %v", err)
	}

	if utilization["total_seats"] != 100 {
		t.Errorf("Total seats mismatch: got %v, want 100", utilization["total_seats"])
	}
	if utilization["allocated_seats"] != 75 {
		t.Errorf("Allocated seats mismatch: got %v, want 75", utilization["allocated_seats"])
	}
	if utilization["available_seats"] != 25 {
		t.Errorf("Available seats mismatch: got %v, want 25", utilization["available_seats"])
	}

	utilizationPercent := utilization["utilization_percent"].(float64)
	if utilizationPercent != 75.0 {
		t.Errorf("Utilization percent mismatch: got %v, want 75.0", utilizationPercent)
	}

	t.Logf("License utilization: %+v", utilization)
}

func TestLicenseAllocationService_ValidateAllocation(t *testing.T) {
	setupAllocationServiceTestDB(t)
	defer teardownAllocationServiceTestDB(t)

	customer := &models.Customer{
		CustomerID:   "alloc-svc-validate-customer",
		Name:         "Validate Customer",
		Email:        "validate@example.com",
		AccountStatus: models.CustomerStatusActive,
	}
	allocationServiceCustomerRepo.Create(allocationServiceTestCtx, customer)

	subscription := &models.Subscription{
		SubscriptionID: "SUB-ALLOC-VALIDATE",
		CustomerID:     customer.ID,
		StartDate:      time.Now(),
		Status:         models.SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	allocationServiceSubscriptionRepo.Create(allocationServiceTestCtx, subscription)

	license := &models.License{
		LicenseID:      "LIC-ALLOC-VALIDATE",
		SubscriptionID: subscription.ID,
		ProductID:      "PROD-VALIDATE",
		LicenseType:    models.LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         models.LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	allocationServiceLicenseRepo.Create(allocationServiceTestCtx, license)

	// Test valid allocation
	err := allocationService.ValidateAllocation(allocationServiceTestCtx, license.LicenseID, 50, "PROD-VALIDATE")
	if err != nil {
		t.Fatalf("Validation should pass: %v", err)
	}

	// Test invalid product
	err = allocationService.ValidateAllocation(allocationServiceTestCtx, license.LicenseID, 50, "PROD-WRONG")
	if err == nil {
		t.Error("Expected error for wrong product")
	}

	// Test insufficient seats
	err = allocationService.ValidateAllocation(allocationServiceTestCtx, license.LicenseID, 150, "PROD-VALIDATE")
	if err == nil {
		t.Error("Expected error for insufficient seats")
	}

	t.Logf("Validation tests passed")
}

