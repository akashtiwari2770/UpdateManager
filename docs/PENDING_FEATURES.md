# Pending Features & Enhancements

This document tracks features and enhancements that are planned but not yet implemented.

## Storage Backend

### S3 Support
**Status:** Pending  
**Priority:** Medium  
**Dependencies:** Package upload functionality (completed)

#### Description
Implement AWS S3 (or S3-compatible storage) as an alternative to local disk storage for package files. This will enable:
- Scalable file storage
- Better reliability and redundancy
- CDN integration
- Cross-region availability

#### Implementation Plan
1. Create storage interface/abstraction layer
2. Implement S3 storage backend
3. Add configuration for storage backend selection (local vs S3)
4. Support for S3-compatible services (MinIO, DigitalOcean Spaces, etc.)
5. Presigned URL generation for secure downloads
6. Migration path from local to S3 storage

#### Configuration
```go
type StorageConfig struct {
    Type          string // "local" or "s3"
    LocalPath     string // For local storage
    S3Bucket      string // For S3 storage
    S3Region      string
    S3AccessKey   string
    S3SecretKey   string
    S3Endpoint    string // For S3-compatible services
    PresignedURLs bool   // Enable presigned URLs for downloads
}
```

#### Files to Create/Modify
- `src/backend/internal/storage/interface.go` - Storage interface
- `src/backend/internal/storage/local.go` - Local storage implementation
- `src/backend/internal/storage/s3.go` - S3 storage implementation
- `src/backend/internal/api/handlers/version_handler.go` - Use storage interface
- `src/backend/cmd/server/main.go` - Initialize storage backend

#### Notes
- Current implementation uses local disk storage at `storage/packages/{version_id}/{package_id}.{ext}`
- S3 implementation should maintain the same directory structure in bucket
- Consider using presigned URLs for downloads to avoid proxying through backend
- Add environment variables for S3 configuration

---

## Other Pending Features

### Additional features will be added here as they are identified.

