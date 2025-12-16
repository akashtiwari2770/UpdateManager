# Release Management

## 1. Overview

This document defines the release management interface and operations for the Accops Update Manager. It provides release managers with the capabilities to manage products and their releases through the management portal. This document focuses on the administrative functions for product lifecycle management and release versioning.

## 2. Release Manager Role

### 2.1 Responsibilities

A **Release Manager** is responsible for:
- Managing product definitions and metadata
- Creating and managing product releases
- Maintaining version information and compatibility matrices
- Publishing release notes
- Managing product lifecycle (EOL dates, deprecation)
- Approving releases for distribution

### 2.2 Access Control

- **FR-RM-1.1**: Release Manager role requires elevated permissions
- **FR-RM-1.2**: All product and release management operations are logged for audit
- **FR-RM-1.3**: Release Manager actions require authentication and authorization

---

## 3. Product Management

### 3.1 Product Overview

Products represent the Accops software components that can be updated through the Update Manager system. Each product has metadata, configuration, and associated releases.

### 3.2 View Products

#### Product List View
- **FR-RM-2.1**: Display all products in a table/grid view with:
  - Product ID
  - Product Name (e.g., "HyWorks", "HySecure", "IRIS", "ARS", "Client for Linux", "Client for Windows", "Client for Mobile")
  - Product Type (Server, Client)
  - Platform (Windows, Linux, Mobile, Cross-Platform)
  - Current Latest Version
  - Total Releases Count
  - Support Status (Active, Deprecated, EOL)
  - Last Updated Date

#### Product Details View
- **FR-RM-2.2**: Clicking a product displays detailed information:
  - Full product metadata
  - Update strategy configuration
  - Multi-tenant support settings
  - Default update policy
  - Associated releases list
  - Compatibility requirements
  - EOL information

#### Product Search and Filtering
- **FR-RM-2.3**: Support filtering products by:
  - Product Type (Server, Client)
  - Platform
  - Support Status
  - Search by product name or ID

### 3.3 Add Product

#### Product Creation Form
- **FR-RM-3.1**: Release Manager can create a new product with:
  - **Product Name** (required, unique, free text - e.g., "HyWorks", "HySecure", "IRIS", "ARS", "Client for Linux", "Client for Windows", "Client for Mobile")
  - **Product ID** (auto-generated or manual, unique identifier)
  - **Product Type** (dropdown: Server or Client)
  - **Platform** (Windows/Linux/Mobile/Cross-Platform)
  - **Update Strategy** (Blue-Green, In-Place, App-Store)
  - **Multi-Tenant Support** (Yes/No)
  - **Default Update Policy** (Auto-Update, Manual Approval, Scheduled)
  - **Description** (optional)
  - **Initial EOL Date** (optional)

#### Product Validation
- **FR-RM-3.2**: System validates:
  - Product name uniqueness
  - Product ID uniqueness
  - Required fields are provided
  - Valid product type (Server or Client)

#### Product Creation Workflow
1. Release Manager navigates to Products page
2. Clicks "Add Product" button
3. Fills in product creation form
4. System validates input
5. Release Manager submits form
6. System creates product record
7. Product appears in product list
8. Release Manager can now add releases for this product

### 3.4 Update Product

#### Product Edit Interface
- **FR-RM-4.1**: Release Manager can update product information:
  - Product Name (with uniqueness check)
  - Product Type (Server or Client, with validation)
  - Platform
  - Update Strategy
  - Multi-Tenant Support
  - Default Update Policy
  - Description
  - EOL Date
  - Support Status

#### Update Restrictions
- **FR-RM-4.2**: System enforces restrictions:
  - Product ID cannot be changed (immutable)
  - Cannot change product type if releases exist (with override option)
  - Cannot set EOL date in the past
  - Cannot deprecate product if active releases exist (with confirmation)

#### Product Update Workflow
1. Release Manager selects product from list
2. Clicks "Edit Product" button
3. System loads current product data into edit form
4. Release Manager modifies fields
5. System validates changes
6. Release Manager saves changes
7. System updates product record
8. Changes are reflected in product list and details

### 3.5 Delete Product

#### Product Deletion Requirements
- **FR-RM-5.1**: Product deletion requires:
  - No active releases associated with the product
  - No endpoints with installed instances of the product
  - Confirmation from Release Manager
  - Reason for deletion (for audit)

