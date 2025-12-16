# Update Manager

Product Release Management System for Accops products.

## Project Structure

```
UpdateManager/
├── docs/                    # Documentation (markdown files only)
│   ├── api-specification.md
│   ├── implementation-guide.md
│   └── product-release-process.md
├── src/
│   ├── backend/            # Go backend application
│   │   ├── cmd/
│   │   │   └── server/     # Main application entry point
│   │   ├── internal/       # Internal packages
│   │   │   ├── api/        # API handlers and routes
│   │   │   ├── models/     # Data models
│   │   │   ├── repository/ # Data access layer
│   │   │   ├── service/    # Business logic
│   │   │   └── config/     # Configuration
│   │   ├── pkg/            # Public packages
│   │   │   ├── database/   # Database utilities
│   │   │   ├── validator/  # Validation utilities
│   │   │   └── logger/     # Logging utilities
│   │   └── go.mod          # Go module definition
│   ├── frontend/           # React frontend application
│   │   └── src/
│   │       ├── components/ # React components
│   │       ├── services/   # API service layer
│   │       ├── types/      # TypeScript types
│   │       └── hooks/      # React hooks
│   └── database/           # Database setup scripts
│       ├── mongodb-indexes.js
│       ├── setup-database.js
│       └── docker-compose.mongodb.yml
├── build/                  # Build artifacts (generated)
├── Makefile               # Build and development commands
└── README.md              # This file
```

## Prerequisites

- Go 1.21 or higher
- MongoDB 6.0 or higher
- Node.js 18+ and npm (for frontend)
- Docker and Docker Compose (optional, for MongoDB)

## Quick Start

### 1. Setup MongoDB

**Using Docker (Recommended):**
```bash
cd src/database
docker-compose -f docker-compose.mongodb.yml up -d
```

**Or using local MongoDB:**
```bash
# Ensure MongoDB is running
make db-setup
```

### 2. Create Database Indexes

```bash
make db-indexes
```

### 3. Install Backend Dependencies

```bash
make install
```

### 4. Build and Run Backend

```bash
make build
make run
```

Or for development:
```bash
make dev
```

## Development

### Backend Development

```bash
# Install dependencies
make install

# Format code
make fmt

# Run linter
make vet

# Run tests
make test

# Run server
make run
```

### Frontend Development

```bash
cd src/frontend
npm install
npm run dev
```

## Makefile Commands

- `make help` - Show available commands
- `make install` - Install Go dependencies
- `make build` - Build backend binary
- `make run` - Run backend server
- `make test` - Run tests
- `make fmt` - Format Go code
- `make vet` - Run go vet
- `make clean` - Clean build artifacts
- `make db-setup` - Setup MongoDB database
- `make db-indexes` - Create MongoDB indexes
- `make docker-build` - Build Docker images
- `make docker-up` - Start Docker containers
- `make docker-down` - Stop Docker containers
- `make setup` - Full setup (install + db setup + indexes)

## Configuration

Create a `.env` file in the backend directory:

```env
# Server
SERVER_PORT=8080
SERVER_HOST=localhost

# MongoDB
MONGODB_URI=mongodb://admin:admin123@localhost:27017/updatemanager?authSource=admin
MONGODB_DATABASE=updatemanager

# JWT
JWT_SECRET=your-secret-key-here
JWT_EXPIRY=24h

# File Storage
STORAGE_PATH=./storage
MAX_FILE_SIZE=1073741824  # 1GB in bytes
```

## API Documentation

See [docs/api-specification.md](docs/api-specification.md) for complete API documentation.

## Implementation Phases

The project is implemented in two phases:

### Phase 1: Foundation and Core Release Management
- Product management
- Version creation and management
- Package upload
- Release approval workflow

### Phase 2: Enhanced Release Management and Distribution
- Compatibility validation
- Notification system
- Update detection
- Update rollout management

See [docs/product-release-process.md](docs/product-release-process.md) for details.

## License

Proprietary - Accops Technologies

