# Product Release Management API Specification

## Overview

This document defines the RESTful API specification for the Product Release Management system. The API supports both Phase 1 and Phase 2 implementation phases as defined in the Product Release Process document.

**Base URL**: `/api/v1`

**Authentication**: Bearer Token (JWT)

**Content-Type**: `application/json`

## Data Models

### Product Model

```go
type Product struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID       string            `bson:"product_id" json:"product_id" validate:"required"`
    Name            string            `bson:"name" json:"name" validate:"required"`
    Type            ProductType       `bson:"type" json:"type" validate:"required"`
    Description     string            `bson:"description" json:"description"`
    Vendor          string            `bson:"vendor" json:"vendor"`
    CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt       time.Time         `bson:"updated_at" json:"updated_at"`
    IsActive        bool              `bson:"is_active" json:"is_active"`
}

type ProductType string

const (
    ProductTypeServer ProductType = "server"
    ProductTypeClient ProductType = "client"
)
```

### Version Model

```go
type Version struct {
    ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID               string            `bson:"product_id" json:"product_id" validate:"required"`
    VersionNumber           string            `bson:"version_number" json:"version_number" validate:"required"`
    ReleaseDate             time.Time         `bson:"release_date" json:"release_date"`
    ReleaseType             ReleaseType       `bson:"release_type" json:"release_type" validate:"required"`
    State                   VersionState      `bson:"state" json:"state" validate:"required"`
    EOLDate                 *time.Time        `bson:"eol_date,omitempty" json:"eol_date,omitempty"`
    
    // Compatibility (for clients)
    MinServerVersion        string            `bson:"min_server_version,omitempty" json:"min_server_version,omitempty"`
    MaxServerVersion        string            `bson:"max_server_version,omitempty" json:"max_server_version,omitempty"`
    RecommendedServerVersion string          `bson:"recommended_server_version,omitempty" json:"recommended_server_version,omitempty"`
    
    // Release Notes
    ReleaseNotes            *ReleaseNotes     `bson:"release_notes,omitempty" json:"release_notes,omitempty"`
    
    // Packages
    Packages                []PackageInfo     `bson:"packages" json:"packages"`
    
    // Approval
    ApprovedBy             string            `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
    ApprovedAt             *time.Time        `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
    CreatedBy               string            `bson:"created_by" json:"created_by"`
    CreatedAt               time.Time         `bson:"created_at" json:"created_at"`
    UpdatedAt               time.Time         `bson:"updated_at" json:"updated_at"`
}

type ReleaseType string

const (
    ReleaseTypeSecurity    ReleaseType = "security"
    ReleaseTypeFeature     ReleaseType = "feature"
    ReleaseTypeMaintenance ReleaseType = "maintenance"
    ReleaseTypeMajor       ReleaseType = "major"
)

type VersionState string

const (
    VersionStateDraft         VersionState = "draft"
    VersionStatePendingReview VersionState = "pending_review"
    VersionStateApproved      VersionState = "approved"
    VersionStateReleased      VersionState = "released"
    VersionStateDeprecated    VersionState = "deprecated"
    VersionStateEOL           VersionState = "eol"
)
```

### Release Notes Model

```go
type ReleaseNotes struct {
    VersionInfo     VersionInfoSection     `bson:"version_info" json:"version_info"`
    WhatsNew        []string               `bson:"whats_new" json:"whats_new"`
    BugFixes        []BugFix               `bson:"bug_fixes" json:"bug_fixes"`
    BreakingChanges []BreakingChange       `bson:"breaking_changes" json:"breaking_changes"`
    Compatibility   CompatibilitySection   `bson:"compatibility" json:"compatibility"`
    UpgradeInstructions string            `bson:"upgrade_instructions" json:"upgrade_instructions"`
    KnownIssues     []KnownIssue           `bson:"known_issues" json:"known_issues"`
}

type VersionInfoSection struct {
    VersionNumber string    `bson:"version_number" json:"version_number"`
    ReleaseDate   time.Time `bson:"release_date" json:"release_date"`
    ReleaseType   ReleaseType `bson:"release_type" json:"release_type"`
}

type BugFix struct {
    ID          string `bson:"id" json:"id"`
    Description string `bson:"description" json:"description"`
    IssueNumber string `bson:"issue_number,omitempty" json:"issue_number,omitempty"`
}

type BreakingChange struct {
    Description      string `bson:"description" json:"description"`
    MigrationSteps   string `bson:"migration_steps" json:"migration_steps"`
    ConfigurationChanges string `bson:"configuration_changes" json:"configuration_changes"`
}

type CompatibilitySection struct {
    ServerVersionRequirements string   `bson:"server_version_requirements" json:"server_version_requirements"`
    ClientVersionRequirements string   `bson:"client_version_requirements" json:"client_version_requirements"`
    OSRequirements           []string `bson:"os_requirements" json:"os_requirements"`
}

type KnownIssue struct {
    ID          string `bson:"id" json:"id"`
    Description string `bson:"description" json:"description"`
    Workaround  string `bson:"workaround,omitempty" json:"workaround,omitempty"`
    PlannedFix  string `bson:"planned_fix,omitempty" json:"planned_fix,omitempty"`
}
```

