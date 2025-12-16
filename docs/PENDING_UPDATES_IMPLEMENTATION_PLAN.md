# Pending Updates Tracking Implementation Plan

**Document Version:** 1.0  
**Date:** 17-Nov-2025  
**Status:** Implementation Plan  
**Related Requirements:** Section 2.6.10 of requirements.md

---

## Overview

This document outlines the implementation plan for Pending Updates Tracking functionality. The feature tracks available updates for each deployment by comparing the currently deployed version with all released versions for the product.

The implementation follows a hierarchical approach:
- **Deployment Level**: Track pending updates for individual deployments
- **Tenant Level**: Aggregate pending updates across deployments in a tenant
- **Customer Level**: Aggregate pending updates across all customer deployments
- **System Level**: Admin view of all pending updates across all customers

The implementation is divided into three main phases:

1. **Phase 1: Backend Services & APIs** - Core business logic and API endpoints
2. **Phase 2: Frontend Components** - UI components for displaying pending updates
3. **Phase 3: Integration & Testing** - Integration with existing features and comprehensive testing

Each phase includes comprehensive test cases to ensure quality and reliability.

---

## Phase 1: Backend Services & APIs

**Duration:** 2-3 weeks  
**Priority:** High  
**Dependencies:** Customer, Tenant, Deployment, Product, Version models (already implemented)

### Goals
- Implement pending updates calculation logic
- Create service methods for retrieving pending updates
- Add API endpoints for pending updates at all levels
- Optimize queries for performance
- Add caching where appropriate

### Tasks

#### 1.1 Pending Updates Service

**File:** `src/backend/internal/service/pending_updates_service.go` (new)

- [ ] Create `PendingUpdatesService` struct with dependencies:
  - `deploymentRepo` - DeploymentRepository
  - `versionRepo` - VersionRepository
  - `productService` - ProductService
  - `customerRepo` - CustomerRepository
  - `tenantRepo` - TenantRepository

- [ ] Implement `GetAvailableUpdatesForDeployment(ctx, deploymentID)` method:
  - Get deployment by ID
  - Get all released versions for the product (state = Released, not deprecated/EOL)
  - Filter versions newer than deployment's installed_version
  - Sort by version number (newest first)
  - Return list of available versions with metadata:
    - Version number
    - Release date
    - Release type (Security, Feature, Maintenance, Major)
    - Compatibility status
    - Upgrade path information

- [ ] Implement `GetPendingUpdatesCount(ctx, deploymentID)` method:
  - Get available updates for deployment
  - Return count and latest version number

- [ ] Implement `GetPendingUpdatesForTenant(ctx, customerID, tenantID, filters)` method:
  - Get all deployments for tenant
  - For each deployment, get pending updates count
  - Aggregate and return:
    - Total deployments with pending updates
    - Total pending update count
    - List of deployments with pending updates (with details)

- [ ] Implement `GetPendingUpdatesForCustomer(ctx, customerID, filters)` method:
  - Get all tenants for customer
  - For each tenant, get pending updates
  - Aggregate across all tenants
  - Return customer-level summary:
    - Total deployments with pending updates
    - Total pending update count
    - Deployments by priority (critical, high, normal)
    - Deployments by product
    - Deployments by tenant

- [ ] Implement `GetAllPendingUpdates(ctx, filters)` method (admin view):
  - Get all active deployments across all customers
  - For each deployment, get pending updates
  - Apply filters (customer, product, tenant, deployment_type, priority)
  - Return paginated list of deployments with pending updates

- [ ] Implement `CalculateUpdatePriority(deployment, availableVersions)` helper:
  - Determine priority based on:
    - Security releases → Critical
    - EOL approaching → Critical
    - Major version updates + Production → High
    - Minor/patch updates → Normal
  - Return priority level

- [ ] Implement `GetVersionGapType(currentVersion, latestVersion)` helper:
  - Compare semantic versions
  - Return: "patch", "minor", "major"

- [ ] Implement caching strategy:
  - Cache pending updates for deployments (TTL: 5 minutes)
  - Invalidate cache when:
    - New version is released
    - Deployment version is updated
    - Version state changes

#### 1.2 Update Deployment Service

**File:** `src/backend/internal/service/deployment_service.go` (modify)

- [ ] Enhance `GetAvailableUpdates` method (if exists) or create new:
  - Use PendingUpdatesService to get available updates
  - Include upgrade path information
  - Return structured response with metadata

- [ ] Add cache invalidation when deployment version is updated:
  - Clear pending updates cache for the deployment
  - Trigger recalculation for parent tenant and customer

