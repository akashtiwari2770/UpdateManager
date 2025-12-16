# Customer Management Implementation Plan

**Document Version:** 2.0  
**Date:** 17-Nov-2025  
**Status:** Implementation Plan  
**Related Requirements:** Section 2.6 of requirements.md

---

## Overview

This document outlines the implementation plan for Customer Management functionality as specified in the requirements document. The implementation follows a hierarchical model:

**Customer → Tenant → Deployment (Product + Type)**

Where:
- **Customer**: Organization/company using the system
- **Tenant**: Independent deployment environment (e.g., data center, region, business unit)
- **Deployment**: Combination of a product and deployment type (UAT/Production)

The implementation is divided into three main phases:

1. **Phase 1: Data Models** - Backend data structures and database schema
2. **Phase 2: API Layer** - Backend API endpoints, services, and repositories
3. **Phase 3: UI Components** - Frontend React components and pages

Each phase includes comprehensive test cases to ensure quality and reliability.

---

## Phase 1: Data Models (Backend)

**Duration:** 1-2 weeks  
**Priority:** Critical  
**Dependencies:** None

### Goals
- Define Customer, Tenant, and Deployment data models
- Create database collections/indexes
- Add validation rules
- Update existing models if needed

### Tasks

#### 1.1 Customer Model
**File:** `src/backend/internal/models/models.go`

- [ ] Define `Customer` struct with fields:
  - `ID` (primitive.ObjectID)
  - `CustomerID` (string, unique, required)
  - `Name` (string, required, max 200)
  - `OrganizationName` (string, optional, max 200)
  - `Email` (string, required, email validation)
  - `Phone` (string, optional, max 50)
  - `Address` (string, optional, max 500)
  - `AccountStatus` (CustomerStatus enum)
  - `NotificationPreferences` (NotificationPreferences struct)
  - `CreatedAt` (time.Time)
  - `UpdatedAt` (time.Time)

- [ ] Define `CustomerStatus` type:
  ```go
  type CustomerStatus string
  const (
      CustomerStatusActive    CustomerStatus = "active"
      CustomerStatusInactive  CustomerStatus = "inactive"
      CustomerStatusSuspended CustomerStatus = "suspended"
  )
  ```

- [ ] Define `NotificationPreferences` struct:
  - `EmailEnabled` (bool)
  - `InAppEnabled` (bool)
  - `UATNotifications` (bool)
  - `ProductionNotifications` (bool)

- [ ] Define `CreateCustomerRequest` struct
- [ ] Define `UpdateCustomerRequest` struct

#### 1.2 Tenant Model (Customer-Managed)
**File:** `src/backend/internal/models/models.go`

- [ ] Define `CustomerTenant` struct with fields:
  - `ID` (primitive.ObjectID)
  - `TenantID` (string, unique, required)
  - `CustomerID` (primitive.ObjectID, required, indexed)
  - `Name` (string, required, max 200)
  - `Description` (string, optional, max 1000)
  - `Status` (TenantStatus enum)
  - `CreatedAt` (time.Time)
  - `UpdatedAt` (time.Time)

- [ ] Define `TenantStatus` type:
  ```go
  type TenantStatus string
  const (
      TenantStatusActive   TenantStatus = "active"
      TenantStatusInactive TenantStatus = "inactive"
  )
  ```

- [ ] Define `CreateTenantRequest` struct
- [ ] Define `UpdateTenantRequest` struct

**Note:** This is different from the existing System Tenant model used for multi-tenant product deployments.

#### 1.3 Deployment Model
**File:** `src/backend/internal/models/models.go`

- [ ] Define `Deployment` struct with fields:
  - `ID` (primitive.ObjectID)
  - `DeploymentID` (string, unique, required)
  - `TenantID` (primitive.ObjectID, required, indexed)
  - `ProductID` (string, required, indexed)
  - `DeploymentType` (DeploymentType enum)
  - `InstalledVersion` (string, required)
  - `NumberOfUsers` (*int, optional)
  - `LicenseInfo` (string, optional, max 1000)
  - `ServerHostname` (string, optional, max 200)
  - `EnvironmentDetails` (string, optional, max 500)
  - `DeploymentDate` (time.Time)
  - `LastUpdatedDate` (time.Time)
  - `Status` (DeploymentStatus enum)

- [ ] Define `DeploymentType` type:
  ```go
  type DeploymentType string
  const (
      DeploymentTypeUAT         DeploymentType = "uat"
      DeploymentTypeTesting     DeploymentType = "testing"
      DeploymentTypeProduction  DeploymentType = "production"
  )
  ```