### Package Model

```go
type PackageInfo struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    PackageType     PackageType        `bson:"package_type" json:"package_type" validate:"required"`
    FileName        string            `bson:"file_name" json:"file_name" validate:"required"`
    FileSize        int64             `bson:"file_size" json:"file_size"`
    DownloadURL     string            `bson:"download_url" json:"download_url"`
    ChecksumSHA256  string            `bson:"checksum_sha256" json:"checksum_sha256" validate:"required"`
    DigitalSignature string           `bson:"digital_signature,omitempty" json:"digital_signature,omitempty"`
    OS              string            `bson:"os,omitempty" json:"os,omitempty"`
    Architecture    string            `bson:"architecture,omitempty" json:"architecture,omitempty"`
    UploadedAt      time.Time         `bson:"uploaded_at" json:"uploaded_at"`
    UploadedBy      string            `bson:"uploaded_by" json:"uploaded_by"`
}

type PackageType string

const (
    PackageTypeFullInstaller PackageType = "full_installer"
    PackageTypeUpdate        PackageType = "update"
    PackageTypeDelta         PackageType = "delta"
    PackageTypeRollback      PackageType = "rollback"
)
```

### Compatibility Model

```go
type CompatibilityMatrix struct {
    ID                      primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID               string            `bson:"product_id" json:"product_id" validate:"required"`
    VersionNumber           string            `bson:"version_number" json:"version_number" validate:"required"`
    MinServerVersion        string            `bson:"min_server_version,omitempty" json:"min_server_version,omitempty"`
    MaxServerVersion        string            `bson:"max_server_version,omitempty" json:"max_server_version,omitempty"`
    RecommendedServerVersion string          `bson:"recommended_server_version,omitempty" json:"recommended_server_version,omitempty"`
    IncompatibleVersions    []string         `bson:"incompatible_versions" json:"incompatible_versions"`
    ValidatedAt             time.Time        `bson:"validated_at" json:"validated_at"`
    ValidatedBy             string           `bson:"validated_by" json:"validated_by"`
    ValidationStatus        ValidationStatus `bson:"validation_status" json:"validation_status"`
    ValidationErrors        []string         `bson:"validation_errors" json:"validation_errors"`
}

type ValidationStatus string

const (
    ValidationStatusPending   ValidationStatus = "pending"
    ValidationStatusPassed    ValidationStatus = "passed"
    ValidationStatusFailed    ValidationStatus = "failed"
    ValidationStatusSkipped   ValidationStatus = "skipped"
)
```

### Upgrade Path Model

```go
type UpgradePath struct {
    ID                  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID           string            `bson:"product_id" json:"product_id" validate:"required"`
    FromVersion         string            `bson:"from_version" json:"from_version" validate:"required"`
    ToVersion           string            `bson:"to_version" json:"to_version" validate:"required"`
    PathType            UpgradePathType   `bson:"path_type" json:"path_type" validate:"required"`
    IntermediateVersions []string         `bson:"intermediate_versions,omitempty" json:"intermediate_versions,omitempty"`
    IsBlocked           bool              `bson:"is_blocked" json:"is_blocked"`
    BlockReason         string            `bson:"block_reason,omitempty" json:"block_reason,omitempty"`
    CreatedAt           time.Time         `bson:"created_at" json:"created_at"`
}

type UpgradePathType string

const (
    UpgradePathTypeDirect    UpgradePathType = "direct"
    UpgradePathTypeMultiStep UpgradePathType = "multi_step"
    UpgradePathTypeBlocked   UpgradePathType = "blocked"
)
```

### Notification Model

