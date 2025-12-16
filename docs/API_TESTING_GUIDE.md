# API Testing Guide

## Overview

This guide covers all methods for testing the Update Manager API:
1. **Automated Unit/Integration Tests** - Existing test suite
2. **Manual API Testing** - Using curl, Postman, or HTTP clients
3. **Load Testing** - Using Artillery

---

## 1. Automated API Tests

### Prerequisites

1. **Start MongoDB** (required for tests):
```bash
make db-start
```

2. **Verify MongoDB is running**:
```bash
make db-status
```

### Running API Tests

#### Run All API Handler Tests
```bash
make test-api
```

This runs all 34 API handler tests covering:
- ProductHandler (9 tests)
- VersionHandler (8 tests)
- CompatibilityHandler (3 tests)
- NotificationHandler (4 tests)
- UpgradePathHandler (3 tests)
- UpdateDetectionHandler (2 tests)
- UpdateRolloutHandler (3 tests)
- AuditLogHandler (2 tests)

#### Run API Tests with Coverage
```bash
make test-api-coverage
```

This generates a coverage report showing which code paths are tested.

#### Run Specific Handler Tests
```bash
# Run only ProductHandler tests
cd src/backend && go test -v ./internal/api/handlers -run TestProductHandler

# Run only VersionHandler tests
cd src/backend && go test -v ./internal/api/handlers -run TestVersionHandler

# Run a specific test
cd src/backend && go test -v ./internal/api/handlers -run TestProductHandler_CreateProduct
```

#### Run All Backend Tests (Repository + Service + API)
```bash
make test-backend
```

### Test Structure

API tests are located in:
```
src/backend/internal/api/handlers/*_handler_test.go
```

Each test:
- Uses a real MongoDB test database (`updatemanager_test`)
- Sets up test data before running
- Cleans up after completion
- Tests HTTP requests/responses using `httptest`

### Example Test Pattern

```go
func TestProductHandler_CreateProduct(t *testing.T) {
    // Setup: Connect to test DB, create handler
    handler, cleanup := setupProductHandlerTest(t)
    defer cleanup()

    // Create request
    req := models.CreateProductRequest{
        ProductID: "test-product-1",
        Name:      "Test Product",
        Type:      models.ProductTypeServer,
    }
    body, _ := json.Marshal(req)
    httpReq := httptest.NewRequest(http.MethodPost, "/api/v1/products", bytes.NewBuffer(body))
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("X-User-ID", "test-user")

    // Execute
    w := httptest.NewRecorder()
    handler.CreateProduct(w, httpReq)

    // Assert
    if w.Code != http.StatusCreated {
        t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
    }
    // ... more assertions
}
```

### Writing New API Tests

1. **Create test file**: `src/backend/internal/api/handlers/your_handler_test.go`
2. **Use setup function** (see existing tests for pattern)
3. **Test all scenarios**:
   - Success cases
   - Error cases (404, 400, 409, 500)
   - Validation errors
   - Edge cases

---

## 2. Manual API Testing

### Prerequisites

1. **Start the backend server**:
```bash
make run
```

Server runs on `http://localhost:8080` by default.

2. **Verify server is running**:
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{"status":"healthy"}
```

### Using curl

#### Create a Product
```bash
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -H "X-User-ID: test-user" \
  -H "X-User-Email: test@example.com" \
  -d '{
    "product_id": "my-product",
    "name": "My Product",
    "type": "server",
    "description": "A test product",
    "vendor": "Test Vendor"
  }'
```

#### Get All Products
```bash
curl http://localhost:8080/api/v1/products?page=1&limit=10
```

#### Get Product by ID
```bash
# Replace {id} with actual product ID
curl http://localhost:8080/api/v1/products/{id}
```

#### Create a Version
```bash
curl -X POST http://localhost:8080/api/v1/products/{product_id}/versions \
  -H "Content-Type: application/json" \
  -H "X-User-ID: test-user" \
  -H "X-User-Email: test@example.com" \
  -d '{
    "version_number": "1.0.0",
    "release_type": "feature",
    "release_date": "2025-01-15T00:00:00Z"
  }'
```

#### Submit Version for Review
```bash
curl -X POST http://localhost:8080/api/v1/versions/{version_id}/submit \
  -H "X-User-ID: reviewer"