- [ ] Define `DeploymentStatus` type:
  ```go
  type DeploymentStatus string
  const (
      DeploymentStatusActive   DeploymentStatus = "active"
      DeploymentStatusInactive DeploymentStatus = "inactive"
  )
  ```

- [ ] Define `CreateDeploymentRequest` struct
- [ ] Define `UpdateDeploymentRequest` struct

#### 1.4 Database Indexes
**File:** `src/backend/internal/repository/customer_repository.go` (new file)

- [ ] Create indexes for Customer collection:
  - Unique index on `customer_id`
  - Index on `email`
  - Index on `account_status`
  - Index on `created_at`

- [ ] Create indexes for Tenant collection:
  - Unique index on `tenant_id`
  - Index on `customer_id`
  - Index on `status`
  - Compound index on `customer_id` + `status`

- [ ] Create indexes for Deployment collection:
  - Unique index on `deployment_id`
  - Index on `tenant_id`
  - Index on `product_id`
  - Index on `deployment_type`
  - Index on `status`
  - Compound index on `tenant_id` + `product_id`
  - Compound index on `tenant_id` + `deployment_type`
  - Compound index on `product_id` + `deployment_type`

#### 1.5 Update Notification Model (if needed)
**File:** `src/backend/internal/models/models.go`

- [ ] Review and update `Notification` model to support customer notifications:
  - Ensure `recipient_id` can be customer ID
  - Add `customer_id` field if needed
  - Add `tenant_id` field for tenant-specific notifications
  - Add `deployment_id` field for deployment-specific notifications

### Test Cases

#### 1.6 Unit Tests for Models
**File:** `src/backend/internal/models/models_test.go`

- [ ] Test Customer model validation:
  - Valid customer creation
  - Invalid email format
  - Missing required fields
  - Field length validation
  - CustomerStatus enum validation

- [ ] Test Tenant model validation:
  - Valid tenant creation
  - Missing required fields
  - TenantStatus enum validation
  - CustomerID reference validation

- [ ] Test Deployment model validation:
  - Valid deployment creation
  - Missing required fields
  - DeploymentType enum validation
  - DeploymentStatus enum validation
  - NumberOfUsers validation (positive integer)
  - Version format validation
  - TenantID reference validation

- [ ] Test Request structs:
  - CreateCustomerRequest validation
  - UpdateCustomerRequest validation
  - CreateTenantRequest validation
  - UpdateTenantRequest validation
  - CreateDeploymentRequest validation
  - UpdateDeploymentRequest validation

#### 1.7 Integration Tests
**Files:**
- `src/backend/internal/repository/customer_repository_test.go` (new file)
- `src/backend/internal/repository/tenant_repository_test.go` (new file)
- `src/backend/internal/repository/deployment_repository_test.go` (new file)

- [ ] Test database indexes creation
- [ ] Test unique constraints
- [ ] Test data persistence
- [ ] Test foreign key relationships (Customer → Tenant → Deployment)

### Acceptance Criteria
- ✅ All data models defined with proper types
- ✅ Validation rules implemented
- ✅ Database indexes created
- ✅ Unit tests passing (100% coverage for models)
- ✅ Integration tests passing

---

## Phase 2: API Layer (Backend)

**Duration:** 3-4 weeks  
**Priority:** Critical  
**Dependencies:** Phase 1 complete

### Goals
- Implement repository layer for database operations
- Implement service layer for business logic
- Implement API handlers for HTTP endpoints
- Add routing for customer management APIs
- Integrate with notification system

### Tasks

#### 2.1 Customer Repository
**File:** `src/backend/internal/repository/customer_repository.go` (new file)

- [ ] Implement `CustomerRepository` interface:
  - `Create(customer *models.Customer) error`
  - `GetByID(id primitive.ObjectID) (*models.Customer, error)`
  - `GetByCustomerID(customerID string) (*models.Customer, error)`
  - `GetByEmail(email string) (*models.Customer, error)`
  - `List(filter *CustomerFilter, pagination *Pagination) ([]*models.Customer, *PaginationInfo, error)`
  - `Update(id primitive.ObjectID, updates *models.Customer) error`
  - `Delete(id primitive.ObjectID) error`
  - `Count(filter *CustomerFilter) (int64, error)`

- [ ] Implement `CustomerFilter` struct:
  - `Search` (string) - search in name, organization, email
  - `Status` (models.CustomerStatus)
  - `Email` (string)

