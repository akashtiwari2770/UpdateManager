# License Management Implementation Plan

**Document Version:** 1.0  
**Date:** 17-Nov-2025  
**Status:** Implementation Plan  
**Related Requirements:** Section 2.7 of requirements.md

---

## Overview

This document outlines the implementation plan for License Management functionality as specified in the requirements document. The implementation follows a hierarchical model:

**Customer → Subscription → License → License Allocation (Tenant/Deployment)**

Where:
- **Customer**: Organization/company using the system
- **Subscription**: A customer can have one or more subscriptions
- **License**: Sales team assigns licenses to subscriptions (for specific products, with user/seats count)
- **License Allocation**: Customer distributes licenses across tenants and deployments

The implementation is divided into three main phases:

1. **Phase 1: Data Models** - Backend data structures and database schema
2. **Phase 2: API Layer** - Backend API endpoints, services, and repositories
3. **Phase 3: UI Components** - Frontend React components and pages

Each phase includes comprehensive test cases to ensure quality and reliability.

---

## Phase 1: Data Models (Backend)

**Duration:** 1-2 weeks  
**Priority:** High  
**Dependencies:** Customer, Tenant, Deployment models (already implemented)

### Goals
- Define Subscription, License, and License Allocation data models
- Create database collections/indexes
- Add validation rules
- Integrate with existing Customer, Tenant, and Deployment models

### Tasks

#### 1.1 Subscription Model
**File:** `src/backend/internal/models/models.go`

- [ ] Define `Subscription` struct with fields:
  - `ID` (primitive.ObjectID)
  - `SubscriptionID` (string, unique, required)
  - `CustomerID` (primitive.ObjectID, required, indexed)
  - `Name` (string, optional, max 200)
  - `Description` (string, optional, max 1000)
  - `StartDate` (time.Time, required)
  - `EndDate` (*time.Time, optional, for time-based subscriptions)
  - `Status` (SubscriptionStatus enum)
  - `CreatedAt` (time.Time)
  - `UpdatedAt` (time.Time)
  - `CreatedBy` (string, user ID)
  - `Notes` (string, optional, max 2000)

- [ ] Define `SubscriptionStatus` type:
  ```go
  type SubscriptionStatus string
  const (
      SubscriptionStatusActive    SubscriptionStatus = "active"
      SubscriptionStatusInactive  SubscriptionStatus = "inactive"
      SubscriptionStatusExpired   SubscriptionStatus = "expired"
      SubscriptionStatusSuspended SubscriptionStatus = "suspended"
  )
  ```

- [ ] Define `CreateSubscriptionRequest` struct
- [ ] Define `UpdateSubscriptionRequest` struct

#### 1.2 License Model
**File:** `src/backend/internal/models/models.go`

- [ ] Define `License` struct with fields:
  - `ID` (primitive.ObjectID)
  - `LicenseID` (string, unique, required)
  - `SubscriptionID` (primitive.ObjectID, required, indexed)
  - `ProductID` (string, required, indexed)
  - `LicenseType` (LicenseType enum)
  - `NumberOfSeats` (int, required, min 1)
  - `StartDate` (time.Time, required)
  - `EndDate` (*time.Time, optional, for time-based licenses)
  - `Status` (LicenseStatus enum)
  - `AssignedBy` (string, user ID, required)
  - `AssignmentDate` (time.Time, required)
  - `Notes` (string, optional, max 2000)
  - `CreatedAt` (time.Time)
  - `UpdatedAt` (time.Time)

- [ ] Define `LicenseType` type:
  ```go
  type LicenseType string
  const (
      LicenseTypePerpetual LicenseType = "perpetual"
      LicenseTypeTimeBased LicenseType = "time_based"
  )
  ```

- [ ] Define `LicenseStatus` type:
  ```go
  type LicenseStatus string
  const (
      LicenseStatusActive   LicenseStatus = "active"
      LicenseStatusInactive LicenseStatus = "inactive"
      LicenseStatusExpired  LicenseStatus = "expired"
      LicenseStatusRevoked  LicenseStatus = "revoked"
  )
  ```

- [ ] Define `CreateLicenseRequest` struct
- [ ] Define `UpdateLicenseRequest` struct