#### 1.3 Update Version Service

**File:** `src/backend/internal/service/version_service.go` (modify)

- [ ] Add cache invalidation when version is released:
  - Clear pending updates cache for all deployments of the product
  - Trigger notification generation (already implemented)

#### 1.4 API Handlers

**File:** `src/backend/internal/api/handlers/pending_updates_handler.go` (new)

- [ ] Create `PendingUpdatesHandler` struct with:
  - `pendingUpdatesService` - PendingUpdatesService
  - `deploymentService` - DeploymentService

- [ ] Implement `GetDeploymentPendingUpdates(w, r)` handler:
  - Route: `GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/updates`
  - Extract customer_id, tenant_id, deployment_id from URL
  - Call service to get available updates
  - Return JSON response with list of available versions

- [ ] Implement `GetTenantPendingUpdates(w, r)` handler:
  - Route: `GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/pending-updates`
  - Extract customer_id, tenant_id from URL
  - Parse query params (product_id, deployment_type, priority)
  - Call service to get tenant-level pending updates
  - Return JSON response with aggregated data

- [ ] Implement `GetCustomerPendingUpdates(w, r)` handler:
  - Route: `GET /api/v1/customers/{customer_id}/deployments/pending-updates`
  - Extract customer_id from URL
  - Parse query params (product_id, deployment_type, priority)
  - Call service to get customer-level pending updates
  - Return JSON response with aggregated data

- [ ] Implement `GetAllPendingUpdates(w, r)` handler (admin):
  - Route: `GET /api/v1/updates/pending`
  - Parse query params (customer_id, product_id, tenant_id, deployment_type, priority, page, limit)
  - Call service to get all pending updates
  - Return paginated JSON response

- [ ] Add error handling and validation:
  - Validate customer/tenant/deployment IDs
  - Handle not found errors
  - Return appropriate HTTP status codes

#### 1.5 Router Configuration

**File:** `src/backend/internal/api/router/router.go` (modify)

- [ ] Add routes for pending updates endpoints:
  ```go
  // Deployment pending updates
  GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/updates
  
  // Tenant pending updates
  GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/pending-updates
  
  // Customer pending updates
  GET /api/v1/customers/{customer_id}/deployments/pending-updates
  
  // All pending updates (admin)
  GET /api/v1/updates/pending
  ```

- [ ] Integrate PendingUpdatesHandler into router
- [ ] Add authentication/authorization middleware

#### 1.6 Service Factory Update

**File:** `src/backend/internal/service/service_factory.go` (modify)

- [ ] Add PendingUpdatesService to ServiceFactory
- [ ] Initialize PendingUpdatesService with required dependencies
- [ ] Wire up dependencies correctly

#### 1.7 Response Models

**File:** `src/backend/internal/models/models.go` (modify)

- [ ] Add `AvailableUpdate` struct:
  ```go
  type AvailableUpdate struct {
      VersionNumber    string    `json:"version_number"`
      ReleaseDate      time.Time `json:"release_date"`
      ReleaseType      string    `json:"release_type"`
      IsSecurityUpdate bool      `json:"is_security_update"`
      CompatibilityStatus string `json:"compatibility_status"`
      UpgradePath      []string  `json:"upgrade_path"` // Intermediate versions if needed
  }
  ```

- [ ] Add `PendingUpdatesResponse` struct:
  ```go
  type PendingUpdatesResponse struct {
      DeploymentID       string            `json:"deployment_id"`
      ProductID           string            `json:"product_id"`
      CurrentVersion      string            `json:"current_version"`
      LatestVersion       string            `json:"latest_version"`
      UpdateCount         int               `json:"update_count"`
      Priority            string            `json:"priority"` // critical, high, normal
      VersionGapType      string            `json:"version_gap_type"` // patch, minor, major
      AvailableUpdates    []AvailableUpdate `json:"available_updates"`
  }
  ```

- [ ] Add `TenantPendingUpdatesSummary` struct:
  ```go
  type TenantPendingUpdatesSummary struct {
      TotalDeployments           int                      `json:"total_deployments"`
      DeploymentsWithUpdates     int                      `json:"deployments_with_updates"`
      TotalPendingUpdateCount    int                      `json:"total_pending_update_count"`
      Deployments                []PendingUpdatesResponse `json:"deployments"`
  }
  ```