- [ ] Implement MongoDB operations
- [ ] Handle database errors appropriately

#### 2.2 Tenant Repository
**File:** `src/backend/internal/repository/tenant_repository.go` (new file)

- [ ] Implement `TenantRepository` interface:
  - `Create(tenant *models.CustomerTenant) error`
  - `GetByID(id primitive.ObjectID) (*models.CustomerTenant, error)`
  - `GetByTenantID(tenantID string) (*models.CustomerTenant, error)`
  - `GetByCustomerID(customerID primitive.ObjectID, filter *TenantFilter, pagination *Pagination) ([]*models.CustomerTenant, *PaginationInfo, error)`
  - `Update(id primitive.ObjectID, updates *models.CustomerTenant) error`
  - `Delete(id primitive.ObjectID) error`
  - `CountByCustomerID(customerID primitive.ObjectID, filter *TenantFilter) (int64, error)`

- [ ] Implement `TenantFilter` struct:
  - `Status` (models.TenantStatus)

- [ ] Implement MongoDB operations
- [ ] Handle database errors appropriately

#### 2.3 Deployment Repository
**File:** `src/backend/internal/repository/deployment_repository.go` (new file)

- [ ] Implement `DeploymentRepository` interface:
  - `Create(deployment *models.Deployment) error`
  - `GetByID(id primitive.ObjectID) (*models.Deployment, error)`
  - `GetByDeploymentID(deploymentID string) (*models.Deployment, error)`
  - `GetByTenantID(tenantID primitive.ObjectID, filter *DeploymentFilter, pagination *Pagination) ([]*models.Deployment, *PaginationInfo, error)`
  - `GetByProductID(productID string, filter *DeploymentFilter, pagination *Pagination) ([]*models.Deployment, *PaginationInfo, error)`
  - `Update(id primitive.ObjectID, updates *models.Deployment) error`
  - `Delete(id primitive.ObjectID) error`
  - `CountByTenantID(tenantID primitive.ObjectID, filter *DeploymentFilter) (int64, error)`
  - `GetDeploymentsForNotification(productID string, version string) ([]*models.Deployment, error)`

- [ ] Implement `DeploymentFilter` struct:
  - `ProductID` (string)
  - `DeploymentType` (models.DeploymentType)
  - `Status` (models.DeploymentStatus)
  - `Version` (string)

- [ ] Implement MongoDB operations
- [ ] Handle database errors appropriately

#### 2.4 Customer Service
**File:** `src/backend/internal/service/customer_service.go` (new file)

- [ ] Implement `CustomerService` struct with dependencies:
  - `customerRepo` (repository.CustomerRepository)
  - `tenantRepo` (repository.TenantRepository)
  - `deploymentRepo` (repository.DeploymentRepository)
  - `auditLogService` (service.AuditLogService)

- [ ] Implement business logic methods:
  - `CreateCustomer(req *models.CreateCustomerRequest) (*models.Customer, error)`
    - Validate request
    - Generate unique customer_id
    - Create customer
    - Log audit event
  - `GetCustomer(id string) (*models.Customer, error)`
  - `ListCustomers(query *ListCustomersQuery) (*CustomerListResponse, error)`
    - Apply filters
    - Apply pagination
    - Return total count
  - `UpdateCustomer(id string, req *models.UpdateCustomerRequest) (*models.Customer, error)`
    - Validate request
    - Update customer
    - Log audit event
  - `DeleteCustomer(id string) error`
    - Check for existing tenants
    - Soft delete (set status to inactive)
    - Log audit event
  - `GetCustomerTenants(customerID string, query *ListTenantsQuery) (*TenantListResponse, error)`
  - `GetCustomerStatistics(customerID string) (*CustomerStatistics, error)`

- [ ] Implement helper methods:
  - `generateCustomerID() string`
  - `validateCustomerData(customer *models.Customer) error`

#### 2.5 Tenant Service
**File:** `src/backend/internal/service/tenant_service.go` (new file)

- [ ] Implement `TenantService` struct with dependencies:
  - `tenantRepo` (repository.TenantRepository)
  - `customerRepo` (repository.CustomerRepository)
  - `deploymentRepo` (repository.DeploymentRepository)
  - `auditLogService` (service.AuditLogService)