#### Soft Delete Option
- **FR-RM-5.2**: System supports soft delete:
  - Mark product as "Deleted" instead of physical removal
  - Preserve historical data and audit logs
  - Hide from active product lists
  - Maintain referential integrity

#### Product Deletion Workflow
1. Release Manager selects product
2. Clicks "Delete Product" button
3. System checks for dependencies:
   - Active releases
   - Installed instances
   - Associated endpoints
4. If dependencies exist:
   - System displays warning with dependency list
   - Release Manager must resolve dependencies first
5. If no dependencies:
   - System prompts for confirmation
   - Release Manager provides deletion reason
   - Release Manager confirms deletion
6. System performs soft delete (or hard delete if configured)
7. Product is removed from active product list
8. Audit log records deletion

---

## 4. Release Management

### 4.1 Release Overview

Releases represent specific versions of a product that can be deployed to endpoints. Each release contains version information, release notes, packages, and compatibility data.

### 4.2 View Releases for a Product

#### Product Releases List
- **FR-RM-6.1**: Display all releases for a selected product with:
  - Version Number (semantic versioning)
  - Release Date
  - Release Type (Security, Feature, Maintenance, Major)
  - Release Status (Draft, Pending Review, Approved, Released, Deprecated, EOL)
  - Download URL
  - Checksum
  - Release Notes Summary
  - Adoption Count (number of endpoints on this version)
  - Last Modified Date

#### Release Filtering and Sorting
- **FR-RM-6.2**: Support filtering releases by:
  - Release Status
  - Release Type
  - Date Range
  - Version Number (search)
  - Sort by: Version (descending), Release Date, Status

#### Release Details View
- **FR-RM-6.3**: Clicking a release displays:
  - Full version metadata
  - Complete release notes
  - Compatibility matrix
  - Upgrade paths
  - Package information
  - Download statistics
  - Adoption metrics
  - Related releases

### 4.3 Add Release

#### Release Creation Form
- **FR-RM-7.1**: Release Manager can create a new release with:
  - **Product** (pre-selected from product context, required)
  - **Version Number** (required, semantic versioning format)
  - **Release Date** (required, default: today)
  - **Release Type** (Security, Feature, Maintenance, Major)
  - **Release Notes** (required, markdown supported)
  - **Package Upload** (required):
    - Full installer file
    - Update package (optional)
    - Delta package (optional)
  - **Checksum** (auto-calculated or manual entry)
  - **Download URL** (auto-generated or manual)
  - **Compatibility Matrix**:
    - Minimum Server Version (for clients)
    - Maximum Server Version (for clients)
    - Compatible Product Versions
    - Incompatible Versions
  - **EOL Date** (optional)
  - **Prerequisites** (optional)
  - **Breaking Changes** (optional)
  - **Known Issues** (optional)

#### Version Number Validation
- **FR-RM-7.2**: System validates version number:
  - Follows semantic versioning (MAJOR.MINOR.PATCH)
  - Uniqueness within product
  - Version number is greater than previous releases (with override option)
  - Pre-release versions allowed (alpha, beta, RC)

#### Package Upload and Validation
- **FR-RM-7.3**: System handles package upload:
  - Secure file upload with progress indicator
  - File size limits and validation
  - Automatic checksum calculation (SHA-256)
  - Package metadata extraction
  - Virus scanning (if configured)
  - Storage in encrypted repository

#### Compatibility Matrix Setup
- **FR-RM-7.4**: Release Manager can define:
  - Minimum compatible server version (for client products)
  - Maximum compatible server version (for client products)
  - Compatible product versions (for dependencies)
  - Incompatible versions (explicit exclusions)
  - Upgrade path requirements

#### Release Creation Workflow
1. Release Manager navigates to product details
2. Clicks "Add Release" button
3. System loads release creation form with product pre-selected
4. Release Manager fills in version information
5. Release Manager uploads package files
6. System calculates checksums automatically
7. Release Manager defines compatibility matrix
8. Release Manager writes release notes
9. System validates all inputs
10. Release Manager saves as "Draft"
11. System creates release record with status "Draft"
12. Release appears in releases list with Draft status

### 4.4 Update Release

#### Release Edit Interface
- **FR-RM-8.1**: Release Manager can update release information:
  - Release Notes (full editing)
  - Release Date (with restrictions)
  - Release Type
  - Compatibility Matrix
  - EOL Date
  - Prerequisites
  - Breaking Changes
  - Known Issues
  - Package files (replace with new version)

