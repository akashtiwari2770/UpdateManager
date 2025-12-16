package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"updatemanager/internal/models"
	"updatemanager/internal/repository"
)

// CustomerService handles customer business logic
type CustomerService struct {
	customerRepo  *repository.CustomerRepository
	tenantRepo    *repository.TenantRepository
	deploymentRepo *repository.DeploymentRepository
	auditRepo     *repository.AuditLogRepository
}

// NewCustomerService creates a new customer service
func NewCustomerService(
	customerRepo *repository.CustomerRepository,
	tenantRepo *repository.TenantRepository,
	deploymentRepo *repository.DeploymentRepository,
	auditRepo *repository.AuditLogRepository,
) *CustomerService {
	return &CustomerService{
		customerRepo:   customerRepo,
		tenantRepo:     tenantRepo,
		deploymentRepo: deploymentRepo,
		auditRepo:      auditRepo,
	}
}

// ListCustomersQuery represents query parameters for listing customers
type ListCustomersQuery struct {
	Search string
	Status models.CustomerStatus
	Email  string
	Page   int
	Limit  int
}

// CustomerListResponse represents a paginated list of customers
type CustomerListResponse struct {
	Customers []*models.Customer
	Pagination *repository.PaginationInfo
}

// CreateCustomer creates a new customer with validation and audit logging
func (s *CustomerService) CreateCustomer(ctx context.Context, req *models.CreateCustomerRequest, userID, userEmail string) (*models.Customer, error) {
	// Validate customer_id uniqueness
	if req.CustomerID != "" {
		existing, err := s.customerRepo.GetByCustomerID(ctx, req.CustomerID)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("customer with customer_id '%s' already exists", req.CustomerID)
		}
	} else {
		// Generate unique customer_id if not provided
		req.CustomerID = s.generateCustomerID()
	}

	// Validate email uniqueness
	if req.Email != "" {
		existing, err := s.customerRepo.GetByEmail(ctx, req.Email)
		if err == nil && existing != nil {
			return nil, fmt.Errorf("customer with email '%s' already exists", req.Email)
		}
	}

	// Create customer
	customer := &models.Customer{
		CustomerID:   req.CustomerID,
		Name:         req.Name,
		OrganizationName: req.OrganizationName,
		Email:        req.Email,
		Phone:        req.Phone,
		Address:      req.Address,
		AccountStatus: req.AccountStatus,
		NotificationPreferences: req.NotificationPreferences,
	}

	if err := s.customerRepo.Create(ctx, customer); err != nil {
		return nil, fmt.Errorf("failed to create customer: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionCreate, "customer", customer.ID.Hex(), userID, userEmail, map[string]interface{}{
		"customer_id": customer.CustomerID,
		"name":        customer.Name,
		"email":       customer.Email,
	})

	return customer, nil
}

// GetCustomer retrieves a customer by ID (string or ObjectID)
func (s *CustomerService) GetCustomer(ctx context.Context, id string) (*models.Customer, error) {
	// Try to parse as ObjectID first
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		customer, err := s.customerRepo.GetByID(ctx, objectID)
		if err != nil {
			return nil, fmt.Errorf("customer not found: %w", err)
		}
		return customer, nil
	}

	// If not ObjectID, try as customer_id
	customer, err := s.customerRepo.GetByCustomerID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}
	return customer, nil
}

// ListCustomers lists customers with filters and pagination
func (s *CustomerService) ListCustomers(ctx context.Context, query *ListCustomersQuery) (*CustomerListResponse, error) {
	filter := &repository.CustomerFilter{}
	if query != nil {
		filter.Search = query.Search
		filter.Status = query.Status
		filter.Email = query.Email
	}

	pagination := &repository.Pagination{}
	if query != nil {
		pagination.Page = query.Page
		pagination.Limit = query.Limit
	}

	customers, paginationInfo, err := s.customerRepo.List(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to list customers: %w", err)
	}

	return &CustomerListResponse{
		Customers:  customers,
		Pagination: paginationInfo,
	}, nil
}

