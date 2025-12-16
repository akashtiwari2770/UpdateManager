package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Product represents a product in the system
type Product struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID   string             `bson:"product_id" json:"product_id" validate:"required,min=1,max=100"`
	Name        string             `bson:"name" json:"name" validate:"required,min=1,max=200"`
	Type        ProductType        `bson:"type" json:"type" validate:"required"`
	Description string             `bson:"description" json:"description" validate:"max=1000"`
	Vendor      string             `bson:"vendor" json:"vendor" validate:"max=100"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	IsActive    bool               `bson:"is_active" json:"is_active"`
}

type ProductType string

const (
	ProductTypeServer ProductType = "server"
	ProductTypeClient ProductType = "client"
)

// Version represents a product version
type Version struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID     string             `bson:"product_id" json:"product_id" validate:"required"`
	VersionNumber string             `bson:"version_number" json:"version_number" validate:"required"`
	ReleaseDate   time.Time          `bson:"release_date" json:"release_date"`
	ReleaseType   ReleaseType        `bson:"release_type" json:"release_type" validate:"required"`
	State         VersionState       `bson:"state" json:"state" validate:"required"`
	EOLDate       *time.Time         `bson:"eol_date,omitempty" json:"eol_date,omitempty"`

	// Compatibility (for clients)
	MinServerVersion         string `bson:"min_server_version,omitempty" json:"min_server_version,omitempty"`
	MaxServerVersion         string `bson:"max_server_version,omitempty" json:"max_server_version,omitempty"`
	RecommendedServerVersion string `bson:"recommended_server_version,omitempty" json:"recommended_server_version,omitempty"`

	// Release Notes
	ReleaseNotes *ReleaseNotes `bson:"release_notes,omitempty" json:"release_notes,omitempty"`

	// Packages
	Packages []PackageInfo `bson:"packages" json:"packages"`

	// Approval
	ApprovedBy string     `bson:"approved_by,omitempty" json:"approved_by,omitempty"`
	ApprovedAt *time.Time `bson:"approved_at,omitempty" json:"approved_at,omitempty"`
	CreatedBy  string     `bson:"created_by" json:"created_by"`
	CreatedAt  time.Time  `bson:"created_at" json:"created_at"`
	UpdatedAt  time.Time  `bson:"updated_at" json:"updated_at"`
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

// ReleaseNotes represents release notes for a version
type ReleaseNotes struct {
	VersionInfo         VersionInfoSection   `bson:"version_info" json:"version_info"`
	WhatsNew            []string             `bson:"whats_new" json:"whats_new"`
	BugFixes            []BugFix             `bson:"bug_fixes" json:"bug_fixes"`
	BreakingChanges     []BreakingChange     `bson:"breaking_changes" json:"breaking_changes"`
	Compatibility       CompatibilitySection `bson:"compatibility" json:"compatibility"`
	UpgradeInstructions string               `bson:"upgrade_instructions" json:"upgrade_instructions"`
	KnownIssues         []KnownIssue         `bson:"known_issues" json:"known_issues"`
}

type VersionInfoSection struct {
	VersionNumber string      `bson:"version_number" json:"version_number"`
	ReleaseDate   time.Time   `bson:"release_date" json:"release_date"`
	ReleaseType   ReleaseType `bson:"release_type" json:"release_type"`
}

type BugFix struct {
	ID          string `bson:"id" json:"id"`
	Description string `bson:"description" json:"description"`
	IssueNumber string `bson:"issue_number,omitempty" json:"issue_number,omitempty"`
}

type BreakingChange struct {
	Description          string `bson:"description" json:"description"`
	MigrationSteps       string `bson:"migration_steps" json:"migration_steps"`
	ConfigurationChanges string `bson:"configuration_changes" json:"configuration_changes"`
}

type CompatibilitySection struct {
	ServerVersionRequirements string   `bson:"server_version_requirements" json:"server_version_requirements"`
	ClientVersionRequirements string   `bson:"client_version_requirements" json:"client_version_requirements"`
	OSRequirements            []string `bson:"os_requirements" json:"os_requirements"`
}

type KnownIssue struct {
	ID          string `bson:"id" json:"id"`
	Description string `bson:"description" json:"description"`
	Workaround  string `bson:"workaround,omitempty" json:"workaround,omitempty"`
	PlannedFix  string `bson:"planned_fix,omitempty" json:"planned_fix,omitempty"`
}

// PackageInfo represents a package file
type PackageInfo struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	PackageType      PackageType        `bson:"package_type" json:"package_type" validate:"required"`
	FileName         string             `bson:"file_name" json:"file_name" validate:"required"`
	FileSize         int64              `bson:"file_size" json:"file_size"`
	DownloadURL      string             `bson:"download_url" json:"download_url"`
	ChecksumSHA256   string             `bson:"checksum_sha256" json:"checksum_sha256" validate:"required"`
	DigitalSignature string             `bson:"digital_signature,omitempty" json:"digital_signature,omitempty"`
	OS               string             `bson:"os,omitempty" json:"os,omitempty"`
	Architecture     string             `bson:"architecture,omitempty" json:"architecture,omitempty"`
	UploadedAt       time.Time          `bson:"uploaded_at" json:"uploaded_at"`
	UploadedBy       string             `bson:"uploaded_by" json:"uploaded_by"`
}

type PackageType string

const (
	PackageTypeFullInstaller PackageType = "full_installer"
	PackageTypeUpdate        PackageType = "update"
	PackageTypeDelta         PackageType = "delta"
	PackageTypeRollback      PackageType = "rollback"
)

// CompatibilityMatrix represents compatibility validation results
type CompatibilityMatrix struct {
	ID                       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID                string             `bson:"product_id" json:"product_id" validate:"required"`
	VersionNumber            string             `bson:"version_number" json:"version_number" validate:"required"`
	MinServerVersion         string             `bson:"min_server_version,omitempty" json:"min_server_version,omitempty"`
	MaxServerVersion         string             `bson:"max_server_version,omitempty" json:"max_server_version,omitempty"`
	RecommendedServerVersion string             `bson:"recommended_server_version,omitempty" json:"recommended_server_version,omitempty"`
	IncompatibleVersions     []string           `bson:"incompatible_versions" json:"incompatible_versions"`
	ValidatedAt              time.Time          `bson:"validated_at" json:"validated_at"`
	ValidatedBy              string             `bson:"validated_by" json:"validated_by"`
	ValidationStatus         ValidationStatus   `bson:"validation_status" json:"validation_status"`
	ValidationErrors         []string           `bson:"validation_errors" json:"validation_errors"`
}

type ValidationStatus string

const (
	ValidationStatusPending ValidationStatus = "pending"
	ValidationStatusPassed  ValidationStatus = "passed"
	ValidationStatusFailed  ValidationStatus = "failed"
	ValidationStatusSkipped ValidationStatus = "skipped"
)

// UpgradePath represents an upgrade path between versions
type UpgradePath struct {
	ID                   primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ProductID            string             `bson:"product_id" json:"product_id" validate:"required"`
	FromVersion          string             `bson:"from_version" json:"from_version" validate:"required"`
	ToVersion            string             `bson:"to_version" json:"to_version" validate:"required"`
	PathType             UpgradePathType    `bson:"path_type" json:"path_type" validate:"required"`
	IntermediateVersions []string           `bson:"intermediate_versions,omitempty" json:"intermediate_versions,omitempty"`
	IsBlocked            bool               `bson:"is_blocked" json:"is_blocked"`
	BlockReason          string             `bson:"block_reason,omitempty" json:"block_reason,omitempty"`
	CreatedAt            time.Time          `bson:"created_at" json:"created_at"`
}

type UpgradePathType string

const (
	UpgradePathTypeDirect    UpgradePathType = "direct"
	UpgradePathTypeMultiStep UpgradePathType = "multi_step"
	UpgradePathTypeBlocked   UpgradePathType = "blocked"
)

// Notification represents a notification
type Notification struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Type        NotificationType     `bson:"type" json:"type" validate:"required"`
	RecipientID string               `bson:"recipient_id" json:"recipient_id" validate:"required"`
	ProductID   string               `bson:"product_id,omitempty" json:"product_id,omitempty"`
	VersionID   string               `bson:"version_id,omitempty" json:"version_id,omitempty"`
	CustomerID  string               `bson:"customer_id,omitempty" json:"customer_id,omitempty"`
	TenantID    string               `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	DeploymentID string             `bson:"deployment_id,omitempty" json:"deployment_id,omitempty"`
	Title       string               `bson:"title" json:"title" validate:"required"`
	Message     string               `bson:"message" json:"message" validate:"required"`
	Priority    NotificationPriority `bson:"priority" json:"priority"`
	IsRead      bool                 `bson:"is_read" json:"is_read"`
	ReadAt      *time.Time           `bson:"read_at,omitempty" json:"read_at,omitempty"`
	CreatedAt   time.Time            `bson:"created_at" json:"created_at"`
}

