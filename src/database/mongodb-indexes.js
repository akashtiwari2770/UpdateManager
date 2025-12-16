// MongoDB Indexes for Product Release Management System
// Run this script to create all necessary indexes

// Products Collection
db.products.createIndex({ "product_id": 1 }, { unique: true });
db.products.createIndex({ "type": 1 });
db.products.createIndex({ "is_active": 1 });
db.products.createIndex({ "created_at": -1 });

// Versions Collection
db.versions.createIndex({ "product_id": 1, "version_number": 1 }, { unique: true });
db.versions.createIndex({ "product_id": 1, "state": 1 });
db.versions.createIndex({ "product_id": 1, "release_type": 1 });
db.versions.createIndex({ "state": 1 });
db.versions.createIndex({ "release_date": -1 });
db.versions.createIndex({ "created_at": -1 });
db.versions.createIndex({ "approved_by": 1 });
db.versions.createIndex({ "eol_date": 1 }, { sparse: true });

// Packages Collection (embedded in versions, but if separate collection)
// db.packages.createIndex({ "version_id": 1, "package_type": 1 });
// db.packages.createIndex({ "checksum_sha256": 1 }, { unique: true });

// Compatibility Matrix Collection
db.compatibility_matrices.createIndex({ "product_id": 1, "version_number": 1 }, { unique: true });
db.compatibility_matrices.createIndex({ "product_id": 1 });
db.compatibility_matrices.createIndex({ "validation_status": 1 });
db.compatibility_matrices.createIndex({ "validated_at": -1 });

// Upgrade Paths Collection
db.upgrade_paths.createIndex({ "product_id": 1, "from_version": 1, "to_version": 1 }, { unique: true });
db.upgrade_paths.createIndex({ "product_id": 1 });
db.upgrade_paths.createIndex({ "from_version": 1 });
db.upgrade_paths.createIndex({ "to_version": 1 });
db.upgrade_paths.createIndex({ "path_type": 1 });
db.upgrade_paths.createIndex({ "is_blocked": 1 });

// Notifications Collection
db.notifications.createIndex({ "recipient_id": 1, "is_read": 1, "created_at": -1 });
db.notifications.createIndex({ "recipient_id": 1, "type": 1 });
db.notifications.createIndex({ "recipient_id": 1, "priority": 1 });
db.notifications.createIndex({ "product_id": 1 });
db.notifications.createIndex({ "version_id": 1 });
db.notifications.createIndex({ "created_at": -1 });
db.notifications.createIndex({ "is_read": 1 });

// Update Detections Collection
db.update_detections.createIndex({ "endpoint_id": 1, "product_id": 1 }, { unique: true });
db.update_detections.createIndex({ "endpoint_id": 1 });
db.update_detections.createIndex({ "product_id": 1 });
db.update_detections.createIndex({ "detected_at": -1 });
db.update_detections.createIndex({ "last_checked_at": -1 });

// Update Rollouts Collection
db.update_rollouts.createIndex({ "endpoint_id": 1, "product_id": 1, "status": 1 });
db.update_rollouts.createIndex({ "endpoint_id": 1 });
db.update_rollouts.createIndex({ "product_id": 1 });
db.update_rollouts.createIndex({ "status": 1 });
db.update_rollouts.createIndex({ "initiated_at": -1 });
db.update_rollouts.createIndex({ "initiated_by": 1 });
db.update_rollouts.createIndex({ "from_version": 1, "to_version": 1 });

// Audit Logs Collection
db.audit_logs.createIndex({ "resource_type": 1, "resource_id": 1 });
db.audit_logs.createIndex({ "user_id": 1 });
db.audit_logs.createIndex({ "action": 1 });
db.audit_logs.createIndex({ "timestamp": -1 });
db.audit_logs.createIndex({ "product_id": 1 }, { sparse: true });
db.audit_logs.createIndex({ "created_at": -1 }); // TTL index for old logs (optional)
// db.audit_logs.createIndex({ "created_at": 1 }, { expireAfterSeconds: 31536000 }); // 1 year TTL

// Compound indexes for common queries

// Versions: Find latest released version for a product
db.versions.createIndex({ "product_id": 1, "state": 1, "release_date": -1 });

// Versions: Find versions pending approval
db.versions.createIndex({ "state": 1, "created_at": 1 });

// Notifications: Unread notifications for a user
db.notifications.createIndex({ "recipient_id": 1, "is_read": 1, "created_at": -1 });

