package repository

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	model "github.com/rem-gestion/api-suite/address/src/models"
	"github.com/rem-gestion/rem-common/config"
	"github.com/rem-gestion/rem-common/db"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// ConnectionEvent representa eventos de conexión
type ConnectionEvent struct {
	Type      string // "connected", "disconnected", "retry"
	Timestamp time.Time
	Error     error
}

// AdaptiveAddressRepo es una versión mejorada con canales para comunicación
type AdaptiveAddressRepo struct {
	mu  sync.RWMutex
	lg  *zap.Logger
	cfg *config.PostgresConfig

	// Estado de la conexión
	dbConnected bool
	db          *gorm.DB

	// Almacenamiento en memoria
	memoryStore map[string]*model.Address
	memoryMutex sync.RWMutex

	// Canales para comunicación asíncrona
	connectionEvents chan ConnectionEvent
	retryTrigger     chan struct{}
	shutdown         chan struct{}

	// Control de concurrencia
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup

	// Métricas
	retryCount    int64
	lastError     error
	lastConnected time.Time
	memoryOps     int64
	dbOps         int64

	// Control de logging para evitar spam
	initialWarningLogged bool
	lastFailureLogged    time.Time
}

// NewAdaptive crea un repositorio adaptativo mejorado con canales
func NewAdaptive(cfg *config.PostgresConfig, lg *zap.Logger) *AdaptiveAddressRepo {
	ctx, cancel := context.WithCancel(context.Background())

	repo := &AdaptiveAddressRepo{
		lg:               lg.Named("adaptive-repo"),
		cfg:              cfg,
		dbConnected:      false,
		memoryStore:      make(map[string]*model.Address),
		connectionEvents: make(chan ConnectionEvent, 100), // Buffer para no bloquear
		retryTrigger:     make(chan struct{}, 1),
		shutdown:         make(chan struct{}),
		ctx:              ctx,
		cancel:           cancel,
	}

	// Iniciar goroutines de fondo
	repo.wg.Add(3)
	go repo.connectionManager() // Maneja la conexión principal
	go repo.retryManager()      // Maneja los reintentos
	go repo.eventProcessor()    // Procesa eventos de conexión

	// Trigger inicial de conexión
	repo.triggerRetry()

	return repo
}

// connectionManager maneja la conexión principal a la base de datos
func (r *AdaptiveAddressRepo) connectionManager() {
	defer r.wg.Done()

	r.lg.Info("connection manager started")

	for {
		select {
		case <-r.ctx.Done():
			r.lg.Info("connection manager shutting down")
			return
		case <-r.retryTrigger:
			r.attemptConnection()
		}
	}
}

// retryManager maneja los reintentos periódicos
func (r *AdaptiveAddressRepo) retryManager() {
	defer r.wg.Done()

	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	r.lg.Info("retry manager started")

	for {
		select {
		case <-r.ctx.Done():
			r.lg.Info("retry manager shutting down")
			return
		case <-ticker.C:
			// Solo reintentar si no estamos conectados
			r.mu.RLock()
			needsRetry := !r.dbConnected
			r.mu.RUnlock()

			if needsRetry {
				r.lg.Debug("periodic retry triggered")
				r.triggerRetry()
			}
		}
	}
}

// eventProcessor procesa eventos de conexión de manera asíncrona
func (r *AdaptiveAddressRepo) eventProcessor() {
	defer r.wg.Done()

	r.lg.Info("event processor started")

	for {
		select {
		case <-r.ctx.Done():
			r.lg.Info("event processor shutting down")
			return
		case event := <-r.connectionEvents:
			r.handleConnectionEvent(event)
		}
	}
}