// UpdateCustomer updates an existing customer
func (s *CustomerService) UpdateCustomer(ctx context.Context, id string, req *models.UpdateCustomerRequest, userID, userEmail string) (*models.Customer, error) {
	// Get existing customer
	customer, err := s.GetCustomer(ctx, id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.OrganizationName != nil {
		customer.OrganizationName = *req.OrganizationName
	}
	if req.Email != nil {
		// Check email uniqueness if changed
		if *req.Email != customer.Email {
			existing, err := s.customerRepo.GetByEmail(ctx, *req.Email)
			if err == nil && existing != nil && existing.ID != customer.ID {
				return nil, fmt.Errorf("customer with email '%s' already exists", *req.Email)
			}
		}
		customer.Email = *req.Email
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.Address != nil {
		customer.Address = *req.Address
	}
	if req.AccountStatus != nil {
		customer.AccountStatus = *req.AccountStatus
	}
	if req.NotificationPreferences != nil {
		customer.NotificationPreferences = *req.NotificationPreferences
	}

	if err := s.customerRepo.Update(ctx, customer.ID, customer); err != nil {
		return nil, fmt.Errorf("failed to update customer: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionUpdate, "customer", customer.ID.Hex(), userID, userEmail, map[string]interface{}{
		"customer_id": customer.CustomerID,
	})

	return customer, nil
}

// DeleteCustomer deletes a customer (soft delete by setting status to inactive)
func (s *CustomerService) DeleteCustomer(ctx context.Context, id string, userID, userEmail string) error {
	customer, err := s.GetCustomer(ctx, id)
	if err != nil {
		return err
	}

	// Check for existing tenants
	tenants, _, err := s.tenantRepo.GetByCustomerID(ctx, customer.ID, nil, &repository.Pagination{Page: 1, Limit: 1})
	if err == nil && len(tenants) > 0 {
		return fmt.Errorf("cannot delete customer: customer has existing tenants")
	}

	// Soft delete (set status to inactive)
	customer.AccountStatus = models.CustomerStatusInactive
	if err := s.customerRepo.Update(ctx, customer.ID, customer); err != nil {
		return fmt.Errorf("failed to delete customer: %w", err)
	}

	// Log audit
	s.logAudit(ctx, models.AuditActionDelete, "customer", customer.ID.Hex(), userID, userEmail, map[string]interface{}{
		"customer_id": customer.CustomerID,
	})

	return nil
}

// GetCustomerTenants retrieves tenants for a customer
func (s *CustomerService) GetCustomerTenants(ctx context.Context, customerID string, query *ListTenantsQuery) (*TenantListResponse, error) {
	customer, err := s.GetCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	filter := &repository.TenantFilter{}
	if query != nil && query.Status != "" {
		filter.Status = query.Status
	}

	pagination := &repository.Pagination{}
	if query != nil {
		pagination.Page = query.Page
		pagination.Limit = query.Limit
	}

	tenants, paginationInfo, err := s.tenantRepo.GetByCustomerID(ctx, customer.ID, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("failed to get customer tenants: %w", err)
	}

	return &TenantListResponse{
		Tenants:    tenants,
		Pagination: paginationInfo,
	}, nil
}

// CustomerStatistics represents statistics for a customer
type CustomerStatistics struct {
	TotalTenants     int64
	TotalDeployments int64
	TotalUsers       int64
	DeploymentsByProduct map[string]int64
	DeploymentsByType    map[string]int64
}

// GetCustomerStatistics retrieves statistics for a customer
func (s *CustomerService) GetCustomerStatistics(ctx context.Context, customerID string) (*CustomerStatistics, error) {
	customer, err := s.GetCustomer(ctx, customerID)
	if err != nil {
		return nil, fmt.Errorf("customer not found: %w", err)
	}

	// Get all tenants for customer
	tenants, _, err := s.tenantRepo.GetByCustomerID(ctx, customer.ID, nil, &repository.Pagination{Page: 1, Limit: 1000})
	if err != nil {
		return nil, fmt.Errorf("failed to get tenants: %w", err)
	}

	stats := &CustomerStatistics{
		TotalTenants:        int64(len(tenants)),
		DeploymentsByProduct: make(map[string]int64),
		DeploymentsByType:    make(map[string]int64),
	}

	var totalUsers int64
	var totalDeployments int64

	// Count deployments per tenant
	for _, tenant := range tenants {
		deployments, _, err := s.deploymentRepo.GetByTenantID(ctx, tenant.ID, nil, &repository.Pagination{Page: 1, Limit: 1000})
		if err != nil {
			continue
		}

		for _, deployment := range deployments {
			totalDeployments++
			stats.DeploymentsByProduct[deployment.ProductID]++
			stats.DeploymentsByType[string(deployment.DeploymentType)]++
			if deployment.NumberOfUsers != nil {
				totalUsers += int64(*deployment.NumberOfUsers)
			}
		}
	}

	stats.TotalDeployments = totalDeployments
	stats.TotalUsers = totalUsers

	return stats, nil
}

// Helper methods

// generateCustomerID generates a unique customer ID
func (s *CustomerService) generateCustomerID() string {
	// Generate a unique ID based on timestamp
	timestamp := time.Now().Unix()
	return fmt.Sprintf("CUST-%d", timestamp)
}

// logAudit logs an audit entry
func (s *CustomerService) logAudit(ctx context.Context, action models.AuditAction, resourceType, resourceID, userID, userEmail string, details map[string]interface{}) {
	if s.auditRepo == nil {
		return
	}

	auditLog := &models.AuditLog{
		Action:       action,
		ResourceType: resourceType,
		ResourceID:   resourceID,
		UserID:       userID,
		UserEmail:    userEmail,
		Details:      details,
		Timestamp:    time.Now(),
	}

	_ = s.auditRepo.Create(ctx, auditLog) // Ignore errors for audit logging
}

// validateCustomerData validates customer data
func (s *CustomerService) validateCustomerData(customer *models.Customer) error {
	if strings.TrimSpace(customer.CustomerID) == "" {
		return fmt.Errorf("customer_id is required")
	}
	if strings.TrimSpace(customer.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if strings.TrimSpace(customer.Email) == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}