#### 1.3 License Allocation Model
**File:** `src/backend/internal/models/models.go`

- [ ] Define `LicenseAllocation` struct with fields:
  - `ID` (primitive.ObjectID)
  - `AllocationID` (string, unique, required)
  - `LicenseID` (primitive.ObjectID, required, indexed)
  - `TenantID` (*primitive.ObjectID, optional, indexed)
  - `DeploymentID` (*primitive.ObjectID, optional, indexed)
  - `NumberOfSeatsAllocated` (int, required, min 1)
  - `AllocationDate` (time.Time, required)
  - `AllocatedBy` (string, user ID, required)
  - `Status` (AllocationStatus enum)
  - `ReleasedDate` (*time.Time, optional)
  - `ReleasedBy` (*string, optional, user ID)
  - `Notes` (string, optional, max 2000)
  - `CreatedAt` (time.Time)
  - `UpdatedAt` (time.Time)

- [ ] Define `AllocationStatus` type:
  ```go
  type AllocationStatus string
  const (
      AllocationStatusActive  AllocationStatus = "active"
      AllocationStatusReleased AllocationStatus = "released"
  )
  ```

- [ ] Define `AllocateLicenseRequest` struct
- [ ] Define `ReleaseLicenseRequest` struct

#### 1.4 Database Indexes
**File:** `src/database/mongodb-indexes.js`

- [ ] Add indexes for `subscriptions` collection:
  - `customer_id` (ascending)
  - `subscription_id` (unique)
  - `status` (ascending)
  - `start_date` and `end_date` (for expiration queries)

- [ ] Add indexes for `licenses` collection:
  - `subscription_id` (ascending)
  - `license_id` (unique)
  - `product_id` (ascending)
  - `license_type` (ascending)
  - `status` (ascending)
  - `end_date` (for expiration queries)
  - Compound index: `subscription_id` + `product_id`

- [ ] Add indexes for `license_allocations` collection:
  - `license_id` (ascending)
  - `allocation_id` (unique)
  - `tenant_id` (ascending, sparse)
  - `deployment_id` (ascending, sparse)
  - `status` (ascending)
  - Compound index: `license_id` + `status`

#### 1.5 Model Validation
- [ ] Add validation tags to all struct fields
- [ ] Create validation functions for:
  - Subscription dates (end date must be after start date)
  - License dates (end date required for time-based, optional for perpetual)
  - License allocation (cannot exceed available seats)
  - Product match validation (license product must match deployment product)

### Testing

#### 1.6 Unit Tests
**File:** `src/backend/internal/models/models_test.go`

- [ ] Test Subscription model validation
- [ ] Test License model validation
- [ ] Test License Allocation model validation
- [ ] Test status enum values
- [ ] Test date validation logic

---

## Phase 2: API Layer (Backend)

**Duration:** 2-3 weeks  
**Priority:** High  
**Dependencies:** Phase 1 (Data Models)

### Goals
- Implement repositories for Subscription, License, and License Allocation
- Create service layer with business logic
- Add API handlers and routes
- Implement validation and error handling
- Add audit logging

### Tasks

#### 2.1 Subscription Repository
**File:** `src/backend/internal/repository/subscription_repository.go`

- [ ] Create `SubscriptionRepository` interface
- [ ] Implement methods:
  - `Create(ctx, subscription)`
  - `GetByID(ctx, id)`
  - `GetBySubscriptionID(ctx, subscriptionID)`
  - `GetByCustomerID(ctx, customerID, filter, pagination)`
  - `Update(ctx, subscription)`
  - `Delete(ctx, id)`
  - `Count(ctx, filter)`
  - `GetExpiringSubscriptions(ctx, days)`

#### 2.2 License Repository
**File:** `src/backend/internal/repository/license_repository.go`

- [ ] Create `LicenseRepository` interface
- [ ] Implement methods:
  - `Create(ctx, license)`
  - `GetByID(ctx, id)`
  - `GetByLicenseID(ctx, licenseID)`
  - `GetBySubscriptionID(ctx, subscriptionID, filter, pagination)`
  - `GetByProductID(ctx, productID, filter, pagination)`
  - `Update(ctx, license)`
  - `Delete(ctx, id)`
  - `Count(ctx, filter)`
  - `GetExpiringLicenses(ctx, days)`
  - `GetAvailableSeats(ctx, licenseID)`

