# Product Field Clarification

## Overview

This document clarifies the Product data model fields to ensure consistency between code and documentation.

## Product Fields

### Product Name
- **Type:** String (free text)
- **Length:** 1-200 characters
- **Required:** Yes
- **Unique:** Yes
- **Purpose:** Human-readable, descriptive name for the product
- **Examples:**
  - "HyWorks"
  - "HySecure"
  - "IRIS"
  - "ARS"
  - "Client for Linux"
  - "Client for Windows"
  - "Client for Mobile"
  - "Accops Application Virtualization Platform"
- **Notes:** This is where specific product names are stored. It's a free-form text field that can be customized by the user.

### Product Type
- **Type:** Enum (string)
- **Values:** 
  - `"server"` - Server-side products
  - `"client"` - Client-side products
- **Required:** Yes
- **Unique:** No (multiple products can have the same type)
- **Purpose:** Categorical classification that determines:
  - Update strategy (Blue-Green for servers, In-Place for clients)
  - Multi-tenant support capabilities
  - Compatibility requirements
  - Deployment behavior
- **Notes:** This is a fixed enum value, not free text. It determines the product's role in the system.

### Product ID
- **Type:** String
- **Length:** 1-100 characters
- **Required:** Yes
- **Unique:** Yes
- **Purpose:** Unique technical identifier for the product
- **Format:** Alphanumeric with hyphens
- **Examples:** "hyworks-prod-001", "hysecure-vpn", "iris-iam"
- **Notes:** Immutable after creation. Used for API references and internal tracking.

## Data Model Structure

```go
type Product struct {
    ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
    ProductID   string             `bson:"product_id" json:"product_id" validate:"required,min=1,max=100"`
    Name        string             `bson:"name" json:"name" validate:"required,min=1,max=200"`
    Type        ProductType        `bson:"type" json:"type" validate:"required"`
    Description string             `bson:"description" json:"description" validate:"max=1000"`
    Vendor      string             `bson:"vendor" json:"vendor" validate:"max=100"`
    CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
    IsActive    bool               `bson:"is_active" json:"is_active"`
}

type ProductType string

const (
    ProductTypeServer ProductType = "server"
    ProductTypeClient ProductType = "client"
)
```

## Example Products

### Server Product Example
```json
{
  "product_id": "hyworks-prod-001",
  "name": "HyWorks",
  "type": "server",
  "description": "Application Virtualization Platform",
  "vendor": "Accops",
  "is_active": true
}
```

### Client Product Example
```json
{
  "product_id": "client-linux-001",
  "name": "Client for Linux",
  "type": "client",
  "description": "Linux client application for HyWorks",
  "vendor": "Accops",
  "is_active": true
}
```

## Key Points

1. **Product Name** is free text and stores specific product names like "HyWorks", "HySecure", etc.
2. **Product Type** is an enum with only two values: "server" or "client"
3. There is **no separate "Product Category" field** - the Type field serves this purpose
4. The specific product names (HyWorks, HySecure, IRIS, ARS) are stored in the **Name** field, not the Type field
5. The Type field determines business logic and behavior (update strategies, compatibility rules, etc.)

## Migration Notes

If you have existing documentation or code that references:
- **"Product Category"** → This should be replaced with **"Product Type"**
- **"Product Type" with values like "HyWorks", "HySecure"** → These should be stored in **"Product Name"**

## Validation Rules

- Product Name must be unique across all products
- Product Type must be either "server" or "client"
- Product ID must be unique and immutable after creation
- Product Name length: 1-200 characters
- Product ID length: 1-100 characters

---

**Last Updated:** 2025-11-13  
**Status:** Current and Accurate

