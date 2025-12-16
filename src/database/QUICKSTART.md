# MongoDB Quick Start Guide

## One-Command Setup

From the project root:

```bash
make db-start
```

That's it! MongoDB will be ready in a few seconds.

## What Gets Created Automatically

When you run `make db-start` for the first time:

✅ MongoDB container starts  
✅ Database `updatemanager` is created  
✅ All 8 collections are created  
✅ Application user `updatemanager` is created  
✅ All indexes are created  
✅ Mongo Express UI starts  

## Verify It's Working

```bash
# Check container status
make db-status

# View logs
make db-logs
```

## Access Points

**MongoDB:**
- Host: `localhost:27017`
- Connection: `mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin`

**Mongo Express (Web UI):**
- URL: http://localhost:8081
- Username: `admin`
- Password: `admin123`

## Common Commands

```bash
make db-start      # Start MongoDB
make db-stop       # Stop MongoDB
make db-status     # Check if running
make db-logs       # View logs
make db-indexes    # Recreate indexes (if needed)
```

## Stop MongoDB

```bash
make db-stop
```

## Reset Everything

If you need to start fresh:

```bash
make db-stop
docker volume rm src_database_mongodb_data src_database_mongodb_config
make db-start
```

## Next Steps

After MongoDB is running, you can:

1. Start building your Go backend
2. Connect using the connection string above
3. Use Mongo Express to browse data
4. Run `make db-indexes` if you need to recreate indexes

