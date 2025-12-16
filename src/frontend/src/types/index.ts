// TypeScript interfaces for React frontend
// Product Release Management System

// Enums
export enum ProductType {
  SERVER = 'server',
  CLIENT = 'client',
}

export enum ReleaseType {
  SECURITY = 'security',
  FEATURE = 'feature',
  MAINTENANCE = 'maintenance',
  MAJOR = 'major',
}

export enum VersionState {
  DRAFT = 'draft',
  PENDING_REVIEW = 'pending_review',
  APPROVED = 'approved',
  RELEASED = 'released',
  DEPRECATED = 'deprecated',
  EOL = 'eol',
}

export enum PackageType {
  FULL_INSTALLER = 'full_installer',
  UPDATE = 'update',
  DELTA = 'delta',
  ROLLBACK = 'rollback',
}

export enum ValidationStatus {
  PENDING = 'pending',
  PASSED = 'passed',
  FAILED = 'failed',
  SKIPPED = 'skipped',
}

export enum UpgradePathType {
  DIRECT = 'direct',
  MULTI_STEP = 'multi_step',
  BLOCKED = 'blocked',
}

export enum NotificationType {
  NEW_VERSION = 'new_version',
  SECURITY_RELEASE = 'security_release',
  EOL_WARNING = 'eol_warning',
  UPDATE_AVAILABLE = 'update_available',
}

export enum NotificationPriority {
  LOW = 'low',
  NORMAL = 'normal',
  HIGH = 'high',
  CRITICAL = 'critical',
}

export enum RolloutStatus {
  PENDING = 'pending',
  IN_PROGRESS = 'in_progress',
  COMPLETED = 'completed',
  FAILED = 'failed',
  CANCELLED = 'cancelled',
}

export enum AuditAction {
  CREATE = 'create',
  UPDATE = 'update',
  DELETE = 'delete',
  APPROVE = 'approve',
  REJECT = 'reject',
  RELEASE = 'release',
  UPLOAD = 'upload',
  DOWNLOAD = 'download',
}

export enum CustomerStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SUSPENDED = 'suspended',
}

export enum TenantStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
}

export enum DeploymentType {
  UAT = 'uat',
  TESTING = 'testing',
  PRODUCTION = 'production',
}

export enum DeploymentStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
}

// License Management Enums
export enum SubscriptionStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  EXPIRED = 'expired',
  SUSPENDED = 'suspended',
}

export enum LicenseType {
  PERPETUAL = 'perpetual',
  TIME_BASED = 'time_based',
}

export enum LicenseStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  EXPIRED = 'expired',
  REVOKED = 'revoked',
}

export enum AllocationStatus {
  ACTIVE = 'active',
  RELEASED = 'released',
}

// Main Models
export interface Product {
  id: string;
  product_id: string;
  name: string;
  type: ProductType;
  description?: string;
  vendor?: string;
  created_at: string;
  updated_at: string;
  is_active: boolean;
}

