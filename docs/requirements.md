# Accops Product Auto Update Functionality - Requirements Document

## 1. Executive Summary

### 1.1 Objective
Implement an automated update system and management portal for Accops products that streamlines version control, rollout of updates, and proactive admin notifications. This ensures product versions are current, compatible, and centrally managed to reduce manual overhead and upgrade errors.

### 1.2 Context
Accops products require continual maintenance and patching to deliver new features, security fixes, and compatibility improvements. Currently, upgrade processes are manual and lack centralized visibility. A modern auto-update system with a management portal and actionable UI elements is critical to enhancing operational efficiency and product lifecycle transparency for administrators.

### 1.3 Key Benefits
- Centralized version control and visibility
- Automated update detection and notifications
- Reduced manual upgrade overhead
- Proactive compatibility checking
- Batch update capabilities
- Customer and installation management
- Automated customer notifications on version releases
- Audit trail for compliance

---

## 2. Functional Requirements

### 2.0 Accops Products

The Update Manager supports the following Accops products:

#### 2.0.1 Product Types

**Server Products:**
- **HyWorks**: Application virtualization and delivery platform
- **HySecure**: Secure remote access and VPN solution
- **IRIS**: Identity and access management platform
- **ARS** (Accops Resource Server): Resource management and provisioning

**Client Products:**
- **Clients for Linux**: Linux client applications
- **Clients for Windows**: Windows client applications
- **Mobile Clients**: iOS and Android client applications

**Other Products:**
- Additional Accops products as they are released

#### 2.0.2 Product Categories

Products are categorized by:
- **Deployment Type**: Server vs. Client
- **Platform**: Windows, Linux, Mobile
- **Update Strategy**: Stateless (blue-green) vs. Stateful (in-place)
- **Tenant Association**: Multi-tenant vs. Single-tenant

#### 2.0.3 Product-Specific Requirements

- **FR-0.1**: Each product type may have specific update procedures
- **FR-0.2**: Server products (HyWorks, HySecure, IRIS, ARS) support multi-tenant deployments
- **FR-0.3**: Client products require endpoint-based update distribution
- **FR-0.4**: Product compatibility matrices must be maintained per product type
- **FR-0.5**: Server products use blue-green deployment for zero-downtime updates
- **FR-0.6**: Client products use in-place updates with service restart

### 2.1 Update Management Portal

#### 2.1.1 Dashboard Overview
- **FR-1.1**: Display all products installed across VMs with current version status
- **FR-1.2**: Show available server/latest version and upgrade paths for each product
- **FR-1.3**: Provide real-time or near real-time status updates
- **FR-1.4**: Support filtering and search capabilities for large deployments
- **FR-1.4a**: Display customer deployments with pending updates in the Updates page
- **FR-1.4b**: Show deployment-level pending updates with customer, tenant, and product context

#### 2.1.2 Visual Indicators
- **FR-1.5**: Display a green dot indicator beside VMs that have pending updates
- **FR-1.6**: Green dot should be clickable and provide contextual information
- **FR-1.7**: Visual distinction between different update types (patch, minor, major)

#### 2.1.3 Update Details View
- **FR-1.8**: Clicking the green dot opens a detailed view containing:
  - Release notes for the available update
  - Current version vs. target version comparison
  - Compatibility information
  - Estimated update time
  - Option to start the update immediately
  - Option to schedule the update

#### 2.1.4 Update Actions
- **FR-1.9**: Support immediate update initiation from the detail view
- **FR-1.10**: Support scheduled updates with date/time selection
- **FR-1.11**: Display update progress in real-time
- **FR-1.12**: Show update history and status for each endpoint

---

### 2.2 Update Agent on Endpoints

#### 2.2.1 Service Architecture
- **FR-2.1**: Run as a background service on each product endpoint (VM or server)
- **FR-2.2**: Support multiple operating systems (Windows, Linux)
- **FR-2.3**: Auto-start on system boot
- **FR-2.4**: Lightweight and resource-efficient

#### 2.2.2 Version Detection
- **FR-2.5**: Periodically check the installed product version
- **FR-2.6**: Detect all Accops products installed on the endpoint
- **FR-2.7**: Report product metadata (name, version, installation path, etc.)

#### 2.2.3 Portal Communication
- **FR-2.8**: Query the central portal to determine if updates are available
- **FR-2.9**: Report status to the portal for visibility and green dot activation
- **FR-2.10**: Support secure communication with mutual authentication
- **FR-2.11**: Handle network failures gracefully with retry logic

#### 2.2.4 Update Execution
- **FR-2.12**: Receive update commands from the portal
- **FR-2.13**: Download update packages securely
- **FR-2.14**: Execute updates with proper error handling
- **FR-2.15**: Support rollback or abort option in case of failure
- **FR-2.16**: Report update status (in-progress, success, failure) back to portal

---

### 2.3 Batch Update Rollout

#### 2.3.1 Selection Interface
- **FR-3.1**: Portal supports selecting multiple VMs or endpoints
- **FR-3.2**: Support bulk selection (select all, select by filter, etc.)
- **FR-3.3**: Display selection count and summary information

#### 2.3.2 Batch Operations
- **FR-3.4**: Admin can initiate simultaneous updates on selected hosts with a single action
- **FR-3.5**: Support bulk upgrades to the latest version
- **FR-3.6**: Support controlled rollout option (staged updates)
- **FR-3.7**: Allow specifying update order or priority

#### 2.3.3 Rollout Management
- **FR-3.8**: Display progress for each endpoint in the batch
- **FR-3.9**: Handle partial failures gracefully
- **FR-3.10**: Report per-host status (success, failure, in-progress, pending)
- **FR-3.11**: Support pausing/resuming batch updates
- **FR-3.12**: Allow canceling individual or all pending updates in a batch