- [ ] Implement business logic methods:
  - `CreateTenant(customerID string, req *models.CreateTenantRequest) (*models.CustomerTenant, error)`
    - Validate customer exists
    - Generate unique tenant_id
    - Create tenant
    - Log audit event
  - `GetTenant(id string) (*models.CustomerTenant, error)`
  - `ListTenants(customerID string, query *ListTenantsQuery) (*TenantListResponse, error)`
  - `UpdateTenant(id string, req *models.UpdateTenantRequest) (*models.CustomerTenant, error)`
    - Update tenant
    - Log audit event
  - `DeleteTenant(id string) error`
    - Check for existing deployments
    - Delete tenant
    - Log audit event
  - `GetTenantDeployments(tenantID string, query *ListDeploymentsQuery) (*DeploymentListResponse, error)`
  - `GetTenantStatistics(tenantID string) (*TenantStatistics, error)`

- [ ] Implement helper methods:
  - `generateTenantID() string`
  - `validateTenantData(tenant *models.CustomerTenant) error`

#### 2.6 Deployment Service
**File:** `src/backend/internal/service/deployment_service.go` (new file)

- [ ] Implement `DeploymentService` struct with dependencies:
  - `deploymentRepo` (repository.DeploymentRepository)
  - `tenantRepo` (repository.TenantRepository)
  - `customerRepo` (repository.CustomerRepository)
  - `productService` (service.ProductService)
  - `versionService` (service.VersionService)
  - `auditLogService` (service.AuditLogService)

- [ ] Implement business logic methods:
  - `CreateDeployment(tenantID string, req *models.CreateDeploymentRequest) (*models.Deployment, error)`
    - Validate tenant exists
    - Validate product exists
    - Validate version exists
    - Check for duplicate deployment (same product + type in tenant)
    - Generate unique deployment_id
    - Create deployment
    - Log audit event
  - `GetDeployment(id string) (*models.Deployment, error)`
  - `ListDeployments(tenantID string, query *ListDeploymentsQuery) (*DeploymentListResponse, error)`
  - `UpdateDeployment(id string, req *models.UpdateDeploymentRequest) (*models.Deployment, error)`
    - Validate version if updated
    - Update deployment
    - Log audit event
  - `DeleteDeployment(id string) error`
    - Delete deployment
    - Log audit event
  - `GetAvailableUpdates(deploymentID string) ([]*models.Version, error)`
    - Get current installed version
    - Find available newer versions
    - Check compatibility
    - Return sorted list

- [ ] Implement helper methods:
  - `generateDeploymentID() string`
  - `validateDeploymentData(deployment *models.Deployment) error`
  - `checkVersionCompatibility(productID, currentVersion, targetVersion string) error`
  - `checkDuplicateDeployment(tenantID, productID string, deploymentType models.DeploymentType) error`

#### 2.7 Customer Notification Service Integration
**File:** `src/backend/internal/service/notification_service.go` (update existing)

- [ ] Add method to generate customer notifications:
  - `NotifyCustomersOnVersionRelease(productID string, versionID string) error`
    - Find all deployments for the product
    - Group by customer and tenant
    - Generate notifications per customer
    - Set priority based on deployment type
    - Include deployment and tenant details in notification

- [ ] Update notification creation to support:
  - Customer as recipient
  - Tenant-specific notifications
  - Deployment-specific notifications
  - Priority based on deployment type

#### 2.8 Customer Handler
**File:** `src/backend/internal/api/handlers/customer_handler.go` (new file)

- [ ] Implement `CustomerHandler` struct:
  - `customerService` (service.CustomerService)

- [ ] Implement HTTP handlers:
  - `CreateCustomer(w http.ResponseWriter, r *http.Request)`
  - `GetCustomer(w http.ResponseWriter, r *http.Request)`
  - `ListCustomers(w http.ResponseWriter, r *http.Request)`
  - `UpdateCustomer(w http.ResponseWriter, r *http.Request)`
  - `DeleteCustomer(w http.ResponseWriter, r *http.Request)`
  - `GetCustomerTenants(w http.ResponseWriter, r *http.Request)`
  - `GetCustomerStatistics(w http.ResponseWriter, r *http.Request)`

#### 2.9 Tenant Handler
**File:** `src/backend/internal/api/handlers/tenant_handler.go` (new file)

- [ ] Implement `TenantHandler` struct:
  - `tenantService` (service.TenantService)

- [ ] Implement HTTP handlers:
  - `CreateTenant(w http.ResponseWriter, r *http.Request)`
  - `GetTenant(w http.ResponseWriter, r *http.Request)`
  - `ListTenants(w http.ResponseWriter, r *http.Request)`
  - `UpdateTenant(w http.ResponseWriter, r *http.Request)`
  - `DeleteTenant(w http.ResponseWriter, r *http.Request)`
  - `GetTenantDeployments(w http.ResponseWriter, r *http.Request)`
  - `GetTenantStatistics(w http.ResponseWriter, r *http.Request)`

