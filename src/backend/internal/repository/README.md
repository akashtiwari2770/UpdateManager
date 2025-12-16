# Repository Layer

This package contains repository implementations for database operations.

## Product Repository

The `ProductRepository` provides CRUD operations for products.

### Usage

```go
import (
    "context"
    "updatemanager/pkg/database"
    "updatemanager/internal/repository"
)

// Connect to database
db, err := database.Connect(ctx, database.DefaultConfig())
if err != nil {
    log.Fatal(err)
}
defer db.Disconnect(ctx)

// Create repository
productRepo := repository.NewProductRepository(db.Collection("products"))

// Create a product
product := &models.Product{
    ProductID: "my-product",
    Name:      "My Product",
    Type:      models.ProductTypeServer,
    IsActive:  true,
}
err = productRepo.Create(ctx, product)

// Get by ID
retrieved, err := productRepo.GetByID(ctx, product.ID)

// Update
product.Name = "Updated Name"
err = productRepo.Update(ctx, product)

// Delete
err = productRepo.Delete(ctx, product.ID)

// List
products, err := productRepo.List(ctx, bson.M{"is_active": true}, nil)

// Count
count, err := productRepo.Count(ctx, bson.M{"type": models.ProductTypeServer})
```

## Testing

Run all repository tests:

```bash
cd src/backend
go test -v ./internal/repository
```

Run a specific test:

```bash
go test -v ./internal/repository -run TestProductCreate
```

### Test Database

Tests use a separate test database (`updatemanager_test`) to avoid affecting production data. The test database is automatically cleaned up after each test run.

## Test Coverage

The Product repository includes comprehensive tests for:

- ✅ Create product
- ✅ Get by ID
- ✅ Get by ProductID
- ✅ Update product
- ✅ Delete product
- ✅ List products (with filters)
- ✅ Count products
- ✅ Error handling (not found cases)