- [ ] Add `CustomerPendingUpdatesSummary` struct:
  ```go
  type CustomerPendingUpdatesSummary struct {
      TotalDeployments           int                      `json:"total_deployments"`
      DeploymentsWithUpdates     int                      `json:"deployments_with_updates"`
      TotalPendingUpdateCount    int                      `json:"total_pending_update_count"`
      ByPriority                 map[string]int           `json:"by_priority"`
      ByProduct                  map[string]int           `json:"by_product"`
      ByTenant                   map[string]int           `json:"by_tenant"`
      Deployments                []PendingUpdatesResponse `json:"deployments"`
  }
  ```

### Testing

#### 1.8 Unit Tests

**File:** `src/backend/internal/service/pending_updates_service_test.go` (new)

- [ ] Test `GetAvailableUpdatesForDeployment`:
  - Deployment with no pending updates
  - Deployment with patch updates
  - Deployment with minor updates
  - Deployment with major updates
  - Deployment with security updates
  - Deployment with deprecated/EOL versions (should be excluded)
  - Invalid deployment ID

- [ ] Test `GetPendingUpdatesCount`:
  - Various update counts
  - Edge cases (zero updates, many updates)

- [ ] Test `GetPendingUpdatesForTenant`:
  - Tenant with multiple deployments
  - Tenant with no deployments
  - Tenant with mixed update statuses
  - Filtering by product, deployment_type, priority

- [ ] Test `GetPendingUpdatesForCustomer`:
  - Customer with multiple tenants
  - Customer with no deployments
  - Aggregation logic
  - Filtering

- [ ] Test `GetAllPendingUpdates`:
  - Pagination
  - Filtering by various criteria
  - Sorting
  - Empty results

- [ ] Test `CalculateUpdatePriority`:
  - Critical priority scenarios
  - High priority scenarios
  - Normal priority scenarios

- [ ] Test `GetVersionGapType`:
  - Patch version gaps
  - Minor version gaps
  - Major version gaps

#### 1.9 Integration Tests

**File:** `src/backend/internal/api/handlers/pending_updates_handler_test.go` (new)

- [ ] Test all API endpoints:
  - GET deployment pending updates
  - GET tenant pending updates
  - GET customer pending updates
  - GET all pending updates (admin)
  - Error cases (not found, invalid IDs)
  - Query parameter validation
  - Pagination

#### 1.10 Performance Tests

- [ ] Test with large number of deployments (1000+)
- [ ] Test with large number of versions per product (100+)
- [ ] Test cache effectiveness
- [ ] Test query optimization

---

## Phase 2: Frontend Components

**Duration:** 2-3 weeks  
**Priority:** High  
**Dependencies:** Phase 1 completion

### Goals
- Create UI components for displaying pending updates
- Integrate pending updates into existing pages
- Add filtering and sorting capabilities
- Implement real-time updates

### Tasks

#### 2.1 TypeScript Types

**File:** `src/frontend/src/types/index.ts` (modify)

- [ ] Add `AvailableUpdate` interface:
  ```typescript
  export interface AvailableUpdate {
    version_number: string;
    release_date: string;
    release_type: string;
    is_security_update: boolean;
    compatibility_status: string;
    upgrade_path: string[];
  }
  ```

- [ ] Add `PendingUpdatesResponse` interface:
  ```typescript
  export interface PendingUpdatesResponse {
    deployment_id: string;
    product_id: string;
    current_version: string;
    latest_version: string;
    update_count: number;
    priority: 'critical' | 'high' | 'normal';
    version_gap_type: 'patch' | 'minor' | 'major';
    available_updates: AvailableUpdate[];
  }
  ```

- [ ] Add `TenantPendingUpdatesSummary` interface
- [ ] Add `CustomerPendingUpdatesSummary` interface
- [ ] Add query parameter interfaces for filtering

#### 2.2 API Services

**File:** `src/frontend/src/services/api/pending-updates.ts` (new)

- [ ] Create `pendingUpdatesApi` service with methods:
  - `getDeploymentPendingUpdates(customerId, tenantId, deploymentId)`
  - `getTenantPendingUpdates(customerId, tenantId, filters?)`
  - `getCustomerPendingUpdates(customerId, filters?)`
  - `getAllPendingUpdates(filters?)`
  - Handle pagination, filtering, error handling

#### 2.3 Deployment Components

**File:** `src/frontend/src/components/deployments/DeploymentPendingUpdates.tsx` (new)

- [ ] Create component to display pending updates for a deployment:
  - Show update count badge
  - List of available updates
  - Version comparison (current vs. latest)
  - Upgrade path visualization
  - Priority indicator
  - Link to view release notes