---

### 2.4 Product Lifecycle Page and API

#### 2.4.1 Lifecycle Display
- **FR-4.1**: Display detailed upgrade paths based on:
  - Product name
  - Installed version
  - Server version for compatibility consideration (optional)
- **FR-4.2**: Show version compatibility matrix
- **FR-4.3**: Display end-of-life (EOL) information
- **FR-4.4**: Show recommended upgrade paths

#### 2.4.2 API Endpoints
- **FR-4.5**: Provide RESTful API returning output in JSON format
- **FR-4.6**: Support integration with automation and audit systems
- **FR-4.7**: API endpoints for:
  - Querying available versions for a product
  - Getting upgrade path between versions
  - Checking compatibility between versions
  - Retrieving release notes

#### 2.4.3 Compatibility Validation
- **FR-4.8**: Validate backward compatibility before approving updates
- **FR-4.9**: Reduce risk of system conflicts
- **FR-4.10**: Check dependencies and prerequisites
- **FR-4.11**: Warn about breaking changes

---

### 2.5 Admin Notifications

#### 2.5.1 Automatic Detection
- **FR-5.1**: IRIS automatically detects when a new product version releases
- **FR-5.2**: Monitor version repositories or release channels
- **FR-5.3**: Support multiple notification triggers (immediate, daily digest, weekly summary)

#### 2.5.2 Notification Content
- **FR-5.4**: Send administrative notifications highlighting:
  - Version details (version number, release date)
  - Release notes summary
  - Affected products and endpoints
  - Urgency level (security, feature, maintenance)

#### 2.5.3 Notification Channels
- **FR-5.5**: In-app notifications within the portal
- **FR-5.6**: Email notifications to registered admins
- **FR-5.7**: Support notification preferences per admin

#### 2.5.4 Notification Actions
- **FR-5.8**: Notifications allow quick access to update portal for immediate action
- **FR-5.9**: Direct links to release notes and update initiation
- **FR-5.10**: Mark notifications as read/dismissed

---

### 2.6 Customer Management

#### 2.6.1 Customer Registration
- **FR-6.1**: Portal supports customer registration and management
- **FR-6.2**: Each customer has a unique customer ID
- **FR-6.3**: Customer information includes:
  - Customer name
  - Organization/Company name
  - Contact information (email, phone)
  - Address
  - Account status (Active, Inactive, Suspended)
  - Created date and last updated date
- **FR-6.4**: Support customer search and filtering
- **FR-6.5**: Customer list view with pagination

#### 2.6.2 Tenant Management
- **FR-6.6**: Each customer can have multiple tenants
- **FR-6.7**: A tenant represents an independent deployment of one or more products
- **FR-6.8**: Each tenant must be associated with:
  - Customer ID
  - Tenant name/identifier
  - Tenant description (optional)
  - Tenant status (Active, Inactive)
  - Created date and last updated date
- **FR-6.9**: Support creating, viewing, updating, and deleting tenants
- **FR-6.10**: Display all tenants for a customer in a list view
- **FR-6.11**: Filter tenants by status

#### 2.6.3 Deployment Management
- **FR-6.12**: Each tenant can have multiple deployments
- **FR-6.13**: A deployment is a combination of a product and deployment type (UAT/Testing, Production)
- **FR-6.14**: Each deployment must be associated with:
  - Tenant ID
  - Product ID
  - Deployment type (UAT/Testing, Production)
  - Installed version number
  - Number of users (integer, optional)
  - License information (free text, optional)
  - Deployment date
  - Server/hostname (optional)
  - Environment details (optional)
  - Status (Active, Inactive)
- **FR-6.15**: Support creating, viewing, updating, and deleting deployments
- **FR-6.16**: Display all deployments for a tenant in a list view
- **FR-6.17**: Filter deployments by product, type (UAT/Production), or version
- **FR-6.18**: Display number of users and license information in deployment list and details

#### 2.6.4 Version Tracking
- **FR-6.19**: Customers can specify and update version details for each deployment
- **FR-6.20**: System tracks current installed version for each deployment
- **FR-6.21**: Display version comparison (current vs. latest available)
- **FR-6.22**: Show version status (up-to-date, update available, outdated)
- **FR-6.23**: Support bulk version updates for multiple deployments

#### 2.6.5 Deployment Types
- **FR-6.24**: Deployment type must be specified as either:
  - **UAT/Testing**: For testing and validation environments
  - **Production**: For live production environments
- **FR-6.25**: Visual distinction between UAT and Production deployments
- **FR-6.26**: Filter and group deployments by type
- **FR-6.27**: Different notification preferences for UAT vs. Production

#### 2.6.6 Notification Generation
- **FR-6.28**: System automatically generates notifications when new versions are rolled out
- **FR-6.29**: Notifications are sent to customers who have deployments of the updated product
- **FR-6.30**: Notification content includes:
  - New version number and release date
  - Release notes summary
  - List of affected deployments (with current versions, tenant, and deployment type)
  - Recommended action (update available)
  - Direct link to view full release notes
- **FR-6.31**: Separate notifications for UAT and Production deployments
- **FR-6.32**: Notification priority based on:
  - Deployment type (Production gets higher priority)
  - Version gap (major updates get higher priority)
  - Security releases (critical priority)
- **FR-6.33**: Support notification preferences per customer
- **FR-6.34**: Track notification delivery and read status

#### 2.6.7 Customer Dashboard
- **FR-6.35**: Customer-specific dashboard showing:
  - Total tenants count
  - Total deployments count
  - Total number of users across all deployments
  - Deployments by product
  - Deployments by type (UAT/Production)
  - Deployments by tenant
  - Pending updates count (aggregate across all deployments)
  - Deployments with pending updates (count and list)
  - Recent version updates
  - License summary information
