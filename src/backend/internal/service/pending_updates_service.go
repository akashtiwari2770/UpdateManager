package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
	"updatemanager/internal/utils"
)

// cacheEntry represents a cached pending updates result
type cacheEntry struct {
	data      interface{}
	expiresAt time.Time
}

// PendingUpdatesService handles pending updates business logic
type PendingUpdatesService struct {
	deploymentRepo *repository.DeploymentRepository
	versionRepo    *repository.VersionRepository
	customerRepo   *repository.CustomerRepository
	tenantRepo     *repository.TenantRepository
	
	// Simple in-memory cache for pending updates
	cache      map[string]*cacheEntry
	cacheMutex sync.RWMutex
	cacheTTL   time.Duration
}

// NewPendingUpdatesService creates a new pending updates service
func NewPendingUpdatesService(
	deploymentRepo *repository.DeploymentRepository,
	versionRepo *repository.VersionRepository,
	customerRepo *repository.CustomerRepository,
	tenantRepo *repository.TenantRepository,
) *PendingUpdatesService {
	return &PendingUpdatesService{
		deploymentRepo: deploymentRepo,
		versionRepo:    versionRepo,
		customerRepo:   customerRepo,
		tenantRepo:     tenantRepo,
		cache:          make(map[string]*cacheEntry),
		cacheTTL:       5 * time.Minute, // Cache for 5 minutes by default
	}
}

// InvalidateCacheForProduct invalidates cache for all deployments of a product
// This should be called when a new version is released for the product
func (s *PendingUpdatesService) InvalidateCacheForProduct(ctx context.Context, productID string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	
	// Remove all cache entries that contain this product ID
	for key := range s.cache {
		// Cache keys format: "deployment:{id}", "tenant:{customerId}:{tenantId}", etc.
		// For simplicity, we'll clear all cache entries when a product version is released
		// In production, you might want more granular invalidation
		delete(s.cache, key)
	}
}

// InvalidateCacheForDeployment invalidates cache for a specific deployment
func (s *PendingUpdatesService) InvalidateCacheForDeployment(deploymentID string) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	
	key := fmt.Sprintf("deployment:%s", deploymentID)
	delete(s.cache, key)
	
	// Also invalidate tenant and customer level caches that might include this deployment
	// For simplicity, we clear all tenant/customer caches
	for k := range s.cache {
		if len(k) > 7 && (k[:7] == "tenant:" || k[:9] == "customer:") {
			delete(s.cache, k)
		}
	}
}

// getCached retrieves a value from cache if it exists and is not expired
func (s *PendingUpdatesService) getCached(key string) (interface{}, bool) {
	s.cacheMutex.RLock()
	entry, exists := s.cache[key]
	if !exists {
		s.cacheMutex.RUnlock()
		return nil, false
	}
	
	expired := time.Now().After(entry.expiresAt)
	s.cacheMutex.RUnlock()
	
	if expired {
		// Entry expired, remove it
		s.cacheMutex.Lock()
		delete(s.cache, key)
		s.cacheMutex.Unlock()
		return nil, false
	}
	
	return entry.data, true
}

// setCached stores a value in cache
func (s *PendingUpdatesService) setCached(key string, data interface{}) {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	
	s.cache[key] = &cacheEntry{
		data:      data,
		expiresAt: time.Now().Add(s.cacheTTL),
	}
}

