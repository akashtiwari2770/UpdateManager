// MongoDB Indexes Creation Script
// This script runs automatically when MongoDB container starts for the first time

// Switch to updatemanager database
db = db.getSiblingDB('updatemanager');

print("Creating indexes for updatemanager database...");

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
db.versions.createIndex({ "product_id": 1, "state": 1, "release_date": -1 });
db.versions.createIndex({ "state": 1, "created_at": 1 });

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
db.update_rollouts.createIndex({ "endpoint_id": 1, "status": 1, "initiated_at": -1 });

// Audit Logs Collection
db.audit_logs.createIndex({ "resource_type": 1, "resource_id": 1 });
db.audit_logs.createIndex({ "user_id": 1 });
db.audit_logs.createIndex({ "action": 1 });
db.audit_logs.createIndex({ "timestamp": -1 });
db.audit_logs.createIndex({ "product_id": 1 }, { sparse: true });
db.audit_logs.createIndex({ "created_at": -1 });
db.audit_logs.createIndex({ "user_id": 1, "timestamp": -1 });
db.audit_logs.createIndex({ "resource_type": 1, "resource_id": 1, "timestamp": -1 });

print("All indexes created successfully!");

