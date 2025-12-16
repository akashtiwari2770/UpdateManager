package service

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// AuditLogService handles audit log business logic
type AuditLogService struct {
	auditRepo *repository.AuditLogRepository
}

// NewAuditLogService creates a new audit log service
func NewAuditLogService(auditRepo *repository.AuditLogRepository) *AuditLogService {
	return &AuditLogService{
		auditRepo: auditRepo,
	}
}

// GetAuditLogsByResource retrieves audit logs for a specific resource
func (s *AuditLogService) GetAuditLogsByResource(ctx context.Context, resourceType, resourceID string, page, limit int) ([]*models.AuditLog, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"timestamp": -1})

	logs, err := s.auditRepo.GetByResource(ctx, resourceType, resourceID, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	filter := bson.M{
		"resource_type": resourceType,
		"resource_id":   resourceID,
	}
	total, err := s.auditRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return logs, total, nil
}

// GetAuditLogsByUser retrieves audit logs for a specific user
func (s *AuditLogService) GetAuditLogsByUser(ctx context.Context, userID string, page, limit int) ([]*models.AuditLog, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"timestamp": -1})

	logs, err := s.auditRepo.GetByUserID(ctx, userID, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	total, err := s.auditRepo.Count(ctx, bson.M{"user_id": userID})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return logs, total, nil
}

// GetAuditLogsByAction retrieves audit logs by action
func (s *AuditLogService) GetAuditLogsByAction(ctx context.Context, action models.AuditAction, page, limit int) ([]*models.AuditLog, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"timestamp": -1})

	logs, err := s.auditRepo.GetByAction(ctx, action, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get audit logs: %w", err)
	}

	total, err := s.auditRepo.Count(ctx, bson.M{"action": action})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return logs, total, nil
}

// ListAuditLogs lists audit logs with filters
func (s *AuditLogService) ListAuditLogs(ctx context.Context, filter bson.M, page, limit int) ([]*models.AuditLog, int64, error) {
	opts := options.Find()
	if limit > 0 {
		opts.SetLimit(int64(limit))
		opts.SetSkip(int64((page - 1) * limit))
	}
	opts.SetSort(bson.M{"timestamp": -1})

	logs, err := s.auditRepo.List(ctx, filter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list audit logs: %w", err)
	}

	total, err := s.auditRepo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return logs, total, nil
}