- **FR-6.36**: Quick view of deployments requiring updates:
  - List of deployments with pending updates
  - Current version vs. latest available version for each deployment
  - Update priority indicator (critical, high, normal)
  - Direct link to deployment details or update view
- **FR-6.37**: Filter deployments by update status:
  - All deployments
  - Deployments with pending updates
  - Deployments up-to-date
  - Deployments with critical updates
  - Deployments with security updates

#### 2.6.8 Tenant Details View
- **FR-6.38**: Detailed view for each tenant showing:
  - Tenant information
  - List of deployments in the tenant
  - Tenant statistics (total deployments, users, etc.)
  - Actions: Edit tenant, Add deployment, Delete tenant

#### 2.6.9 Deployment Details View
- **FR-6.39**: Detailed view for each deployment showing:
  - Deployment metadata (tenant, product, type)
  - Number of users
  - License information
  - Current version information
  - Available updates
  - Update history
  - Compatibility information
  - Related notifications
- **FR-6.40**: Action buttons to:
  - Update version information
  - Update number of users and license information
  - View available updates
  - Initiate update (if applicable)
  - View release notes

#### 2.6.10 Pending Updates Tracking
- **FR-6.41**: System automatically tracks pending updates for each deployment by comparing:
  - Currently deployed version (installed_version in deployment)
  - Latest released version for the product
  - All released versions between current and latest
- **FR-6.42**: Pending update calculation logic:
  - Identify all released versions for the deployment's product
  - Filter versions that are newer than the deployment's installed_version
  - Consider only versions with state "Released" or "Available"
  - Exclude versions that are deprecated or EOL
  - Sort by version number (newest first)
- **FR-6.43**: Pending update information includes:
  - Count of available updates (number of newer versions)
  - Latest available version number
  - Version gap type (patch, minor, major)
  - Security update indicator (if any pending updates are security releases)
  - Release dates of pending updates
  - Compatibility status with current deployment
- **FR-6.44**: Display pending updates at deployment level:
  - Show pending update count badge on deployment list
  - Display pending updates section in deployment details view
  - Show version comparison (current vs. latest available)
  - List all intermediate versions between current and latest
  - Indicate if direct upgrade path exists or intermediate upgrades required
- **FR-6.45**: Display pending updates at tenant level:
  - Aggregate pending updates count across all deployments in a tenant
  - Show tenant-level summary of deployments requiring updates
  - Filter deployments by update status (up-to-date, updates available, critical updates)
- **FR-6.46**: Display pending updates at customer level:
  - Aggregate pending updates count across all deployments for a customer
  - Show customer-level summary of:
    - Total deployments with pending updates
    - Total pending update count (sum of all available updates)
    - Deployments by update priority (critical, high, normal)
    - Deployments by product requiring updates
    - Deployments by tenant requiring updates
- **FR-6.47**: Display pending updates in Updates page:
  - Show all deployments with pending updates across all customers
  - Filter by product, customer, tenant, or deployment type
  - Group by customer or product
  - Sort by update priority, version gap, or deployment type
  - Show deployment details: customer, tenant, product, current version, latest version
  - Display update path (direct upgrade or intermediate steps required)
- **FR-6.48**: Update priority calculation:
  - **Critical**: Security releases or EOL approaching
  - **High**: Major version updates or Production deployments with significant version gap
  - **Normal**: Minor or patch updates for non-critical deployments
- **FR-6.49**: Real-time update tracking:
  - Recalculate pending updates when new versions are released
  - Recalculate pending updates when deployment version is updated
  - Update pending counts immediately in all views (customer, tenant, deployment, updates page)
- **FR-6.50**: Pending updates API endpoints:
  - GET `/api/v1/customers/{customer_id}/deployments/pending-updates` - Get all pending updates for customer
  - GET `/api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/pending-updates` - Get pending updates for tenant
  - GET `/api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/updates` - Get available updates for deployment
  - GET `/api/v1/updates/pending` - Get all pending updates across all customers (admin view)

### 2.7 License Management

#### 2.7.1 Subscription Management
- **FR-7.1**: Each customer can have one or more subscriptions
- **FR-7.2**: Each subscription must be associated with:
  - Customer ID
  - Subscription ID (unique identifier)
  - Subscription name/description (optional)
  - Start date
  - End date (optional, for time-based subscriptions)
  - Status (Active, Inactive, Expired, Suspended)
  - Created date and last updated date
- **FR-7.3**: Support creating, viewing, updating, and deleting subscriptions
- **FR-7.4**: Display all subscriptions for a customer in a list view
- **FR-7.5**: Filter subscriptions by status, date range, or product
- **FR-7.6**: Support subscription renewal and extension
- **FR-7.7**: Track subscription history and changes

#### 2.7.2 License Assignment
- **FR-7.8**: Sales team can assign licenses to subscriptions
- **FR-7.9**: Each license must be associated with:
  - Subscription ID
  - License ID (unique identifier)
  - Product ID (specific product the license is for)
  - License type (Perpetual, Time-based)
  - Number of users/seats (integer)
  - Start date
  - End date (for time-based licenses, optional for perpetual)
  - Status (Active, Inactive, Expired, Revoked)
  - Assigned by (sales team member/user ID)
  - Assignment date
  - Notes (optional, free text)
- **FR-7.10**: Support creating, viewing, updating, and deleting licenses
- **FR-7.11**: Display all licenses for a subscription in a list view
- **FR-7.12**: Filter licenses by product, type, status, or date range
- **FR-7.13**: Track license assignment history and audit trail
- **FR-7.14**: Support bulk license assignment for multiple products