type NotificationType string

const (
	NotificationTypeNewVersion      NotificationType = "new_version"
	NotificationTypeSecurityRelease NotificationType = "security_release"
	NotificationTypeEOLWarning      NotificationType = "eol_warning"
	NotificationTypeUpdateAvailable NotificationType = "update_available"
)

type NotificationPriority string

const (
	NotificationPriorityLow      NotificationPriority = "low"
	NotificationPriorityNormal   NotificationPriority = "normal"
	NotificationPriorityHigh     NotificationPriority = "high"
	NotificationPriorityCritical NotificationPriority = "critical"
)

// UpdateDetection represents update detection for an endpoint
type UpdateDetection struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EndpointID       string             `bson:"endpoint_id" json:"endpoint_id" validate:"required"`
	ProductID        string             `bson:"product_id" json:"product_id" validate:"required"`
	CurrentVersion   string             `bson:"current_version" json:"current_version" validate:"required"`
	AvailableVersion string             `bson:"available_version" json:"available_version" validate:"required"`
	DetectedAt       time.Time          `bson:"detected_at" json:"detected_at"`
	LastCheckedAt    time.Time          `bson:"last_checked_at" json:"last_checked_at"`
}

// UpdateRollout represents an update rollout
type UpdateRollout struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	EndpointID   string             `bson:"endpoint_id" json:"endpoint_id" validate:"required"`
	ProductID    string             `bson:"product_id" json:"product_id" validate:"required"`
	FromVersion  string             `bson:"from_version" json:"from_version" validate:"required"`
	ToVersion    string             `bson:"to_version" json:"to_version" validate:"required"`
	Status       RolloutStatus      `bson:"status" json:"status" validate:"required"`
	InitiatedBy  string             `bson:"initiated_by" json:"initiated_by" validate:"required"`
	InitiatedAt  time.Time          `bson:"initiated_at" json:"initiated_at"`
	StartedAt    *time.Time         `bson:"started_at,omitempty" json:"started_at,omitempty"`
	CompletedAt  *time.Time         `bson:"completed_at,omitempty" json:"completed_at,omitempty"`
	FailedAt     *time.Time         `bson:"failed_at,omitempty" json:"failed_at,omitempty"`
	ErrorMessage string             `bson:"error_message,omitempty" json:"error_message,omitempty"`
	Progress     int                `bson:"progress" json:"progress"` // 0-100
}

