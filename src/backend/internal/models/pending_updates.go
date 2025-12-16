package models

import "time"

// AvailableUpdate represents an available update for a deployment
type AvailableUpdate struct {
	VersionNumber      string    `json:"version_number"`
	ReleaseDate        time.Time `json:"release_date"`
	ReleaseType        string    `json:"release_type"`
	IsSecurityUpdate   bool      `json:"is_security_update"`
	CompatibilityStatus string   `json:"compatibility_status"`
	UpgradePath        []string  `json:"upgrade_path"` // Intermediate versions if needed
}

// PendingUpdatesResponse represents pending updates for a deployment
type PendingUpdatesResponse struct {
	DeploymentID     string            `json:"deployment_id"`
	ProductID         string            `json:"product_id"`
	CurrentVersion    string            `json:"current_version"`
	LatestVersion     string            `json:"latest_version"`
	UpdateCount       int               `json:"update_count"`
	Priority          string            `json:"priority"` // critical, high, normal
	VersionGapType    string            `json:"version_gap_type"` // patch, minor, major
	AvailableUpdates  []AvailableUpdate `json:"available_updates"`
	// Additional context
	TenantID          string            `json:"tenant_id,omitempty"`
	TenantName        string            `json:"tenant_name,omitempty"`
	CustomerID        string            `json:"customer_id,omitempty"`
	CustomerName      string            `json:"customer_name,omitempty"`
	DeploymentType    DeploymentType    `json:"deployment_type,omitempty"`
}

// TenantPendingUpdatesSummary represents aggregated pending updates for a tenant
type TenantPendingUpdatesSummary struct {
	TenantID                string                 `json:"tenant_id"`
	TenantName              string                 `json:"tenant_name"`
	TotalDeployments        int                    `json:"total_deployments"`
	DeploymentsWithUpdates  int                    `json:"deployments_with_updates"`
	TotalPendingUpdateCount int                    `json:"total_pending_update_count"`
	ByPriority              map[string]int          `json:"by_priority"`
	ByProduct               map[string]int          `json:"by_product"`
	Deployments             []PendingUpdatesResponse `json:"deployments"`
}

// CustomerPendingUpdatesSummary represents aggregated pending updates for a customer
type CustomerPendingUpdatesSummary struct {
	CustomerID              string                 `json:"customer_id"`
	CustomerName            string                 `json:"customer_name"`
	TotalDeployments        int                    `json:"total_deployments"`
	DeploymentsWithUpdates  int                    `json:"deployments_with_updates"`
	TotalPendingUpdateCount int                    `json:"total_pending_update_count"`
	ByPriority              map[string]int          `json:"by_priority"`
	ByProduct               map[string]int          `json:"by_product"`
	ByTenant                map[string]int          `json:"by_tenant"`
	Deployments             []PendingUpdatesResponse `json:"deployments"`
}

// PendingUpdatesFilter represents filters for pending updates queries
type PendingUpdatesFilter struct {
	ProductID      string
	DeploymentType DeploymentType
	Priority       string // critical, high, normal
	TenantID       string
	CustomerID     string
}