// Update Rollouts: Active rollouts for an endpoint
db.update_rollouts.createIndex({ "endpoint_id": 1, "status": 1, "initiated_at": -1 });

// Audit Logs: Recent actions by user
db.audit_logs.createIndex({ "user_id": 1, "timestamp": -1 });

// Audit Logs: Actions on a specific resource
db.audit_logs.createIndex({ "resource_type": 1, "resource_id": 1, "timestamp": -1 });

// Customer Management Indexes

// Customers Collection
db.customers.createIndex({ "customer_id": 1 }, { unique: true });
db.customers.createIndex({ "email": 1 });
db.customers.createIndex({ "account_status": 1 });
db.customers.createIndex({ "created_at": -1 });

// Customer Tenants Collection
db.customer_tenants.createIndex({ "tenant_id": 1 }, { unique: true });
db.customer_tenants.createIndex({ "customer_id": 1 });
db.customer_tenants.createIndex({ "status": 1 });
db.customer_tenants.createIndex({ "customer_id": 1, "status": 1 });
db.customer_tenants.createIndex({ "created_at": -1 });

// Deployments Collection
db.deployments.createIndex({ "deployment_id": 1 }, { unique: true });
db.deployments.createIndex({ "tenant_id": 1 });
db.deployments.createIndex({ "product_id": 1 });
db.deployments.createIndex({ "deployment_type": 1 });
db.deployments.createIndex({ "status": 1 });
db.deployments.createIndex({ "tenant_id": 1, "product_id": 1 });
db.deployments.createIndex({ "tenant_id": 1, "deployment_type": 1 });
db.deployments.createIndex({ "product_id": 1, "deployment_type": 1 });
db.deployments.createIndex({ "deployment_date": -1 });
db.deployments.createIndex({ "last_updated_date": -1 });

// Compound indexes for common queries

// Deployments: Find deployments for a tenant by product and type
db.deployments.createIndex({ "tenant_id": 1, "product_id": 1, "deployment_type": 1 }, { unique: true });

// Deployments: Find all deployments for a product (for notifications)
db.deployments.createIndex({ "product_id": 1, "status": 1, "deployment_type": 1 });

// Notifications: Customer-specific notifications
db.notifications.createIndex({ "customer_id": 1, "is_read": 1, "created_at": -1 });
db.notifications.createIndex({ "tenant_id": 1 });
db.notifications.createIndex({ "deployment_id": 1 });

// License Management Indexes

// Subscriptions Collection
db.subscriptions.createIndex({ "subscription_id": 1 }, { unique: true });
db.subscriptions.createIndex({ "customer_id": 1 });
db.subscriptions.createIndex({ "status": 1 });
db.subscriptions.createIndex({ "start_date": 1 });
db.subscriptions.createIndex({ "end_date": 1 }, { sparse: true });
db.subscriptions.createIndex({ "customer_id": 1, "status": 1 });
db.subscriptions.createIndex({ "end_date": 1 }, { sparse: true }); // For expiration queries

// Licenses Collection
db.licenses.createIndex({ "license_id": 1 }, { unique: true });
db.licenses.createIndex({ "subscription_id": 1 });
db.licenses.createIndex({ "product_id": 1 });
db.licenses.createIndex({ "license_type": 1 });
db.licenses.createIndex({ "status": 1 });
db.licenses.createIndex({ "end_date": 1 }, { sparse: true }); // For expiration queries
db.licenses.createIndex({ "subscription_id": 1, "product_id": 1 }); // Compound index
db.licenses.createIndex({ "subscription_id": 1, "status": 1 });
db.licenses.createIndex({ "product_id": 1, "status": 1 });

// License Allocations Collection
db.license_allocations.createIndex({ "allocation_id": 1 }, { unique: true });
db.license_allocations.createIndex({ "license_id": 1 });
db.license_allocations.createIndex({ "tenant_id": 1 }, { sparse: true });
db.license_allocations.createIndex({ "deployment_id": 1 }, { sparse: true });
db.license_allocations.createIndex({ "status": 1 });
db.license_allocations.createIndex({ "license_id": 1, "status": 1 }); // Compound index
db.license_allocations.createIndex({ "tenant_id": 1, "status": 1 }, { sparse: true });
db.license_allocations.createIndex({ "deployment_id": 1, "status": 1 }, { sparse: true });
db.license_allocations.createIndex({ "allocation_date": -1 });

