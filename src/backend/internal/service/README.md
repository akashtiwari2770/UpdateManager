# Service Layer

The service layer contains business logic and orchestrates repository operations.

## Services

### 1. ProductService
- **File**: `product_service.go`
- **Dependencies**: ProductRepository, AuditLogRepository
- **Methods**:
  - `CreateProduct()` - Creates product with validation
  - `GetProduct()` - Retrieves product by ID
  - `GetProductByProductID()` - Retrieves by product_id
  - `UpdateProduct()` - Updates product with validation
  - `DeleteProduct()` - Soft delete (sets IsActive=false)
  - `ListProducts()` - Lists products with pagination
  - `GetActiveProducts()` - Gets only active products

### 2. VersionService
- **File**: `version_service.go`
- **Dependencies**: VersionRepository, ProductRepository, AuditLogRepository
- **Methods**:
  - `CreateVersion()` - Creates version with validation
  - `GetVersion()` - Retrieves version by ID
  - `GetVersionByProductAndVersion()` - Retrieves by product_id and version_number
  - `GetVersionsByProduct()` - Lists versions for a product
  - `UpdateVersion()` - Updates draft versions only
  - `ApproveVersion()` - Approves pending versions
  - `SubmitForReview()` - Submits draft for review
  - `ReleaseVersion()` - Releases approved versions
  - `GetVersionsByState()` - Gets versions by state

### 3. CompatibilityService
- **File**: `compatibility_service.go`
- **Dependencies**: CompatibilityRepository, VersionRepository, AuditLogRepository
- **Methods**:
  - `ValidateCompatibility()` - Validates/creates compatibility matrix
  - `GetCompatibility()` - Retrieves compatibility matrix
  - `ListCompatibility()` - Lists compatibility matrices

### 4. UpgradePathService
- **File**: `upgrade_path_service.go`
- **Dependencies**: UpgradePathRepository, VersionRepository
- **Methods**:
  - `CreateUpgradePath()` - Creates upgrade path with validation
  - `GetUpgradePath()` - Retrieves upgrade path
  - `GetUpgradePathsByProduct()` - Lists upgrade paths for product
  - `BlockUpgradePath()` - Blocks an upgrade path

### 5. NotificationService
- **File**: `notification_service.go`
- **Dependencies**: NotificationRepository
- **Methods**:
  - `CreateNotification()` - Creates notification
  - `GetNotifications()` - Gets notifications for recipient (with unread filter)
  - `MarkAsRead()` - Marks notification as read
  - `MarkAllAsRead()` - Marks all notifications as read
  - `GetUnreadCount()` - Gets count of unread notifications

### 6. UpdateDetectionService
- **File**: `update_detection_service.go`
- **Dependencies**: UpdateDetectionRepository, VersionRepository
- **Methods**:
  - `DetectUpdate()` - Creates or updates detection
  - `GetDetection()` - Retrieves detection
  - `UpdateAvailableVersion()` - Updates available version
  - `ListDetections()` - Lists detections with filters

### 7. UpdateRolloutService
- **File**: `update_rollout_service.go`
- **Dependencies**: UpdateRolloutRepository, UpdateDetectionRepository
- **Methods**:
  - `InitiateRollout()` - Initiates new rollout
  - `ListRollouts()` - Lists rollouts with filters
  - (Other methods need ObjectID conversion - to be implemented)

### 8. AuditLogService
- **File**: `audit_log_service.go`
- **Dependencies**: AuditLogRepository
- **Methods**:
  - `GetAuditLogsByResource()` - Gets logs for a resource
  - `GetAuditLogsByUser()` - Gets logs for a user
  - `GetAuditLogsByAction()` - Gets logs by action
  - `ListAuditLogs()` - Lists logs with filters

## Service Factory

The `ServiceFactory` initializes all services with their dependencies:

```go
import (
    "updatemanager/pkg/database"
    "updatemanager/internal/service"
)

// Connect to database
db, err := database.Connect(ctx, database.DefaultConfig())
if err != nil {
    log.Fatal(err)
}

// Create service factory
services := service.NewServiceFactory(db.Database)

// Use services
product, err := services.ProductService.CreateProduct(ctx, req, userID, userEmail)
```

## Business Logic Features

### Validation
- Product ID uniqueness
- Version uniqueness per product
- Version state transitions
- Upgrade path validation

### Audit Logging
- Automatic audit logging for create/update/delete operations
- Tracks user actions with details

### State Management
- Version state machine (draft → pending_review → approved → released)
- Rollout status management
- Notification read/unread tracking

### Soft Deletes
- Products are soft-deleted (IsActive=false) instead of hard delete

## Usage Example

```go
// Create a product
req := &models.CreateProductRequest{
    ProductID: "my-product",
    Name:      "My Product",
    Type:      models.ProductTypeServer,
}
product, err := services.ProductService.CreateProduct(ctx, req, "user-123", "user@example.com")

// Create a version
versionReq := &models.CreateVersionRequest{
    VersionNumber: "1.0.0",
    ReleaseType:   models.ReleaseTypeFeature,
    ReleaseDate:   time.Now(),
}
version, err := services.VersionService.CreateVersion(ctx, product.ProductID, versionReq, "user-123")

// Submit for review
version, err = services.VersionService.SubmitForReview(ctx, version.ID, "user-123")

// Approve version
approveReq := &models.ApproveVersionRequest{ApprovedBy: "admin-123"}
version, err = services.VersionService.ApproveVersion(ctx, version.ID, approveReq)
```