#### 2.10 Deployment Handler
**File:** `src/backend/internal/api/handlers/deployment_handler.go` (new file)

- [ ] Implement `DeploymentHandler` struct:
  - `deploymentService` (service.DeploymentService)

- [ ] Implement HTTP handlers:
  - `CreateDeployment(w http.ResponseWriter, r *http.Request)`
  - `GetDeployment(w http.ResponseWriter, r *http.Request)`
  - `ListDeployments(w http.ResponseWriter, r *http.Request)`
  - `UpdateDeployment(w http.ResponseWriter, r *http.Request)`
  - `DeleteDeployment(w http.ResponseWriter, r *http.Request)`
  - `GetAvailableUpdates(w http.ResponseWriter, r *http.Request)`

#### 2.11 API Routing
**File:** `src/backend/internal/api/router/router.go` (update existing)

- [ ] Add customer routes:
  - `GET /api/v1/customers` - List customers
  - `POST /api/v1/customers` - Create customer
  - `GET /api/v1/customers/:id` - Get customer
  - `PUT /api/v1/customers/:id` - Update customer
  - `DELETE /api/v1/customers/:id` - Delete customer
  - `GET /api/v1/customers/:id/tenants` - List customer tenants
  - `GET /api/v1/customers/:id/statistics` - Get customer statistics

- [ ] Add tenant routes:
  - `POST /api/v1/customers/:customer_id/tenants` - Create tenant
  - `GET /api/v1/customers/:customer_id/tenants/:id` - Get tenant
  - `PUT /api/v1/customers/:customer_id/tenants/:id` - Update tenant
  - `DELETE /api/v1/customers/:customer_id/tenants/:id` - Delete tenant
  - `GET /api/v1/customers/:customer_id/tenants/:id/deployments` - List tenant deployments
  - `GET /api/v1/customers/:customer_id/tenants/:id/statistics` - Get tenant statistics

- [ ] Add deployment routes:
  - `POST /api/v1/customers/:customer_id/tenants/:tenant_id/deployments` - Create deployment
  - `GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id` - Get deployment
  - `PUT /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id` - Update deployment
  - `DELETE /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id` - Delete deployment
  - `GET /api/v1/customers/:customer_id/tenants/:tenant_id/deployments/:id/updates` - Get available updates

- [ ] Register handlers in router

#### 2.12 Service Factory Update
**File:** `src/backend/internal/service/service_factory.go` (update existing)

- [ ] Add customer service to factory
- [ ] Add tenant service to factory
- [ ] Add deployment service to factory
- [ ] Initialize repositories
- [ ] Wire dependencies

### Test Cases

#### 2.13 Repository Tests
**Files:**
- `src/backend/internal/repository/customer_repository_test.go`
- `src/backend/internal/repository/tenant_repository_test.go`
- `src/backend/internal/repository/deployment_repository_test.go`

- [ ] Test Customer Repository (CRUD operations, filters, pagination)
- [ ] Test Tenant Repository (CRUD operations, filters, pagination)
- [ ] Test Deployment Repository (CRUD operations, filters, pagination, notification queries)
- [ ] Test error handling

#### 2.14 Service Tests
**Files:**
- `src/backend/internal/service/customer_service_test.go`
- `src/backend/internal/service/tenant_service_test.go`
- `src/backend/internal/service/deployment_service_test.go`

- [ ] Test Customer Service (business logic, validation, audit logging)
- [ ] Test Tenant Service (business logic, validation, audit logging)
- [ ] Test Deployment Service (business logic, validation, version checking, audit logging)
- [ ] Test duplicate deployment prevention
- [ ] Test version compatibility checking

#### 2.15 Handler Tests
**Files:**
- `src/backend/internal/api/handlers/customer_handler_test.go`
- `src/backend/internal/api/handlers/tenant_handler_test.go`
- `src/backend/internal/api/handlers/deployment_handler_test.go`

- [ ] Test all HTTP handlers (success cases, validation errors, not found, etc.)
- [ ] Test request/response formats
- [ ] Test error responses

#### 2.16 Integration Tests
**File:** `src/backend/tests/integration/customer_management_test.go` (new file)