**File:** `src/frontend/src/components/deployments/DeploymentsList.tsx` (modify)

- [ ] Add pending updates count badge to each deployment row
- [ ] Add filter by update status (up-to-date, updates available, critical)
- [ ] Add sort by update count or priority

**File:** `src/frontend/src/pages/DeploymentDetails.tsx` (new, if needed)

- [ ] Create deployment details page with:
  - Pending updates section
  - Available updates list
  - Version comparison
  - Update history

#### 2.4 Tenant Components

**File:** `src/frontend/src/components/tenants/TenantPendingUpdates.tsx` (new)

- [ ] Create component to display tenant-level pending updates:
  - Summary statistics (total deployments, deployments with updates, total update count)
  - List of deployments with pending updates
  - Filtering and sorting

**File:** `src/frontend/src/pages/TenantDetails.tsx` (modify)

- [ ] Add pending updates section to overview tab
- [ ] Display tenant-level summary
- [ ] Link to deployment details

#### 2.5 Customer Components

**File:** `src/frontend/src/components/customers/CustomerPendingUpdates.tsx` (new)

- [ ] Create component to display customer-level pending updates:
  - Summary statistics
  - Breakdown by priority, product, tenant
  - List of deployments with pending updates
  - Quick actions (view deployment, update version)

**File:** `src/frontend/src/pages/CustomerDetails.tsx` (modify)

- [ ] Add pending updates section to overview tab
- [ ] Display customer-level summary
- [ ] Add filter by update status
- [ ] Link to deployments with updates

#### 2.6 Updates Page

**File:** `src/frontend/src/pages/Updates.tsx` (modify)

- [ ] Add new section/tab for "Deployment Updates" or "Customer Updates"
- [ ] Display all deployments with pending updates across all customers
- [ ] Add filters:
  - Customer
  - Product
  - Tenant
  - Deployment type (UAT/Production)
  - Priority (Critical/High/Normal)
  - Update status
- [ ] Add grouping options:
  - By customer
  - By product
  - By tenant
- [ ] Add sorting:
  - By priority
  - By version gap
  - By deployment type
  - By update count
- [ ] Display deployment details:
  - Customer name
  - Tenant name
  - Product
  - Current version
  - Latest version
  - Update count
  - Priority
  - Link to deployment details

**File:** `src/frontend/src/components/updates/PendingUpdatesList.tsx` (new)

- [ ] Create reusable component for listing pending updates
- [ ] Support different views (deployment, tenant, customer, system-wide)
- [ ] Include filtering, sorting, pagination
- [ ] Display update badges and priority indicators

**File:** `src/frontend/src/components/updates/PendingUpdatesFilters.tsx` (new)

- [ ] Create filter component for pending updates
- [ ] Support multiple filter types
- [ ] Clear filters functionality

#### 2.7 Dashboard Updates

**File:** `src/frontend/src/pages/Dashboard.tsx` (modify)

- [ ] Add statistics card for "Deployments with Pending Updates"
- [ ] Add link to Updates page filtered by pending updates

**File:** `src/frontend/src/services/api/dashboard.ts` (modify)

- [ ] Add pending updates count to dashboard statistics
- [ ] Aggregate from customer management APIs

#### 2.8 UI Components

**File:** `src/frontend/src/components/ui/UpdateBadge.tsx` (new)

- [ ] Create badge component for displaying update counts
- [ ] Color coding based on priority (red for critical, orange for high, blue for normal)
- [ ] Show count or "Up to date" indicator

**File:** `src/frontend/src/components/ui/PriorityBadge.tsx` (new)

- [ ] Create badge component for priority levels
- [ ] Color coding and icons

**File:** `src/frontend/src/components/ui/VersionGapBadge.tsx` (new)

- [ ] Create badge component for version gap types
- [ ] Visual indicators for patch/minor/major

### Testing

#### 2.9 Component Tests

- [ ] Test all new components:
  - Rendering with various data states
  - Empty states
  - Loading states
  - Error states
  - User interactions (filtering, sorting, clicking)

#### 2.10 Integration Tests

- [ ] Test integration with existing pages
- [ ] Test navigation flows
- [ ] Test data refresh after version updates

#### 2.11 E2E Tests

**File:** `src/frontend/tests/e2e/pending-updates.spec.ts` (new)

- [ ] Test viewing pending updates at deployment level
- [ ] Test viewing pending updates at tenant level
- [ ] Test viewing pending updates at customer level
- [ ] Test viewing all pending updates in Updates page
- [ ] Test filtering and sorting
- [ ] Test updating deployment version and seeing pending updates recalculate
- [ ] Test real-time updates when new version is released