type RolloutStatus string

const (
	RolloutStatusPending    RolloutStatus = "pending"
	RolloutStatusInProgress RolloutStatus = "in_progress"
	RolloutStatusCompleted  RolloutStatus = "completed"
	RolloutStatusFailed     RolloutStatus = "failed"
	RolloutStatusCancelled  RolloutStatus = "cancelled"
)

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Action       AuditAction            `bson:"action" json:"action" validate:"required"`
	ResourceType string                 `bson:"resource_type" json:"resource_type" validate:"required"`
	ResourceID   string                 `bson:"resource_id" json:"resource_id" validate:"required"`
	UserID       string                 `bson:"user_id" json:"user_id" validate:"required"`
	UserEmail    string                 `bson:"user_email" json:"user_email"`
	Details      map[string]interface{} `bson:"details" json:"details"`
	IPAddress    string                 `bson:"ip_address" json:"ip_address"`
	UserAgent    string                 `bson:"user_agent" json:"user_agent"`
	Timestamp    time.Time              `bson:"timestamp" json:"timestamp"`
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

// Request/Response DTOs

// CreateProductRequest represents a request to create a product
type CreateProductRequest struct {
	ProductID   string      `json:"product_id" validate:"required,min=1,max=100"`
	Name        string      `json:"name" validate:"required,min=1,max=200"`
	Type        ProductType `json:"type" validate:"required"`
	Description string      `json:"description" validate:"max=1000"`
	Vendor      string      `json:"vendor" validate:"max=100"`
}

