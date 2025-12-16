## Quick Testing Guide

### 1. Backend Unit Tests
```bash
cd src/backend
go test ./internal/utils -v -run TestCompareVersions
go test ./internal/service -v -run TestPendingUpdatesService
```

### 2. API Testing (requires backend running)
```bash
# Make sure backend is running: make run

# Test individual endpoints
make test-pending-updates

# Or use the test script
./test-pending-updates.sh
```

### 3. Frontend Testing
```bash
# Start frontend
make frontend-dev

# Then open browser and test:
# - http://localhost:5173/updates (Deployment Updates tab)
# - http://localhost:5173/customers/{id} (Pending Updates section)
# - http://localhost:5173/customers/{id}/tenants/{id} (Pending Updates section)
```

### 4. Manual Test Data Setup
See docs/TESTING_PENDING_UPDATES.md for detailed setup instructions.

### 5. Quick Test Checklist
- [ ] Backend API returns data
- [ ] Frontend displays badges
- [ ] Filters work
- [ ] Priority calculation correct
- [ ] Aggregation at all levels works