#### 2.7.3 License Types
- **FR-7.15**: License type must be specified as either:
  - **Perpetual**: No expiration date, valid indefinitely
  - **Time-based**: Has start and end dates, expires after end date
- **FR-7.16**: Visual distinction between Perpetual and Time-based licenses
- **FR-7.17**: Automatic status update for time-based licenses:
  - Active: Current date is between start and end date
  - Expired: Current date is after end date
  - Pending: Current date is before start date
- **FR-7.18**: Display expiration warnings for licenses expiring within 30/60/90 days
- **FR-7.19**: Support license renewal and extension for time-based licenses

#### 2.7.4 License Distribution
- **FR-7.20**: Customer can distribute licenses across multiple tenants and deployments
- **FR-7.21**: License allocation must track:
  - License ID
  - Tenant ID (optional, if allocated to tenant)
  - Deployment ID (optional, if allocated to specific deployment)
  - Number of users/seats allocated
  - Allocation date
  - Allocated by (customer user ID)
  - Status (Active, Released)
- **FR-7.22**: Support allocating licenses from subscription to tenants/deployments
- **FR-7.23**: Support releasing allocated licenses back to subscription pool
- **FR-7.24**: Track total allocated vs. available licenses per subscription
- **FR-7.25**: Prevent over-allocation (cannot allocate more seats than available in license)
- **FR-7.26**: Display license utilization:
  - Total licenses per subscription
  - Allocated licenses count
  - Available licenses count
  - Utilization percentage
- **FR-7.27**: Support partial allocation (allocate subset of seats from a license)
- **FR-7.28**: Display allocation history and changes

#### 2.7.5 License Validation
- **FR-7.29**: System validates license availability before allocation
- **FR-7.30**: Check license expiration before allowing new allocations
- **FR-7.31**: Validate product match (license must be for the product being deployed)
- **FR-7.32**: Enforce seat limits (cannot exceed licensed user count)
- **FR-7.33**: Display validation errors and warnings
- **FR-7.34**: Support license compliance reporting

#### 2.7.6 License Dashboard
- **FR-7.35**: Customer-specific license dashboard showing:
  - Total subscriptions count
  - Total licenses count (by type: Perpetual, Time-based)
  - Active licenses count
  - Expired licenses count
  - Licenses expiring soon (30/60/90 days)
  - Total licensed users/seats
  - Allocated users/seats
  - Available users/seats
  - Utilization percentage
  - Licenses by product
  - Licenses by subscription
- **FR-7.36**: Subscription details view showing:
  - Subscription information
  - List of licenses in the subscription
  - License allocation summary
  - Available vs. allocated seats
  - Expiration timeline
- **FR-7.37**: License details view showing:
  - License information (product, type, seats, dates)
  - Allocation history
  - Current allocations (tenants/deployments)
  - Utilization metrics
  - Expiration status and warnings

#### 2.7.7 License Reporting
- **FR-7.38**: Generate license reports:
  - License inventory by customer
  - License utilization by product
  - Expiring licenses report
  - License allocation by tenant/deployment
  - License compliance report
- **FR-7.39**: Export license data (CSV, JSON)
- **FR-7.40**: Filter and search licenses across all customers (admin view)
- **FR-7.41**: Track license usage trends over time

#### 2.7.8 Integration with Deployments
- **FR-7.42**: Link deployments to allocated licenses
- **FR-7.43**: Display license information in deployment details
- **FR-7.44**: Show license compliance status for each deployment
- **FR-7.45**: Validate deployment user count against allocated license seats
- **FR-7.46**: Warn when deployment user count exceeds allocated license seats
- **FR-7.47**: Support license reallocation when deployment user count changes

#### 2.7.9 License API Endpoints
- **FR-7.48**: Subscription API endpoints:
  - GET `/api/v1/customers/{customer_id}/subscriptions` - List all subscriptions for customer
  - POST `/api/v1/customers/{customer_id}/subscriptions` - Create new subscription
  - GET `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}` - Get subscription details
  - PUT `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}` - Update subscription
  - DELETE `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}` - Delete subscription
  - GET `/api/v1/subscriptions` - List all subscriptions (admin view)
- **FR-7.49**: License API endpoints:
  - GET `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses` - List licenses in subscription
  - POST `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses` - Assign license to subscription
  - GET `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}` - Get license details
  - PUT `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}` - Update license
  - DELETE `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}` - Revoke license
  - GET `/api/v1/licenses` - List all licenses (admin view)
- **FR-7.50**: License Allocation API endpoints:
  - POST `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}/allocate` - Allocate license to tenant/deployment
  - POST `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}/release` - Release allocated license
  - GET `/api/v1/customers/{customer_id}/subscriptions/{subscription_id}/licenses/{license_id}/allocations` - Get allocation history
  - GET `/api/v1/customers/{customer_id}/tenants/{tenant_id}/licenses` - Get licenses allocated to tenant
  - GET `/api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/licenses` - Get licenses allocated to deployment

---

## 3. Multi-Tenant Requirements

### 3.1 Blue-Green Stateless Edge Upgrade
- **FR-7.1**: Tenant edge deployments are stateless and updated by deploying a parallel environment ("Blue")
- **FR-7.2**: Switch tenant traffic post-validation to the new environment
- **FR-7.3**: Rollback is automatic and near-instant if issues arise
- **FR-7.4**: Ensure zero downtime during upgrades
- **FR-7.5**: Support traffic validation before full cutover

