# Testing Pending Updates Feature

This guide covers how to test the pending updates tracking feature across backend, frontend, and end-to-end scenarios.

## Prerequisites

1. **Backend running**: `make run` or `cd src/backend && go run cmd/main.go`
2. **Frontend running**: `make frontend-dev` or `cd src/frontend && npm run dev`
3. **MongoDB running**: `make db-start` or ensure MongoDB is accessible
4. **Test data**: Create products, versions, customers, tenants, and deployments

## 1. Setting Up Test Data

### Step 1: Create a Product

```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{
    "product_id": "test-product",
    "name": "Test Product",
    "type": "server",
    "description": "Test product for pending updates"
  }'
```

### Step 2: Create Multiple Versions

```bash
# Version 1.0.0 (oldest)
curl -X POST http://localhost:8080/api/v1/products/test-product/versions \
  -H "Content-Type: application/json" \
  -d '{
    "version_number": "1.0.0",
    "release_date": "2024-01-01T00:00:00Z",
    "release_type": "feature",
    "state": "released"
  }'

# Version 1.1.0 (minor update)
curl -X POST http://localhost:8080/api/v1/products/test-product/versions \
  -H "Content-Type: application/json" \
  -d '{
    "version_number": "1.1.0",
    "release_date": "2024-02-01T00:00:00Z",
    "release_type": "feature",
    "state": "released"
  }'

# Version 1.2.0 (minor update with security)
curl -X POST http://localhost:8080/api/v1/products/test-product/versions \
  -H "Content-Type: application/json" \
  -d '{
    "version_number": "1.2.0",
    "release_date": "2024-03-01T00:00:00Z",
    "release_type": "security",
    "state": "released"
  }'

# Version 2.0.0 (major update)
curl -X POST http://localhost:8080/api/v1/products/test-product/versions \
  -H "Content-Type: application/json" \
  -d '{
    "version_number": "2.0.0",
    "release_date": "2024-04-01T00:00:00Z",
    "release_type": "major",
    "state": "released"
  }'
```

**Note**: You'll need to approve and release these versions. Check the version IDs from the responses and update their state:

```bash
# Get version ID from response, then:
curl -X PUT http://localhost:8080/api/v1/versions/{version_id} \
  -H "Content-Type: application/json" \
  -d '{
    "state": "released"
  }'
```

### Step 3: Create Customer, Tenant, and Deployment

```bash
# Create Customer
curl -X POST http://localhost:8080/api/v1/customers \
  -H "Content-Type: application/json" \
  -d '{
    "customer_id": "test-customer",
    "name": "Test Customer",
    "email": "test@example.com",
    "account_status": "active"
  }'

# Create Tenant (use customer_id from response)
curl -X POST http://localhost:8080/api/v1/customers/test-customer/tenants \
  -H "Content-Type: application/json" \
  -d '{
    "tenant_id": "test-tenant",
    "name": "Test Tenant",
    "status": "active"
  }'

# Create Deployment with old version (use tenant_id from response)
curl -X POST http://localhost:8080/api/v1/customers/test-customer/tenants/test-tenant/deployments \
  -H "Content-Type: application/json" \
  -d '{
    "deployment_id": "test-deployment",
    "product_id": "test-product",
    "deployment_type": "production",
    "installed_version": "1.0.0",
    "status": "active"
  }'
```

## 2. Backend API Testing

### Test 1: Get Pending Updates for a Deployment

```bash
curl http://localhost:8080/api/v1/customers/test-customer/tenants/test-tenant/deployments/test-deployment/updates
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "deployment_id": "test-deployment",
    "product_id": "test-product",
    "current_version": "1.0.0",
    "latest_version": "2.0.0",
    "update_count": 3,
    "priority": "critical",
    "version_gap_type": "major",
    "available_updates": [
      {
        "version_number": "2.0.0",
        "release_date": "2024-04-01T00:00:00Z",
        "release_type": "major",
        "is_security_update": false,
        "compatibility_status": "compatible"
      },
      {
        "version_number": "1.2.0",
        "release_date": "2024-03-01T00:00:00Z",
        "release_type": "security",
        "is_security_update": true,
        "compatibility_status": "compatible"
      },
      {
        "version_number": "1.1.0",
        "release_date": "2024-02-01T00:00:00Z",
        "release_type": "feature",
        "is_security_update": false,
        "compatibility_status": "compatible"
      }
    ]
  }
}
```

