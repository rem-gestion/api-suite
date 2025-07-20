package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"

	"github.com/rem-gestion/api-suite/organization/src/models"
)

// CacheService define la interfaz para operaciones de cache
type CacheService interface {
	// Organization cache operations
	GetOrganization(orgID string) (*models.Organization, error)
	SetOrganization(org *models.Organization, ttl time.Duration) error
	DeleteOrganization(orgID string) error

	// Organization settings cache operations
	GetOrganizationSettings(orgID string) (map[string]interface{}, error)
	SetOrganizationSettings(orgID string, settings map[string]interface{}, ttl time.Duration) error
	DeleteOrganizationSettings(orgID string) error

	// Employee cache operations
	GetEmployee(empID string) (*models.Employee, error)
	SetEmployee(emp *models.Employee, ttl time.Duration) error
	DeleteEmployee(empID string) error

	// Branch cache operations
	GetBranch(branchID string) (*models.OrganizationBranch, error)
	SetBranch(branch *models.OrganizationBranch, ttl time.Duration) error
	DeleteBranch(branchID string) error

	// Role cache operations
	GetRole(roleID string) (*models.OrganizationRole, error)
	SetRole(role *models.OrganizationRole, ttl time.Duration) error
	DeleteRole(roleID string) error

	// List cache operations (for pagination)
	GetEmployeesList(orgID string, page, perPage int) ([]models.Employee, error)
	SetEmployeesList(orgID string, page, perPage int, employees []models.Employee, ttl time.Duration) error

	// Invalidation operations
	InvalidateOrganizationCache(orgID string) error
}

// RedisOrganizationCacheService implementa CacheService usando Redis
type RedisOrganizationCacheService struct {
	redis  *redis.Client
	logger *zap.Logger
}

// NewRedisOrganizationCacheService crea una nueva instancia del servicio de cache
func NewRedisOrganizationCacheService(redis *redis.Client, logger *zap.Logger) CacheService {
	return &RedisOrganizationCacheService{
		redis:  redis,
		logger: logger,
	}
}

// NewCacheService crea una nueva instancia de CacheService
// Si redis es nil, retorna un NoOpCacheService
func NewCacheService(redis *redis.Client, logger *zap.Logger) CacheService {
	if redis == nil {
		return &NoOpCacheService{logger: logger}
	}
	return NewRedisOrganizationCacheService(redis, logger)
}

// Organization cache operations implementation

func (c *RedisOrganizationCacheService) GetOrganization(orgID string) (*models.Organization, error) {
	key := c.organizationCacheKey(orgID)

	data, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		c.logger.Error("Failed to get organization from cache", zap.Error(err))
		return nil, err
	}

	var org models.Organization
	if err := json.Unmarshal([]byte(data), &org); err != nil {
		c.logger.Error("Failed to unmarshal organization from cache", zap.Error(err))
		return nil, err
	}

	return &org, nil
}

func (c *RedisOrganizationCacheService) SetOrganization(org *models.Organization, ttl time.Duration) error {
	key := c.organizationCacheKey(org.ID.String())

	data, err := json.Marshal(org)
	if err != nil {
		c.logger.Error("Failed to marshal organization for cache", zap.Error(err))
		return err
	}

	if err := c.redis.Set(context.Background(), key, string(data), ttl).Err(); err != nil {
		c.logger.Error("Failed to set organization in cache", zap.Error(err))
		return err
	}

	return nil
}

func (c *RedisOrganizationCacheService) DeleteOrganization(orgID string) error {
	key := c.organizationCacheKey(orgID)

	if err := c.redis.Del(context.Background(), key).Err(); err != nil {
		c.logger.Error("Failed to delete organization from cache", zap.Error(err))
		return err
	}

	return nil
}

// Organization settings cache operations implementation

func (c *RedisOrganizationCacheService) GetOrganizationSettings(orgID string) (map[string]interface{}, error) {
	key := c.organizationSettingsCacheKey(orgID)

	data, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		c.logger.Error("Failed to get organization settings from cache", zap.Error(err))
		return nil, err
	}

	var settings map[string]interface{}
	if err := json.Unmarshal([]byte(data), &settings); err != nil {
		c.logger.Error("Failed to unmarshal organization settings from cache", zap.Error(err))
		return nil, err
	}

	return settings, nil
}

func (c *RedisOrganizationCacheService) SetOrganizationSettings(orgID string, settings map[string]interface{}, ttl time.Duration) error {
	key := c.organizationSettingsCacheKey(orgID)

	data, err := json.Marshal(settings)
	if err != nil {
		c.logger.Error("Failed to marshal organization settings for cache", zap.Error(err))
		return err
	}

	if err := c.redis.Set(context.Background(), key, string(data), ttl).Err(); err != nil {
		c.logger.Error("Failed to set organization settings in cache", zap.Error(err))
		return err
	}

	return nil
}