### 3.2 Forced Upgrade Policy & Version Enforcement
- **FR-7.6**: Admins are notified of new releases and must upgrade before reaching last supported version
- **FR-7.7**: No tenants remain on end-of-life versions
- **FR-7.8**: System enforces maximum supported versions for all active tenants
- **FR-7.9**: Display warnings for approaching EOL versions
- **FR-7.10**: Support grace period before forced upgrades

---

## 4. Non-Functional Requirements

### 4.1 Security
- **NFR-1**: Secure communication between update agents and portal with mutual authentication
- **NFR-2**: Encrypt update packages in transit and at rest
- **NFR-3**: Support certificate-based authentication
- **NFR-4**: Role-based access control (RBAC) for portal access
- **NFR-5**: Audit logging for all update operations

### 4.2 Performance
- **NFR-6**: Dashboard must refresh in near real-time to reflect upgrade status changes
- **NFR-7**: Support large-scale deployments (1000+ endpoints)
- **NFR-8**: Efficient batch update processing
- **NFR-9**: Minimal impact on endpoint performance during version checks

### 4.3 Reliability
- **NFR-10**: Update agent must support rollback or abort option in case of failure
- **NFR-11**: Handle network interruptions gracefully
- **NFR-12**: Validate update integrity before installation
- **NFR-13**: Support resume of interrupted downloads

### 4.4 Usability
- **NFR-14**: Intuitive user interface with clear visual indicators
- **NFR-15**: Comprehensive release notes display
- **NFR-16**: Clear error messages and troubleshooting guidance
- **NFR-17**: Responsive design for various screen sizes

### 4.5 Scalability
- **NFR-18**: Support horizontal scaling of portal services
- **NFR-19**: Efficient database queries for large datasets
- **NFR-20**: Support concurrent batch updates

---

## 5. Architectural Components

### 5.1 Update Manager Portal
- **Description**: Central UI for status, release notes, batch updates, customer management
- **Technology**: Web-based application (React/Vue/Angular)
- **Responsibilities**:
  - Display product and version status
  - Manage batch update operations
  - Show release notes and compatibility information
  - Handle admin notifications
  - Manage customers and installations
  - Customer dashboard and installation tracking

### 5.2 Endpoint Update Agent
- **Description**: Lightweight client checks version, communicates with portal
- **Technology**: Cross-platform service (Python/Go)
- **Responsibilities**:
  - Detect installed product versions
  - Query portal for available updates
  - Execute updates with rollback support
  - Report status to portal

### 5.3 Version Control API
- **Description**: Stores and exposes product lifecycle and version compatibility
- **Technology**: RESTful API (FastAPI/Flask/Express)
- **Responsibilities**:
  - Manage product version metadata
  - Provide upgrade path calculations
  - Validate compatibility
  - Serve release notes

### 5.4 Notification System
- **Description**: Pushes in-app and email alerts on new releases
- **Technology**: Message queue and email service
- **Responsibilities**:
  - Detect new version releases
  - Generate and send notifications to admins and customers
  - Identify affected customers based on installations
  - Manage notification preferences per customer
  - Track notification delivery and read status
  - Support different notification priorities for UAT vs. Production

### 5.5 Audit & Logging
- **Description**: Tracks update status, success, failure for all endpoints
- **Technology**: Logging service and database
- **Responsibilities**:
  - Log all update operations
  - Track endpoint status changes
  - Generate audit reports
  - Support compliance requirements

---

## 6. Data Models

### 6.1 Core Entities

#### Product
- Product ID
- Product Name (free text, unique - e.g., "HyWorks", "HySecure", "IRIS", "ARS", "Client for Linux", "Client for Windows", "Client for Mobile")
- Product Type (Server, Client)
- Platform (Windows, Linux, Mobile, Cross-Platform)
- Update Strategy (Blue-Green, In-Place, App-Store)
- Multi-Tenant Support (Yes/No)
- Current Latest Version
- EOL Date
- Support Status
- Default Update Policy

#### Version
- Version ID
- Product ID
- Version Number (semantic versioning)
- Release Date
- Release Notes
- Compatibility Matrix
- Download URL
- Checksum

#### Customer
- Customer ID
- Customer Name
- Organization/Company Name
- Contact Information (email, phone, address)
- Account Status (Active, Inactive, Suspended)
- Notification Preferences
- Created Date
- Last Updated Date

#### Customer Tenant (Customer-Managed)
- Tenant ID
- Customer ID
- Tenant Name/Identifier
- Tenant Description (optional)
- Tenant Status (Active, Inactive)
- Created Date
- Last Updated Date

#### Deployment
- Deployment ID
- Tenant ID
- Product ID
- Deployment Type (UAT/Testing, Production)
- Installed Version Number
- Number of Users (integer, optional)
- License Information (free text, optional)
- Server/Hostname (optional)
- Environment Details (optional)
- Deployment Date
- Last Updated Date
- Status (Active, Inactive)

#### Subscription
- Subscription ID
- Customer ID
- Subscription Name/Description (optional)
- Start Date
- End Date (optional, for time-based subscriptions)
- Status (Active, Inactive, Expired, Suspended)
- Created Date
- Last Updated Date
- Created By (user ID)
- Notes (optional, free text)

#### License
- License ID
- Subscription ID
- Product ID
- License Type (Perpetual, Time-based)
- Number of Users/Seats (integer)
- Start Date
- End Date (for time-based licenses, optional for perpetual)
- Status (Active, Inactive, Expired, Revoked)
- Assigned By (sales team member/user ID)
- Assignment Date
- Notes (optional, free text)
- Created Date
- Last Updated Date

#### License Allocation
- Allocation ID
- License ID
- Tenant ID (optional, if allocated to tenant)
- Deployment ID (optional, if allocated to specific deployment)
- Number of Users/Seats Allocated (integer)
- Allocation Date
- Allocated By (customer user ID)
- Status (Active, Released)
- Released Date (optional, when released)
- Released By (optional, user ID who released)
- Notes (optional, free text)