#### 2.3 License Allocation Repository
**File:** `src/backend/internal/repository/license_allocation_repository.go`

- [ ] Create `LicenseAllocationRepository` interface
- [ ] Implement methods:
  - `Create(ctx, allocation)`
  - `GetByID(ctx, id)`
  - `GetByAllocationID(ctx, allocationID)`
  - `GetByLicenseID(ctx, licenseID, filter, pagination)`
  - `GetByTenantID(ctx, tenantID, filter, pagination)`
  - `GetByDeploymentID(ctx, deploymentID, filter, pagination)`
  - `GetActiveAllocationsByLicenseID(ctx, licenseID)`
  - `Update(ctx, allocation)`
  - `Release(ctx, allocationID, releasedBy)`
  - `Count(ctx, filter)`
  - `GetTotalAllocatedSeats(ctx, licenseID)`

#### 2.4 Subscription Service
**File:** `src/backend/internal/service/subscription_service.go`

- [ ] Create `SubscriptionService` struct
- [ ] Implement methods:
  - `CreateSubscription(ctx, customerID, req, userID)`
  - `GetSubscription(ctx, customerID, subscriptionID)`
  - `ListSubscriptions(ctx, customerID, filter, pagination)`
  - `UpdateSubscription(ctx, customerID, subscriptionID, req, userID)`
  - `DeleteSubscription(ctx, customerID, subscriptionID, userID)`
  - `GetSubscriptionStatistics(ctx, customerID, subscriptionID)`
  - `RenewSubscription(ctx, customerID, subscriptionID, newEndDate, userID)`
  - `ValidateSubscriptionStatus(ctx, subscriptionID)`

#### 2.5 License Service
**File:** `src/backend/internal/service/license_service.go`

- [ ] Create `LicenseService` struct
- [ ] Implement methods:
  - `AssignLicense(ctx, customerID, subscriptionID, req, userID)`
  - `GetLicense(ctx, customerID, subscriptionID, licenseID)`
  - `ListLicenses(ctx, customerID, subscriptionID, filter, pagination)`
  - `UpdateLicense(ctx, customerID, subscriptionID, licenseID, req, userID)`
  - `RevokeLicense(ctx, customerID, subscriptionID, licenseID, userID)`
  - `GetLicenseStatistics(ctx, customerID, subscriptionID, licenseID)`
  - `ValidateLicenseStatus(ctx, licenseID)`
  - `GetAvailableSeats(ctx, licenseID)`
  - `CheckLicenseExpiration(ctx, licenseID)`
  - `RenewLicense(ctx, customerID, subscriptionID, licenseID, newEndDate, userID)`

#### 2.6 License Allocation Service
**File:** `src/backend/internal/service/license_allocation_service.go`

- [ ] Create `LicenseAllocationService` struct
- [ ] Implement methods:
  - `AllocateLicense(ctx, customerID, subscriptionID, licenseID, req, userID)`
  - `ReleaseAllocation(ctx, customerID, subscriptionID, licenseID, allocationID, userID)`
  - `GetAllocations(ctx, customerID, subscriptionID, licenseID, filter, pagination)`
  - `GetAllocationsByTenant(ctx, customerID, tenantID, filter, pagination)`
  - `GetAllocationsByDeployment(ctx, customerID, tenantID, deploymentID, filter, pagination)`
  - `GetLicenseUtilization(ctx, licenseID)`
  - `ValidateAllocation(ctx, licenseID, seats, productID)`
  - `GetTotalAllocatedSeats(ctx, licenseID)`

#### 2.7 API Handlers
**Files:** 
- `src/backend/internal/api/handlers/subscription_handler.go`
- `src/backend/internal/api/handlers/license_handler.go`
- `src/backend/internal/api/handlers/license_allocation_handler.go`

- [ ] Subscription handlers:
  - `CreateSubscription`
  - `GetSubscription`
  - `ListSubscriptions`
  - `UpdateSubscription`
  - `DeleteSubscription`
  - `GetSubscriptionStatistics`