// CreateVersionRequest represents a request to create a version
type CreateVersionRequest struct {
	VersionNumber            string        `json:"version_number" validate:"required"`
	ReleaseDate              time.Time     `json:"release_date"`
	ReleaseType              ReleaseType   `json:"release_type" validate:"required"`
	EOLDate                  *time.Time    `json:"eol_date,omitempty"`
	MinServerVersion         string        `json:"min_server_version,omitempty"`
	MaxServerVersion         string        `json:"max_server_version,omitempty"`
	RecommendedServerVersion string        `json:"recommended_server_version,omitempty"`
	ReleaseNotes             *ReleaseNotes `json:"release_notes,omitempty"`
}

// UpdateVersionRequest represents a request to update a version
type UpdateVersionRequest struct {
	ReleaseDate              *time.Time    `json:"release_date,omitempty"`
	ReleaseType              *ReleaseType  `json:"release_type,omitempty"`
	EOLDate                  *time.Time    `json:"eol_date,omitempty"`
	MinServerVersion         *string       `json:"min_server_version,omitempty"`
	MaxServerVersion         *string       `json:"max_server_version,omitempty"`
	RecommendedServerVersion *string       `json:"recommended_server_version,omitempty"`
	ReleaseNotes             *ReleaseNotes `json:"release_notes,omitempty"`
}

// ApproveVersionRequest represents a request to approve a version
type ApproveVersionRequest struct {
	ApprovedBy string `json:"approved_by" validate:"required"`
}

// ValidateCompatibilityRequest represents a request to validate compatibility
type ValidateCompatibilityRequest struct {
	MinServerVersion         string   `json:"min_server_version,omitempty"`
	MaxServerVersion         string   `json:"max_server_version,omitempty"`
	RecommendedServerVersion string   `json:"recommended_server_version,omitempty"`
	IncompatibleVersions     []string `json:"incompatible_versions"`
}

// InitiateRolloutRequest represents a request to initiate an update rollout
type InitiateRolloutRequest struct {
	ToVersion string `json:"to_version" validate:"required"`
}

// PaginationResponse represents a paginated response
type PaginationResponse struct {
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Details map[string]interface{} `json:"details,omitempty"`
}

// Customer Management Models

// Customer represents a customer in the system
type Customer struct {
	ID                    primitive.ObjectID      `bson:"_id,omitempty" json:"id"`
	CustomerID            string                  `bson:"customer_id" json:"customer_id" validate:"required,min=1,max=100"`
	Name                  string                  `bson:"name" json:"name" validate:"required,min=1,max=200"`
	OrganizationName      string                  `bson:"organization_name,omitempty" json:"organization_name,omitempty" validate:"max=200"`
	Email                 string                  `bson:"email" json:"email" validate:"required,email"`
	Phone                 string                  `bson:"phone,omitempty" json:"phone,omitempty" validate:"max=50"`
	Address               string                  `bson:"address,omitempty" json:"address,omitempty" validate:"max=500"`
	AccountStatus         CustomerStatus          `bson:"account_status" json:"account_status" validate:"required"`
	NotificationPreferences NotificationPreferences `bson:"notification_preferences" json:"notification_preferences"`
	CreatedAt             time.Time               `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time               `bson:"updated_at" json:"updated_at"`
}

// CustomerStatus represents the status of a customer account
type CustomerStatus string

const (
	CustomerStatusActive    CustomerStatus = "active"
	CustomerStatusInactive  CustomerStatus = "inactive"
	CustomerStatusSuspended CustomerStatus = "suspended"
)

// NotificationPreferences represents notification preferences for a customer
type NotificationPreferences struct {
	EmailEnabled           bool `bson:"email_enabled" json:"email_enabled"`
	InAppEnabled           bool `bson:"in_app_enabled" json:"in_app_enabled"`
	UATNotifications       bool `bson:"uat_notifications" json:"uat_notifications"`
	ProductionNotifications bool `bson:"production_notifications" json:"production_notifications"`
}

// CustomerTenant represents a tenant (independent deployment) for a customer
type CustomerTenant struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	TenantID    string             `bson:"tenant_id" json:"tenant_id" validate:"required,min=1,max=100"`
	CustomerID  primitive.ObjectID `bson:"customer_id" json:"customer_id" validate:"required"`
	Name        string             `bson:"name" json:"name" validate:"required,min=1,max=200"`
	Description string             `bson:"description,omitempty" json:"description,omitempty" validate:"max=1000"`
	Status      TenantStatus       `bson:"status" json:"status" validate:"required"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// TenantStatus represents the status of a tenant
type TenantStatus string

const (
	TenantStatusActive   TenantStatus = "active"
	TenantStatusInactive TenantStatus = "inactive"
)

// Deployment represents a deployment (product + type) for a tenant
type Deployment struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	DeploymentID     string             `bson:"deployment_id" json:"deployment_id" validate:"required,min=1,max=100"`
	TenantID         primitive.ObjectID `bson:"tenant_id" json:"tenant_id" validate:"required"`
	ProductID        string             `bson:"product_id" json:"product_id" validate:"required"`
	DeploymentType   DeploymentType     `bson:"deployment_type" json:"deployment_type" validate:"required"`
	InstalledVersion string             `bson:"installed_version" json:"installed_version" validate:"required"`
	NumberOfUsers    *int               `bson:"number_of_users,omitempty" json:"number_of_users,omitempty"`
	LicenseInfo      string             `bson:"license_info,omitempty" json:"license_info,omitempty" validate:"max=1000"`
	ServerHostname   string             `bson:"server_hostname,omitempty" json:"server_hostname,omitempty" validate:"max=200"`
	EnvironmentDetails string           `bson:"environment_details,omitempty" json:"environment_details,omitempty" validate:"max=500"`
	DeploymentDate   time.Time          `bson:"deployment_date" json:"deployment_date"`
	LastUpdatedDate  time.Time          `bson:"last_updated_date" json:"last_updated_date"`
	Status           DeploymentStatus   `bson:"status" json:"status" validate:"required"`
}