#### System Tenant (Multi-Tenant Product Deployment)
- Tenant ID
- Tenant Name
- Organization Name
- Contact Information
- Subscription/License Information
- Update Policy Preferences
- Notification Preferences
- Maintenance Windows
- Compliance Status
- Created Date
- Last Updated Date

**Note:** This Tenant model is for multi-tenant product deployments (Section 3). The Customer Tenant model above is for customer-managed independent deployments (Section 2.6).

#### Endpoint
- Endpoint ID
- Tenant ID (for multi-tenant deployments)
- Endpoint Name
- IP Address / Hostname
- Operating System
- Endpoint Type (Server, Client, Edge)
- Blue-Green Environment (Blue/Green/None)
- Last Check-in Time
- Status (Online/Offline)
- Agent Version

#### Installed Product
- Installation ID
- Endpoint ID (optional, for agent-based installations)
- Tenant ID (optional, for tenant-managed deployments)
- Customer ID (optional, for customer-managed deployments)
- Product ID
- Installed Version
- Installation Date
- Installation Path
- Deployment Type (UAT/Testing, Production)

#### Update Job
- Job ID
- Tenant ID (optional, for tenant-scoped updates)
- Endpoint ID
- Product ID
- Source Version
- Target Version
- Status (Pending/In-Progress/Success/Failed/Rolled-Back)
- Update Type (Immediate, Scheduled, Forced)
- Blue-Green Status (Blue/Green/Switching)
- Start Time
- End Time
- Error Message
- Rollback Status
- Scheduled Time (if scheduled)

#### Update History
- History ID
- Job ID
- Tenant ID
- Endpoint ID
- Product ID
- Action (Update/Rollback)
- Status
- Timestamp
- Admin User
- Notes

#### Tenant Policy
- Policy ID
- Tenant ID
- Product ID (optional, for product-specific policies)
- Policy Type (Forced Upgrade, Maximum Version, Minimum Version)
- Version Constraints
- Grace Period (days)
- Enforcement Date
- Status (Active, Suspended)

---

## 7. API Specifications

### 7.1 Agent API Endpoints

#### POST /api/v1/agent/register
- Register a new endpoint with the portal
- Request: Endpoint metadata, installed products
- Response: Agent ID, configuration

#### POST /api/v1/agent/checkin
- Periodic check-in from agent
- Request: Endpoint ID, installed products, status
- Response: Available updates, commands

#### GET /api/v1/agent/updates/{endpoint_id}
- Get available updates for an endpoint
- Response: List of available updates

#### POST /api/v1/agent/update/status
- Report update status from agent
- Request: Job ID, status, progress, error details

### 7.2 Portal API Endpoints

#### GET /api/v1/portal/dashboard
- Get dashboard data
- Response: Endpoints, products, update status

#### GET /api/v1/portal/endpoints
- List all endpoints with status
- Query params: filter, search, pagination

#### GET /api/v1/portal/endpoints/{endpoint_id}/updates
- Get available updates for specific endpoint

#### POST /api/v1/portal/updates/batch
- Initiate batch update
- Request: Endpoint IDs, target version, schedule

#### GET /api/v1/portal/updates/jobs/{job_id}
- Get update job status

#### GET /api/v1/portal/updates/history
- Get update history
- Query params: endpoint_id, product_id, date_range

### 7.3 Lifecycle API Endpoints

#### GET /api/v1/lifecycle/products
- List all products

#### GET /api/v1/lifecycle/products/{product_id}/versions
- Get all versions for a product

#### GET /api/v1/lifecycle/upgrade-path
- Get upgrade path between versions
- Query params: product_id, from_version, to_version, server_version (optional)

#### GET /api/v1/lifecycle/compatibility
- Check version compatibility
- Query params: product_id, version1, version2

#### GET /api/v1/lifecycle/release-notes/{version_id}
- Get release notes for a version

### 7.4 Customer Management API Endpoints

#### GET /api/v1/customers
- List all customers
- Query params: search, status, pagination
- Response: List of customers

#### POST /api/v1/customers
- Create a new customer
- Request: Customer information
- Response: Created customer

#### GET /api/v1/customers/{customer_id}
- Get customer details
- Response: Customer information with installations

#### PUT /api/v1/customers/{customer_id}
- Update customer information
- Request: Updated customer data
- Response: Updated customer

#### DELETE /api/v1/customers/{customer_id}
- Delete customer (soft delete)
- Response: Success confirmation

#### GET /api/v1/customers/{customer_id}/tenants
- List all tenants for a customer
- Query params: status, pagination
- Response: List of tenants

#### POST /api/v1/customers/{customer_id}/tenants
- Create a new tenant for a customer
- Request: Tenant details (name, description, etc.)
- Response: Created tenant

#### GET /api/v1/customers/{customer_id}/tenants/{tenant_id}
- Get tenant details
- Response: Tenant information with deployments

#### PUT /api/v1/customers/{customer_id}/tenants/{tenant_id}
- Update tenant details
- Request: Updated tenant data
- Response: Updated tenant

#### DELETE /api/v1/customers/{customer_id}/tenants/{tenant_id}
- Delete tenant
- Response: Success confirmation

#### GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments
- List all deployments for a tenant
- Query params: product_id, deployment_type, status, pagination
- Response: List of deployments

#### POST /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments
- Create a new deployment for a tenant
- Request: Deployment details (product_id, deployment_type, version, number_of_users, license_info, etc.)
- Response: Created deployment

#### GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}
- Get deployment details
- Response: Deployment information with version details

#### PUT /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}
- Update deployment details (version, number of users, license information, etc.)
- Request: Updated deployment data
- Response: Updated deployment

#### DELETE /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}
- Delete deployment
- Response: Success confirmation

#### GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/{deployment_id}/updates
- Get available updates for a deployment
- Response: List of available versions (newer than installed_version)
- Includes: version number, release date, release type, compatibility status, upgrade path

#### GET /api/v1/customers/{customer_id}/deployments/pending-updates
- Get all pending updates for a customer (aggregated across all tenants)
- Query params: product_id, deployment_type, priority
- Response: List of deployments with pending updates, including update counts and latest versions

#### GET /api/v1/customers/{customer_id}/tenants/{tenant_id}/deployments/pending-updates
- Get pending updates for a tenant (aggregated across all deployments)
- Query params: product_id, deployment_type, priority
- Response: List of deployments with pending updates for the tenant

#### GET /api/v1/updates/pending
- Get all pending updates across all customers (admin view)
- Query params: customer_id, product_id, tenant_id, deployment_type, priority, pagination
- Response: List of all deployments with pending updates across the system
- Includes: customer, tenant, deployment, product, current version, latest version, update count, priority

#### POST /api/v1/customers/{customer_id}/notifications
- Manually trigger notification for customer
- Request: Notification details
- Response: Notification sent confirmation

---

## 8. User Stories

### 8.1 Admin User Stories

**US-1**: As an admin, I want to see all VMs with pending updates at a glance, so I can prioritize update activities.

**US-2**: As an admin, I want to click on a green dot to see release notes and update details, so I can make informed decisions about when to update.

**US-3**: As an admin, I want to select multiple VMs and update them simultaneously, so I can efficiently manage large deployments.

**US-4**: As an admin, I want to receive notifications when new versions are released, so I can stay informed about available updates.

**US-5**: As an admin, I want to query the lifecycle API to check upgrade paths, so I can integrate with my automation systems.

**US-6**: As an admin, I want to see update history and audit logs, so I can track compliance and troubleshoot issues.

**US-7**: As an admin, I want to schedule updates for off-peak hours, so I can minimize disruption to users.

**US-8**: As an admin, I want to rollback failed updates, so I can quickly restore system stability.

**US-9**: As an admin, I want to manage customers and their installations, so I can track which customers have which product versions.

**US-10**: As an admin, I want to see all tenants for a customer, so I can understand their deployment landscape.

**US-11**: As an admin, I want to see all deployments for a tenant, so I can manage product installations.

**US-12**: As an admin, I want to distinguish between UAT and Production deployments, so I can manage updates appropriately.

**US-13**: As an admin, I want customers to receive automatic notifications when new versions are released, so they stay informed about available updates.

**US-14**: As an admin, I want to view number of users and license information for customer deployments, so I can track usage and compliance.

### 8.2 Customer User Stories

**US-15**: As a customer, I want to create and manage tenants, so I can organize my independent deployments.

**US-16**: As a customer, I want to register my deployments (product + type) in the portal, so I can track my product versions.

**US-17**: As a customer, I want to specify whether my deployment is UAT or Production, so I receive appropriate notifications.

**US-18**: As a customer, I want to update version information for my deployments, so the system has accurate version tracking.

**US-19**: As a customer, I want to receive notifications when new versions are available for my deployed products, so I can plan updates.

**US-20**: As a customer, I want to see all my tenants and deployments in one place, so I have visibility into my deployment status.

**US-21**: As a customer, I want to specify the number of users for each deployment, so the system can track usage and licensing.

**US-22**: As a customer, I want to add license information for each deployment, so I can maintain license records in the portal.

### 8.3 System User Stories

**US-21**: As the system, I want to automatically detect new product versions, so admins and customers are notified promptly.

**US-22**: As the system, I want to validate compatibility before allowing updates, so I can prevent system conflicts.

**US-23**: As the system, I want to enforce version policies for multi-tenant environments, so no tenants remain on unsupported versions.

**US-24**: As the system, I want to automatically generate notifications to customers when new versions are rolled out, so customers are informed about updates for their deployments.

---

## 9. Example Use Cases

### Use Case 1: Single VM Update
1. Admin logs into update portal
2. Sees 5 VMs with green dots indicating pending upgrades
3. Clicks on one VM's indicator
4. Views detailed release notes and compatibility information
5. Clicks "Update Now"
6. System initiates update and shows progress
7. Update completes successfully, green dot disappears

### Use Case 2: Batch Update
1. Admin selects three VMs with pending updates
2. Initiates batch update to latest version
3. System processes updates in parallel
4. Dashboard shows progress for each VM
5. Two updates succeed, one fails
6. Admin reviews failure details and retries failed update

### Use Case 3: Notification and Proactive Update
1. System detects new major version release
2. Admin receives notification with version details and release notes summary
3. Admin clicks notification link to view full details
4. Admin reviews affected endpoints
5. Admin schedules batch update for maintenance window

### Use Case 4: Lifecycle API Integration
1. Admin's automation system queries lifecycle API
2. API returns upgrade path for specific product and version
3. Automation system validates compatibility
4. Automation system initiates update via portal API
5. System tracks update in audit logs

### Use Case 5: Multi-Tenant Blue-Green Upgrade
1. System identifies tenant edge deployment requiring update
2. System deploys parallel "Blue" environment with new version
3. System validates new environment
4. System switches tenant traffic to Blue environment
5. System monitors for issues
6. If issues detected, system automatically rolls back to Green environment