- [ ] End-to-end API tests:
  - Create customer → Create tenant → Create deployment → Get deployment
  - List operations with pagination and filters
  - Update operations
  - Delete operations (with cascade checks)
  - Get available updates
  - Notification generation on version release

- [ ] Database transaction tests
- [ ] Concurrent request handling
- [ ] Error scenarios

### Acceptance Criteria
- ✅ All repository methods implemented and tested
- ✅ All service methods implemented with business logic
- ✅ All API handlers implemented
- ✅ Routes registered and working
- ✅ Unit tests passing (80%+ coverage)
- ✅ Integration tests passing
- ✅ API documentation updated
- ✅ Error handling implemented
- ✅ Audit logging integrated
- ✅ Duplicate deployment prevention working

---

## Phase 3: UI Components (Frontend)

**Duration:** 4-5 weeks  
**Priority:** High  
**Dependencies:** Phase 2 complete

### Goals
- Create customer management UI components
- Create tenant management UI components
- Create deployment management UI components
- Implement customer dashboard
- Integrate with notification system
- Add E2E tests

### Tasks

#### 3.1 TypeScript Types
**File:** `src/frontend/src/types/index.ts` (update existing)

- [ ] Add Customer type
- [ ] Add CustomerTenant type
- [ ] Add Deployment type
- [ ] Add request/response types
- [ ] Add filter types

#### 3.2 API Services
**File:** `src/frontend/src/services/api/customers.ts` (new file)

- [ ] Implement `customersApi`:
  - `getAll(query?: ListCustomersQuery): Promise<PaginatedResponse<Customer>>`
  - `getById(id: string): Promise<Customer>`
  - `create(data: CreateCustomerRequest): Promise<Customer>`
  - `update(id: string, data: UpdateCustomerRequest): Promise<Customer>`
  - `delete(id: string): Promise<void>`
  - `getTenants(customerId: string, query?: ListTenantsQuery): Promise<PaginatedResponse<CustomerTenant>>`
  - `getStatistics(customerId: string): Promise<CustomerStatistics>`

**File:** `src/frontend/src/services/api/tenants.ts` (new file)

- [ ] Implement `tenantsApi`:
  - `create(customerId: string, data: CreateTenantRequest): Promise<CustomerTenant>`
  - `getById(customerId: string, tenantId: string): Promise<CustomerTenant>`
  - `update(customerId: string, tenantId: string, data: UpdateTenantRequest): Promise<CustomerTenant>`
  - `delete(customerId: string, tenantId: string): Promise<void>`
  - `getDeployments(customerId: string, tenantId: string, query?: ListDeploymentsQuery): Promise<PaginatedResponse<Deployment>>`
  - `getStatistics(customerId: string, tenantId: string): Promise<TenantStatistics>`

**File:** `src/frontend/src/services/api/deployments.ts` (new file)

- [ ] Implement `deploymentsApi`:
  - `create(customerId: string, tenantId: string, data: CreateDeploymentRequest): Promise<Deployment>`
  - `getById(customerId: string, tenantId: string, deploymentId: string): Promise<Deployment>`
  - `update(customerId: string, tenantId: string, deploymentId: string, data: UpdateDeploymentRequest): Promise<Deployment>`
  - `delete(customerId: string, tenantId: string, deploymentId: string): Promise<void>`
  - `getAvailableUpdates(customerId: string, tenantId: string, deploymentId: string): Promise<Version[]>`

#### 3.3 Customer List Page
**File:** `src/frontend/src/pages/Customers.tsx` (new file)

- [ ] Create customers list page with:
  - Search bar
  - Filter by status
  - Table showing:
    - Customer ID
    - Name
    - Organization
    - Email
    - Status badge
    - Number of tenants
    - Actions (View, Edit, Delete)
  - Pagination
  - "Create Customer" button
  - Empty state

#### 3.4 Customer Details Page
**File:** `src/frontend/src/pages/CustomerDetails.tsx` (new file)

- [ ] Create customer details page with tabs:
  - **Overview Tab:**
    - Customer information
    - Account status
    - Notification preferences
    - Statistics (tenants, deployments, users)
    - Edit button
  - **Tenants Tab:**
    - List of tenants
    - Filter by status
    - "Add Tenant" button
    - Tenant cards/list
  - **Activity Tab:**
    - Recent notifications
    - Audit logs

#### 3.5 Customer Form Component
**File:** `src/frontend/src/components/customers/CustomerForm.tsx` (new file)

- [ ] Create form for create/edit customer (same as before)

#### 3.6 Tenant List Component
**File:** `src/frontend/src/components/tenants/TenantsList.tsx` (new file)