```

#### Approve Version
```bash
curl -X POST http://localhost:8080/api/v1/versions/{version_id}/approve \
  -H "X-User-ID: approver" \
  -H "X-User-Email: approver@example.com"
```

#### Release Version
```bash
curl -X POST http://localhost:8080/api/v1/versions/{version_id}/release \
  -H "X-User-ID: releaser"
```

### Using Postman

1. **Import Collection** (create from API spec):
   - Base URL: `http://localhost:8080`
   - API prefix: `/api/v1`

2. **Set Headers**:
   - `Content-Type: application/json`
   - `X-User-ID: your-user-id`
   - `X-User-Email: your-email@example.com`

3. **Test Endpoints**:
   - Products: `GET/POST /api/v1/products`
   - Versions: `GET/POST /api/v1/products/{id}/versions`
   - Compatibility: `POST /api/v1/products/{id}/versions/{version}/compatibility`
   - Notifications: `GET/POST /api/v1/notifications`
   - Upgrade Paths: `POST /api/v1/products/{id}/upgrade-paths`
   - Update Detection: `POST /api/v1/update-detections`
   - Update Rollouts: `POST /api/v1/update-rollouts`
   - Audit Logs: `GET /api/v1/audit-logs`

### Using HTTPie

```bash
# Install: pip install httpie

# Create product
http POST localhost:8080/api/v1/products \
  X-User-ID:test-user \
  X-User-Email:test@example.com \
  product_id=my-product \
  name="My Product" \
  type=server

# Get products
http GET localhost:8080/api/v1/products
```

### Using JavaScript/Node.js

```javascript
const axios = require('axios');

const API_BASE = 'http://localhost:8080/api/v1';

async function testAPI() {
  // Create product
  const product = await axios.post(`${API_BASE}/products`, {
    product_id: 'test-product',
    name: 'Test Product',
    type: 'server'
  }, {
    headers: {
      'X-User-ID': 'test-user',
      'X-User-Email': 'test@example.com'
    }
  });
  
  console.log('Created product:', product.data);
  
  // Get products
  const products = await axios.get(`${API_BASE}/products`);
  console.log('Products:', products.data);
}

testAPI();
```

---

## 3. Load Testing with Artillery

### Prerequisites

Install Artillery:
```bash
npm install -g artillery
```

### Available Load Tests

#### Mixed Load Test
```bash
make load-test
```
Runs a balanced mix of read and write operations.

#### Read-Heavy Load Test
```bash
make load-test-read
```
Focuses on GET requests (products, versions, notifications).

#### Write-Heavy Load Test
```bash
make load-test-write
```
Focuses on POST/PUT requests (create, update operations).

#### Spike Test
```bash
make load-test-spike
```
Tests system behavior under sudden traffic spikes.

### Load Test Configuration

Test configurations are in `load-tests/`:
- `artillery-config.yml` - Mixed load test
- `artillery-read-heavy.yml` - Read-heavy test
- `artillery-write-heavy.yml` - Write-heavy test
- `artillery-spike-test.yml` - Spike test

### Customizing Load Tests

Edit the YAML files to adjust:
- **Duration**: How long the test runs
- **Arrival rate**: Requests per second
- **Scenarios**: Which endpoints to test
- **Payloads**: Test data to use

Example:
```yaml
config:
  target: 'http://localhost:8080'
  phases:
    - duration: 60
      arrivalRate: 10
      name: "Warm up"
    - duration: 300
      arrivalRate: 50
      name: "Sustained load"
scenarios:
  - name: "Get products"
    flow:
      - get:
          url: "/api/v1/products"
```

### Running Custom Load Test

```bash
artillery run load-tests/your-custom-test.yml
```

---

## 4. API Endpoints Reference

### Products
- `GET /api/v1/products` - List products (with pagination)
- `POST /api/v1/products` - Create product
- `GET /api/v1/products/{id}` - Get product by ID
- `PUT /api/v1/products/{id}` - Update product
- `DELETE /api/v1/products/{id}` - Delete product
- `GET /api/v1/products/active` - Get active products
- `GET /api/v1/products/by-product-id/{product_id}` - Get by product ID