### Test 2: Get Tenant Pending Updates

```bash
curl http://localhost:8080/api/v1/customers/test-customer/tenants/test-tenant/deployments/pending-updates
```

**Expected Response:**
```json
{
  "success": true,
  "data": {
    "tenant_id": "test-tenant",
    "tenant_name": "Test Tenant",
    "total_deployments": 1,
    "deployments_with_updates": 1,
    "total_pending_update_count": 3,
    "by_priority": {
      "critical": 1,
      "high": 0,
      "normal": 0
    },
    "by_product": {
      "test-product": 1
    },
    "deployments": [...]
  }
}
```

### Test 3: Get Customer Pending Updates

```bash
curl http://localhost:8080/api/v1/customers/test-customer/deployments/pending-updates
```

### Test 4: Get All Pending Updates (Admin View)

```bash
curl "http://localhost:8080/api/v1/updates/pending?page=1&limit=20"
```

### Test 5: Filter by Priority

```bash
curl "http://localhost:8080/api/v1/updates/pending?priority=critical"
```

### Test 6: Filter by Product

```bash
curl "http://localhost:8080/api/v1/updates/pending?product_id=test-product"
```

## 3. Backend Unit Tests

### Run Version Utility Tests

```bash
cd src/backend
go test ./internal/utils -v -run TestCompareVersions
```

### Run Pending Updates Service Tests

```bash
cd src/backend
go test ./internal/service -v -run TestPendingUpdatesService
```

**Note**: These tests require MongoDB to be running. Make sure your test database is configured.

## 4. Frontend Manual Testing

### Test 1: View Pending Updates in Deployment List

1. Navigate to: `http://localhost:5173/customers/{customer_id}/tenants/{tenant_id}`
2. Go to the "Deployments" tab
3. Verify that deployments with pending updates show an `UpdateBadge` with:
   - Update count
   - Color-coded priority (red for critical, orange for high, blue for normal)

### Test 2: View Deployment Details with Pending Updates

1. Navigate to a deployment details page
2. Verify the `DeploymentPendingUpdates` component shows:
   - Current version vs. latest version
   - Update count
   - Priority badge
   - Version gap type badge
   - List of available updates with details

### Test 3: View Customer Pending Updates Summary

1. Navigate to: `http://localhost:5173/customers/{customer_id}`
2. Go to the "Overview" tab
3. Verify the "Pending Updates" section shows:
   - Deployments with updates count
   - Total pending update count
   - Breakdown by priority (Critical, High, Normal)

### Test 4: View Tenant Pending Updates Summary

1. Navigate to: `http://localhost:5173/customers/{customer_id}/tenants/{tenant_id}`
2. Go to the "Overview" tab
3. Verify the "Pending Updates" section shows tenant-level aggregation

### Test 5: View All Pending Updates (Admin View)

1. Navigate to: `http://localhost:5173/updates`
2. Click on the "Deployment Updates" tab
3. Verify:
   - Table shows all deployments with pending updates
   - Filters work (Priority, Product ID, Deployment Type)
   - Pagination works
   - "View" button navigates to deployment details

### Test 6: Test Priority Calculation

Create test scenarios:

1. **Critical Priority**: Deployment with security update available
   - Create version with `release_type: "security"`
   - Verify priority is "critical"

2. **High Priority**: Production deployment with major version update
   - Create deployment with `deployment_type: "production"` and `installed_version: "1.0.0"`
   - Create version `2.0.0` (major)
   - Verify priority is "high"

3. **Normal Priority**: UAT deployment with minor update
   - Create deployment with `deployment_type: "uat"` and `installed_version: "1.0.0"`
   - Create version `1.1.0` (minor)
   - Verify priority is "normal"

## 5. End-to-End Testing Scenarios

### Scenario 1: New Version Release Triggers Pending Updates

