package circuit

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// HealthChecker define la interfaz para verificar salud de servicios
type HealthChecker interface {
	Check(ctx context.Context) error
}

// ServiceHealthChecker verifica la salud de múltiples servicios
type ServiceHealthChecker struct {
	services map[string]HealthChecker
	mu       sync.RWMutex
}

// NewServiceHealthChecker crea un nuevo health checker
func NewServiceHealthChecker() *ServiceHealthChecker {
	return &ServiceHealthChecker{
		services: make(map[string]HealthChecker),
	}
}

// RegisterService registra un servicio para health checking
func (shc *ServiceHealthChecker) RegisterService(name string, checker HealthChecker) {
	shc.mu.Lock()
	defer shc.mu.Unlock()
	shc.services[name] = checker
}

// CheckService verifica la salud de un servicio específico
func (shc *ServiceHealthChecker) CheckService(ctx context.Context, serviceName string) error {
	shc.mu.RLock()
	checker, exists := shc.services[serviceName]
	shc.mu.RUnlock()

	if !exists {
		return fmt.Errorf("service %s not registered", serviceName)
	}

	return checker.Check(ctx)
}

// CheckAllServices verifica la salud de todos los servicios
func (shc *ServiceHealthChecker) CheckAllServices(ctx context.Context) map[string]error {
	shc.mu.RLock()
	services := make(map[string]HealthChecker)
	for name, checker := range shc.services {
		services[name] = checker
	}
	shc.mu.RUnlock()

	results := make(map[string]error)
	var wg sync.WaitGroup
	var mu sync.Mutex

	for name, checker := range services {
		wg.Add(1)
		go func(serviceName string, serviceChecker HealthChecker) {
			defer wg.Done()

			checkCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err := serviceChecker.Check(checkCtx)

			mu.Lock()
			results[serviceName] = err
			mu.Unlock()
		}(name, checker)
	}

	wg.Wait()
	return results
}

// RequiredServicesHealthy verifica que todos los servicios requeridos estén saludables
func (shc *ServiceHealthChecker) RequiredServicesHealthy(ctx context.Context, requiredServices []string) error {
	var unhealthyServices []string

	for _, serviceName := range requiredServices {
		if err := shc.CheckService(ctx, serviceName); err != nil {
			unhealthyServices = append(unhealthyServices, fmt.Sprintf("%s: %v", serviceName, err))
		}
	}

	if len(unhealthyServices) > 0 {
		return fmt.Errorf("required services are unhealthy: %v", unhealthyServices)
	}

	return nil
}

// SimpleHealthChecker implementación simple para testing
type SimpleHealthChecker struct {
	healthFunc func(ctx context.Context) error
}

// NewSimpleHealthChecker crea un health checker simple
func NewSimpleHealthChecker(healthFunc func(ctx context.Context) error) *SimpleHealthChecker {
	return &SimpleHealthChecker{
		healthFunc: healthFunc,
	}
}

// Check ejecuta la función de health check
func (shc *SimpleHealthChecker) Check(ctx context.Context) error {
	if shc.healthFunc == nil {
		return nil // Siempre saludable si no hay función definida
	}
	return shc.healthFunc(ctx)
}
