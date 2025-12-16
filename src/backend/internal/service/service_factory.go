package service

import (
	"go.mongodb.org/mongo-driver/mongo"

	"updatemanager/internal/repository"
)

// ServiceFactory holds all services
type ServiceFactory struct {
	ProductService            *ProductService
	VersionService            *VersionService
	CompatibilityService      *CompatibilityService
	UpgradePathService        *UpgradePathService
	NotificationService       *NotificationService
	UpdateDetectionService    *UpdateDetectionService
	UpdateRolloutService      *UpdateRolloutService
	AuditLogService           *AuditLogService
	CustomerService          *CustomerService
	TenantService            *TenantService
	DeploymentService        *DeploymentService
	PendingUpdatesService     *PendingUpdatesService
	SubscriptionService      *SubscriptionService
	LicenseService            *LicenseService
	LicenseAllocationService  *LicenseAllocationService
}

// NewServiceFactory creates all services with their dependencies
func NewServiceFactory(db *mongo.Database) *ServiceFactory {
	// Initialize repositories
	productRepo := repository.NewProductRepository(db.Collection("products"))
	versionRepo := repository.NewVersionRepository(db.Collection("versions"))
	compatibilityRepo := repository.NewCompatibilityRepository(db.Collection("compatibility_matrices"))
	upgradePathRepo := repository.NewUpgradePathRepository(db.Collection("upgrade_paths"))
	notificationRepo := repository.NewNotificationRepository(db.Collection("notifications"))
	detectionRepo := repository.NewUpdateDetectionRepository(db.Collection("update_detections"))
	rolloutRepo := repository.NewUpdateRolloutRepository(db.Collection("update_rollouts"))
	auditRepo := repository.NewAuditLogRepository(db.Collection("audit_logs"))
	customerRepo := repository.NewCustomerRepository(db.Collection("customers"))
	tenantRepo := repository.NewTenantRepository(db.Collection("customer_tenants"))
	deploymentRepo := repository.NewDeploymentRepository(db.Collection("deployments"))
	subscriptionRepo := repository.NewSubscriptionRepository(db.Collection("subscriptions"))
	licenseRepo := repository.NewLicenseRepository(db.Collection("licenses"))
	allocationRepo := repository.NewLicenseAllocationRepository(db.Collection("license_allocations"))

	// Initialize services
	productService := NewProductService(productRepo, auditRepo)
	versionService := NewVersionService(versionRepo, productRepo, auditRepo)
	compatibilityService := NewCompatibilityService(compatibilityRepo, versionRepo, auditRepo)
	upgradePathService := NewUpgradePathService(upgradePathRepo, versionRepo)
	notificationService := NewNotificationService(notificationRepo)
	detectionService := NewUpdateDetectionService(detectionRepo, versionRepo, productRepo)
	rolloutService := NewUpdateRolloutService(rolloutRepo, detectionRepo, versionRepo, productRepo)
	auditLogService := NewAuditLogService(auditRepo)
	customerService := NewCustomerService(customerRepo, tenantRepo, deploymentRepo, auditRepo)
	tenantService := NewTenantService(tenantRepo, customerRepo, deploymentRepo, auditRepo)
	deploymentService := NewDeploymentService(deploymentRepo, tenantRepo, customerRepo, productService, versionService, auditRepo)
	pendingUpdatesService := NewPendingUpdatesService(deploymentRepo, versionRepo, customerRepo, tenantRepo)
	subscriptionService := NewSubscriptionService(subscriptionRepo, customerRepo, licenseRepo, auditRepo)
	licenseService := NewLicenseService(licenseRepo, subscriptionRepo, customerRepo, allocationRepo, auditRepo)
	licenseAllocationService := NewLicenseAllocationService(allocationRepo, licenseRepo, subscriptionRepo, customerRepo, tenantRepo, deploymentRepo, auditRepo)

	return &ServiceFactory{
		ProductService:           productService,
		VersionService:           versionService,
		CompatibilityService:     compatibilityService,
		UpgradePathService:       upgradePathService,
		NotificationService:      notificationService,
		UpdateDetectionService:   detectionService,
		UpdateRolloutService:     rolloutService,
		AuditLogService:          auditLogService,
		CustomerService:          customerService,
		TenantService:            tenantService,
		DeploymentService:        deploymentService,
		PendingUpdatesService:    pendingUpdatesService,
		SubscriptionService:      subscriptionService,
		LicenseService:           licenseService,
		LicenseAllocationService: licenseAllocationService,
	}
}
