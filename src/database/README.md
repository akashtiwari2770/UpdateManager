# MongoDB Setup

This directory contains MongoDB setup scripts and configuration for Docker-based MongoDB.

## Quick Start

The easiest way to set up MongoDB is using Docker:

```bash
# From project root
make db-start
```

This will:
- Start MongoDB container
- Create database and collections automatically
- Create all indexes automatically
- Start Mongo Express UI

## Files

- `docker-compose.mongodb.yml` - Docker Compose configuration for MongoDB and Mongo Express
- `init-scripts/01-setup-database.js` - Runs automatically on first container start
- `init-scripts/02-create-indexes.js` - Creates indexes automatically on first container start
- `mongodb-indexes.js` - Manual index creation script (for recreating indexes)
- `setup-database.js` - Manual database setup script

## Makefile Commands

From the project root, use these commands:

```bash
make db-start      # Start MongoDB Docker container
make db-stop       # Stop MongoDB Docker container
make db-status     # Check container status
make db-logs       # View container logs
make db-setup      # Manually run database setup (usually automatic)
make db-indexes    # Manually recreate indexes
```

## Connection Information

See [CONNECTION.md](CONNECTION.md) for detailed connection strings and credentials.

**Quick Reference:**
- Connection: `mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin`
- Mongo Express: http://localhost:8081 (admin/admin123)

## Automatic Setup

When you run `make db-start` for the first time:

1. MongoDB container starts
2. Database `updatemanager` is created
3. All collections are created
4. Application user `updatemanager` is created
5. All indexes are created automatically

**Note:** The init scripts in `init-scripts/` run automatically only on the **first** container start. If you need to recreate them, remove the volume and restart.

## Manual Setup (if needed)

If you need to manually run setup scripts:

```bash
# Ensure container is running
make db-start

# Run setup script
docker exec -i updatemanager-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin updatemanager < setup-database.js

# Create indexes
docker exec -i updatemanager-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin updatemanager < mongodb-indexes.js
```

## Collections

The following collections are automatically created:

1. `products` - Product definitions
2. `versions` - Product versions
3. `compatibility_matrices` - Compatibility validation results
4. `upgrade_paths` - Upgrade path definitions
5. `notifications` - User notifications
6. `update_detections` - Update detection records
7. `update_rollouts` - Update rollout records
8. `audit_logs` - Audit log entries

## Indexes

All necessary indexes are created automatically. See `mongodb-indexes.js` for the complete list, including:
- Unique indexes for primary keys
- Compound indexes for common queries
- Sparse indexes for optional fields

## Data Persistence

MongoDB data is stored in a Docker volume `mongodb_data`. To reset the database:

```bash
make db-stop
docker volume rm src_database_mongodb_data
make db-start
```

## Troubleshooting

**Container won't start:**
```bash
make db-logs  # Check logs for errors
```

**Connection refused:**
- Ensure container is running: `make db-status`
- Check if port 27017 is already in use

**Need to reset everything:**
```bash
make db-stop
docker volume rm src_database_mongodb_data src_database_mongodb_config
make db-start
```

