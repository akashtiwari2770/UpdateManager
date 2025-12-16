#!/bin/bash

# Database cleanup script for Update Manager
# This script cleans up test data, invalid products, and load test artifacts

set -e

# Database connection details
DB_NAME="${DB_NAME:-updatemanager}"
CONTAINER_NAME="${CONTAINER_NAME:-updatemanager-mongodb}"

echo "========================================="
echo "Database Cleanup Script"
echo "========================================="
echo "Database: $DB_NAME"
echo "Container: $CONTAINER_NAME"
echo ""

# Check if Docker is available
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH"
    exit 1
fi

# Determine Docker command (with or without sudo)
HAS_DOCKER_ACCESS=$(test -r /var/run/docker.sock 2>/dev/null && echo "yes" || echo "no")
DOCKER_PATH=$(which docker)
if [ "$HAS_DOCKER_ACCESS" = "yes" ]; then
    DOCKER_CMD="docker"
else
    if [ -n "$DOCKER_PATH" ]; then
        DOCKER_CMD="sudo $DOCKER_PATH"
    else
        DOCKER_CMD="sudo docker"
    fi
fi

# Check if container is running
if ! $DOCKER_CMD ps --format "{{.Names}}" 2>/dev/null | grep -q "^${CONTAINER_NAME}$"; then
    echo "Error: MongoDB container '$CONTAINER_NAME' is not running."
    echo "Please start it with: make db-start"
    exit 1
fi

# Function to execute MongoDB command
mongo_exec() {
    $DOCKER_CMD exec -i "$CONTAINER_NAME" mongosh -u admin -p admin123 --authenticationDatabase admin "$DB_NAME" --quiet --eval "$1"
}

# Ask for confirmation
read -p "This will delete all data from the database. Are you sure? (yes/no): " confirm
if [ "$confirm" != "yes" ]; then
    echo "Cleanup cancelled."
    exit 0
fi

echo ""
echo "Starting cleanup..."

# Option 1: Full cleanup - drop all collections
read -p "Do you want to drop ALL collections? (yes/no): " drop_all
if [ "$drop_all" = "yes" ]; then
    echo "Dropping all collections..."
    mongo_exec "
    db.products.drop();
    db.versions.drop();
    db.compatibility_matrices.drop();
    db.upgrade_paths.drop();
    db.notifications.drop();
    db.update_detections.drop();
    db.update_rollouts.drop();
    db.audit_logs.drop();
    print('All collections dropped successfully.');
    "
    echo "✓ All collections dropped."
else
    # Option 2: Selective cleanup
    echo ""
    echo "Selective cleanup options:"
    echo "1. Delete products with empty product_id or name"
    echo "2. Delete load test products (containing 'undefined' or 'Load Test')"
    echo "3. Delete test products (starting with 'test-')"
    echo "4. Delete all products"
    echo "5. Delete all versions"
    echo "6. Delete all data (products, versions, etc.)"
    echo ""
    read -p "Enter option (1-6) or 'all' for all options: " option
    
    case $option in
        1)
            echo "Deleting products with empty product_id or name..."
            count=$(mongo_exec "db.products.countDocuments({\$or: [{product_id: ''}, {product_id: null}, {name: ''}, {name: null}]})")
            mongo_exec "db.products.deleteMany({\$or: [{product_id: ''}, {product_id: null}, {name: ''}, {name: null}]})"
            echo "✓ Deleted $count products with empty fields."
            ;;
        2)
            echo "Deleting load test products..."
            count=$(mongo_exec "db.products.countDocuments({\$or: [{product_id: /undefined/}, {name: /Load Test/}, {name: /undefined/}]})")
            mongo_exec "db.products.deleteMany({\$or: [{product_id: /undefined/}, {name: /Load Test/}, {name: /undefined/}]})"
            echo "✓ Deleted $count load test products."
            ;;
        3)
            echo "Deleting test products..."
            count=$(mongo_exec "db.products.countDocuments({product_id: /^test-/})")
            mongo_exec "db.products.deleteMany({product_id: /^test-/})"
            echo "✓ Deleted $count test products."
            ;;
        4)
            echo "Deleting all products..."
            count=$(mongo_exec "db.products.countDocuments({})")
            mongo_exec "db.products.deleteMany({})"
            echo "✓ Deleted $count products."
            ;;
        5)
            echo "Deleting all versions..."
            count=$(mongo_exec "db.versions.countDocuments({})")
            mongo_exec "db.versions.deleteMany({})"
            echo "✓ Deleted $count versions."
            ;;
        6|all)
            echo "Deleting all data..."
            products=$(mongo_exec "db.products.countDocuments({})")
            versions=$(mongo_exec "db.versions.countDocuments({})")
            compat=$(mongo_exec "db.compatibility_matrices.countDocuments({})")
            paths=$(mongo_exec "db.upgrade_paths.countDocuments({})")
            notifications=$(mongo_exec "db.notifications.countDocuments({})")
            detections=$(mongo_exec "db.update_detections.countDocuments({})")
            rollouts=$(mongo_exec "db.update_rollouts.countDocuments({})")
            audit=$(mongo_exec "db.audit_logs.countDocuments({})")
            
            mongo_exec "
            db.products.deleteMany({});
            db.versions.deleteMany({});
            db.compatibility_matrices.deleteMany({});
            db.upgrade_paths.deleteMany({});
            db.notifications.deleteMany({});
            db.update_detections.deleteMany({});
            db.update_rollouts.deleteMany({});
            db.audit_logs.deleteMany({});
            "
            echo "✓ Deleted all data:"
            echo "  - $products products"
            echo "  - $versions versions"
            echo "  - $compat compatibility matrices"
            echo "  - $paths upgrade paths"
            echo "  - $notifications notifications"
            echo "  - $detections update detections"
            echo "  - $rollouts update rollouts"
            echo "  - $audit audit logs"
            ;;
        *)
            echo "Invalid option. Cleanup cancelled."
            exit 1
            ;;
    esac
fi

echo ""
echo "Cleanup completed successfully!"
echo ""
echo "Current database stats:"
mongo_exec "
print('Products: ' + db.products.countDocuments({}));
print('Versions: ' + db.versions.countDocuments({}));
print('Compatibility Matrices: ' + db.compatibility_matrices.countDocuments({}));
print('Upgrade Paths: ' + db.upgrade_paths.countDocuments({}));
print('Notifications: ' + db.notifications.countDocuments({}));
print('Update Detections: ' + db.update_detections.countDocuments({}));
print('Update Rollouts: ' + db.update_rollouts.countDocuments({}));
print('Audit Logs: ' + db.audit_logs.countDocuments({}));
"