### Versions
- `GET /api/v1/products/{product_id}/versions` - List versions
- `POST /api/v1/products/{product_id}/versions` - Create version
- `GET /api/v1/versions/{id}` - Get version
- `PUT /api/v1/versions/{id}` - Update version
- `POST /api/v1/versions/{id}/submit` - Submit for review
- `POST /api/v1/versions/{id}/approve` - Approve version
- `POST /api/v1/versions/{id}/release` - Release version

### Compatibility
- `POST /api/v1/products/{product_id}/versions/{version}/compatibility` - Validate compatibility
- `GET /api/v1/products/{product_id}/versions/{version}/compatibility` - Get compatibility
- `GET /api/v1/compatibility` - List compatibility matrices

### Notifications
- `GET /api/v1/notifications?recipient_id={id}` - Get notifications
- `POST /api/v1/notifications` - Create notification
- `GET /api/v1/notifications/unread-count?recipient_id={id}` - Get unread count
- `POST /api/v1/notifications/mark-all-read` - Mark all as read

### Upgrade Paths
- `POST /api/v1/products/{product_id}/upgrade-paths` - Create upgrade path
- `GET /api/v1/products/{product_id}/upgrade-paths/{from}/{to}` - Get upgrade path
- `POST /api/v1/products/{product_id}/upgrade-paths/{from}/{to}/block` - Block upgrade path

### Update Detection
- `POST /api/v1/update-detections` - Detect/register update
- `PUT /api/v1/update-detections/{id}/available-version` - Update available version

### Update Rollouts
- `POST /api/v1/update-rollouts` - Initiate rollout
- `PUT /api/v1/update-rollouts/{id}/status` - Update rollout status
- `PUT /api/v1/update-rollouts/{id}/progress` - Update rollout progress

### Audit Logs
- `GET /api/v1/audit-logs` - Get audit logs (with filters)

---

## 5. Testing Best Practices

### 1. Test Data Management
- Use separate test database (`updatemanager_test`)
- Clean up test data after tests
- Use unique identifiers to avoid conflicts

### 2. Test Coverage
- Test success cases
- Test error cases (400, 404, 409, 500)
- Test validation errors
- Test edge cases (empty data, large payloads, etc.)

### 3. Headers
Always include:
- `Content-Type: application/json` (for POST/PUT)
- `X-User-ID: {user_id}` (for audit logging)
- `X-User-Email: {email}` (optional, for notifications)

### 4. Response Format
All responses follow this structure:
```json
{
  "success": true,
  "data": { ... },
  "meta": { ... },  // For paginated responses
  "error": null
}
```

Error responses:
```json
{
  "success": false,
  "data": null,
  "error": {
    "message": "Error description",
    "code": "ERROR_CODE"
  }
}
```

### 5. Status Codes
- `200 OK` - Successful GET, PUT, DELETE
- `201 Created` - Successful POST
- `400 Bad Request` - Validation error
- `404 Not Found` - Resource not found
- `409 Conflict` - Duplicate resource
- `500 Internal Server Error` - Server error

---

## 6. Troubleshooting

### Tests Fail with Database Connection Error
```bash
# Ensure MongoDB is running
make db-start

# Check MongoDB status
make db-status
```

### Tests Fail with Port Already in Use
```bash
# Check what's using port 8080
lsof -i :8080

# Kill the process or change PORT in .env
```

### Load Tests Fail
```bash
# Ensure server is running
make run

# Check server logs for errors
# Verify MongoDB can handle the load
```

### API Returns 500 Errors
- Check MongoDB connection
- Check server logs
- Verify request format matches API spec
- Check required headers are present

---

## 7. Quick Reference

```bash
# Start services
make db-start          # Start MongoDB
make run              # Start API server

# Run tests
make test-api         # Run API tests
make test-api-coverage # API tests with coverage
make test-backend     # All backend tests

# Load testing
make load-test        # Mixed load
make load-test-read   # Read-heavy
make load-test-write  # Write-heavy
make load-test-spike  # Spike test

# Health check
curl http://localhost:8080/health
```

---

## Additional Resources

- [API Specification](api-specification.md) - Complete API documentation
- [API Layer Tests](api-layer-tests.md) - Detailed test documentation
- [Test Counts](test-counts.md) - Test statistics

