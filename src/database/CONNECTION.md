# MongoDB Connection Information

## Docker Setup (Recommended)

When using Docker, MongoDB is automatically configured with:

### Connection Details

**Root User (Admin):**
- Username: `admin`
- Password: `admin123`
- Database: `admin` (for authentication)

**Application User:**
- Username: `updatemanager`
- Password: `updatemanager123`
- Database: `updatemanager`

### Connection Strings

**For Go Backend (.env file):**
```env
MONGODB_URI=mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin
```

**Or using application user:**
```env
MONGODB_URI=mongodb://updatemanager:updatemanager123@localhost:27017/updatemanager?authSource=admin
```

**For MongoDB Shell:**
```bash
mongosh "mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin"
```

### Web UI

**Mongo Express:**
- URL: http://localhost:8081
- Username: `admin`
- Password: `admin123`

## Quick Start

1. Start MongoDB:
   ```bash
   make db-start
   ```

2. Check status:
   ```bash
   make db-status
   ```

3. View logs:
   ```bash
   make db-logs
   ```

4. Stop MongoDB:
   ```bash
   make db-stop
   ```

## Manual Connection

If you need to connect manually:

```bash
# Using mongosh
docker exec -it updatemanager-mongodb mongosh -u admin -p admin123 --authenticationDatabase admin

# Or from host
mongosh "mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin"
```

## Notes

- Database and indexes are created automatically on first container start
- Data persists in Docker volume `mongodb_data`
- To reset database, stop container and remove volume:
  ```bash
  make db-stop
  docker volume rm src_database_mongodb_data
  make db-start
  ```

