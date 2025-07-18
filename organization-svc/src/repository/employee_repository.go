package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/rem-gestion/api-suite/organization/src/dto"
	"github.com/rem-gestion/api-suite/organization/src/models"
)

// employeeRepository implements EmployeeRepository interface
type employeeRepository struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewEmployeeRepository creates a new employee repository instance
func NewEmployeeRepository(db *gorm.DB, logger *zap.Logger) EmployeeRepository {
	return &employeeRepository{db: db, logger: logger}
}

// ===============================
// BASIC CRUD OPERATIONS
// ===============================

func (r *employeeRepository) Create(ctx context.Context, employee *models.Employee) (*models.Employee, error) {
	if err := r.db.WithContext(ctx).Create(employee).Error; err != nil {
		return nil, fmt.Errorf("failed to create employee: %w", err)
	}
	return employee, nil
}

func (r *employeeRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}
	return &employee, nil
}

func (r *employeeRepository) GetByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.WithContext(ctx).Where("user_id = ? AND organization_id = ?", userID, orgID).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}
	return &employee, nil
}

func (r *employeeRepository) Update(ctx context.Context, id uuid.UUID, updates map[string]interface{}) (*models.Employee, error) {
	var employee models.Employee

	// First check if employee exists
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee: %w", err)
	}

	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	// Perform update
	if err := r.db.WithContext(ctx).Model(&employee).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update employee: %w", err)
	}

	// Return updated employee
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&employee).Error; err != nil {
		return nil, fmt.Errorf("failed to get updated employee: %w", err)
	}

	return &employee, nil
}

func (r *employeeRepository) Delete(ctx context.Context, id uuid.UUID) error {
	result := r.db.WithContext(ctx).Unscoped().Delete(&models.Employee{}, id)
	if result.Error != nil {
		return fmt.Errorf("failed to delete employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *employeeRepository) SoftDelete(ctx context.Context, id uuid.UUID, deletedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"deleted_at": time.Now(),
		"deleted_by": deletedBy,
	}

	result := r.db.WithContext(ctx).Model(&models.Employee{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to soft delete employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

// ===============================
// LISTING AND FILTERING
// ===============================

func (r *employeeRepository) List(ctx context.Context, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Employee{})

	// Apply filters
	query = r.applyEmployeeFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyEmployeePaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employees: %w", err)
	}

	return employees, total, nil
}

func (r *employeeRepository) ListByOrganization(ctx context.Context, orgID uuid.UUID, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Employee{}).Where("organization_id = ?", orgID)

	// Apply filters
	query = r.applyEmployeeFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees by organization: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyEmployeePaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employees by organization: %w", err)
	}

	return employees, total, nil
}

func (r *employeeRepository) ListByRole(ctx context.Context, roleID uuid.UUID, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Employee{}).
		Joins("JOIN employee_roles ON employees.id = employee_roles.employee_id").
		Where("employee_roles.role_id = ?", roleID)

	// Apply filters
	query = r.applyEmployeeFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees by role: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyEmployeePaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employees by role: %w", err)
	}

	return employees, total, nil
}

func (r *employeeRepository) ListByBranch(ctx context.Context, branchID uuid.UUID, filters dto.EmployeeFiltersRequest) ([]models.Employee, int64, error) {
	var employees []models.Employee
	var total int64

	query := r.db.WithContext(ctx).Model(&models.Employee{}).Where("branch_id = ?", branchID)

	// Apply filters
	query = r.applyEmployeeFilters(query, filters)

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count employees by branch: %w", err)
	}

	// Apply pagination and sorting
	query = r.applyEmployeePaginationAndSorting(query, filters.PaginationRequest, filters.FilterRequest)

	// Execute query
	if err := query.Find(&employees).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to list employees by branch: %w", err)
	}

	return employees, total, nil
}

// ===============================
// STATUS MANAGEMENT
// ===============================