export interface Version {
  id: string;
  product_id: string;
  version_number: string;
  release_date: string;
  release_type: ReleaseType;
  state: VersionState;
  eol_date?: string;
  min_server_version?: string;
  max_server_version?: string;
  recommended_server_version?: string;
  release_notes?: ReleaseNotes;
  packages: PackageInfo[];
  approved_by?: string;
  approved_at?: string;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface ReleaseNotes {
  version_info: VersionInfoSection;
  whats_new: string[];
  bug_fixes: BugFix[];
  breaking_changes: BreakingChange[];
  compatibility: CompatibilitySection;
  upgrade_instructions: string;
  known_issues: KnownIssue[];
}

export interface VersionInfoSection {
  version_number: string;
  release_date: string;
  release_type: ReleaseType;
}

export interface BugFix {
  id: string;
  description: string;
  issue_number?: string;
}

export interface BreakingChange {
  description: string;
  migration_steps: string;
  configuration_changes: string;
}

export interface CompatibilitySection {
  server_version_requirements: string;
  client_version_requirements: string;
  os_requirements: string[];
}

export interface KnownIssue {
  id: string;
  description: string;
  workaround?: string;
  planned_fix?: string;
}

export interface PackageInfo {
  id: string;
  package_type: PackageType;
  file_name: string;
  file_size: number;
  download_url: string;
  checksum_sha256: string;
  digital_signature?: string;
  os?: string;
  architecture?: string;
  uploaded_at: string;
  uploaded_by: string;
}

export interface CompatibilityMatrix {
  id: string;
  product_id: string;
  version_number: string;
  min_server_version?: string;
  max_server_version?: string;
  recommended_server_version?: string;
  incompatible_versions: string[];
  validated_at: string;
  validated_by: string;
  validation_status: ValidationStatus;
  validation_errors: string[];
}

export interface UpgradePath {
  id: string;
  product_id: string;
  from_version: string;
  to_version: string;
  path_type: UpgradePathType;
  intermediate_versions?: string[];
  is_blocked: boolean;
  block_reason?: string;
  created_at: string;
}

export interface Notification {
  id: string;
  type: NotificationType;
  recipient_id: string;
  product_id?: string;
  version_id?: string;
  customer_id?: string;
  tenant_id?: string;
  deployment_id?: string;
  title: string;
  message: string;
  priority: NotificationPriority;
  is_read: boolean;
  read_at?: string;
  created_at: string;
}

export interface UpdateDetection {
  id: string;
  endpoint_id: string;
  product_id: string;
  current_version: string;
  available_version: string;
  detected_at: string;
  last_checked_at: string;
}

export interface UpdateRollout {
  id: string;
  endpoint_id: string;
  product_id: string;
  from_version: string;
  to_version: string;
  status: RolloutStatus;
  initiated_by: string;
  initiated_at: string;
  started_at?: string;
  completed_at?: string;
  failed_at?: string;
  error_message?: string;
  progress: number; // 0-100
}

export interface AuditLog {
  id: string;
  action: AuditAction;
  resource_type: string;
  resource_id: string;
  user_id: string;
  user_email: string;
  details: Record<string, any>;
  ip_address: string;
  user_agent: string;
  timestamp: string;
}

// Customer Management Models
export interface NotificationPreferences {
  email_enabled: boolean;
  in_app_enabled: boolean;
  uat_notifications: boolean;
  production_notifications: boolean;
}

export interface Customer {
  id: string;
  customer_id: string;
  name: string;
  organization_name?: string;
  email: string;
  phone?: string;
  address?: string;
  account_status: CustomerStatus;
  notification_preferences: NotificationPreferences;
  created_at: string;
  updated_at: string;
}

export interface CustomerTenant {
  id: string;
  tenant_id: string;
  customer_id: string;
  name: string;
  description?: string;
  status: TenantStatus;
  created_at: string;
  updated_at: string;
}

export interface Deployment {
  id: string;
  deployment_id: string;
  tenant_id: string;
  product_id: string;
  deployment_type: DeploymentType;
  installed_version: string;
  number_of_users?: number;
  license_info?: string;
  server_hostname?: string;
  environment_details?: string;
  deployment_date: string;
  last_updated_date: string;
  status: DeploymentStatus;
}

// Request DTOs
export interface CreateProductRequest {
  product_id: string;
  name: string;
  type: ProductType;
  description?: string;
  vendor?: string;
}

export interface CreateVersionRequest {
  version_number: string;
  release_date: string;
  release_type: ReleaseType;
  eol_date?: string;
  min_server_version?: string;
  max_server_version?: string;
  recommended_server_version?: string;
  release_notes?: ReleaseNotes;
}

export interface UpdateVersionRequest {
  release_date?: string;
  release_type?: ReleaseType;
  eol_date?: string;
  min_server_version?: string;
  max_server_version?: string;
  recommended_server_version?: string;
  release_notes?: ReleaseNotes;
}

export interface ApproveVersionRequest {
  approved_by: string;
}

export interface ValidateCompatibilityRequest {
  min_server_version?: string;
  max_server_version?: string;
  recommended_server_version?: string;
  incompatible_versions?: string[];
}

export interface InitiateRolloutRequest {
  to_version: string;
}

// Customer Management Request DTOs
export interface CreateCustomerRequest {
  customer_id?: string;
  name: string;
  organization_name?: string;
  email: string;
  phone?: string;
  address?: string;
  account_status: CustomerStatus;
  notification_preferences: NotificationPreferences;
}

export interface UpdateCustomerRequest {
  name?: string;
  organization_name?: string;
  email?: string;
  phone?: string;
  address?: string;
  account_status?: CustomerStatus;
  notification_preferences?: NotificationPreferences;
}

export interface CreateTenantRequest {
  tenant_id?: string;
  name: string;
  description?: string;
  status: TenantStatus;
}

export interface UpdateTenantRequest {
  name?: string;
  description?: string;
  status?: TenantStatus;
}

export interface CreateDeploymentRequest {
  deployment_id?: string;
  product_id: string;
  deployment_type: DeploymentType;
  installed_version: string;
  number_of_users?: number;
  license_info?: string;
  server_hostname?: string;
  environment_details?: string;
  status: DeploymentStatus;
}

export interface UpdateDeploymentRequest {
  deployment_type?: DeploymentType;
  installed_version?: string;
  number_of_users?: number;
  license_info?: string;
  server_hostname?: string;
  environment_details?: string;
  status?: DeploymentStatus;
}

// Response DTOs
export interface PaginationResponse {
  page: number;
  limit: number;
  total: number;
  total_pages: number;
}

export interface PaginatedResponse<T> {
  data: T[];
  pagination: PaginationResponse;
}

export interface ErrorResponse {
  error: {
    code: string;
    message: string;
    details?: Record<string, any>;
  };
}

// API Response Types
export interface ProductsResponse extends PaginatedResponse<Product> {
  products: Product[];
}

export interface VersionsResponse extends PaginatedResponse<Version> {
  versions: Version[];
}

export interface NotificationsResponse extends PaginatedResponse<Notification> {
  notifications: Notification[];
}

export interface UpdateDetectionsResponse extends PaginatedResponse<UpdateDetection> {
  updates: UpdateDetection[];
}

export interface RolloutsResponse extends PaginatedResponse<UpdateRollout> {
  rollouts: UpdateRollout[];
}

export interface AuditLogsResponse extends PaginatedResponse<AuditLog> {
  audit_logs: AuditLog[];
}

// Query Parameters
export interface ListProductsQuery {
  type?: ProductType;
  is_active?: boolean;
  page?: number;
  limit?: number;
}

export interface ListVersionsQuery {
  product_id?: string;
  state?: VersionState;
  release_type?: ReleaseType;
  page?: number;
  limit?: number;
}

export interface ListNotificationsQuery {
  is_read?: boolean;
  type?: NotificationType;
  priority?: NotificationPriority;
  page?: number;
  limit?: number;
}

export interface ListAuditLogsQuery {
  resource_type?: string;
  resource_id?: string;
  user_id?: string;
  action?: AuditAction;
  start_date?: string;
  end_date?: string;
  page?: number;
  limit?: number;
}

export interface UpgradePathsQuery {
  from_version?: string;
  to_version?: string;
}

// Customer Management Query Parameters
export interface ListCustomersQuery {
  search?: string;
  status?: CustomerStatus;
  email?: string;
  page?: number;
  limit?: number;
}

export interface ListTenantsQuery {
  status?: TenantStatus;
  page?: number;
  limit?: number;
}

export interface ListDeploymentsQuery {
  product_id?: string;
  deployment_type?: DeploymentType;
  status?: DeploymentStatus;
  version?: string;
  page?: number;
  limit?: number;
}

// Customer Management Statistics
export interface CustomerStatistics {
  total_tenants: number;
  total_deployments: number;
  total_users: number;
  deployments_by_product: Record<string, number>;
  deployments_by_type: Record<string, number>;
}

export interface TenantStatistics {
  total_deployments: number;
  total_users: number;
  deployments_by_product: Record<string, number>;
  deployments_by_type: Record<string, number>;
}

// Pending Updates Types
export interface AvailableUpdate {
  version_number: string;
  release_date: string;
  release_type: string;
  is_security_update: boolean;
  compatibility_status: string;
  upgrade_path: string[];
}

export interface PendingUpdatesResponse {
  deployment_id: string;
  product_id: string;
  current_version: string;
  latest_version: string;
  update_count: number;
  priority: 'critical' | 'high' | 'normal';
  version_gap_type: 'patch' | 'minor' | 'major';
  available_updates: AvailableUpdate[];
  tenant_id?: string;
  tenant_name?: string;
  customer_id?: string;
  customer_name?: string;
  deployment_type?: DeploymentType;
}

export interface TenantPendingUpdatesSummary {
  tenant_id: string;
  tenant_name: string;
  total_deployments: number;
  deployments_with_updates: number;
  total_pending_update_count: number;
  by_priority: Record<string, number>;
  by_product: Record<string, number>;
  deployments: PendingUpdatesResponse[];
}

export interface CustomerPendingUpdatesSummary {
  customer_id: string;
  customer_name: string;
  total_deployments: number;
  deployments_with_updates: number;
  total_pending_update_count: number;
  by_priority: Record<string, number>;
  by_product: Record<string, number>;
  by_tenant: Record<string, number>;
  deployments: PendingUpdatesResponse[];
}

export interface PendingUpdatesQuery {
  product_id?: string;
  deployment_type?: DeploymentType;
  priority?: 'critical' | 'high' | 'normal';
  tenant_id?: string;
  customer_id?: string;
  page?: number;
  limit?: number;
}

// License Management Models
export interface Subscription {
  id: string;
  subscription_id: string;
  customer_id: string;
  name?: string;
  description?: string;
  start_date: string;
  end_date?: string;
  status: SubscriptionStatus;
  created_at: string;
  updated_at: string;
  created_by: string;
  notes?: string;
}

export interface License {
  id: string;
  license_id: string;
  subscription_id: string;
  product_id: string;
  license_type: LicenseType;
  number_of_seats: number;
  start_date: string;
  end_date?: string;
  status: LicenseStatus;
  assigned_by: string;
  assignment_date: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

export interface LicenseAllocation {
  id: string;
  allocation_id: string;
  license_id: string;
  tenant_id?: string;
  deployment_id?: string;
  number_of_seats_allocated: number;
  allocation_date: string;
  allocated_by: string;
  status: AllocationStatus;
  released_date?: string;
  released_by?: string;
  notes?: string;
  created_at: string;
  updated_at: string;
}

// License Management Request DTOs
export interface CreateSubscriptionRequest {
  subscription_id: string;
  name?: string;
  description?: string;
  start_date: string;
  end_date?: string;
  status: SubscriptionStatus;
  notes?: string;
}

export interface UpdateSubscriptionRequest {
  name?: string;
  description?: string;
  start_date?: string;
  end_date?: string;
  status?: SubscriptionStatus;
  notes?: string;
}

export interface CreateLicenseRequest {
  license_id: string;
  product_id: string;
  license_type: LicenseType;
  number_of_seats: number;
  start_date: string;
  end_date?: string;
  status: LicenseStatus;
  notes?: string;
}

export interface UpdateLicenseRequest {
  license_type?: LicenseType;
  number_of_seats?: number;
  start_date?: string;
  end_date?: string;
  status?: LicenseStatus;
  notes?: string;
}

export interface AllocateLicenseRequest {
  tenant_id?: string;
  deployment_id?: string;
  number_of_seats_allocated: number;
  notes?: string;
}

export interface ReleaseLicenseRequest {
  notes?: string;
}

// License Management Query Types
export interface ListSubscriptionsQuery {
  status?: SubscriptionStatus;
  page?: number;
  limit?: number;
}

export interface ListLicensesQuery {
  product_id?: string;
  license_type?: LicenseType;
  status?: LicenseStatus;
  page?: number;
  limit?: number;
}

export interface ListLicenseAllocationsQuery {
  status?: AllocationStatus;
  page?: number;
  limit?: number;
}

// License Management Statistics
export interface SubscriptionStatistics {
  total_licenses: number;
  active_licenses: number;
  expired_licenses: number;
  total_seats: number;
  perpetual_licenses: number;
  time_based_licenses: number;
}

export interface LicenseStatistics {
  total_seats: number;
  allocated_seats: number;
  available_seats: number;
  utilization_percent: number;
  active_allocations: number;
  license_type: LicenseType;
  status: LicenseStatus;
}

export interface LicenseUtilization {
  total_seats: number;
  allocated_seats: number;
  available_seats: number;
  utilization_percent: number;
  active_allocations: number;
}

// Form Types
export interface ProductFormData {
  product_id: string;
  name: string;
  type: ProductType;
  description: string;
  vendor: string;
}

export interface VersionFormData {
  version_number: string;
  release_date: string;
  release_type: ReleaseType;
  eol_date?: string;
  min_server_version?: string;
  max_server_version?: string;
  recommended_server_version?: string;
  release_notes: {
    version_info: {
      version_number: string;
      release_date: string;
      release_type: ReleaseType;
    };
    whats_new: string[];
    bug_fixes: Array<{
      id: string;
      description: string;
      issue_number?: string;
    }>;
    breaking_changes: Array<{
      description: string;
      migration_steps: string;
      configuration_changes: string;
    }>;
    compatibility: {
      server_version_requirements: string;
      client_version_requirements: string;
      os_requirements: string[];
    };
    upgrade_instructions: string;
    known_issues: Array<{
      id: string;
      description: string;
      workaround?: string;
      planned_fix?: string;
    }>;
  };
}

export interface PackageUploadFormData {
  file: File;
  package_type: PackageType;
  os?: string;
  architecture?: string;
}

// UI State Types
export interface VersionStateTransition {
  from: VersionState;
  to: VersionState;
  label: string;
  action: string;
  requiresApproval: boolean;
}

export interface VersionStateConfig {
  state: VersionState;
  label: string;
  color: string;
  icon: string;
  allowedTransitions: VersionState[];
}

export const VERSION_STATE_CONFIGS: Record<VersionState, VersionStateConfig> = {
  [VersionState.DRAFT]: {
    state: VersionState.DRAFT,
    label: 'Draft',
    color: 'gray',
    icon: 'edit',
    allowedTransitions: [VersionState.PENDING_REVIEW],
  },
  [VersionState.PENDING_REVIEW]: {
    state: VersionState.PENDING_REVIEW,
    label: 'Pending Review',
    color: 'yellow',
    icon: 'clock',
    allowedTransitions: [VersionState.APPROVED, VersionState.DRAFT],
  },
  [VersionState.APPROVED]: {
    state: VersionState.APPROVED,
    label: 'Approved',
    color: 'green',
    icon: 'check',
    allowedTransitions: [VersionState.RELEASED, VersionState.DRAFT],
  },
  [VersionState.RELEASED]: {
    state: VersionState.RELEASED,
    label: 'Released',
    color: 'blue',
    icon: 'rocket',
    allowedTransitions: [VersionState.DEPRECATED, VersionState.EOL],
  },
  [VersionState.DEPRECATED]: {
    state: VersionState.DEPRECATED,
    label: 'Deprecated',
    color: 'orange',
    icon: 'warning',
    allowedTransitions: [VersionState.EOL],
  },
  [VersionState.EOL]: {
    state: VersionState.EOL,
    label: 'End of Life',
    color: 'red',
    icon: 'stop',
    allowedTransitions: [],
  },
};

// Utility Types
export type ApiResponse<T> = T | ErrorResponse;

export type AsyncState<T> = {
  data: T | null;
  loading: boolean;
  error: ErrorResponse | null;
};

export type TableColumn<T> = {
  key: keyof T | string;
  label: string;
  sortable?: boolean;
  render?: (value: any, row: T) => React.ReactNode;
};