func (c *RedisOrganizationCacheService) DeleteOrganizationSettings(orgID string) error {
	key := c.organizationSettingsCacheKey(orgID)

	if err := c.redis.Del(context.Background(), key).Err(); err != nil {
		c.logger.Error("Failed to delete organization settings from cache", zap.Error(err))
		return err
	}

	return nil
}

// Employee cache operations implementation

func (c *RedisOrganizationCacheService) GetEmployee(empID string) (*models.Employee, error) {
	key := c.employeeCacheKey(empID)

	data, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		c.logger.Error("Failed to get employee from cache", zap.Error(err))
		return nil, err
	}

	var emp models.Employee
	if err := json.Unmarshal([]byte(data), &emp); err != nil {
		c.logger.Error("Failed to unmarshal employee from cache", zap.Error(err))
		return nil, err
	}

	return &emp, nil
}

func (c *RedisOrganizationCacheService) SetEmployee(emp *models.Employee, ttl time.Duration) error {
	key := c.employeeCacheKey(emp.ID.String())

	data, err := json.Marshal(emp)
	if err != nil {
		c.logger.Error("Failed to marshal employee for cache", zap.Error(err))
		return err
	}

	if err := c.redis.Set(context.Background(), key, string(data), ttl).Err(); err != nil {
		c.logger.Error("Failed to set employee in cache", zap.Error(err))
		return err
	}

	return nil
}

func (c *RedisOrganizationCacheService) DeleteEmployee(empID string) error {
	key := c.employeeCacheKey(empID)

	if err := c.redis.Del(context.Background(), key).Err(); err != nil {
		c.logger.Error("Failed to delete employee from cache", zap.Error(err))
		return err
	}

	return nil
}

// Branch cache operations implementation

func (c *RedisOrganizationCacheService) GetBranch(branchID string) (*models.OrganizationBranch, error) {
	key := c.branchCacheKey(branchID)

	data, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		c.logger.Error("Failed to get branch from cache", zap.Error(err))
		return nil, err
	}

	var branch models.OrganizationBranch
	if err := json.Unmarshal([]byte(data), &branch); err != nil {
		c.logger.Error("Failed to unmarshal branch from cache", zap.Error(err))
		return nil, err
	}

	return &branch, nil
}

func (c *RedisOrganizationCacheService) SetBranch(branch *models.OrganizationBranch, ttl time.Duration) error {
	key := c.branchCacheKey(branch.ID.String())

	data, err := json.Marshal(branch)
	if err != nil {
		c.logger.Error("Failed to marshal branch for cache", zap.Error(err))
		return err
	}

	if err := c.redis.Set(context.Background(), key, string(data), ttl).Err(); err != nil {
		c.logger.Error("Failed to set branch in cache", zap.Error(err))
		return err
	}

	return nil
}

func (c *RedisOrganizationCacheService) DeleteBranch(branchID string) error {
	key := c.branchCacheKey(branchID)

	if err := c.redis.Del(context.Background(), key).Err(); err != nil {
		c.logger.Error("Failed to delete branch from cache", zap.Error(err))
		return err
	}

	return nil
}

// Role cache operations implementation

func (c *RedisOrganizationCacheService) GetRole(roleID string) (*models.OrganizationRole, error) {
	key := c.roleCacheKey(roleID)

	data, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		c.logger.Error("Failed to get role from cache", zap.Error(err))
		return nil, err
	}

	var role models.OrganizationRole
	if err := json.Unmarshal([]byte(data), &role); err != nil {
		c.logger.Error("Failed to unmarshal role from cache", zap.Error(err))
		return nil, err
	}

	return &role, nil
}

func (c *RedisOrganizationCacheService) SetRole(role *models.OrganizationRole, ttl time.Duration) error {
	key := c.roleCacheKey(role.ID.String())

	data, err := json.Marshal(role)
	if err != nil {
		c.logger.Error("Failed to marshal role for cache", zap.Error(err))
		return err
	}

	if err := c.redis.Set(context.Background(), key, string(data), ttl).Err(); err != nil {
		c.logger.Error("Failed to set role in cache", zap.Error(err))
		return err
	}

	return nil
}

func (c *RedisOrganizationCacheService) DeleteRole(roleID string) error {
	key := c.roleCacheKey(roleID)

	if err := c.redis.Del(context.Background(), key).Err(); err != nil {
		c.logger.Error("Failed to delete role from cache", zap.Error(err))
		return err
	}

	return nil
}

// List cache operations implementation

func (c *RedisOrganizationCacheService) GetEmployeesList(orgID string, page, perPage int) ([]models.Employee, error) {
	key := c.employeesListCacheKey(orgID, page, perPage)

	data, err := c.redis.Get(context.Background(), key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil // Not found in cache
		}
		c.logger.Error("Failed to get employees list from cache", zap.Error(err))
		return nil, err
	}

	var employees []models.Employee
	if err := json.Unmarshal([]byte(data), &employees); err != nil {
		c.logger.Error("Failed to unmarshal employees list from cache", zap.Error(err))
		return nil, err
	}

	return employees, nil
}