---

## Phase 3: Integration & Testing

**Duration:** 1-2 weeks  
**Priority:** Medium  
**Dependencies:** Phase 1 and Phase 2 completion

### Goals
- Integrate pending updates with existing features
- Optimize performance
- Add caching
- Comprehensive end-to-end testing
- Documentation

### Tasks

#### 3.1 Integration with Version Release

- [ ] When new version is released:
  - Trigger pending updates recalculation for all affected deployments
  - Invalidate cache
  - Update UI in real-time (if WebSocket/SSE available)

#### 3.2 Integration with Deployment Version Update

- [ ] When deployment version is updated:
  - Recalculate pending updates for that deployment
  - Update parent tenant and customer summaries
  - Invalidate cache
  - Update UI

#### 3.3 Integration with Notifications

- [ ] Enhance notification system to include pending updates count
- [ ] Add notification when critical updates are available
- [ ] Link notifications to pending updates view

#### 3.4 Performance Optimization

- [ ] Implement caching strategy:
  - Cache pending updates for 5 minutes
  - Cache aggregated summaries for 2 minutes
  - Invalidate on version release or deployment update

- [ ] Optimize database queries:
  - Add indexes for common queries
  - Use aggregation pipelines where appropriate
  - Batch operations where possible

- [ ] Implement pagination for large result sets
- [ ] Add lazy loading for deployment lists

#### 3.5 Real-time Updates

- [ ] Implement WebSocket or Server-Sent Events for real-time updates
- [ ] Push updates when:
  - New version is released
  - Deployment version is updated
  - Version state changes

#### 3.6 Documentation

- [ ] Update API documentation with new endpoints
- [ ] Add code comments and documentation
- [ ] Create user guide for pending updates feature
- [ ] Update architecture documentation

#### 3.7 Monitoring & Logging

- [ ] Add logging for pending updates calculations
- [ ] Monitor performance metrics
- [ ] Track cache hit rates
- [ ] Alert on slow queries

### Testing

#### 3.8 End-to-End Testing

- [ ] Test complete flow:
  1. Create customer, tenant, deployment
  2. Release new version
  3. View pending updates at all levels
  4. Update deployment version
  5. Verify pending updates recalculate
  6. Filter and sort in Updates page

#### 3.9 Performance Testing

- [ ] Load testing with large datasets
- [ ] Stress testing cache invalidation
- [ ] Test concurrent updates

#### 3.10 Security Testing

- [ ] Verify authorization checks
- [ ] Test customer data isolation
- [ ] Test admin vs. customer access

---

## Risk Mitigation

### Performance Risks
- **Risk**: Slow queries with large number of deployments/versions
- **Mitigation**: 
  - Implement caching
  - Optimize database queries
  - Add indexes
  - Use pagination

### Data Consistency Risks
- **Risk**: Stale pending updates data
- **Mitigation**:
  - Implement cache invalidation strategy
  - Real-time updates when possible
  - Reasonable cache TTL

### Scalability Risks
- **Risk**: System performance degrades with many customers/deployments
- **Mitigation**:
  - Implement efficient aggregation
  - Use background jobs for heavy calculations
  - Consider read replicas for queries

---

## Timeline

- **Phase 1 (Backend)**: 2-3 weeks
- **Phase 2 (Frontend)**: 2-3 weeks
- **Phase 3 (Integration)**: 1-2 weeks
- **Total**: 5-8 weeks

---

## Success Criteria

1. ✅ Pending updates are accurately calculated for all deployments
2. ✅ Pending updates are visible at deployment, tenant, customer, and system levels
3. ✅ Updates page displays all deployments with pending updates
4. ✅ Filtering and sorting work correctly
5. ✅ Pending updates recalculate in real-time when versions are released or updated
6. ✅ Performance is acceptable with large datasets (1000+ deployments)
7. ✅ All API endpoints are tested and documented
8. ✅ UI components are responsive and user-friendly
9. ✅ E2E tests pass
10. ✅ Documentation is complete

---

## Dependencies

- Customer Management (Phase 1-3) - ✅ Completed
- Product and Version Management - ✅ Completed
- Notification System - ✅ Completed
- Audit Logging - ✅ Completed

---

## Notes

- Consider implementing background jobs for heavy calculations if performance becomes an issue
- May need to add database indexes for optimal query performance
- Consider implementing WebSocket/SSE for real-time updates in future iterations
- Cache invalidation strategy should be carefully designed to balance performance and data freshness