```go
type Notification struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Type            NotificationType   `bson:"type" json:"type" validate:"required"`
    RecipientID     string            `bson:"recipient_id" json:"recipient_id" validate:"required"`
    ProductID       string            `bson:"product_id,omitempty" json:"product_id,omitempty"`
    VersionID       string            `bson:"version_id,omitempty" json:"version_id,omitempty"`
    Title           string            `bson:"title" json:"title" validate:"required"`
    Message         string            `bson:"message" json:"message" validate:"required"`
    Priority        NotificationPriority `bson:"priority" json:"priority"`
    IsRead          bool              `bson:"is_read" json:"is_read"`
    ReadAt          *time.Time        `bson:"read_at,omitempty" json:"read_at,omitempty"`
    CreatedAt       time.Time         `bson:"created_at" json:"created_at"`
}

type NotificationType string

const (
    NotificationTypeNewVersion      NotificationType = "new_version"
    NotificationTypeSecurityRelease  NotificationType = "security_release"
    NotificationTypeEOLWarning     NotificationType = "eol_warning"
    NotificationTypeUpdateAvailable NotificationType = "update_available"
)

type NotificationPriority string

const (
    NotificationPriorityLow      NotificationPriority = "low"
    NotificationPriorityNormal   NotificationPriority = "normal"
    NotificationPriorityHigh     NotificationPriority = "high"
    NotificationPriorityCritical NotificationPriority = "critical"
)
```

### Update Detection Model

```go
type UpdateDetection struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    EndpointID      string            `bson:"endpoint_id" json:"endpoint_id" validate:"required"`
    ProductID       string            `bson:"product_id" json:"product_id" validate:"required"`
    CurrentVersion  string            `bson:"current_version" json:"current_version" validate:"required"`
    AvailableVersion string           `bson:"available_version" json:"available_version" validate:"required"`
    DetectedAt      time.Time         `bson:"detected_at" json:"detected_at"`
    LastCheckedAt   time.Time         `bson:"last_checked_at" json:"last_checked_at"`
}
```

### Update Rollout Model

```go
type UpdateRollout struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    EndpointID      string            `bson:"endpoint_id" json:"endpoint_id" validate:"required"`
    ProductID       string            `bson:"product_id" json:"product_id" validate:"required"`
    FromVersion     string            `bson:"from_version" json:"from_version" validate:"required"`
    ToVersion       string            `bson:"to_version" json:"to_version" validate:"required"`
    Status          RolloutStatus     `bson:"status" json:"status" validate:"required"`
    InitiatedBy     string            `bson:"initiated_by" json:"initiated_by" validate:"required"`
    InitiatedAt     time.Time         `bson:"initiated_at" json:"initiated_at"`
    StartedAt       *time.Time        `bson:"started_at,omitempty" json:"started_at,omitempty"`
    CompletedAt     *time.Time        `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
    FailedAt        *time.Time        `bson:"failed_at,omitempty" json:"failed_at,omitempty"`
    ErrorMessage    string            `bson:"error_message,omitempty" json:"error_message,omitempty"`
    Progress        int               `bson:"progress" json:"progress"` // 0-100
}

type RolloutStatus string

const (
    RolloutStatusPending   RolloutStatus = "pending"
    RolloutStatusInProgress RolloutStatus = "in_progress"
    RolloutStatusCompleted RolloutStatus = "completed"
    RolloutStatusFailed    RolloutStatus = "failed"
    RolloutStatusCancelled RolloutStatus = "cancelled"
)
```

### Audit Log Model

```go
type AuditLog struct {
    ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    Action          AuditAction       `bson:"action" json:"action" validate:"required"`
    ResourceType    string            `bson:"resource_type" json:"resource_type" validate:"required"`
    ResourceID      string            `bson:"resource_id" json:"resource_id" validate:"required"`
    UserID          string            `bson:"user_id" json:"user_id" validate:"required"`
    UserEmail       string            `bson:"user_email" json:"user_email"`
    Details         map[string]interface{} `bson:"details" json:"details"`
    IPAddress       string            `bson:"ip_address" json:"ip_address"`
    UserAgent       string            `bson:"user_agent" json:"user_agent"`
    Timestamp       time.Time         `bson:"timestamp" json:"timestamp"`
}

type AuditAction string

