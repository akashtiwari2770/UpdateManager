// MongoDB Database Setup Script
// Creates database and collections for Update Manager
// This script can be run manually if needed

// Switch to updatemanager database
db = db.getSiblingDB('updatemanager');

// Create collections (MongoDB creates collections automatically, but we can pre-create them)
db.createCollection("products");
db.createCollection("versions");
db.createCollection("compatibility_matrices");
db.createCollection("upgrade_paths");
db.createCollection("notifications");
db.createCollection("update_detections");
db.createCollection("update_rollouts");
db.createCollection("audit_logs");

print("Database 'updatemanager' setup complete!");
print("Collections created successfully.");