func (r *employeeRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status models.EmployeeStatus, updatedBy uuid.UUID) error {
	updates := map[string]interface{}{
		"status":     status,
		"updated_by": updatedBy,
		"updated_at": time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.Employee{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update employee status: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

func (r *employeeRepository) GetByStatus(ctx context.Context, orgID uuid.UUID, status models.EmployeeStatus) ([]models.Employee, error) {
	var employees []models.Employee
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND status = ?", orgID, status).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get employees by status: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) TerminateEmployee(ctx context.Context, id uuid.UUID, terminatedBy uuid.UUID, reason string) error {
	updates := map[string]interface{}{
		"status":             models.EmpStatusTerminated,
		"terminated_at":      time.Now(),
		"terminated_by":      terminatedBy,
		"termination_reason": reason,
		"updated_by":         terminatedBy,
		"updated_at":         time.Now(),
	}

	result := r.db.WithContext(ctx).Model(&models.Employee{}).Where("id = ?", id).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to terminate employee: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("employee not found")
	}
	return nil
}

// ===============================
// RELATIONSHIPS
// ===============================

func (r *employeeRepository) GetWithRoles(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.WithContext(ctx).Preload("EmployeeRoles").Preload("EmployeeRoles.Role").Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee with roles: %w", err)
	}
	return &employee, nil
}

func (r *employeeRepository) GetWithOrganization(ctx context.Context, id uuid.UUID) (*models.Employee, error) {
	var employee models.Employee
	if err := r.db.WithContext(ctx).Preload("Organization").Where("id = ?", id).First(&employee).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("employee not found")
		}
		return nil, fmt.Errorf("failed to get employee with organization: %w", err)
	}
	return &employee, nil
}

// ===============================
// BUSINESS LOGIC SUPPORT
// ===============================

func (r *employeeRepository) ExistsByUserAndOrg(ctx context.Context, userID, orgID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("user_id = ? AND organization_id = ?", userID, orgID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check employee existence: %w", err)
	}
	return count > 0, nil
}

func (r *employeeRepository) CountByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).Where("organization_id = ?", orgID).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count employees by organization: %w", err)
	}
	return count, nil
}

func (r *employeeRepository) CountActiveByOrganization(ctx context.Context, orgID uuid.UUID) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("organization_id = ? AND status = ?", orgID, models.EmpStatusActive).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count active employees by organization: %w", err)
	}
	return count, nil
}

func (r *employeeRepository) GetActiveEmployeesByOrg(ctx context.Context, orgID uuid.UUID) ([]models.Employee, error) {
	var employees []models.Employee
	if err := r.db.WithContext(ctx).Where("organization_id = ? AND status = ?", orgID, models.EmpStatusActive).Find(&employees).Error; err != nil {
		return nil, fmt.Errorf("failed to get active employees by organization: %w", err)
	}
	return employees, nil
}

func (r *employeeRepository) IsUserEmployeeOfOrg(ctx context.Context, userID, orgID uuid.UUID) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.Employee{}).
		Where("user_id = ? AND organization_id = ? AND status = ?", userID, orgID, models.EmpStatusActive).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if user is employee of organization: %w", err)
	}
	return count > 0, nil
}

// ===============================
// HELPER METHODS
// ===============================

func (r *employeeRepository) applyEmployeeFilters(query *gorm.DB, filters dto.EmployeeFiltersRequest) *gorm.DB {
	// Search filter
	if filters.Search != "" {
		searchTerm := "%" + filters.Search + "%"
		query = query.Where("notes ILIKE ?", searchTerm)
	}

	// Status filter
	if filters.Status != "" {
		query = query.Where("status = ?", filters.Status)
	}

	// Role filter
	if filters.RoleID != nil {
		query = query.Joins("JOIN employee_roles ON employees.id = employee_roles.employee_id").
			Where("employee_roles.role_id = ?", *filters.RoleID)
	}

	// Branch filter
	if filters.BranchID != nil {
		query = query.Where("branch_id = ?", *filters.BranchID)
	}

	// User filter
	if filters.UserID != nil {
		query = query.Where("user_id = ?", *filters.UserID)
	}

	// Hired date range filters
	if filters.HiredAfter != nil {
		query = query.Where("hired_at >= ?", *filters.HiredAfter)
	}

	if filters.HiredBefore != nil {
		query = query.Where("hired_at <= ?", *filters.HiredBefore)
	}

	// Date filters
	if filters.CreatedAt != nil {
		query = query.Where("created_at >= ?", *filters.CreatedAt)
	}

	if filters.UpdatedAt != nil {
		query = query.Where("updated_at >= ?", *filters.UpdatedAt)
	}

	return query
}

func (r *employeeRepository) applyEmployeePaginationAndSorting(query *gorm.DB, pagination dto.PaginationRequest, filter dto.FilterRequest) *gorm.DB {
	// Apply default sorting
	query = query.Order("created_at DESC")

	// Apply pagination
	if pagination.Page > 0 && pagination.PerPage > 0 {
		offset := (pagination.Page - 1) * pagination.PerPage
		query = query.Offset(offset).Limit(pagination.PerPage)
	}

	return query
}