- [ ] License handlers:
  - `AssignLicense`
  - `GetLicense`
  - `ListLicenses`
  - `UpdateLicense`
  - `RevokeLicense`
  - `GetLicenseStatistics`

- [ ] License Allocation handlers:
  - `AllocateLicense`
  - `ReleaseAllocation`
  - `GetAllocations`
  - `GetAllocationsByTenant`
  - `GetAllocationsByDeployment`
  - `GetLicenseUtilization`

#### 2.8 Router Integration
**File:** `src/backend/internal/api/router/router.go`

- [ ] Add subscription routes:
  - `GET /api/v1/customers/{customer_id}/subscriptions`
  - `POST /api/v1/customers/{customer_id}/subscriptions`
  - `GET /api/v1/customers/{customer_id}/subscriptions/{subscription_id}`
  - `PUT /api/v1/customers/{customer_id}/subscriptions/{subscription_id}`
  - `DELETE /api/v1/customers/{customer_id}/subscriptions/{subscription_id}`
  - `GET /api/v1/subscriptions` (admin view)

- [ ] Add license routes:
  - `GET /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses`
  - `POST /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses`
  - `GET /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}`
  - `PUT /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}`
  - `DELETE /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}`
  - `GET /api/v1/licenses` (admin view)

- [ ] Add license allocation routes:
  - `POST /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}/allocate`
  - `POST /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}/release`
  - `GET /api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}/allocations`
  - `GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/licenses`
  - `GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/licenses`

#### 2.9 Service Factory Integration
**File:** `src/backend/internal/service/service_factory.go`

- [ ] Add SubscriptionService to ServiceFactory
- [ ] Add LicenseService to ServiceFactory
- [ ] Add LicenseAllocationService to ServiceFactory
- [ ] Wire up dependencies correctly

### Testing

#### 2.10 Repository Tests
- [ ] Test SubscriptionRepository CRUD operations
- [ ] Test LicenseRepository CRUD operations
- [ ] Test LicenseAllocationRepository CRUD operations
- [ ] Test query filters and pagination
- [ ] Test expiration queries

#### 2.11 Service Tests
- [ ] Test SubscriptionService business logic
- [ ] Test LicenseService business logic
- [ ] Test LicenseAllocationService business logic
- [ ] Test validation logic
- [ ] Test seat allocation limits
- [ ] Test expiration handling

#### 2.12 API Tests
- [ ] Test all subscription endpoints
- [ ] Test all license endpoints
- [ ] Test all license allocation endpoints
- [ ] Test error handling
- [ ] Test authentication/authorization

---

## Phase 3: UI Components (Frontend)

**Duration:** 3-4 weeks  
**Priority:** High  
**Dependencies:** Phase 2 (API Layer)

### Goals
- Create TypeScript types and interfaces
- Build API service layer
- Create React components for subscriptions, licenses, and allocations
- Integrate with existing customer, tenant, and deployment pages
- Add license dashboard and reporting

### Tasks

#### 3.1 TypeScript Types
**File:** `src/frontend/src/types/index.ts`

- [ ] Define `Subscription` interface
- [ ] Define `SubscriptionStatus` enum
- [ ] Define `License` interface
- [ ] Define `LicenseType` enum
- [ ] Define `LicenseStatus` enum
- [ ] Define `LicenseAllocation` interface
- [ ] Define `AllocationStatus` enum
- [ ] Define request/response DTOs:
  - `CreateSubscriptionRequest`
  - `UpdateSubscriptionRequest`
  - `CreateLicenseRequest`
  - `UpdateLicenseRequest`
  - `AllocateLicenseRequest`
  - `ReleaseLicenseRequest`
- [ ] Define query parameters and filters

#### 3.2 API Services
**Files:**
- `src/frontend/src/services/api/subscriptions.ts`
- `src/frontend/src/services/api/licenses.ts`
- `src/frontend/src/services/api/license-allocations.ts`

- [ ] Subscription API methods:
  - `create(customerId, data)`
  - `getById(customerId, subscriptionId)`
  - `list(customerId, filters, pagination)`
  - `update(customerId, subscriptionId, data)`
  - `delete(customerId, subscriptionId)`
  - `getStatistics(customerId, subscriptionId)`