- [ ] Create tenants list component:
  - Filter by status
  - Table/cards showing:
    - Tenant name
    - Description
    - Status badge
    - Number of deployments
    - Actions
  - Pagination
  - Empty state

#### 3.7 Tenant Details Component
**File:** `src/frontend/src/components/tenants/TenantDetails.tsx` (new file)

- [ ] Create tenant details component:
  - Tenant information
  - Statistics (deployments, users)
  - List of deployments
  - Edit button
  - "Add Deployment" button

#### 3.8 Tenant Form Component
**File:** `src/frontend/src/components/tenants/TenantForm.tsx` (new file)

- [ ] Create form for create/edit tenant:
  - Tenant name (required)
  - Description (optional)
  - Status (dropdown)
  - Form validation
  - Submit/Cancel buttons

#### 3.9 Deployment List Component
**File:** `src/frontend/src/components/deployments/DeploymentsList.tsx` (new file)

- [ ] Create deployments list component:
  - Filter by product, type, status
  - Search functionality
  - Table/cards showing:
    - Product name
    - Deployment type badge (UAT/Production)
    - Installed version
    - Number of users
    - Status
    - Actions
  - Pagination
  - Empty state

#### 3.10 Deployment Form Component
**File:** `src/frontend/src/components/deployments/DeploymentForm.tsx` (new file)

- [ ] Create form for create/edit deployment:
  - Product (dropdown, required)
  - Deployment type (UAT/Testing/Production, required)
  - Installed version (dropdown, required)
  - Number of users (number input, optional)
  - License information (textarea, optional)
  - Server hostname (optional)
  - Environment details (textarea, optional)
  - Form validation
  - Duplicate deployment check (same product + type)
  - Submit/Cancel buttons

#### 3.11 Deployment Details Component
**File:** `src/frontend/src/components/deployments/DeploymentDetails.tsx` (new file)

- [ ] Create deployment details component:
  - Deployment information (tenant, product, type)
  - Version information with comparison
  - Available updates section
  - Update history
  - Edit button
  - Actions:
    - Update version
    - View release notes
    - Check for updates

#### 3.12 Customer Dashboard Component
**File:** `src/frontend/src/components/customers/CustomerDashboard.tsx` (new file)

- [ ] Create customer dashboard:
  - Statistics cards:
    - Total tenants
    - Total deployments
    - Total users
    - Deployments by product (chart/list)
    - Deployments by type (UAT/Production)
    - Deployments by tenant
    - Pending updates count
  - Recent updates section
  - Quick actions
  - License summary

#### 3.13 Deployment Type Badge Component
**File:** `src/frontend/src/components/deployments/DeploymentTypeBadge.tsx` (new file)

- [ ] Create badge component:
  - Visual distinction for UAT/Testing/Production
  - Color coding
  - Icons (optional)

#### 3.14 Update Notification Integration
**File:** `src/frontend/src/services/api/notifications.ts` (update existing)

- [ ] Update notification service to handle customer notifications
- [ ] Add method to fetch customer-specific notifications

**File:** `src/frontend/src/hooks/useNotifications.ts` (update existing)

- [ ] Update to support customer notifications
- [ ] Filter by deployment type

#### 3.15 Routing
**File:** `src/frontend/src/router/index.tsx` (update existing)

- [ ] Add routes:
  - `/customers` - Customers list
  - `/customers/new` - Create customer
  - `/customers/:id` - Customer details
  - `/customers/:id/edit` - Edit customer
  - `/customers/:id/tenants/new` - Create tenant
  - `/customers/:id/tenants/:tenantId` - Tenant details
  - `/customers/:id/tenants/:tenantId/edit` - Edit tenant
  - `/customers/:id/tenants/:tenantId/deployments/new` - Create deployment
  - `/customers/:id/tenants/:tenantId/deployments/:deploymentId` - Deployment details
  - `/customers/:id/tenants/:tenantId/deployments/:deploymentId/edit` - Edit deployment

#### 3.16 Navigation Updates
**File:** `src/frontend/src/components/layout/Sidebar.tsx` (update existing)

- [ ] Add "Customers" menu item
- [ ] Add submenu for customer management

### Test Cases

#### 3.17 Component Unit Tests
**Files:**
- `src/frontend/src/components/customers/__tests__/CustomerForm.test.tsx`
- `src/frontend/src/components/tenants/__tests__/TenantForm.test.tsx`
- `src/frontend/src/components/deployments/__tests__/DeploymentForm.test.tsx`
- `src/frontend/src/components/customers/__tests__/CustomerDashboard.test.tsx`