// attemptConnection intenta conectar con la base de datos
func (r *AdaptiveAddressRepo) attemptConnection() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.retryCount++

	// Solo loggear el primer intento o cada 5 minutos para evitar spam
	shouldLog := !r.initialWarningLogged || time.Since(r.lastFailureLogged) > 5*time.Minute

	r.lg.Debug("attempting database connection", zap.Int64("attempt", r.retryCount))

	pg, err := db.NewPostgres(*r.cfg)
	if err != nil {
		r.lastError = err
		r.dbConnected = false
		r.db = nil

		if shouldLog {
			r.lg.Warn("database connection failed - using memory fallback",
				zap.Error(err),
				zap.Int64("attempt", r.retryCount))
			r.lastFailureLogged = time.Now()
			if !r.initialWarningLogged {
				r.initialWarningLogged = true
			}
		}

		// Enviar evento de error de manera no bloqueante
		select {
		case r.connectionEvents <- ConnectionEvent{
			Type:      "disconnected",
			Timestamp: time.Now(),
			Error:     err,
		}:
		default:
			// Canal lleno, skipear evento
		}
		return
	}

	// Verificar que la conexión funcione
	sqlDB, err := pg.DB()
	if err != nil {
		r.lastError = err
		r.dbConnected = false
		r.db = nil

		if shouldLog {
			r.lg.Warn("database connection validation failed - using memory fallback",
				zap.Error(err),
				zap.Int64("attempt", r.retryCount))
			r.lastFailureLogged = time.Now()
			if !r.initialWarningLogged {
				r.initialWarningLogged = true
			}
		}

		select {
		case r.connectionEvents <- ConnectionEvent{
			Type:      "disconnected",
			Timestamp: time.Now(),
			Error:     err,
		}:
		default:
		}
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		r.lastError = err
		r.dbConnected = false
		r.db = nil

		if shouldLog {
			r.lg.Warn("database ping failed - using memory fallback",
				zap.Error(err),
				zap.Int64("attempt", r.retryCount))
			r.lastFailureLogged = time.Now()
			if !r.initialWarningLogged {
				r.initialWarningLogged = true
			}
		}

		select {
		case r.connectionEvents <- ConnectionEvent{
			Type:      "disconnected",
			Timestamp: time.Now(),
			Error:     err,
		}:
		default:
		}
		return
	}

	// Conexión exitosa
	wasDisconnected := !r.dbConnected
	r.db = pg
	r.dbConnected = true
	r.lastConnected = time.Now()
	r.lastError = nil

	// Enviar evento de conexión exitosa
	select {
	case r.connectionEvents <- ConnectionEvent{
		Type:      "connected",
		Timestamp: time.Now(),
		Error:     nil,
	}:
	default:
	}

	if wasDisconnected {
		r.lg.Info("database connection restored",
			zap.Int64("retry_count", r.retryCount),
			zap.Duration("downtime", time.Since(r.lastConnected)))

		// Sincronizar datos de memoria a DB en background
		go r.syncMemoryToDB()
	}
}

// handleConnectionEvent maneja eventos de conexión
func (r *AdaptiveAddressRepo) handleConnectionEvent(event ConnectionEvent) {
	switch event.Type {
	case "connected":
		r.lg.Info("database connection established", zap.Time("timestamp", event.Timestamp))

	case "disconnected":
		// Solo loggear si es la primera desconexión o ha pasado tiempo suficiente
		r.mu.RLock()
		shouldLog := !r.initialWarningLogged || time.Since(r.lastFailureLogged) > 5*time.Minute
		r.mu.RUnlock()

		if shouldLog {
			r.lg.Warn("database connection lost - fallback to memory store active",
				zap.Time("timestamp", event.Timestamp),
				zap.Error(event.Error))
		}

	case "retry":
		r.lg.Debug("retry triggered", zap.Time("timestamp", event.Timestamp))
	}
}

// triggerRetry dispara un intento de reconexión de manera no bloqueante
func (r *AdaptiveAddressRepo) triggerRetry() {
	select {
	case r.retryTrigger <- struct{}{}:
		// Retry triggereado exitosamente
	default:
		// Ya hay un retry pendiente, no es necesario agregar otro
	}
}

// syncMemoryToDB sincroniza datos de memoria a la base de datos
func (r *AdaptiveAddressRepo) syncMemoryToDB() {
	r.lg.Info("starting memory to database migration process")

	r.memoryMutex.RLock()
	addresses := make([]*model.Address, 0, len(r.memoryStore))
	for _, addr := range r.memoryStore {
		// Crear copia para evitar race conditions
		addrCopy := *addr
		addresses = append(addresses, &addrCopy)
	}
	r.memoryMutex.RUnlock()

	if len(addresses) == 0 {
		r.lg.Info("no addresses in memory to migrate")
		return
	}

	r.lg.Info("found addresses in memory to migrate", zap.Int("count", len(addresses)))

	r.mu.RLock()
	db := r.db
	connected := r.dbConnected
	r.mu.RUnlock()

	if !connected || db == nil {
		r.lg.Warn("cannot migrate: database not available")
		return
	}

	successCount := 0
	errorCount := 0
	duplicateCount := 0

	for _, addr := range addresses {
		// Crear una nueva instancia sin ID para que GORM genere uno nuevo
		dbAddr := &model.Address{
			Street:  addr.Street,
			Number:  addr.Number,
			Floor:   addr.Floor,
			Unit:    addr.Unit,
			City:    addr.City,
			State:   addr.State,
			Zip:     addr.Zip,
			Country: addr.Country,
		}

		if err := db.Create(dbAddr).Error; err != nil {
			// Verificar si es error de duplicado
			if strings.Contains(strings.ToLower(err.Error()), "duplicate") ||
				strings.Contains(strings.ToLower(err.Error()), "unique") {
				duplicateCount++
				r.lg.Debug("address already exists in database during migration",
					zap.String("memory_id", addr.ID),
					zap.String("address", fmt.Sprintf("%s %d, %s", addr.Street, addr.Number, addr.City)))
			} else {
				errorCount++
				r.lg.Warn("failed to migrate address to database",
					zap.String("memory_id", addr.ID),
					zap.String("address", fmt.Sprintf("%s %d, %s", addr.Street, addr.Number, addr.City)),
					zap.Error(err))
			}
		} else {
			successCount++
			r.lg.Debug("successfully migrated address to database",
				zap.String("memory_id", addr.ID),
				zap.String("db_id", dbAddr.ID),
				zap.String("address", fmt.Sprintf("%s %d, %s", addr.Street, addr.Number, addr.City)))
		}
	}

	// Log del resultado de la migración
	if successCount > 0 || duplicateCount > 0 {
		r.lg.Info("memory to database migration completed",
			zap.Int("total_addresses", len(addresses)),
			zap.Int("migrated", successCount),
			zap.Int("duplicates_skipped", duplicateCount),
			zap.Int("errors", errorCount))

		// Limpiar memoria después de migración (exitosa o duplicados son ambos ok)
		r.memoryMutex.Lock()
		clearedCount := len(r.memoryStore)
		r.memoryStore = make(map[string]*model.Address)
		r.memoryMutex.Unlock()

		r.lg.Info("memory cache cleared after migration",
			zap.Int("cleared_addresses", clearedCount))
	} else {
		r.lg.Error("memory to database migration failed completely",
			zap.Int("total_addresses", len(addresses)),
			zap.Int("errors", errorCount))
	}
}