#### Update Restrictions
- **FR-RM-8.2**: System enforces restrictions:
  - Version Number cannot be changed (immutable)
  - Product cannot be changed (immutable)
  - Cannot modify release if status is "Released" and endpoints are using it (with override)
  - Cannot set EOL date in the past
  - Package replacement requires new checksum

#### Release Status Transitions
- **FR-RM-8.3**: Release Manager can change release status:
  - **Draft** → **Pending Review**: Submit for approval
  - **Pending Review** → **Approved**: Approve for release
  - **Approved** → **Released**: Make available to endpoints
  - **Released** → **Deprecated**: Mark as no longer recommended
  - **Deprecated** → **EOL**: Mark as end of life
  - **Any** → **Draft**: Revert to draft (with restrictions)

#### Release Update Workflow
1. Release Manager selects release from list
2. Clicks "Edit Release" button
3. System loads current release data into edit form
4. Release Manager modifies fields
5. System validates changes based on current status
6. Release Manager can change status (if allowed)
7. Release Manager saves changes
8. System updates release record
9. If status changed to "Released", system:
   - Makes release available to endpoints
   - Triggers notification generation
   - Updates product's "Current Latest Version"

### 4.5 Delete Release

#### Release Deletion Requirements
- **FR-RM-9.1**: Release deletion requires:
  - Release status is "Draft" or "Pending Review" (not approved/released)
  - No endpoints with this version installed
  - Confirmation from Release Manager
  - Reason for deletion (for audit)

#### Soft Delete for Released Versions
- **FR-RM-9.2**: For released versions:
  - Cannot be deleted (immutable for audit trail)
  - Can be marked as "Deprecated" or "EOL"
  - Historical data must be preserved

#### Release Deletion Workflow
1. Release Manager selects release
2. Clicks "Delete Release" button
3. System checks:
   - Release status
   - Installed instances
4. If release is approved/released:
   - System prevents deletion
   - Suggests deprecation instead
5. If release is draft/pending:
   - System prompts for confirmation
   - Release Manager provides deletion reason
   - Release Manager confirms deletion
6. System performs deletion (soft or hard based on configuration)
7. Release is removed from releases list
8. Audit log records deletion

---

## 5. Release Approval Workflow

### 5.1 Approval States

- **FR-RM-10.1**: Releases progress through states:
  1. **Draft**: Being prepared by Release Manager
  2. **Pending Review**: Submitted for approval
  3. **Approved**: Approved for release
  4. **Released**: Available to endpoints
  5. **Deprecated**: No longer recommended
  6. **EOL**: End of Life

### 5.2 Approval Actions

#### Submit for Review
- **FR-RM-10.2**: Release Manager can submit draft release:
  - Validates all required fields are complete
  - Checks package uploads are present
  - Verifies release notes are provided
  - Changes status to "Pending Review"
  - Notifies approvers (if configured)

#### Approve Release
- **FR-RM-10.3**: Release Manager (with approval permissions) can:
  - Review release details
  - Approve release (status → "Approved")
  - Reject release with comments (status → "Draft")
  - Request changes with comments

#### Publish Release
- **FR-RM-10.4**: Release Manager can publish approved release:
  - Changes status to "Released"
  - Makes release available to endpoints
  - Triggers notification system
  - Updates product's latest version
  - Generates upgrade paths

---

## 6. Compatibility Management

### 6.1 Compatibility Matrix

#### Define Compatibility
- **FR-RM-11.1**: Release Manager can define for each release:
  - **Server-Client Compatibility** (for client products):
    - Minimum server version required
    - Maximum server version supported
    - Recommended server version
  - **Product Dependencies**:
    - Required product versions
    - Compatible product versions
    - Incompatible product versions
  - **Operating System Compatibility**:
    - Supported OS versions
    - Minimum OS version
    - Architecture requirements (x86, x64, ARM)

#### Compatibility Validation
- **FR-RM-11.2**: System validates compatibility:
  - Checks for conflicts with existing releases
  - Validates version ranges are logical
  - Warns about potential compatibility issues
  - Suggests compatible version combinations

### 6.2 Upgrade Path Management

#### Define Upgrade Paths
- **FR-RM-11.3**: Release Manager can specify:
  - Direct upgrade allowed (from version X to Y)
  - Multi-step upgrade required (must go through intermediate versions)
  - Blocked upgrades (cannot upgrade directly)
  - Recommended upgrade path