// DeploymentType represents the type of deployment
type DeploymentType string

const (
	DeploymentTypeUAT        DeploymentType = "uat"
	DeploymentTypeTesting    DeploymentType = "testing"
	DeploymentTypeProduction DeploymentType = "production"
)

// DeploymentStatus represents the status of a deployment
type DeploymentStatus string

const (
	DeploymentStatusActive   DeploymentStatus = "active"
	DeploymentStatusInactive DeploymentStatus = "inactive"
)

// Request DTOs for Customer Management

// CreateCustomerRequest represents a request to create a customer
type CreateCustomerRequest struct {
	CustomerID            string                  `json:"customer_id" validate:"required,min=1,max=100"`
	Name                  string                  `json:"name" validate:"required,min=1,max=200"`
	OrganizationName      string                  `json:"organization_name,omitempty" validate:"max=200"`
	Email                 string                  `json:"email" validate:"required,email"`
	Phone                 string                  `json:"phone,omitempty" validate:"max=50"`
	Address               string                  `json:"address,omitempty" validate:"max=500"`
	AccountStatus         CustomerStatus          `json:"account_status" validate:"required"`
	NotificationPreferences NotificationPreferences `json:"notification_preferences"`
}

// UpdateCustomerRequest represents a request to update a customer
type UpdateCustomerRequest struct {
	Name                  *string                 `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	OrganizationName      *string                 `json:"organization_name,omitempty" validate:"omitempty,max=200"`
	Email                 *string                 `json:"email,omitempty" validate:"omitempty,email"`
	Phone                 *string                 `json:"phone,omitempty" validate:"omitempty,max=50"`
	Address               *string                 `json:"address,omitempty" validate:"omitempty,max=500"`
	AccountStatus         *CustomerStatus         `json:"account_status,omitempty"`
	NotificationPreferences *NotificationPreferences `json:"notification_preferences,omitempty"`
}

// CreateTenantRequest represents a request to create a tenant
type CreateTenantRequest struct {
	TenantID    string `json:"tenant_id" validate:"required,min=1,max=100"`
	Name        string `json:"name" validate:"required,min=1,max=200"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	Status      TenantStatus `json:"status" validate:"required"`
}

