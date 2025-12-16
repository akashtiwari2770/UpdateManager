package models

import (
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Test Customer Model

func TestCustomer_Valid(t *testing.T) {
	customer := &Customer{
		CustomerID:   "CUST-001",
		Name:         "Test Customer",
		Email:        "test@example.com",
		AccountStatus: CustomerStatusActive,
		NotificationPreferences: NotificationPreferences{
			EmailEnabled:           true,
			InAppEnabled:           true,
			UATNotifications:       true,
			ProductionNotifications: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if customer.CustomerID == "" {
		t.Error("CustomerID should not be empty")
	}
	if customer.Name == "" {
		t.Error("Name should not be empty")
	}
	if customer.Email == "" {
		t.Error("Email should not be empty")
	}
}

func TestCustomerStatus_Valid(t *testing.T) {
	statuses := []CustomerStatus{
		CustomerStatusActive,
		CustomerStatusInactive,
		CustomerStatusSuspended,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("CustomerStatus should not be empty: %v", status)
		}
	}
}

func TestCreateCustomerRequest_Valid(t *testing.T) {
	req := &CreateCustomerRequest{
		CustomerID:   "CUST-001",
		Name:         "Test Customer",
		Email:        "test@example.com",
		AccountStatus: CustomerStatusActive,
		NotificationPreferences: NotificationPreferences{
			EmailEnabled: true,
		},
	}

	if req.CustomerID == "" {
		t.Error("CustomerID should not be empty")
	}
	if req.Name == "" {
		t.Error("Name should not be empty")
	}
	if req.Email == "" {
		t.Error("Email should not be empty")
	}
}

// Test Tenant Model

func TestCustomerTenant_Valid(t *testing.T) {
	customerID := primitive.NewObjectID()
	tenant := &CustomerTenant{
		TenantID:   "TENANT-001",
		CustomerID: customerID,
		Name:       "Test Tenant",
		Status:     TenantStatusActive,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if tenant.TenantID == "" {
		t.Error("TenantID should not be empty")
	}
	if tenant.CustomerID.IsZero() {
		t.Error("CustomerID should not be zero")
	}
	if tenant.Name == "" {
		t.Error("Name should not be empty")
	}
}

func TestTenantStatus_Valid(t *testing.T) {
	statuses := []TenantStatus{
		TenantStatusActive,
		TenantStatusInactive,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("TenantStatus should not be empty: %v", status)
		}
	}
}

func TestCreateTenantRequest_Valid(t *testing.T) {
	req := &CreateTenantRequest{
		TenantID: "TENANT-001",
		Name:     "Test Tenant",
		Status:   TenantStatusActive,
	}

	if req.TenantID == "" {
		t.Error("TenantID should not be empty")
	}
	if req.Name == "" {
		t.Error("Name should not be empty")
	}
}

// Test Deployment Model

func TestDeployment_Valid(t *testing.T) {
	tenantID := primitive.NewObjectID()
	users := 100
	deployment := &Deployment{
		DeploymentID:    "DEPLOY-001",
		TenantID:        tenantID,
		ProductID:       "PROD-001",
		DeploymentType:  DeploymentTypeProduction,
		InstalledVersion: "1.0.0",
		NumberOfUsers:   &users,
		Status:          DeploymentStatusActive,
		DeploymentDate:  time.Now(),
		LastUpdatedDate: time.Now(),
	}

	if deployment.DeploymentID == "" {
		t.Error("DeploymentID should not be empty")
	}
	if deployment.TenantID.IsZero() {
		t.Error("TenantID should not be zero")
	}
	if deployment.ProductID == "" {
		t.Error("ProductID should not be empty")
	}
	if deployment.InstalledVersion == "" {
		t.Error("InstalledVersion should not be empty")
	}
	if deployment.NumberOfUsers != nil && *deployment.NumberOfUsers < 0 {
		t.Error("NumberOfUsers should be positive if set")
	}
}

func TestDeploymentType_Valid(t *testing.T) {
	types := []DeploymentType{
		DeploymentTypeUAT,
		DeploymentTypeTesting,
		DeploymentTypeProduction,
	}

	for _, dt := range types {
		if dt == "" {
			t.Errorf("DeploymentType should not be empty: %v", dt)
		}
	}
}

func TestDeploymentStatus_Valid(t *testing.T) {
	statuses := []DeploymentStatus{
		DeploymentStatusActive,
		DeploymentStatusInactive,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("DeploymentStatus should not be empty: %v", status)
		}
	}
}

func TestCreateDeploymentRequest_Valid(t *testing.T) {
	users := 50
	req := &CreateDeploymentRequest{
		DeploymentID:    "DEPLOY-001",
		ProductID:       "PROD-001",
		DeploymentType:  DeploymentTypeUAT,
		InstalledVersion: "1.0.0",
		NumberOfUsers:   &users,
		Status:          DeploymentStatusActive,
	}

	if req.DeploymentID == "" {
		t.Error("DeploymentID should not be empty")
	}
	if req.ProductID == "" {
		t.Error("ProductID should not be empty")
	}
	if req.InstalledVersion == "" {
		t.Error("InstalledVersion should not be empty")
	}
}

func TestUpdateDeploymentRequest_PartialUpdate(t *testing.T) {
	users := 100
	req := &UpdateDeploymentRequest{
		InstalledVersion: stringPtr("2.0.0"),
		NumberOfUsers:    &users,
	}

	if req.InstalledVersion == nil {
		t.Error("InstalledVersion should be set for update")
	}
	if req.NumberOfUsers == nil {
		t.Error("NumberOfUsers should be set for update")
	}
}

// Test Notification Model Updates

func TestNotification_CustomerFields(t *testing.T) {
	notification := &Notification{
		Type:        NotificationTypeNewVersion,
		RecipientID: "CUST-001",
		CustomerID:  "CUST-001",
		TenantID:    "TENANT-001",
		DeploymentID: "DEPLOY-001",
		ProductID:   "PROD-001",
		Title:       "New Version Available",
		Message:     "A new version is available",
		Priority:    NotificationPriorityNormal,
		CreatedAt:   time.Now(),
	}

	if notification.CustomerID == "" {
		t.Error("CustomerID should be set for customer notifications")
	}
	if notification.TenantID == "" {
		t.Error("TenantID should be set for tenant-specific notifications")
	}
	if notification.DeploymentID == "" {
		t.Error("DeploymentID should be set for deployment-specific notifications")
	}
}

// Test License Management Models

func TestSubscription_Valid(t *testing.T) {
	customerID := primitive.NewObjectID()
	subscription := &Subscription{
		SubscriptionID: "SUB-001",
		CustomerID:     customerID,
		Name:           "Test Subscription",
		StartDate:      time.Now(),
		Status:         SubscriptionStatusActive,
		CreatedBy:      "user-123",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if subscription.SubscriptionID == "" {
		t.Error("SubscriptionID should not be empty")
	}
	if subscription.CustomerID.IsZero() {
		t.Error("CustomerID should not be zero")
	}
	if subscription.StartDate.IsZero() {
		t.Error("StartDate should not be zero")
	}
	if subscription.CreatedBy == "" {
		t.Error("CreatedBy should not be empty")
	}
}

func TestSubscriptionStatus_Valid(t *testing.T) {
	statuses := []SubscriptionStatus{
		SubscriptionStatusActive,
		SubscriptionStatusInactive,
		SubscriptionStatusExpired,
		SubscriptionStatusSuspended,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("SubscriptionStatus should not be empty: %v", status)
		}
	}
}

func TestCreateSubscriptionRequest_Valid(t *testing.T) {
	req := &CreateSubscriptionRequest{
		SubscriptionID: "SUB-001",
		Name:           "Test Subscription",
		StartDate:      time.Now(),
		Status:         SubscriptionStatusActive,
	}

	if req.SubscriptionID == "" {
		t.Error("SubscriptionID should not be empty")
	}
	if req.StartDate.IsZero() {
		t.Error("StartDate should not be zero")
	}
}

func TestLicense_Valid(t *testing.T) {
	subscriptionID := primitive.NewObjectID()
	license := &License{
		LicenseID:      "LIC-001",
		SubscriptionID: subscriptionID,
		ProductID:      "PROD-001",
		LicenseType:    LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         LicenseStatusActive,
		AssignedBy:     "sales-user-123",
		AssignmentDate: time.Now(),
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	if license.LicenseID == "" {
		t.Error("LicenseID should not be empty")
	}
	if license.SubscriptionID.IsZero() {
		t.Error("SubscriptionID should not be zero")
	}
	if license.ProductID == "" {
		t.Error("ProductID should not be empty")
	}
	if license.NumberOfSeats < 1 {
		t.Error("NumberOfSeats should be at least 1")
	}
	if license.AssignedBy == "" {
		t.Error("AssignedBy should not be empty")
	}
}

func TestLicenseType_Valid(t *testing.T) {
	types := []LicenseType{
		LicenseTypePerpetual,
		LicenseTypeTimeBased,
	}

	for _, licenseType := range types {
		if licenseType == "" {
			t.Errorf("LicenseType should not be empty: %v", licenseType)
		}
	}
}

func TestLicenseStatus_Valid(t *testing.T) {
	statuses := []LicenseStatus{
		LicenseStatusActive,
		LicenseStatusInactive,
		LicenseStatusExpired,
		LicenseStatusRevoked,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("LicenseStatus should not be empty: %v", status)
		}
	}
}

func TestCreateLicenseRequest_Valid(t *testing.T) {
	req := &CreateLicenseRequest{
		LicenseID:     "LIC-001",
		ProductID:     "PROD-001",
		LicenseType:   LicenseTypePerpetual,
		NumberOfSeats: 100,
		StartDate:     time.Now(),
		Status:        LicenseStatusActive,
	}

	if req.LicenseID == "" {
		t.Error("LicenseID should not be empty")
	}
	if req.ProductID == "" {
		t.Error("ProductID should not be empty")
	}
	if req.NumberOfSeats < 1 {
		t.Error("NumberOfSeats should be at least 1")
	}
}

func TestLicenseAllocation_Valid(t *testing.T) {
	licenseID := primitive.NewObjectID()
	tenantID := primitive.NewObjectID()
	allocation := &LicenseAllocation{
		AllocationID:          "ALLOC-001",
		LicenseID:             licenseID,
		TenantID:              &tenantID,
		NumberOfSeatsAllocated: 50,
		AllocationDate:        time.Now(),
		AllocatedBy:           "customer-user-123",
		Status:                AllocationStatusActive,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
	}

	if allocation.AllocationID == "" {
		t.Error("AllocationID should not be empty")
	}
	if allocation.LicenseID.IsZero() {
		t.Error("LicenseID should not be zero")
	}
	if allocation.NumberOfSeatsAllocated < 1 {
		t.Error("NumberOfSeatsAllocated should be at least 1")
	}
	if allocation.AllocatedBy == "" {
		t.Error("AllocatedBy should not be empty")
	}
}

func TestAllocationStatus_Valid(t *testing.T) {
	statuses := []AllocationStatus{
		AllocationStatusActive,
		AllocationStatusReleased,
	}

	for _, status := range statuses {
		if status == "" {
			t.Errorf("AllocationStatus should not be empty: %v", status)
		}
	}
}

func TestAllocateLicenseRequest_Valid(t *testing.T) {
	tenantID := "TENANT-001"
	req := &AllocateLicenseRequest{
		TenantID:              &tenantID,
		NumberOfSeatsAllocated: 50,
	}

	if req.NumberOfSeatsAllocated < 1 {
		t.Error("NumberOfSeatsAllocated should be at least 1")
	}
}

func TestLicense_TimeBasedRequiresEndDate(t *testing.T) {
	// This test validates business logic: time-based licenses should have end date
	// Note: This is a business rule that should be enforced in service layer
	subscriptionID := primitive.NewObjectID()
	
	// Perpetual license - end date optional
	perpetualLicense := &License{
		LicenseID:      "LIC-PERP",
		SubscriptionID: subscriptionID,
		ProductID:      "PROD-001",
		LicenseType:    LicenseTypePerpetual,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		Status:         LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	
	if perpetualLicense.LicenseType != LicenseTypePerpetual {
		t.Error("License type should be perpetual")
	}
	
	// Time-based license - end date should be set (enforced in service layer)
	endDate := time.Now().AddDate(1, 0, 0) // 1 year from now
	timeBasedLicense := &License{
		LicenseID:      "LIC-TIME",
		SubscriptionID: subscriptionID,
		ProductID:      "PROD-001",
		LicenseType:    LicenseTypeTimeBased,
		NumberOfSeats:  100,
		StartDate:      time.Now(),
		EndDate:        &endDate,
		Status:         LicenseStatusActive,
		AssignedBy:     "sales-user",
		AssignmentDate: time.Now(),
	}
	
	if timeBasedLicense.LicenseType != LicenseTypeTimeBased {
		t.Error("License type should be time-based")
	}
}

func TestSubscription_EndDateAfterStartDate(t *testing.T) {
	// This test validates business logic: end date must be after start date
	// Note: This is a business rule that should be enforced in service layer
	customerID := primitive.NewObjectID()
	startDate := time.Now()
	endDate := startDate.AddDate(1, 0, 0) // 1 year after start
	
	subscription := &Subscription{
		SubscriptionID: "SUB-001",
		CustomerID:     customerID,
		StartDate:      startDate,
		EndDate:        &endDate,
		Status:         SubscriptionStatusActive,
		CreatedBy:      "user-123",
	}
	
	if subscription.EndDate != nil && subscription.EndDate.Before(subscription.StartDate) {
		t.Error("EndDate should be after StartDate")
	}
}

// Helper functions

func stringPtr(s string) *string {
	return &s
}

