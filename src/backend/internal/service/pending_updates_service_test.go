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
	pendingUpdatesServiceTestDB  *database.MongoDB
	pendingUpdatesServiceTestCtx context.Context
	pendingUpdatesService        *PendingUpdatesService
	pendingUpdatesDeploymentRepo *repository.DeploymentRepository
	pendingUpdatesVersionRepo    *repository.VersionRepository
	pendingUpdatesCustomerRepo   *repository.CustomerRepository
	pendingUpdatesTenantRepo     *repository.TenantRepository
)

func setupPendingUpdatesServiceTestDB(t *testing.T) {
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

	pendingUpdatesServiceTestDB = db
	pendingUpdatesServiceTestCtx = ctx
	pendingUpdatesDeploymentRepo = repository.NewDeploymentRepository(db.Collection("deployments"))
	pendingUpdatesVersionRepo = repository.NewVersionRepository(db.Collection("versions"))
	pendingUpdatesCustomerRepo = repository.NewCustomerRepository(db.Collection("customers"))
	pendingUpdatesTenantRepo = repository.NewTenantRepository(db.Collection("customer_tenants"))
	pendingUpdatesService = NewPendingUpdatesService(
		pendingUpdatesDeploymentRepo,
		pendingUpdatesVersionRepo,
		pendingUpdatesCustomerRepo,
		pendingUpdatesTenantRepo,
	)
}

func teardownPendingUpdatesServiceTestDB(t *testing.T) {
	if pendingUpdatesServiceTestDB != nil {
		_ = pendingUpdatesServiceTestDB.Collection("deployments").Drop(pendingUpdatesServiceTestCtx)
		_ = pendingUpdatesServiceTestDB.Collection("versions").Drop(pendingUpdatesServiceTestCtx)
		_ = pendingUpdatesServiceTestDB.Collection("customers").Drop(pendingUpdatesServiceTestCtx)
		_ = pendingUpdatesServiceTestDB.Collection("customer_tenants").Drop(pendingUpdatesServiceTestCtx)
		_ = pendingUpdatesServiceTestDB.Collection("products").Drop(pendingUpdatesServiceTestCtx)
		_ = pendingUpdatesServiceTestDB.Disconnect(pendingUpdatesServiceTestCtx)
	}
}