1. Create a deployment with version `1.0.0`
2. Release a new version `1.1.0`
3. Verify:
   - Deployment shows 1 pending update
   - Customer summary shows updated count
   - Admin view shows the deployment in pending updates list

### Scenario 2: Multiple Deployments with Different Versions

1. Create 3 deployments:
   - Deployment A: `1.0.0` (should have 3 updates)
   - Deployment B: `1.1.0` (should have 2 updates)
   - Deployment C: `1.2.0` (should have 1 update)
2. Verify aggregation at tenant/customer level shows correct totals

### Scenario 3: Update Deployment Version

1. Update a deployment's `installed_version` to the latest
2. Verify:
   - Pending updates count becomes 0
   - Badge shows "Up to date"
   - Deployment removed from pending updates lists

### Scenario 4: Filtering and Search

1. Create multiple deployments across different products
2. Test filters:
   - Filter by priority (critical/high/normal)
   - Filter by product ID
   - Filter by deployment type (UAT/Production)
3. Verify results are correctly filtered

## 6. Performance Testing

### Test with Large Dataset

```bash
# Create multiple deployments (script example)
for i in {1..100}; do
  curl -X POST http://localhost:8080/api/v1/customers/test-customer/tenants/test-tenant/deployments \
    -H "Content-Type: application/json" \
    -d "{
      \"deployment_id\": \"deployment-$i\",
      \"product_id\": \"test-product\",
      \"deployment_type\": \"production\",
      \"installed_version\": \"1.0.0\",
      \"status\": \"active\"
    }"
done

# Then test aggregation performance
time curl "http://localhost:8080/api/v1/customers/test-customer/deployments/pending-updates"
```

## 7. Edge Cases to Test

1. **No Pending Updates**: Deployment with latest version should show "Up to date"
2. **Deprecated Versions**: Versions with `state: "deprecated"` should be excluded
3. **EOL Versions**: Versions past EOL date should be excluded
4. **Invalid Deployment ID**: Should return 404
5. **Empty Tenant**: Tenant with no deployments should show 0 updates
6. **Multiple Products**: Customer with deployments across multiple products

## 8. Browser Console Testing

Open browser DevTools and check:

1. **Network Tab**: Verify API calls are made correctly
2. **Console**: Check for any JavaScript errors
3. **React DevTools**: Inspect component state and props

## 9. Common Issues and Solutions

### Issue: "No pending updates" when updates should exist

**Check:**
- Version state is `"released"` (not `"draft"` or `"pending_review"`)
- Version is newer than installed version (semantic versioning)
- Version is not deprecated or EOL

### Issue: Priority not calculated correctly

**Check:**
- Security updates should be "critical"
- Major updates on production should be "high"
- Other updates should be "normal"

### Issue: Badge not showing in deployment list

**Check:**
- `loadPendingUpdates()` is called after deployments load
- API response structure matches expected format
- Deployment ID matches between list and API call

## 10. Automated E2E Tests (Future)

To add Playwright E2E tests:

```typescript
// tests/e2e/pending-updates.spec.ts
import { test, expect } from '@playwright/test';

test('should display pending updates in deployment list', async ({ page }) => {
  await page.goto('/customers/test-customer/tenants/test-tenant');
  await page.click('text=Deployments');
  
  // Wait for pending updates to load
  await page.waitForSelector('[data-testid="update-badge"]');
  
  // Verify badge shows update count
  const badge = await page.locator('[data-testid="update-badge"]').first();
  await expect(badge).toContainText('Updates');
});
```

## Quick Test Checklist

- [ ] Backend API returns pending updates for deployment
- [ ] Backend API aggregates at tenant level
- [ ] Backend API aggregates at customer level
- [ ] Backend API returns all pending updates (admin view)
- [ ] Frontend shows update badges in deployment list
- [ ] Frontend shows pending updates summary on customer page
- [ ] Frontend shows pending updates summary on tenant page
- [ ] Frontend "Deployment Updates" tab displays correctly
- [ ] Filters work (priority, product, type)
- [ ] Pagination works
- [ ] Priority calculation is correct
- [ ] Version gap type is correct
- [ ] Security updates show as critical priority