### Use Case 6: Customer Tenant and Deployment Management
1. Admin creates a new customer in the portal
2. Customer creates a tenant "US Data Center"
3. Customer adds deployments for the tenant:
   - Product: HyWorks, Type: Production, Version: 2.5.0, Users: 500, License: "Enterprise License - 500 users"
   - Product: HySecure, Type: Production, Version: 3.1.0, Users: 500
4. Customer creates another tenant "EU Data Center"
5. Customer adds deployments for the EU tenant:
   - Product: HyWorks, Type: Production, Version: 2.4.0, Users: 300
   - Product: HyWorks, Type: UAT, Version: 2.4.0, Users: 50, License: "Test License"
6. System tracks all deployments across tenants
7. When new version 2.6.0 is released for HyWorks, system generates notifications for all affected deployments
8. Customer receives notifications indicating:
   - US Data Center: Production deployment (2.5.0) can upgrade to 2.6.0
   - EU Data Center: Production deployment (2.4.0) can upgrade to 2.6.0
   - EU Data Center: UAT deployment (2.4.0) can upgrade to 2.6.0
9. Customer views deployment details including tenant, user count, and license information
10. Customer updates deployment version information after upgrading

### Use Case 7: Customer Notification on Version Release
1. New version 3.0.0 is released for HySecure product
2. System identifies all customers with HySecure deployments
3. System generates notifications for each affected customer
4. Notification includes:
   - New version details
   - List of customer's deployments (with tenant and type) that can be updated
   - Current version vs. new version for each deployment
   - Release notes summary
5. Customer receives notification and can view full details
6. Customer updates deployment version information after upgrading

### Use Case 8: Pending Updates Tracking and Visibility
1. Customer has multiple deployments:
   - Tenant "US Data Center": HyWorks Production (v2.5.0), HySecure Production (v3.1.0)
   - Tenant "EU Data Center": HyWorks Production (v2.4.0), HyWorks UAT (v2.4.0)
2. System releases new versions:
   - HyWorks v2.6.0 (released)
   - HyWorks v2.7.0 (released)
   - HySecure v3.2.0 (released)
3. System automatically calculates pending updates for each deployment:
   - US Data Center - HyWorks: 2 pending updates (v2.6.0, v2.7.0)
   - EU Data Center - HyWorks Production: 2 pending updates (v2.6.0, v2.7.0)
   - EU Data Center - HyWorks UAT: 2 pending updates (v2.6.0, v2.7.0)
   - US Data Center - HySecure: 1 pending update (v3.2.0)
4. Customer views customer dashboard and sees:
   - Total pending updates: 7 (aggregated across all deployments)
   - Deployments with pending updates: 4
   - List of deployments requiring updates with priority indicators
5. Customer navigates to Updates page and sees:
   - All deployments with pending updates across all customers
   - Filtered view showing only their deployments
   - Current version vs. latest version for each deployment
   - Update path information (direct upgrade or intermediate steps)
6. Customer views specific deployment details:
   - Sees pending updates section with list of available versions
   - Views version comparison (current v2.5.0 vs. latest v2.7.0)
   - Sees intermediate version (v2.6.0) if direct upgrade not possible
   - Reviews compatibility information for each pending update
7. Admin views Updates page (admin view):
   - Sees all pending updates across all customers
   - Filters by product, customer, tenant, or deployment type
   - Groups by customer or product for better visibility
   - Sorts by priority to identify critical updates first
8. When customer updates a deployment version:
   - System recalculates pending updates for that deployment
   - Updates pending counts in customer, tenant, and deployment views
   - Updates Updates page immediately

---

## 10. Technical Considerations

### 10.1 Communication Protocol
- RESTful API for portal-agent communication
- WebSocket or Server-Sent Events for real-time updates
- Message queue for async operations

### 10.2 Authentication & Authorization
- Mutual TLS for agent-portal communication
- OAuth 2.0 / JWT for portal access
- API keys for programmatic access

### 10.3 Update Package Distribution
- CDN or distributed storage for update packages
- Support for delta updates to reduce bandwidth
- Checksum validation for package integrity

### 10.4 Database Design
- Relational database for structured data (products, versions, endpoints)
- Time-series database for metrics and logs
- Caching layer for frequently accessed data

### 10.5 Error Handling
- Comprehensive error codes and messages
- Retry logic with exponential backoff
- Dead letter queue for failed operations

---

## 11. Success Criteria

1. ✅ Portal displays all endpoints with accurate version status
2. ✅ Green dot indicators correctly identify pending updates
3. ✅ Batch updates successfully process multiple endpoints
4. ✅ Update agent reliably detects and reports versions
5. ✅ Lifecycle API returns accurate upgrade paths
6. ✅ Notifications are delivered promptly upon new releases
7. ✅ Multi-tenant blue-green upgrades complete with zero downtime
8. ✅ Version enforcement prevents EOL versions in production
9. ✅ Audit logs capture all update operations
10. ✅ System handles failures gracefully with rollback support
11. ✅ Customer management portal supports CRUD operations for customers and installations
12. ✅ System automatically generates notifications to customers when new versions are released
13. ✅ Installation version tracking accurately reflects customer-reported versions
14. ✅ UAT and Production installations are properly distinguished and managed
15. ✅ Customer dashboard provides clear visibility into installation status and updates

---

## 12. Out of Scope (Future Enhancements)

- Automatic update scheduling based on policies
- A/B testing of new versions
- Update preview/sandbox environment
- Integration with third-party monitoring tools
- Mobile app for update management
- Update analytics and reporting dashboards
- Custom update scripts and hooks

---

## Document Version
- **Version**: 1.2
- **Date**: 17-Nov-2025
- **Status**: Active Development
- **Last Updated**: 17-Nov-2025 - Added Pending Updates Tracking requirements (Section 2.6.10)