#### Automatic Upgrade Path Calculation
- **FR-RM-11.4**: System can automatically calculate:
  - Shortest upgrade path between versions
  - All possible upgrade paths
  - Blocked paths with reasons
  - Recommended path based on compatibility

---

## 7. Release Notes Management

### 7.1 Release Notes Editor

#### Rich Text Editor
- **FR-RM-12.1**: Release Manager can create release notes with:
  - Markdown support
  - Rich text formatting
  - Section templates (What's New, Bug Fixes, Breaking Changes, etc.)
  - Preview functionality
  - Version history (track changes)

#### Release Notes Sections
- **FR-RM-12.2**: Standard sections include:
  - Version Information
  - What's New (features, enhancements)
  - Bug Fixes
  - Breaking Changes
  - Compatibility Information
  - Upgrade Instructions
  - Known Issues
  - Security Advisories (if applicable)

### 7.2 Release Notes Templates

#### Template Support
- **FR-RM-12.3**: System provides templates:
  - Security Release Template
  - Feature Release Template
  - Maintenance Release Template
  - Major Release Template
  - Custom templates per product

---

## 8. Package Management

### 8.1 Package Upload

#### Upload Interface
- **FR-RM-13.1**: Release Manager can upload:
  - Full installer packages
  - Update packages (incremental)
  - Delta packages (minimal changes)
  - Rollback packages (for reverting)
  - Multiple packages per release (different platforms)

#### Package Validation
- **FR-RM-13.2**: System validates packages:
  - File format validation
  - Size limits
  - Checksum verification
  - Virus scanning (if configured)
  - Package metadata extraction

### 8.2 Package Storage

#### Secure Storage
- **FR-RM-13.3**: Packages are stored:
  - In encrypted storage
  - With version control
  - With access logging
  - In CDN for distribution
  - With geographic distribution

#### Package Metadata
- **FR-RM-13.4**: Each package stores:
  - File name and size
  - Checksum (SHA-256)
  - Upload date and user
  - Download count
  - Platform and architecture
  - Package type

---

## 9. EOL (End of Life) Management

### 9.1 Set EOL Date

#### EOL Configuration
- **FR-RM-14.1**: Release Manager can set:
  - EOL date for a release
  - EOL date for a product
  - Grace period before EOL
  - EOL notification schedule

#### EOL Workflow
- **FR-RM-14.2**: When EOL date is set:
  - System schedules notifications (90 days, 30 days, 7 days before)
  - System marks version as approaching EOL
  - System enforces upgrade policies
  - System blocks new installations after EOL

### 9.2 Deprecation Management

#### Deprecate Release
- **FR-RM-14.3**: Release Manager can deprecate releases:
  - Mark as "Deprecated" status
  - Provide deprecation reason
  - Set deprecation date
  - Recommend alternative version
  - Maintain availability for existing installations

---

## 10. Bulk Operations

### 10.1 Bulk Release Management

#### Multi-Select Operations
- **FR-RM-15.1**: Release Manager can:
  - Select multiple releases
  - Bulk change status (approve, deprecate, etc.)
  - Bulk set EOL dates
  - Bulk export release information
  - Bulk delete (draft releases only)

### 10.2 Import/Export

#### Export Releases
- **FR-RM-15.2**: Release Manager can export:
  - Release data (CSV, JSON)
  - Release notes (markdown, PDF)
  - Compatibility matrices
  - Product catalogs

#### Import Releases
- **FR-RM-15.3**: Release Manager can import:
  - Release data from CSV/JSON
  - Bulk release creation
  - Compatibility matrix updates
  - Product definitions

---

## 11. Search and Filtering

### 11.1 Product Search

- **FR-RM-16.1**: Search products by:
  - Product name
  - Product ID
  - Product type
  - Support status
  - Platform

### 11.2 Release Search

- **FR-RM-16.2**: Search releases by:
  - Version number
  - Release date range
  - Release type
  - Release status
  - Product name
  - Release notes content

---

## 12. Audit and Logging

### 12.1 Operation Logging

- **FR-RM-17.1**: System logs all Release Manager operations:
  - Product creation, update, deletion
  - Release creation, update, deletion
  - Status changes
  - Package uploads
  - Approval actions
  - EOL date changes

### 12.2 Audit Trail

- **FR-RM-17.2**: Audit trail includes:
  - User who performed action
  - Timestamp
  - Action type
  - Before/after values (for updates)
  - IP address
  - Reason (for deletions)

### 12.3 Audit Reports

- **FR-RM-17.3**: Release Manager can view:
  - Activity logs
  - Change history for products/releases
  - Approval history
  - Deletion history

---

## 13. API for Release Management

### 13.1 Product Management API

#### GET /api/v1/release-management/products
- List all products
- Query params: filter, search, pagination

#### GET /api/v1/release-management/products/{product_id}
- Get product details

#### POST /api/v1/release-management/products
- Create new product
- Request: Product metadata
- Response: Created product

#### PUT /api/v1/release-management/products/{product_id}
- Update product
- Request: Updated product metadata
- Response: Updated product

#### DELETE /api/v1/release-management/products/{product_id}
- Delete product
- Request: Deletion reason
- Response: Deletion confirmation

### 13.2 Release Management API

#### GET /api/v1/release-management/products/{product_id}/releases
- List all releases for a product
- Query params: status, type, date_range, pagination

#### GET /api/v1/release-management/releases/{release_id}
- Get release details

#### POST /api/v1/release-management/products/{product_id}/releases
- Create new release
- Request: Release metadata, package files
- Response: Created release

#### PUT /api/v1/release-management/releases/{release_id}
- Update release
- Request: Updated release metadata
- Response: Updated release

#### DELETE /api/v1/release-management/releases/{release_id}
- Delete release
- Request: Deletion reason
- Response: Deletion confirmation

#### POST /api/v1/release-management/releases/{release_id}/approve
- Approve release
- Request: Approval comments
- Response: Updated release status

#### POST /api/v1/release-management/releases/{release_id}/publish
- Publish release (make available)
- Response: Updated release status

#### POST /api/v1/release-management/releases/{release_id}/deprecate
- Deprecate release
- Request: Deprecation reason, date
- Response: Updated release status

#### POST /api/v1/release-management/releases/{release_id}/upload-package
- Upload package file for release
- Request: Package file, package type
- Response: Package metadata

---

## 14. User Interface Requirements

### 14.1 Product Management UI

- **FR-RM-18.1**: Product list page with:
  - Table view of all products
  - Add Product button
  - Search and filter controls
  - Actions column (View, Edit, Delete)
  - Pagination for large lists

- **FR-RM-18.2**: Product details page with:
  - Product information display
  - Edit button
  - Delete button
  - Associated releases section
  - Activity log section

- **FR-RM-18.3**: Product creation/edit form with:
  - All product fields
  - Validation messages
  - Save/Cancel buttons
  - Field help text

### 14.2 Release Management UI

- **FR-RM-18.4**: Release list page with:
  - Table view of releases for a product
  - Add Release button
  - Filter by status, type, date
  - Actions column (View, Edit, Delete, Approve, Publish)
  - Version sorting

- **FR-RM-18.5**: Release details page with:
  - Release information display
  - Release notes viewer
  - Compatibility matrix display
  - Package download links
  - Status change controls
  - Activity log

- **FR-RM-18.6**: Release creation/edit form with:
  - All release fields
  - Package upload interface
  - Release notes editor
  - Compatibility matrix editor
  - Validation messages
  - Save/Cancel buttons

---

## 15. Notifications and Alerts

### 15.1 Release Manager Notifications

- **FR-RM-19.1**: Release Manager receives notifications for:
  - Pending approvals (releases submitted for review)
  - Release approval requests
  - Package upload completion
  - Validation errors
  - EOL date approaching
  - Release adoption milestones

### 15.2 System Alerts

- **FR-RM-19.2**: System alerts Release Manager about:
  - Incomplete release information
  - Missing compatibility data
  - Package upload failures
  - Validation errors
  - Duplicate version numbers

---

## 16. Data Validation and Constraints

### 16.1 Product Validation

- **FR-RM-20.1**: Product validation rules:
  - Product name: Required, unique, 1-200 characters (free text, e.g., "HyWorks", "HySecure", "IRIS", "ARS", "Client for Linux")
  - Product ID: Required, unique, alphanumeric with hyphens
  - Product type: Required, must be "server" or "client"
  - Platform: Required for client products

### 16.2 Release Validation

- **FR-RM-20.2**: Release validation rules:
  - Version number: Required, semantic versioning format, unique per product
  - Release date: Required, cannot be in future (with override)
  - Release notes: Required, minimum 50 characters
  - Package: Required, valid file format, size limits
  - Checksum: Required, valid SHA-256 format
  - Compatibility: Required for client products (server version range)

---

## Document Version
- **Version**: 1.0
- **Date**: 2025
- **Status**: Draft