- [ ] Test all form components (validation, field updates, submit handling)
- [ ] Test dashboard components (statistics calculation, data display)
- [ ] Test list components (filtering, pagination)

#### 3.18 Integration Tests
**Files:**
- `src/frontend/src/pages/__tests__/Customers.test.tsx`
- `src/frontend/src/pages/__tests__/CustomerDetails.test.tsx`

- [ ] Test customer list page
- [ ] Test customer details page
- [ ] Test tenant management
- [ ] Test deployment management

#### 3.19 E2E Tests
**File:** `src/frontend/tests/e2e/customer-management.spec.ts` (new file)

- [ ] Test customer management flow:
  - Navigate to customers page
  - Create new customer
  - View customer details
  - Edit customer
  - Delete customer
  - Search customers
  - Filter by status

- [ ] Test tenant management flow:
  - Create tenant for customer
  - View tenant details
  - Edit tenant
  - Delete tenant
  - View tenant deployments

- [ ] Test deployment management flow:
  - Create deployment for tenant
  - View deployment details
  - Update deployment version
  - Update number of users and license
  - View available updates
  - Delete deployment
  - Filter deployments
  - Test duplicate deployment prevention

- [ ] Test customer dashboard:
  - View statistics
  - View deployments by product
  - View deployments by tenant
  - View pending updates

- [ ] Test notification integration:
  - Receive notification on version release
  - View deployment-specific notifications
  - Notification priority display

### Acceptance Criteria
- ✅ All UI components created and styled
- ✅ Forms with validation
- ✅ Customer CRUD operations working
- ✅ Tenant CRUD operations working
- ✅ Deployment CRUD operations working
- ✅ Duplicate deployment prevention in UI
- ✅ Customer dashboard displaying data
- ✅ Integration with notification system
- ✅ Responsive design
- ✅ Unit tests passing (80%+ coverage)
- ✅ E2E tests passing
- ✅ Error handling and loading states
- ✅ Accessibility (ARIA labels, keyboard navigation)

---

## Integration Points

### Notification System Integration
- When a new version is released, the system should:
  1. Find all deployments for that product
  2. Group by customer and tenant
  3. Generate notifications per customer
  4. Include deployment and tenant details in notification
  5. Set priority based on deployment type

### Audit Logging Integration
- All customer, tenant, and deployment operations should be logged:
  - Create, update, delete operations
  - Version updates
  - Status changes

### Version Service Integration
- Deployment service should:
  - Validate versions exist
  - Check version compatibility
  - Get available updates

---

## Dependencies

### External Dependencies
- MongoDB (database)
- React 19+ (frontend)
- Go 1.21+ (backend)

### Internal Dependencies
- Product Service (for product validation)
- Version Service (for version validation and updates)
- Notification Service (for customer notifications)
- Audit Log Service (for logging operations)

---

## Timeline Estimate

- **Phase 1 (Data Models):** 1-2 weeks
- **Phase 2 (API Layer):** 3-4 weeks
- **Phase 3 (UI Components):** 4-5 weeks
- **Total:** 8-11 weeks

---

## Risk Mitigation

### Risks
1. **Database performance** with large number of deployments
   - Mitigation: Proper indexing, pagination, query optimization

2. **Notification system overload** when many customers have deployments
   - Mitigation: Batch processing, queue system, rate limiting

3. **Version compatibility checking** complexity
   - Mitigation: Reuse existing compatibility service, thorough testing

4. **UI performance** with large customer/tenant/deployment lists
   - Mitigation: Virtual scrolling, pagination, lazy loading

5. **Duplicate deployment prevention** (same product + type in tenant)
   - Mitigation: Database unique constraint, UI validation, clear error messages

---

## Success Metrics

- ✅ All requirements from Section 2.6 implemented
- ✅ 80%+ test coverage for backend
- ✅ 80%+ test coverage for frontend
- ✅ All E2E tests passing
- ✅ API response times < 200ms for list operations
- ✅ Zero critical bugs in production
- ✅ Documentation complete
- ✅ Duplicate deployment prevention working

---

## Document Version History

- **v2.0** (17-Nov-2025): Updated to reflect Customer → Tenant → Deployment model structure
- **v1.0** (17-Nov-2025): Initial implementation plan created

---

**Status:** Ready for Implementation  
**Next Steps:** Begin Phase 1 - Data Models implementation
