// MongoDB Database Setup Script
// This script runs automatically when MongoDB container starts for the first time
// It creates the database, collections, and a dedicated user

// Switch to admin database
db = db.getSiblingDB('admin');

// Authenticate as root user (if needed)
// db.auth('admin', 'admin123');

// Switch to updatemanager database
db = db.getSiblingDB('updatemanager');

// Create collections
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

// Create a dedicated user for the application
db = db.getSiblingDB('admin');
db.createUser({
  user: "updatemanager",
  pwd: "updatemanager123",
  roles: [
    {
      role: "readWrite",
      db: "updatemanager"
    }
  ]
});

print("User 'updatemanager' created successfully!");