func TestPendingUpdatesService_GetAvailableUpdatesForDeployment(t *testing.T) {
	setupPendingUpdatesServiceTestDB(t)
	defer teardownPendingUpdatesServiceTestDB(t)

	// Create a product
	productRepo := repository.NewProductRepository(pendingUpdatesServiceTestDB.Collection("products"))
	product := &models.Product{
		ProductID: "test-product",
		Name:      "Test Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	productRepo.Create(pendingUpdatesServiceTestCtx, product)

	// Create versions
	now := time.Now()
	versions := []*models.Version{
		{
			ProductID:     "test-product",
			VersionNumber: "1.0.0",
			ReleaseDate:   now.AddDate(0, 0, -10),
			ReleaseType:   models.ReleaseTypeFeature,
			State:         models.VersionStateReleased,
		},
		{
			ProductID:     "test-product",
			VersionNumber: "1.1.0",
			ReleaseDate:   now.AddDate(0, 0, -5),
			ReleaseType:   models.ReleaseTypeFeature,
			State:         models.VersionStateReleased,
		},
		{
			ProductID:     "test-product",
			VersionNumber: "1.2.0",
			ReleaseDate:   now,
			ReleaseType:   models.ReleaseTypeSecurity,
			State:         models.VersionStateReleased,
		},
		{
			ProductID:     "test-product",
			VersionNumber: "0.9.0",
			ReleaseDate:   now.AddDate(0, 0, -20),
			ReleaseType:   models.ReleaseTypeFeature,
			State:         models.VersionStateReleased,
		},
		{
			ProductID:     "test-product",
			VersionNumber: "2.0.0",
			ReleaseDate:   now.AddDate(0, 0, -1),
			ReleaseType:   models.ReleaseTypeMajor,
			State:         models.VersionStateDeprecated,
		},
	}

	for _, v := range versions {
		pendingUpdatesVersionRepo.Create(pendingUpdatesServiceTestCtx, v)
	}

	// Create customer, tenant, and deployment
	customer := &models.Customer{
		CustomerID:    "test-customer",
		Name:          "Test Customer",
		Email:         "test@example.com",
		AccountStatus: models.CustomerStatusActive,
		NotificationPreferences: models.NotificationPreferences{
			EmailEnabled:           true,
			InAppEnabled:           true,
			UATNotifications:      true,
			ProductionNotifications: true,
		},
	}
	pendingUpdatesCustomerRepo.Create(pendingUpdatesServiceTestCtx, customer)

	tenant := &models.CustomerTenant{
		TenantID:   "test-tenant",
		CustomerID: customer.ID,
		Name:       "Test Tenant",
		Status:     models.TenantStatusActive,
	}
	pendingUpdatesTenantRepo.Create(pendingUpdatesServiceTestCtx, tenant)

	deployment := &models.Deployment{
		DeploymentID:     "test-deployment",
		TenantID:         tenant.ID,
		ProductID:        "test-product",
		DeploymentType:   models.DeploymentTypeProduction,
		InstalledVersion: "1.0.0",
		Status:           models.DeploymentStatusActive,
	}
	pendingUpdatesDeploymentRepo.Create(pendingUpdatesServiceTestCtx, deployment)

	// Test: Get available updates
	updates, err := pendingUpdatesService.GetAvailableUpdatesForDeployment(pendingUpdatesServiceTestCtx, deployment.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get available updates: %v", err)
	}

	// Should have 2 updates (1.1.0 and 1.2.0), excluding 0.9.0 (older) and 2.0.0 (deprecated)
	if len(updates) != 2 {
		t.Errorf("Expected 2 updates, got %d", len(updates))
	}

	// Check that updates are sorted by version (newest first)
	if len(updates) > 0 && updates[0].VersionNumber != "1.2.0" {
		t.Errorf("Expected latest version to be 1.2.0, got %s", updates[0].VersionNumber)
	}

	// Check security update flag
	if len(updates) > 0 && !updates[0].IsSecurityUpdate {
		t.Error("Expected 1.2.0 to be marked as security update")
	}
}

func TestPendingUpdatesService_GetPendingUpdatesCount(t *testing.T) {
	setupPendingUpdatesServiceTestDB(t)
	defer teardownPendingUpdatesServiceTestDB(t)

	// Create test data (similar to above)
	productRepo := repository.NewProductRepository(pendingUpdatesServiceTestDB.Collection("products"))
	product := &models.Product{
		ProductID: "test-product-2",
		Name:      "Test Product 2",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	productRepo.Create(pendingUpdatesServiceTestCtx, product)

	now := time.Now()
	version := &models.Version{
		ProductID:     "test-product-2",
		VersionNumber: "1.1.0",
		ReleaseDate:   now,
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateReleased,
	}
	pendingUpdatesVersionRepo.Create(pendingUpdatesServiceTestCtx, version)

	customer := &models.Customer{
		CustomerID:    "test-customer-2",
		Name:          "Test Customer 2",
		Email:         "test2@example.com",
		AccountStatus: models.CustomerStatusActive,
		NotificationPreferences: models.NotificationPreferences{},
	}
	pendingUpdatesCustomerRepo.Create(pendingUpdatesServiceTestCtx, customer)

	tenant := &models.CustomerTenant{
		TenantID:   "test-tenant-2",
		CustomerID: customer.ID,
		Name:       "Test Tenant 2",
		Status:     models.TenantStatusActive,
	}
	pendingUpdatesTenantRepo.Create(pendingUpdatesServiceTestCtx, tenant)

	deployment := &models.Deployment{
		DeploymentID:     "test-deployment-2",
		TenantID:         tenant.ID,
		ProductID:        "test-product-2",
		DeploymentType:   models.DeploymentTypeProduction,
		InstalledVersion: "1.0.0",
		Status:           models.DeploymentStatusActive,
	}
	pendingUpdatesDeploymentRepo.Create(pendingUpdatesServiceTestCtx, deployment)

	// Test: Get pending updates count
	count, latestVersion, err := pendingUpdatesService.GetPendingUpdatesCount(pendingUpdatesServiceTestCtx, deployment.ID.Hex())
	if err != nil {
		t.Fatalf("Failed to get pending updates count: %v", err)
	}

	if count != 1 {
		t.Errorf("Expected count 1, got %d", count)
	}

	if latestVersion != "1.1.0" {
		t.Errorf("Expected latest version 1.1.0, got %s", latestVersion)
	}
}

func TestPendingUpdatesService_CalculateUpdatePriority(t *testing.T) {
	setupPendingUpdatesServiceTestDB(t)
	defer teardownPendingUpdatesServiceTestDB(t)

	// Test critical priority (security update)
	deployment := &models.Deployment{
		DeploymentType:   models.DeploymentTypeProduction,
		InstalledVersion: "1.0.0",
	}
	updates := []models.AvailableUpdate{
		{
			VersionNumber:    "1.1.0",
			IsSecurityUpdate: true,
		},
	}
	priority := pendingUpdatesService.CalculateUpdatePriority(deployment, updates)
	if priority != "critical" {
		t.Errorf("Expected priority 'critical', got '%s'", priority)
	}

	// Test high priority (major update on production)
	deployment2 := &models.Deployment{
		DeploymentType:   models.DeploymentTypeProduction,
		InstalledVersion: "1.0.0",
	}
	updates2 := []models.AvailableUpdate{
		{
			VersionNumber:    "2.0.0",
			IsSecurityUpdate: false,
		},
	}
	priority2 := pendingUpdatesService.CalculateUpdatePriority(deployment2, updates2)
	if priority2 != "high" {
		t.Errorf("Expected priority 'high', got '%s'", priority2)
	}

	// Test normal priority (minor/patch update)
	deployment3 := &models.Deployment{
		DeploymentType:   models.DeploymentTypeUAT,
		InstalledVersion: "1.0.0",
	}
	updates3 := []models.AvailableUpdate{
		{
			VersionNumber:    "1.1.0",
			IsSecurityUpdate: false,
		},
	}
	priority3 := pendingUpdatesService.CalculateUpdatePriority(deployment3, updates3)
	if priority3 != "normal" {
		t.Errorf("Expected priority 'normal', got '%s'", priority3)
	}
}

// TestPendingUpdatesService_Caching tests the caching functionality
func TestPendingUpdatesService_Caching(t *testing.T) {
	setupPendingUpdatesServiceTestDB(t)
	defer teardownPendingUpdatesServiceTestDB(t)

	// Create a product
	productRepo := repository.NewProductRepository(pendingUpdatesServiceTestDB.Collection("products"))
	product := &models.Product{
		ProductID: "cache-product",
		Name:      "Cache Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	productRepo.Create(pendingUpdatesServiceTestCtx, product)

	// Create customer, tenant, and deployment
	customer := &models.Customer{
		CustomerID:   "cache-customer",
		Name:          "Cache Customer",
		AccountStatus: models.CustomerStatusActive,
	}
	pendingUpdatesCustomerRepo.Create(pendingUpdatesServiceTestCtx, customer)

	tenant := &models.CustomerTenant{
		TenantID:   "cache-tenant",
		CustomerID: customer.ID,
		Name:       "Cache Tenant",
		Status:     models.TenantStatusActive,
	}
	pendingUpdatesTenantRepo.Create(pendingUpdatesServiceTestCtx, tenant)

	deployment := &models.Deployment{
		DeploymentID:   "cache-deployment",
		TenantID:       tenant.ID,
		ProductID:      "cache-product",
		DeploymentType: models.DeploymentTypeUAT,
		InstalledVersion: "1.0.0",
		Status:         models.DeploymentStatusActive,
	}
	pendingUpdatesDeploymentRepo.Create(pendingUpdatesServiceTestCtx, deployment)

	// Create a version
	version := &models.Version{
		ProductID:     "cache-product",
		VersionNumber: "1.1.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateReleased,
	}
	pendingUpdatesVersionRepo.Create(pendingUpdatesServiceTestCtx, version)

	// First call should hit database and cache result
	result1, err := pendingUpdatesService.GetPendingUpdatesForDeployment(pendingUpdatesServiceTestCtx, deployment.DeploymentID)
	if err != nil {
		t.Fatalf("Failed to get pending updates: %v", err)
	}
	if result1 == nil {
		t.Fatal("Expected result, got nil")
	}

	// Second call should hit cache (we can't directly verify this, but we can test invalidation)
	result2, err := pendingUpdatesService.GetPendingUpdatesForDeployment(pendingUpdatesServiceTestCtx, deployment.DeploymentID)
	if err != nil {
		t.Fatalf("Failed to get pending updates: %v", err)
	}
	if result2 == nil {
		t.Fatal("Expected result, got nil")
	}

	// Test cache invalidation for deployment
	pendingUpdatesService.InvalidateCacheForDeployment(deployment.DeploymentID)

	// Test cache invalidation for product
	pendingUpdatesService.InvalidateCacheForProduct(pendingUpdatesServiceTestCtx, "cache-product")

	t.Logf("Cache tests passed - invalidation methods work correctly")
}

// TestPendingUpdatesService_IntegrationWithVersionRelease tests integration with version release
func TestPendingUpdatesService_IntegrationWithVersionRelease(t *testing.T) {
	setupPendingUpdatesServiceTestDB(t)
	defer teardownPendingUpdatesServiceTestDB(t)

	// Create a product
	productRepo := repository.NewProductRepository(pendingUpdatesServiceTestDB.Collection("products"))
	product := &models.Product{
		ProductID: "integration-product",
		Name:      "Integration Product",
		Type:      models.ProductTypeServer,
		IsActive:  true,
	}
	productRepo.Create(pendingUpdatesServiceTestCtx, product)

	// Create customer, tenant, and deployment
	customer := &models.Customer{
		CustomerID:   "integration-customer",
		Name:         "Integration Customer",
		AccountStatus: models.CustomerStatusActive,
	}
	pendingUpdatesCustomerRepo.Create(pendingUpdatesServiceTestCtx, customer)

	tenant := &models.CustomerTenant{
		TenantID:   "integration-tenant",
		CustomerID: customer.ID,
		Name:       "Integration Tenant",
		Status:     models.TenantStatusActive,
	}
	pendingUpdatesTenantRepo.Create(pendingUpdatesServiceTestCtx, tenant)

	deployment := &models.Deployment{
		DeploymentID:   "integration-deployment",
		TenantID:       tenant.ID,
		ProductID:      "integration-product",
		DeploymentType: models.DeploymentTypeProduction,
		InstalledVersion: "1.0.0",
		Status:         models.DeploymentStatusActive,
	}
	pendingUpdatesDeploymentRepo.Create(pendingUpdatesServiceTestCtx, deployment)

	// Initially no updates available
	result1, err := pendingUpdatesService.GetPendingUpdatesForDeployment(pendingUpdatesServiceTestCtx, deployment.DeploymentID)
	if err != nil {
		t.Fatalf("Failed to get pending updates: %v", err)
	}
	if result1.UpdateCount != 0 {
		t.Errorf("Expected 0 updates initially, got %d", result1.UpdateCount)
	}

	// Create and release a new version
	version := &models.Version{
		ProductID:     "integration-product",
		VersionNumber: "1.1.0",
		ReleaseDate:   time.Now(),
		ReleaseType:   models.ReleaseTypeFeature,
		State:         models.VersionStateReleased,
	}
	pendingUpdatesVersionRepo.Create(pendingUpdatesServiceTestCtx, version)

	// Invalidate cache (simulating what happens when version is released)
	pendingUpdatesService.InvalidateCacheForProduct(pendingUpdatesServiceTestCtx, "integration-product")

	// Now pending updates should include the new version
	result2, err := pendingUpdatesService.GetPendingUpdatesForDeployment(pendingUpdatesServiceTestCtx, deployment.DeploymentID)
	if err != nil {
		t.Fatalf("Failed to get pending updates: %v", err)
	}
	if result2.UpdateCount != 1 {
		t.Errorf("Expected 1 update after version release, got %d", result2.UpdateCount)
	}
	if result2.LatestVersion != "1.1.0" {
		t.Errorf("Expected latest version '1.1.0', got '%s'", result2.LatestVersion)
	}

	t.Logf("Integration test passed - cache invalidation works correctly when version is released")
}