// UpdateTenantRequest represents a request to update a tenant
type UpdateTenantRequest struct {
	Name        *string      `json:"name,omitempty" validate:"omitempty,min=1,max=200"`
	Description *string      `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      *TenantStatus `json:"status,omitempty"`
}

// CreateDeploymentRequest represents a request to create a deployment
type CreateDeploymentRequest struct {
	DeploymentID      string         `json:"deployment_id" validate:"required,min=1,max=100"`
	ProductID         string         `json:"product_id" validate:"required"`
	DeploymentType    DeploymentType `json:"deployment_type" validate:"required"`
	InstalledVersion  string         `json:"installed_version" validate:"required"`
	NumberOfUsers     *int           `json:"number_of_users,omitempty"`
	LicenseInfo       string         `json:"license_info,omitempty" validate:"max=1000"`
	ServerHostname    string         `json:"server_hostname,omitempty" validate:"max=200"`
	EnvironmentDetails string        `json:"environment_details,omitempty" validate:"max=500"`
	Status            DeploymentStatus `json:"status" validate:"required"`
}

// UpdateDeploymentRequest represents a request to update a deployment
type UpdateDeploymentRequest struct {
	DeploymentType    *DeploymentType  `json:"deployment_type,omitempty"`
	InstalledVersion  *string          `json:"installed_version,omitempty"`
	NumberOfUsers     *int             `json:"number_of_users,omitempty"`
	LicenseInfo       *string          `json:"license_info,omitempty" validate:"omitempty,max=1000"`
	ServerHostname    *string          `json:"server_hostname,omitempty" validate:"omitempty,max=200"`
	EnvironmentDetails *string        `json:"environment_details,omitempty" validate:"omitempty,max=500"`
	Status            *DeploymentStatus `json:"status,omitempty"`
}

// License Management Models

// Subscription represents a subscription for a customer
type Subscription struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	SubscriptionID string           `bson:"subscription_id" json:"subscription_id" validate:"required,min=1,max=100"`
	CustomerID  primitive.ObjectID  `bson:"customer_id" json:"customer_id" validate:"required"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty" validate:"max=200"`
	Description string             `bson:"description,omitempty" json:"description,omitempty" validate:"max=1000"`
	StartDate   time.Time          `bson:"start_date" json:"start_date" validate:"required"`
	EndDate     *time.Time         `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Status      SubscriptionStatus `bson:"status" json:"status" validate:"required"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
	CreatedBy   string             `bson:"created_by" json:"created_by" validate:"required"`
	Notes       string             `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=2000"`
}

// SubscriptionStatus represents the status of a subscription
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusInactive  SubscriptionStatus = "inactive"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusSuspended SubscriptionStatus = "suspended"
)

// License represents a license assigned to a subscription
type License struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	LicenseID      string             `bson:"license_id" json:"license_id" validate:"required,min=1,max=100"`
	SubscriptionID primitive.ObjectID `bson:"subscription_id" json:"subscription_id" validate:"required"`
	ProductID      string             `bson:"product_id" json:"product_id" validate:"required"`
	LicenseType    LicenseType        `bson:"license_type" json:"license_type" validate:"required"`
	NumberOfSeats  int                `bson:"number_of_seats" json:"number_of_seats" validate:"required,min=1"`
	StartDate      time.Time          `bson:"start_date" json:"start_date" validate:"required"`
	EndDate        *time.Time         `bson:"end_date,omitempty" json:"end_date,omitempty"`
	Status         LicenseStatus      `bson:"status" json:"status" validate:"required"`
	AssignedBy     string             `bson:"assigned_by" json:"assigned_by" validate:"required"`
	AssignmentDate time.Time          `bson:"assignment_date" json:"assignment_date" validate:"required"`
	Notes          string             `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=2000"`
	CreatedAt      time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt      time.Time          `bson:"updated_at" json:"updated_at"`
}

// LicenseType represents the type of license
type LicenseType string

const (
	LicenseTypePerpetual LicenseType = "perpetual"
	LicenseTypeTimeBased LicenseType = "time_based"
)

// LicenseStatus represents the status of a license
type LicenseStatus string

const (
	LicenseStatusActive   LicenseStatus = "active"
	LicenseStatusInactive LicenseStatus = "inactive"
	LicenseStatusExpired  LicenseStatus = "expired"
	LicenseStatusRevoked  LicenseStatus = "revoked"
)