// IsConnected retorna si la base de datos está conectada
func (r *AdaptiveAddressRepo) IsConnected() bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.dbConnected
}

// ForceRetry fuerza un intento de reconexión inmediato
func (r *AdaptiveAddressRepo) ForceRetry() {
	r.lg.Info("forcing database reconnection attempt")
	r.triggerRetry()
}

// ForceSyncMemoryToDB fuerza la sincronización de memoria a DB si hay conexión
func (r *AdaptiveAddressRepo) ForceSyncMemoryToDB() error {
	r.mu.RLock()
	connected := r.dbConnected
	r.mu.RUnlock()

	if !connected {
		return fmt.Errorf("database not connected - cannot sync")
	}

	r.memoryMutex.RLock()
	hasData := len(r.memoryStore) > 0
	r.memoryMutex.RUnlock()

	if !hasData {
		r.lg.Info("no data in memory to sync")
		return nil
	}

	r.lg.Info("forcing memory to database sync")
	go r.syncMemoryToDB()
	return nil
}

// GetStats retorna estadísticas detalladas
func (r *AdaptiveAddressRepo) GetStats() map[string]interface{} {
	r.mu.RLock()
	r.memoryMutex.RLock()
	defer r.mu.RUnlock()
	defer r.memoryMutex.RUnlock()

	stats := map[string]interface{}{
		"db_connected":           r.dbConnected,
		"retry_count":            r.retryCount,
		"memory_operations":      r.memoryOps,
		"db_operations":          r.dbOps,
		"memory_store_size":      len(r.memoryStore),
		"last_connected":         r.lastConnected,
		"events_buffered":        len(r.connectionEvents),
		"initial_warning_logged": r.initialWarningLogged,
		"last_failure_logged":    r.lastFailureLogged,
		"mode":                   r.getCurrentMode(),
	}

	if r.lastError != nil {
		stats["last_error"] = r.lastError.Error()
	}

	// Agregar detalles de las direcciones en memoria si las hay
	if len(r.memoryStore) > 0 {
		memoryAddresses := make([]map[string]interface{}, 0, len(r.memoryStore))
		for id, addr := range r.memoryStore {
			memoryAddresses = append(memoryAddresses, map[string]interface{}{
				"id":      id,
				"address": fmt.Sprintf("%s %d, %s", addr.Street, addr.Number, addr.City),
				"created": addr.CreatedAt,
			})
		}
		stats["memory_addresses"] = memoryAddresses
	}

	return stats
}

// getCurrentMode retorna el modo actual de operación
func (r *AdaptiveAddressRepo) getCurrentMode() string {
	if r.dbConnected {
		if len(r.memoryStore) > 0 {
			return "database_with_pending_migration"
		}
		return "database"
	}
	return "memory_fallback"
}

// Close cierra el repositorio de manera limpia
func (r *AdaptiveAddressRepo) Close() {
	r.lg.Info("shutting down adaptive repository")

	// Cerrar canales y cancelar contexto
	r.cancel()
	close(r.shutdown)

	// Esperar a que todas las goroutines terminen
	r.wg.Wait()

	// Cerrar conexión a DB si existe
	r.mu.Lock()
	if r.db != nil {
		if sqlDB, err := r.db.DB(); err == nil {
			sqlDB.Close()
		}
	}
	r.mu.Unlock()

	r.lg.Info("adaptive repository shutdown complete")
}

// generateMemoryID genera un ID único para almacenamiento en memoria
func (r *AdaptiveAddressRepo) generateMemoryID() string {
	return fmt.Sprintf("mem_%d_%d", time.Now().UnixNano(), r.memoryOps)
}