func (c *RedisOrganizationCacheService) SetEmployeesList(orgID string, page, perPage int, employees []models.Employee, ttl time.Duration) error {
	key := c.employeesListCacheKey(orgID, page, perPage)

	data, err := json.Marshal(employees)
	if err != nil {
		c.logger.Error("Failed to marshal employees list for cache", zap.Error(err))
		return err
	}

	if err := c.redis.Set(context.Background(), key, string(data), ttl).Err(); err != nil {
		c.logger.Error("Failed to set employees list in cache", zap.Error(err))
		return err
	}

	return nil
}

// Invalidation operations implementation

func (c *RedisOrganizationCacheService) InvalidateOrganizationCache(orgID string) error {
	// Delete organization cache
	orgKey := c.organizationCacheKey(orgID)
	settingsKey := c.organizationSettingsCacheKey(orgID)

	// Use a pipeline for efficiency
	pipe := c.redis.Pipeline()
	pipe.Del(context.Background(), orgKey)
	pipe.Del(context.Background(), settingsKey)

	// Also invalidate employees list cache for this organization
	pattern := fmt.Sprintf("org:employees:%s:*", orgID)
	keys, err := c.redis.Keys(context.Background(), pattern).Result()
	if err != nil {
		c.logger.Error("Failed to get cache keys for invalidation", zap.Error(err))
		return err
	}

	if len(keys) > 0 {
		pipe.Del(context.Background(), keys...)
	}

	_, err = pipe.Exec(context.Background())
	if err != nil {
		c.logger.Error("Failed to invalidate organization cache", zap.Error(err))
		return err
	}

	return nil
}

// Cache key generation helpers

func (c *RedisOrganizationCacheService) organizationCacheKey(orgID string) string {
	return fmt.Sprintf("org:%s", orgID)
}

func (c *RedisOrganizationCacheService) organizationSettingsCacheKey(orgID string) string {
	return fmt.Sprintf("org:settings:%s", orgID)
}

func (c *RedisOrganizationCacheService) employeeCacheKey(empID string) string {
	return fmt.Sprintf("emp:%s", empID)
}

func (c *RedisOrganizationCacheService) branchCacheKey(branchID string) string {
	return fmt.Sprintf("branch:%s", branchID)
}

func (c *RedisOrganizationCacheService) roleCacheKey(roleID string) string {
	return fmt.Sprintf("role:%s", roleID)
}

func (c *RedisOrganizationCacheService) employeesListCacheKey(orgID string, page, perPage int) string {
	return fmt.Sprintf("org:employees:%s:%d:%d", orgID, page, perPage)
}

// NoOpCacheService implements CacheService but does nothing (for when Redis is not available)
type NoOpCacheService struct {
	logger *zap.Logger
}

func (n *NoOpCacheService) GetOrganization(orgID string) (*models.Organization, error) {
	return nil, nil
}
func (n *NoOpCacheService) SetOrganization(org *models.Organization, ttl time.Duration) error {
	return nil
}
func (n *NoOpCacheService) DeleteOrganization(orgID string) error { return nil }
func (n *NoOpCacheService) GetOrganizationSettings(orgID string) (map[string]interface{}, error) {
	return nil, nil
}
func (n *NoOpCacheService) SetOrganizationSettings(orgID string, settings map[string]interface{}, ttl time.Duration) error {
	return nil
}
func (n *NoOpCacheService) DeleteOrganizationSettings(orgID string) error             { return nil }
func (n *NoOpCacheService) GetEmployee(empID string) (*models.Employee, error)        { return nil, nil }
func (n *NoOpCacheService) SetEmployee(emp *models.Employee, ttl time.Duration) error { return nil }
func (n *NoOpCacheService) DeleteEmployee(empID string) error                         { return nil }
func (n *NoOpCacheService) GetBranch(branchID string) (*models.OrganizationBranch, error) {
	return nil, nil
}
func (n *NoOpCacheService) SetBranch(branch *models.OrganizationBranch, ttl time.Duration) error {
	return nil
}
func (n *NoOpCacheService) DeleteBranch(branchID string) error                      { return nil }
func (n *NoOpCacheService) GetRole(roleID string) (*models.OrganizationRole, error) { return nil, nil }
func (n *NoOpCacheService) SetRole(role *models.OrganizationRole, ttl time.Duration) error {
	return nil
}
func (n *NoOpCacheService) DeleteRole(roleID string) error { return nil }
func (n *NoOpCacheService) GetEmployeesList(orgID string, page, perPage int) ([]models.Employee, error) {
	return nil, nil
}
func (n *NoOpCacheService) SetEmployeesList(orgID string, page, perPage int, employees []models.Employee, ttl time.Duration) error {
	return nil
}
func (n *NoOpCacheService) InvalidateOrganizationCache(orgID string) error { return nil }