// LicenseAllocation represents an allocation of a license to a tenant or deployment
type LicenseAllocation struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AllocationID          string             `bson:"allocation_id" json:"allocation_id" validate:"required,min=1,max=100"`
	LicenseID             primitive.ObjectID `bson:"license_id" json:"license_id" validate:"required"`
	TenantID              *primitive.ObjectID `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	DeploymentID          *primitive.ObjectID `bson:"deployment_id,omitempty" json:"deployment_id,omitempty"`
	NumberOfSeatsAllocated int                `bson:"number_of_seats_allocated" json:"number_of_seats_allocated" validate:"required,min=1"`
	AllocationDate        time.Time          `bson:"allocation_date" json:"allocation_date" validate:"required"`
	AllocatedBy           string             `bson:"allocated_by" json:"allocated_by" validate:"required"`
	Status                AllocationStatus   `bson:"status" json:"status" validate:"required"`
	ReleasedDate          *time.Time         `bson:"released_date,omitempty" json:"released_date,omitempty"`
	ReleasedBy            *string            `bson:"released_by,omitempty" json:"released_by,omitempty"`
	Notes                 string             `bson:"notes,omitempty" json:"notes,omitempty" validate:"max=2000"`
	CreatedAt             time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt             time.Time          `bson:"updated_at" json:"updated_at"`
}

// AllocationStatus represents the status of a license allocation
type AllocationStatus string

const (
	AllocationStatusActive  AllocationStatus = "active"
	AllocationStatusReleased AllocationStatus = "released"
)

// Request DTOs for License Management

// CreateSubscriptionRequest represents a request to create a subscription
type CreateSubscriptionRequest struct {
	SubscriptionID string     `json:"subscription_id" validate:"required,min=1,max=100"`
	Name           string     `json:"name,omitempty" validate:"max=200"`
	Description    string     `json:"description,omitempty" validate:"max=1000"`
	StartDate      time.Time  `json:"start_date" validate:"required"`
	EndDate        *time.Time `json:"end_date,omitempty"`
	Status         SubscriptionStatus `json:"status" validate:"required"`
	Notes          string     `json:"notes,omitempty" validate:"max=2000"`
}

// UpdateSubscriptionRequest represents a request to update a subscription
type UpdateSubscriptionRequest struct {
	Name        *string             `json:"name,omitempty" validate:"omitempty,max=200"`
	Description *string             `json:"description,omitempty" validate:"omitempty,max=1000"`
	StartDate   *time.Time          `json:"start_date,omitempty"`
	EndDate     *time.Time          `json:"end_date,omitempty"`
	Status      *SubscriptionStatus `json:"status,omitempty"`
	Notes       *string             `json:"notes,omitempty" validate:"omitempty,max=2000"`
}

// CreateLicenseRequest represents a request to assign a license to a subscription
type CreateLicenseRequest struct {
	LicenseID     string      `json:"license_id" validate:"required,min=1,max=100"`
	ProductID     string      `json:"product_id" validate:"required"`
	LicenseType   LicenseType `json:"license_type" validate:"required"`
	NumberOfSeats int         `json:"number_of_seats" validate:"required,min=1"`
	StartDate     time.Time   `json:"start_date" validate:"required"`
	EndDate       *time.Time  `json:"end_date,omitempty"`
	Status        LicenseStatus `json:"status" validate:"required"`
	Notes         string      `json:"notes,omitempty" validate:"max=2000"`
}

// UpdateLicenseRequest represents a request to update a license
type UpdateLicenseRequest struct {
	LicenseType   *LicenseType   `json:"license_type,omitempty"`
	NumberOfSeats *int           `json:"number_of_seats,omitempty" validate:"omitempty,min=1"`
	StartDate     *time.Time     `json:"start_date,omitempty"`
	EndDate       *time.Time     `json:"end_date,omitempty"`
	Status        *LicenseStatus `json:"status,omitempty"`
	Notes         *string        `json:"notes,omitempty" validate:"omitempty,max=2000"`
}

// AllocateLicenseRequest represents a request to allocate a license to a tenant or deployment
type AllocateLicenseRequest struct {
	TenantID              *string `json:"tenant_id,omitempty"`
	DeploymentID          *string `json:"deployment_id,omitempty"`
	NumberOfSeatsAllocated int    `json:"number_of_seats_allocated" validate:"required,min=1"`
	Notes                 string  `json:"notes,omitempty" validate:"max=2000"`
}

// ReleaseLicenseRequest represents a request to release a license allocation
type ReleaseLicenseRequest struct {
	Notes string `json:"notes,omitempty" validate:"max=2000"`
}