const (
    AuditActionCreate   AuditAction = "create"
    AuditActionUpdate   AuditAction = "update"
    AuditActionDelete   AuditAction = "delete"
    AuditActionApprove  AuditAction = "approve"
    AuditActionReject   AuditAction = "reject"
    AuditActionRelease  AuditAction = "release"
    AuditActionUpload   AuditAction = "upload"
    AuditActionDownload AuditAction = "download"
)
```

## API Endpoints

### Products API

#### GET /products
List all products

**Query Parameters:**
- `type` (optional): Filter by product type (server/client)
- `is_active` (optional): Filter by active status
- `page` (optional): Page number (default: 1)
- `limit` (optional): Items per page (default: 20)

**Response:**
```json
{
  "products": [
    {
      "id": "507f1f77bcf86cd799439011",
      "product_id": "hyworks",
      "name": "HyWorks",
      "type": "server",
      "description": "Application Virtualization Platform",
      "vendor": "Accops",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-15T10:00:00Z",
      "is_active": true
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 100,
    "total_pages": 5
  }
}
```

#### GET /products/{product_id}
Get product details

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439011",
  "product_id": "hyworks",
  "name": "HyWorks",
  "type": "server",
  "description": "Application Virtualization Platform",
  "vendor": "Accops",
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-15T10:00:00Z",
  "is_active": true
}
```

#### POST /products
Create a new product

**Request Body:**
```json
{
  "product_id": "hyworks",
  "name": "HyWorks",
  "type": "server",
  "description": "Application Virtualization Platform",
  "vendor": "Accops"
}
```

**Response:** 201 Created with product object

### Versions API

#### GET /products/{product_id}/versions
List all versions for a product

**Query Parameters:**
- `state` (optional): Filter by version state
- `release_type` (optional): Filter by release type
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response:**
```json
{
  "versions": [
    {
      "id": "507f1f77bcf86cd799439012",
      "product_id": "hyworks",
      "version_number": "2.1.0",
      "release_date": "2025-01-20T10:00:00Z",
      "release_type": "feature",
      "state": "released",
      "created_at": "2025-01-15T10:00:00Z",
      "updated_at": "2025-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 50,
    "total_pages": 3
  }
}
```

#### GET /products/{product_id}/versions/{version_number}
Get version details

**Response:**
```json
{
  "id": "507f1f77bcf86cd799439012",
  "product_id": "hyworks",
  "version_number": "2.1.0",
  "release_date": "2025-01-20T10:00:00Z",
  "release_type": "feature",
  "state": "released",
  "release_notes": {
    "version_info": {
      "version_number": "2.1.0",
      "release_date": "2025-01-20T10:00:00Z",
      "release_type": "feature"
    },
    "whats_new": ["New feature 1", "New feature 2"],
    "bug_fixes": [
      {
        "id": "BUG-123",
        "description": "Fixed issue with authentication",
        "issue_number": "GH-456"
      }
    ]
  },
  "packages": [
    {
      "id": "507f1f77bcf86cd799439013",
      "package_type": "full_installer",
      "file_name": "hyworks-2.1.0-installer.deb",
      "file_size": 52428800,
      "download_url": "https://cdn.accops.com/packages/hyworks-2.1.0-installer.deb",
      "checksum_sha256": "abc123...",
      "os": "linux",
      "architecture": "amd64"
    }
  ],
  "created_at": "2025-01-15T10:00:00Z",
  "updated_at": "2025-01-20T10:00:00Z"
}
```

#### POST /products/{product_id}/versions
Create a new version (Phase 1)

**Request Body:**
```json
{
  "version_number": "2.1.0",
  "release_date": "2025-01-20T10:00:00Z",
  "release_type": "feature",
  "eol_date": null,
  "min_server_version": "2.0.0",
  "max_server_version": "3.0.0",
  "recommended_server_version": "2.5.0",
  "release_notes": {
    "version_info": {
      "version_number": "2.1.0",
      "release_date": "2025-01-20T10:00:00Z",
      "release_type": "feature"
    },
    "whats_new": ["New feature 1"],
    "bug_fixes": [],
    "breaking_changes": [],
    "compatibility": {
      "server_version_requirements": ">= 2.0.0",
      "client_version_requirements": ">= 1.5.0",
      "os_requirements": ["linux", "windows"]
    },
    "upgrade_instructions": "Follow the upgrade guide...",
    "known_issues": []
  }
}
```

**Response:** 201 Created with version object

#### PUT /products/{product_id}/versions/{version_number}
Update version (Phase 1)

**Request Body:** Same as POST, all fields optional

**Response:** 200 OK with updated version object

#### POST /products/{product_id}/versions/{version_number}/submit
Submit version for review (Phase 1)

**Response:**
```json
{
  "message": "Version submitted for review",
  "version": {
    "state": "pending_review",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

#### POST /products/{product_id}/versions/{version_number}/approve
Approve version for release (Phase 1)

**Request Body:**
```json
{
  "approved_by": "user@example.com"
}
```

**Response:**
```json
{
  "message": "Version approved for release",
  "version": {
    "state": "approved",
    "approved_by": "user@example.com",
    "approved_at": "2025-01-15T10:00:00Z"
  }
}
```

#### POST /products/{product_id}/versions/{version_number}/release
Release version (Phase 1)

**Response:**
```json
{
  "message": "Version released",
  "version": {
    "state": "released",
    "updated_at": "2025-01-15T10:00:00Z"
  }
}
```

### Packages API

#### POST /products/{product_id}/versions/{version_number}/packages
Upload a package (Phase 1)

**Request:** Multipart form data
- `file`: Package file
- `package_type`: full_installer | update | delta | rollback
- `os`: linux | windows | macos (optional)
- `architecture`: amd64 | arm64 | x86_64 (optional)

**Response:**
```json
{
  "message": "Package uploaded successfully",
  "package": {
    "id": "507f1f77bcf86cd799439013",
    "package_type": "full_installer",
    "file_name": "hyworks-2.1.0-installer.deb",
    "file_size": 52428800,
    "download_url": "https://cdn.accops.com/packages/hyworks-2.1.0-installer.deb",
    "checksum_sha256": "abc123...",
    "os": "linux",
    "architecture": "amd64",
    "uploaded_at": "2025-01-15T10:00:00Z"
  }
}
```

#### GET /products/{product_id}/versions/{version_number}/packages
List packages for a version

**Response:**
```json
{
  "packages": [
    {
      "id": "507f1f77bcf86cd799439013",
      "package_type": "full_installer",
      "file_name": "hyworks-2.1.0-installer.deb",
      "file_size": 52428800,
      "download_url": "https://cdn.accops.com/packages/hyworks-2.1.0-installer.deb",
      "checksum_sha256": "abc123...",
      "os": "linux",
      "architecture": "amd64"
    }
  ]
}
```

#### DELETE /products/{product_id}/versions/{version_number}/packages/{package_id}
Delete a package

**Response:** 204 No Content

### Compatibility API (Phase 2)

#### POST /products/{product_id}/versions/{version_number}/validate
Validate compatibility (Phase 2)

**Request Body:**
```json
{
  "min_server_version": "2.0.0",
  "max_server_version": "3.0.0",
  "recommended_server_version": "2.5.0",
  "incompatible_versions": ["1.0.0", "1.5.0"]
}
```

**Response:**
```json
{
  "validation_status": "passed",
  "compatibility_matrix": {
    "id": "507f1f77bcf86cd799439014",
    "product_id": "hyworks",
    "version_number": "2.1.0",
    "min_server_version": "2.0.0",
    "max_server_version": "3.0.0",
    "recommended_server_version": "2.5.0",
    "incompatible_versions": ["1.0.0", "1.5.0"],
    "validated_at": "2025-01-15T10:00:00Z",
    "validated_by": "system",
    "validation_status": "passed",
    "validation_errors": []
  }
}
```

#### GET /products/{product_id}/versions/{version_number}/compatibility
Get compatibility matrix

**Response:** Compatibility matrix object

#### GET /products/{product_id}/upgrade-paths
Get upgrade paths for a product

**Query Parameters:**
- `from_version` (optional): Filter by from version
- `to_version` (optional): Filter by to version

**Response:**
```json
{
  "upgrade_paths": [
    {
      "id": "507f1f77bcf86cd799439015",
      "product_id": "hyworks",
      "from_version": "2.0.0",
      "to_version": "2.1.0",
      "path_type": "direct",
      "is_blocked": false
    }
  ]
}
```

### Notifications API (Phase 2)

#### GET /notifications
Get user notifications

**Query Parameters:**
- `is_read` (optional): Filter by read status
- `type` (optional): Filter by notification type
- `priority` (optional): Filter by priority
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response:**
```json
{
  "notifications": [
    {
      "id": "507f1f77bcf86cd799439016",
      "type": "new_version",
      "product_id": "hyworks",
      "version_id": "2.1.0",
      "title": "New Version Available",
      "message": "HyWorks 2.1.0 is now available",
      "priority": "normal",
      "is_read": false,
      "created_at": "2025-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "total_pages": 1
  }
}
```

#### POST /notifications/{notification_id}/read
Mark notification as read

**Response:**
```json
{
  "message": "Notification marked as read",
  "notification": {
    "id": "507f1f77bcf86cd799439016",
    "is_read": true,
    "read_at": "2025-01-20T10:00:00Z"
  }
}
```

#### POST /notifications/read-all
Mark all notifications as read

**Response:**
```json
{
  "message": "All notifications marked as read",
  "count": 5
}
```

### Update Detection API (Phase 2)

#### GET /endpoints/{endpoint_id}/updates
Get available updates for an endpoint

**Response:**
```json
{
  "updates": [
    {
      "id": "507f1f77bcf86cd799439017",
      "endpoint_id": "endpoint-123",
      "product_id": "hyworks",
      "current_version": "2.0.0",
      "available_version": "2.1.0",
      "detected_at": "2025-01-20T10:00:00Z",
      "last_checked_at": "2025-01-20T10:00:00Z"
    }
  ]
}
```

#### POST /endpoints/{endpoint_id}/updates/check
Trigger update check for an endpoint

**Response:**
```json
{
  "message": "Update check initiated",
  "updates": []
}
```

### Update Rollout API (Phase 2)

#### POST /endpoints/{endpoint_id}/updates/{update_id}/rollout
Initiate update rollout

**Request Body:**
```json
{
  "to_version": "2.1.0"
}
```

**Response:**
```json
{
  "message": "Update rollout initiated",
  "rollout": {
    "id": "507f1f77bcf86cd799439018",
    "endpoint_id": "endpoint-123",
    "product_id": "hyworks",
    "from_version": "2.0.0",
    "to_version": "2.1.0",
    "status": "pending",
    "initiated_by": "user@example.com",
    "initiated_at": "2025-01-20T10:00:00Z",
    "progress": 0
  }
}
```

#### GET /endpoints/{endpoint_id}/rollouts
Get rollout history for an endpoint

**Response:**
```json
{
  "rollouts": [
    {
      "id": "507f1f77bcf86cd799439018",
      "endpoint_id": "endpoint-123",
      "product_id": "hyworks",
      "from_version": "2.0.0",
      "to_version": "2.1.0",
      "status": "completed",
      "initiated_by": "user@example.com",
      "initiated_at": "2025-01-20T10:00:00Z",
      "completed_at": "2025-01-20T10:05:00Z",
      "progress": 100
    }
  ]
}
```

#### GET /rollouts/{rollout_id}
Get rollout details

**Response:** Rollout object

#### POST /rollouts/{rollout_id}/cancel
Cancel a rollout

**Response:**
```json
{
  "message": "Rollout cancelled",
  "rollout": {
    "status": "cancelled"
  }
}
```

### Audit Logs API

#### GET /audit-logs
Get audit logs

**Query Parameters:**
- `resource_type` (optional): Filter by resource type
- `resource_id` (optional): Filter by resource ID
- `user_id` (optional): Filter by user ID
- `action` (optional): Filter by action
- `start_date` (optional): Start date filter
- `end_date` (optional): End date filter
- `page` (optional): Page number
- `limit` (optional): Items per page

**Response:**
```json
{
  "audit_logs": [
    {
      "id": "507f1f77bcf86cd799439019",
      "action": "approve",
      "resource_type": "version",
      "resource_id": "2.1.0",
      "user_id": "user-123",
      "user_email": "user@example.com",
      "details": {
        "product_id": "hyworks",
        "version_number": "2.1.0"
      },
      "ip_address": "192.168.1.1",
      "user_agent": "Mozilla/5.0...",
      "timestamp": "2025-01-20T10:00:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 1000,
    "total_pages": 50
  }
}
```

## Error Responses

All error responses follow this format:

```json
{
  "error": {
    "code": "ERROR_CODE",
    "message": "Human-readable error message",
    "details": {
      "field": "Additional error details"
    }
  }
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource conflict (e.g., duplicate version)
- `422 Unprocessable Entity`: Validation error
- `500 Internal Server Error`: Server error

### Example Error Response

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Validation failed",
    "details": {
      "version_number": "Version number is required",
      "release_type": "Invalid release type"
    }
  }
}
```

## Authentication

All API requests require authentication via Bearer token:

```
Authorization: Bearer <jwt_token>
```

## Rate Limiting

- Standard endpoints: 100 requests per minute
- File upload endpoints: 10 requests per minute
- Update rollout endpoints: 20 requests per minute

## Versioning

API version is specified in the URL path: `/api/v1`

Future versions will use `/api/v2`, etc.