- [ ] License API methods:
  - `assign(customerId, subscriptionId, data)`
  - `getById(customerId, subscriptionId, licenseId)`
  - `list(customerId, subscriptionId, filters, pagination)`
  - `update(customerId, subscriptionId, licenseId, data)`
  - `revoke(customerId, subscriptionId, licenseId)`
  - `getStatistics(customerId, subscriptionId, licenseId)`

- [ ] License Allocation API methods:
  - `allocate(customerId, subscriptionId, licenseId, data)`
  - `release(customerId, subscriptionId, licenseId, allocationId)`
  - `getAllocations(customerId, subscriptionId, licenseId, filters, pagination)`
  - `getByTenant(customerId, tenantId, filters, pagination)`
  - `getByDeployment(customerId, tenantId, deploymentId, filters, pagination)`
  - `getUtilization(customerId, subscriptionId, licenseId)`

#### 3.3 Subscription Components
**Files:**
- `src/frontend/src/components/subscriptions/SubscriptionsList.tsx`
- `src/frontend/src/components/subscriptions/SubscriptionForm.tsx`
- `src/frontend/src/components/subscriptions/SubscriptionDetails.tsx`
- `src/frontend/src/components/subscriptions/SubscriptionStatistics.tsx`

- [ ] SubscriptionsList component:
  - Display paginated list of subscriptions
  - Filter by status, date range
  - Search functionality
  - Actions: View, Edit, Delete

- [ ] SubscriptionForm component:
  - Create/Edit subscription form
  - Validation
  - Date pickers for start/end dates
  - Status selection

- [ ] SubscriptionDetails component:
  - Display subscription information
  - List of licenses in subscription
  - License allocation summary
  - Available vs. allocated seats
  - Expiration timeline
  - Actions: Edit, Delete, Renew

- [ ] SubscriptionStatistics component:
  - Total licenses count
  - Active/Expired licenses
  - Utilization metrics
  - Expiring licenses count

#### 3.4 License Components
**Files:**
- `src/frontend/src/components/licenses/LicensesList.tsx`
- `src/frontend/src/components/licenses/LicenseForm.tsx`
- `src/frontend/src/components/licenses/LicenseDetails.tsx`
- `src/frontend/src/components/licenses/LicenseTypeBadge.tsx`
- `src/frontend/src/components/licenses/LicenseStatusBadge.tsx`

- [ ] LicensesList component:
  - Display paginated list of licenses
  - Filter by product, type, status, date range
  - Search functionality
  - Expiration warnings
  - Actions: View, Edit, Revoke

- [ ] LicenseForm component:
  - Create/Edit license form
  - Product selection
  - License type selection (Perpetual/Time-based)
  - Number of seats input
  - Date pickers
  - Validation

- [ ] LicenseDetails component:
  - Display license information
  - Allocation history
  - Current allocations (tenants/deployments)
  - Utilization metrics
  - Expiration status and warnings
  - Actions: Edit, Revoke, Renew, Allocate

- [ ] LicenseTypeBadge component:
  - Visual badge for Perpetual/Time-based

- [ ] LicenseStatusBadge component:
  - Color-coded status badges

#### 3.5 License Allocation Components
**Files:**
- `src/frontend/src/components/license-allocations/AllocateLicenseForm.tsx`
- `src/frontend/src/components/license-allocations/AllocationsList.tsx`
- `src/frontend/src/components/license-allocations/LicenseUtilization.tsx`

- [ ] AllocateLicenseForm component:
  - Tenant/Deployment selection
  - Number of seats input
  - Validation (available seats check)
  - Product match validation

- [ ] AllocationsList component:
  - Display allocations for a license
  - Filter by tenant, deployment, status
  - Actions: Release allocation

- [ ] LicenseUtilization component:
  - Total vs. allocated seats
  - Utilization percentage
  - Visual progress bar/chart

#### 3.6 Customer Integration
**File:** `src/frontend/src/pages/CustomerDetails.tsx`

- [ ] Add "Subscriptions" tab to CustomerDetails
- [ ] Display subscription summary in overview
- [ ] Link to subscription details
- [ ] Show license statistics