// GetAvailableUpdatesForDeployment retrieves available updates for a deployment
func (s *PendingUpdatesService) GetAvailableUpdatesForDeployment(ctx context.Context, deploymentID string) ([]models.AvailableUpdate, error) {
	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(deploymentID)
	if err != nil {
		// Try getting by deployment_id instead
		deployment, err2 := s.deploymentRepo.GetByDeploymentID(ctx, deploymentID)
		if err2 != nil {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
		return s.getAvailableUpdatesForDeployment(ctx, deployment)
	}

	// Get deployment
	deployment, err := s.deploymentRepo.GetByID(ctx, objectID)
	if err != nil {
		return nil, fmt.Errorf("deployment not found: %w", err)
	}

	return s.getAvailableUpdatesForDeployment(ctx, deployment)
}

// getAvailableUpdatesForDeployment is the internal method that does the actual work
func (s *PendingUpdatesService) getAvailableUpdatesForDeployment(ctx context.Context, deployment *models.Deployment) ([]models.AvailableUpdate, error) {

	// Get all versions for the product
	opts := options.Find()
	opts.SetSort(bson.M{"release_date": -1})
	versions, err := s.versionRepo.GetByProductID(ctx, deployment.ProductID, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get versions: %w", err)
	}

	// Filter versions that are newer than installed version and are released
	var availableUpdates []models.AvailableUpdate
	now := time.Now()

	for _, version := range versions {
		// Only consider released versions
		if version.State != models.VersionStateReleased {
			continue
		}

		// Exclude deprecated and EOL versions
		if version.State == models.VersionStateDeprecated || version.State == models.VersionStateEOL {
			continue
		}

		// Check if version has passed EOL date
		if version.EOLDate != nil && version.EOLDate.Before(now) {
			continue
		}

		// Check if version is newer than installed version
		if !utils.IsVersionNewer(version.VersionNumber, deployment.InstalledVersion) {
			continue
		}

		// Determine if this is a security update
		isSecurityUpdate := version.ReleaseType == models.ReleaseTypeSecurity

		// Build upgrade path (simplified - could be enhanced with compatibility matrix)
		upgradePath := []string{version.VersionNumber}

		availableUpdate := models.AvailableUpdate{
			VersionNumber:      version.VersionNumber,
			ReleaseDate:        version.ReleaseDate,
			ReleaseType:        string(version.ReleaseType),
			IsSecurityUpdate:   isSecurityUpdate,
			CompatibilityStatus: "compatible", // Default, could be enhanced
			UpgradePath:        upgradePath,
		}

		availableUpdates = append(availableUpdates, availableUpdate)
	}

	return availableUpdates, nil
}

// GetPendingUpdatesCount retrieves pending updates count for a deployment
func (s *PendingUpdatesService) GetPendingUpdatesCount(ctx context.Context, deploymentID string) (int, string, error) {
	updates, err := s.GetAvailableUpdatesForDeployment(ctx, deploymentID)
	if err != nil {
		return 0, "", err
	}

	count := len(updates)
	latestVersion := ""
	if count > 0 {
		latestVersion = updates[0].VersionNumber
	}

	return count, latestVersion, nil
}

// GetPendingUpdatesForDeployment retrieves full pending updates response for a deployment
func (s *PendingUpdatesService) GetPendingUpdatesForDeployment(ctx context.Context, deploymentID string) (*models.PendingUpdatesResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("deployment:%s", deploymentID)
	if cached, found := s.getCached(cacheKey); found {
		if result, ok := cached.(*models.PendingUpdatesResponse); ok {
			return result, nil
		}
	}

	// Convert string ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(deploymentID)
	var deployment *models.Deployment
	if err != nil {
		// Try getting by deployment_id instead
		deployment, err = s.deploymentRepo.GetByDeploymentID(ctx, deploymentID)
		if err != nil {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
	} else {
		deployment, err = s.deploymentRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("deployment not found: %w", err)
		}
	}

	// Get available updates
	availableUpdates, err := s.getAvailableUpdatesForDeployment(ctx, deployment)
	if err != nil {
		return nil, err
	}

	// Get tenant and customer for context
	tenant, err := s.tenantRepo.GetByID(ctx, deployment.TenantID)
	if err == nil {
		customer, err := s.customerRepo.GetByID(ctx, tenant.CustomerID)
		if err == nil {
			// Build response with context
			response := &models.PendingUpdatesResponse{
				DeploymentID:    deployment.DeploymentID,
				ProductID:       deployment.ProductID,
				CurrentVersion:   deployment.InstalledVersion,
				UpdateCount:      len(availableUpdates),
				AvailableUpdates: availableUpdates,
				TenantID:         tenant.TenantID,
				TenantName:       tenant.Name,
				CustomerID:       customer.CustomerID,
				CustomerName:     customer.Name,
				DeploymentType:   deployment.DeploymentType,
			}

			// Set latest version
			if len(availableUpdates) > 0 {
				response.LatestVersion = availableUpdates[0].VersionNumber
				response.VersionGapType = utils.GetVersionGapType(deployment.InstalledVersion, response.LatestVersion)
			}

			// Calculate priority
			response.Priority = s.CalculateUpdatePriority(deployment, availableUpdates)

			// Cache the result
			s.setCached(cacheKey, response)

			return response, nil
		}
	}

	// Fallback without tenant/customer context
	response := &models.PendingUpdatesResponse{
		DeploymentID:    deployment.DeploymentID,
		ProductID:       deployment.ProductID,
		CurrentVersion:   deployment.InstalledVersion,
		UpdateCount:      len(availableUpdates),
		AvailableUpdates: availableUpdates,
		DeploymentType:   deployment.DeploymentType,
	}

	if len(availableUpdates) > 0 {
		response.LatestVersion = availableUpdates[0].VersionNumber
		response.VersionGapType = utils.GetVersionGapType(deployment.InstalledVersion, response.LatestVersion)
	}

	response.Priority = s.CalculateUpdatePriority(deployment, availableUpdates)

	// Cache the result
	s.setCached(cacheKey, response)

	return response, nil
}

// GetPendingUpdatesForTenant retrieves pending updates for all deployments in a tenant
func (s *PendingUpdatesService) GetPendingUpdatesForTenant(ctx context.Context, customerID, tenantID string, filter *models.PendingUpdatesFilter) (*models.TenantPendingUpdatesSummary, error) {
	// Get tenant
	tenant, err := s.tenantRepo.GetByTenantID(ctx, tenantID)
	if err != nil {
		// Try as ObjectID
		objectID, parseErr := primitive.ObjectIDFromHex(tenantID)
		if parseErr != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
		tenant, err = s.tenantRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("tenant not found: %w", err)
		}
	}

	// Build deployment filter
	deploymentFilter := &repository.DeploymentFilter{}
	if filter != nil {
		if filter.ProductID != "" {
			deploymentFilter.ProductID = filter.ProductID
		}
		if filter.DeploymentType != "" {
			deploymentFilter.DeploymentType = filter.DeploymentType
		}
	}

	// Get all deployments for tenant
	deployments, _, err := s.deploymentRepo.GetByTenantID(ctx, tenant.ID, deploymentFilter, &repository.Pagination{Page: 1, Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	summary := &models.TenantPendingUpdatesSummary{
		TenantID:               tenant.TenantID,
		TenantName:             tenant.Name,
		TotalDeployments:       len(deployments),
		ByPriority:              make(map[string]int),
		ByProduct:              make(map[string]int),
		Deployments:            []models.PendingUpdatesResponse{},
	}

	// Process each deployment
	for _, deployment := range deployments {
		// Get pending updates for deployment
		pendingUpdates, err := s.GetPendingUpdatesForDeployment(ctx, deployment.ID.Hex())
		if err != nil {
			continue // Skip deployments with errors
		}

		// Apply priority filter if specified
		if filter != nil && filter.Priority != "" && pendingUpdates.Priority != filter.Priority {
			continue
		}

		if pendingUpdates.UpdateCount > 0 {
			summary.DeploymentsWithUpdates++
			summary.TotalPendingUpdateCount += pendingUpdates.UpdateCount
			summary.ByPriority[pendingUpdates.Priority]++
			summary.ByProduct[pendingUpdates.ProductID]++
		}

		summary.Deployments = append(summary.Deployments, *pendingUpdates)
	}

	return summary, nil
}

// GetPendingUpdatesForCustomer retrieves pending updates for all deployments of a customer
func (s *PendingUpdatesService) GetPendingUpdatesForCustomer(ctx context.Context, customerID string, filter *models.PendingUpdatesFilter) (*models.CustomerPendingUpdatesSummary, error) {
	// Get customer
	customer, err := s.customerRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		// Try as ObjectID
		objectID, parseErr := primitive.ObjectIDFromHex(customerID)
		if parseErr != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
		customer, err = s.customerRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
	}

	// Get all tenants for customer
	tenants, _, err := s.tenantRepo.GetByCustomerID(ctx, customer.ID, nil, &repository.Pagination{Page: 1, Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get tenants: %w", err)
	}

	summary := &models.CustomerPendingUpdatesSummary{
		CustomerID:              customer.CustomerID,
		CustomerName:            customer.Name,
		ByPriority:              make(map[string]int),
		ByProduct:               make(map[string]int),
		ByTenant:                make(map[string]int),
		Deployments:             []models.PendingUpdatesResponse{},
	}

	// Process each tenant
	for _, tenant := range tenants {
		tenantFilter := filter
		if tenantFilter == nil {
			tenantFilter = &models.PendingUpdatesFilter{}
		}
		tenantFilter.TenantID = tenant.TenantID

		tenantSummary, err := s.GetPendingUpdatesForTenant(ctx, customer.CustomerID, tenant.TenantID, tenantFilter)
		if err != nil {
			continue // Skip tenants with errors
		}

		summary.TotalDeployments += tenantSummary.TotalDeployments
		summary.DeploymentsWithUpdates += tenantSummary.DeploymentsWithUpdates
		summary.TotalPendingUpdateCount += tenantSummary.TotalPendingUpdateCount

		// Aggregate by priority, product, tenant
		for priority, count := range tenantSummary.ByPriority {
			summary.ByPriority[priority] += count
		}
		for product, count := range tenantSummary.ByProduct {
			summary.ByProduct[product] += count
		}
		summary.ByTenant[tenant.TenantID] = tenantSummary.DeploymentsWithUpdates

		// Add deployments
		summary.Deployments = append(summary.Deployments, tenantSummary.Deployments...)
	}

	return summary, nil
}

// GetAllPendingUpdates retrieves all pending updates across all customers (admin view)
func (s *PendingUpdatesService) GetAllPendingUpdates(ctx context.Context, filter *models.PendingUpdatesFilter, pagination *repository.Pagination) ([]models.PendingUpdatesResponse, *repository.PaginationInfo, error) {
	// Build deployment filter
	deploymentFilter := &repository.DeploymentFilter{}
	if filter != nil {
		if filter.ProductID != "" {
			deploymentFilter.ProductID = filter.ProductID
		}
		if filter.DeploymentType != "" {
			deploymentFilter.DeploymentType = filter.DeploymentType
		}
	}

	// Get all deployments (with optional filters)
	deployments, paginationResult, err := s.deploymentRepo.GetAll(ctx, deploymentFilter, pagination)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get deployments: %w", err)
	}

	return s.processDeploymentsForPendingUpdates(ctx, deployments, filter, paginationResult)
}

// processDeploymentsForPendingUpdates processes a list of deployments and returns pending updates
func (s *PendingUpdatesService) processDeploymentsForPendingUpdates(ctx context.Context, deployments []*models.Deployment, filter *models.PendingUpdatesFilter, paginationResult *repository.PaginationInfo) ([]models.PendingUpdatesResponse, *repository.PaginationInfo, error) {

	var results []models.PendingUpdatesResponse

	// Process each deployment
	for _, deployment := range deployments {
		// Get pending updates for deployment
		pendingUpdates, err := s.GetPendingUpdatesForDeployment(ctx, deployment.ID.Hex())
		if err != nil {
			continue // Skip deployments with errors
		}

		// Apply filters
		if filter != nil {
			if filter.CustomerID != "" && pendingUpdates.CustomerID != filter.CustomerID {
				continue
			}
			if filter.TenantID != "" && pendingUpdates.TenantID != filter.TenantID {
				continue
			}
			if filter.Priority != "" && pendingUpdates.Priority != filter.Priority {
				continue
			}
		}

		// Only include deployments with pending updates
		if pendingUpdates.UpdateCount > 0 {
			results = append(results, *pendingUpdates)
		}
	}

	// Use the original pagination info but update total to reflect filtered results
	// Note: This is an approximation - the actual total might be different
	// For accurate totals, we'd need to count deployments with pending updates separately
	paginationInfo := &repository.PaginationInfo{
		Page:       paginationResult.Page,
		Limit:      paginationResult.Limit,
		Total:      int64(len(results)), // Approximate - could be enhanced
		TotalPages: (int64(len(results)) + int64(paginationResult.Limit) - 1) / int64(paginationResult.Limit),
	}

	return results, paginationInfo, nil
}

// CalculateUpdatePriority calculates the priority level for pending updates
func (s *PendingUpdatesService) CalculateUpdatePriority(deployment *models.Deployment, availableUpdates []models.AvailableUpdate) string {
	// Check for security updates
	for _, update := range availableUpdates {
		if update.IsSecurityUpdate {
			return "critical"
		}
	}

	// Check for EOL approaching (if we have version info)
	// This could be enhanced with actual EOL date checking

	// Check for major version updates on production deployments
	if len(availableUpdates) > 0 {
		latestVersion := availableUpdates[0].VersionNumber
		gapType := utils.GetVersionGapType(deployment.InstalledVersion, latestVersion)

		if gapType == "major" && deployment.DeploymentType == models.DeploymentTypeProduction {
			return "high"
		}
	}

	// Default to normal
	return "normal"
}