#### 3.7 Deployment Integration
**File:** `src/frontend/src/components/deployments/DeploymentDetails.tsx`

- [ ] Add "Licenses" section to deployment details
- [ ] Display allocated licenses
- [ ] Show license compliance status
- [ ] Warn when user count exceeds allocated seats
- [ ] Link to license details

#### 3.8 License Dashboard
**File:** `src/frontend/src/pages/Licenses.tsx`

- [ ] Create main Licenses page
- [ ] Display license dashboard with:
  - Total subscriptions count
  - Total licenses count (by type)
  - Active/Expired licenses
  - Licenses expiring soon
  - Utilization metrics
  - Licenses by product
- [ ] Filter and search capabilities
- [ ] Export functionality

#### 3.9 Routing
**File:** `src/frontend/src/router/index.tsx`

- [ ] Add routes:
  - `/customers/:customerId/subscriptions`
  - `/customers/:customerId/subscriptions/new`
  - `/customers/:customerId/subscriptions/:subscriptionId`
  - `/customers/:customerId/subscriptions/:subscriptionId/edit`
  - `/customers/:customerId/subscriptions/:subscriptionId/licenses`
  - `/customers/:customerId/subscriptions/:subscriptionId/licenses/new`
  - `/customers/:customerId/subscriptions/:subscriptionId/licenses/:licenseId`
  - `/licenses` (admin view)

#### 3.10 Navigation
**File:** `src/frontend/src/components/layout/Sidebar.tsx`

- [ ] Add "Licenses" menu item (admin view)
- [ ] Add "Subscriptions" submenu under Customers

### Testing

#### 3.11 Component Tests
- [ ] Test all subscription components
- [ ] Test all license components
- [ ] Test all license allocation components
- [ ] Test form validation
- [ ] Test error handling
- [ ] Test loading states

#### 3.12 Integration Tests
- [ ] Test subscription workflow (create → view → edit → delete)
- [ ] Test license assignment workflow
- [ ] Test license allocation workflow
- [ ] Test integration with customer/deployment pages

#### 3.13 E2E Tests
**File:** `src/frontend/tests/e2e/license-management.spec.ts`

- [ ] Test subscription CRUD operations
- [ ] Test license assignment
- [ ] Test license allocation to tenant/deployment
- [ ] Test license release
- [ ] Test expiration warnings
- [ ] Test utilization display
- [ ] Test license compliance validation

---

## Implementation Notes

### Business Logic Considerations

1. **License Expiration**: Time-based licenses should automatically update status based on current date vs. end date
2. **Seat Allocation**: Must prevent over-allocation (total allocated seats cannot exceed license seats)
3. **Product Matching**: License must be for the same product as the deployment
4. **Partial Allocation**: Support allocating subset of seats from a license
5. **License Release**: When releasing allocation, seats become available again
6. **Status Calculation**: Subscription and license status should be calculated based on dates and allocations

### Integration Points

1. **Customer Management**: Subscriptions belong to customers
2. **Tenant Management**: Licenses can be allocated to tenants
3. **Deployment Management**: Licenses can be allocated to deployments, validate user count
4. **Product Management**: Licenses are for specific products
5. **Audit Logging**: All license operations should be logged

### Performance Considerations

1. **License Utilization Calculation**: Cache utilization metrics for frequently accessed licenses
2. **Expiration Queries**: Use indexes for efficient expiration date queries
3. **Allocation Aggregation**: Optimize queries for total allocated seats calculation

---

## Success Criteria

- ✅ All data models defined and tested
- ✅ All API endpoints implemented and tested
- ✅ All UI components created and tested
- ✅ License allocation validation working correctly
- ✅ Expiration warnings displayed appropriately
- ✅ Utilization metrics accurate
- ✅ Integration with customer/tenant/deployment pages complete
- ✅ E2E tests passing
- ✅ Documentation complete

---

## Timeline Estimate

- **Phase 1 (Data Models)**: 1-2 weeks
- **Phase 2 (API Layer)**: 2-3 weeks
- **Phase 3 (UI Components)**: 3-4 weeks
- **Total**: 6-9 weeks

---

**Status**: Ready for Implementation

