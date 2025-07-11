package repository

import (
	"strings"
	"time"

	model "github.com/rem-gestion/api-suite/address/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
)

// Create crea una nueva dirección usando DB o memoria según disponibilidad
func (r *AdaptiveAddressRepo) Create(a *model.Address) (*model.Address, error) {
	if r.IsConnected() {
		return r.createInDB(a)
	}
	return r.createInMemory(a)
}

// createInDB crea la dirección en base de datos
func (r *AdaptiveAddressRepo) createInDB(a *model.Address) (*model.Address, error) {
	r.mu.RLock()
	db := r.db
	r.mu.RUnlock()

	if db == nil {
		return r.createInMemory(a) // Fallback a memoria
	}

	r.mu.Lock()
	r.dbOps++
	r.mu.Unlock()

	// 1) intento de inserción
	if err := db.Create(a).Error; err != nil {
		r.lg.Warn("insert failed, trying dedup", zap.Error(err))

		f, u, s, c, st, z, co := r.normalize(a)
		var existing model.Address

		findErr := db.
			Where(`lower(trim(floor))   = ? AND
				   lower(trim(unit))    = ? AND
				   lower(trim(street))  = ? AND
				   number              = ? AND
				   lower(trim(city))    = ? AND
				   lower(trim(state))   = ? AND
				   trim(zip)           = ? AND
				   upper(trim(country)) = ?`,
				f, u, s, a.Number, c, st, z, co).
			First(&existing).Error

		if findErr != nil {
			r.lg.Error("dedup query failed", zap.Error(findErr))
			// Si falla la deduplicación, intentar con memoria
			return r.createInMemory(a)
		}

		r.lg.Info("address already exists, returning existing",
			zap.String("existing_id", existing.ID))
		return &existing, nil
	}

	r.lg.Info("address created in database", zap.String("id", a.ID))
	return a, nil
}

// createInMemory crea la dirección en memoria
func (r *AdaptiveAddressRepo) createInMemory(a *model.Address) (*model.Address, error) {
	r.memoryMutex.Lock()
	defer r.memoryMutex.Unlock()

	r.memoryOps++

	// Generar nuevo ID para memoria
	memAddr := &model.Address{
		ID:        r.generateMemoryID(),
		Street:    a.Street,
		Number:    a.Number,
		Floor:     a.Floor,
		Unit:      a.Unit,
		City:      a.City,
		State:     a.State,
		Zip:       a.Zip,
		Country:   a.Country,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verificar duplicados en memoria usando normalización
	f, u, s, c, st, z, co := r.normalize(memAddr)
	for _, existing := range r.memoryStore {
		ef, eu, es, ec, est, ez, eco := r.normalize(existing)
		if ef == f && eu == u && es == s && existing.Number == memAddr.Number &&
			ec == c && est == st && ez == z && eco == co {
			r.lg.Info("address already exists in memory, returning existing",
				zap.String("existing_id", existing.ID))
			return existing, nil
		}
	}

	r.memoryStore[memAddr.ID] = memAddr
	r.lg.Info("address created in memory", zap.String("id", memAddr.ID))
	return memAddr, nil
}

// Get obtiene una dirección por ID
func (r *AdaptiveAddressRepo) Get(id string) (*model.Address, error) {
	if r.IsConnected() {
		return r.getFromDB(id)
	}
	return r.getFromMemory(id)
}

// GetByID es un alias de Get para compatibilidad
func (r *AdaptiveAddressRepo) GetByID(id string) (*model.Address, error) {
	return r.Get(id)
}

// getFromDB obtiene la dirección de la base de datos
func (r *AdaptiveAddressRepo) getFromDB(id string) (*model.Address, error) {
	r.mu.RLock()
	db := r.db
	r.mu.RUnlock()

	if db == nil {
		return r.getFromMemory(id) // Fallback a memoria
	}

	r.mu.Lock()
	r.dbOps++
	r.mu.Unlock()

	var a model.Address
	if err := db.First(&a, "id = ?", id).Error; err != nil {
		r.lg.Warn("get address from db failed", zap.String("id", id), zap.Error(err))
		// Intentar en memoria como fallback
		return r.getFromMemory(id)
	}
	return &a, nil
}

// getFromMemory obtiene la dirección de la memoria
func (r *AdaptiveAddressRepo) getFromMemory(id string) (*model.Address, error) {
	r.memoryMutex.RLock()
	defer r.memoryMutex.RUnlock()

	if addr, exists := r.memoryStore[id]; exists {
		r.lg.Debug("address found in memory", zap.String("id", id))
		return addr, nil
	}

	r.lg.Warn("address not found in memory", zap.String("id", id))
	return nil, &rerrors.NotFoundError{Msg: "address not found"}
}

// List lista direcciones con filtros opcionales
func (r *AdaptiveAddressRepo) List(city string, limit, offset int) ([]model.Address, error) {
	if r.IsConnected() {
		return r.listFromDB(city, limit, offset)
	}
	return r.listFromMemory(city, limit, offset)
}

// listFromDB lista direcciones de la base de datos
func (r *AdaptiveAddressRepo) listFromDB(city string, limit, offset int) ([]model.Address, error) {
	r.mu.RLock()
	db := r.db
	r.mu.RUnlock()

	if db == nil {
		return r.listFromMemory(city, limit, offset) // Fallback a memoria
	}

	r.mu.Lock()
	r.dbOps++
	r.mu.Unlock()

	var res []model.Address
	q := db.Limit(limit).Offset(offset)
	if city != "" {
		q = q.Where("city = ?", city)
	}
	if err := q.Find(&res).Error; err != nil {
		r.lg.Warn("list addresses from db failed", zap.Error(err))
		// Fallback a memoria
		return r.listFromMemory(city, limit, offset)
	}
	return res, nil
}

// listFromMemory lista direcciones de la memoria
func (r *AdaptiveAddressRepo) listFromMemory(city string, limit, offset int) ([]model.Address, error) {
	r.memoryMutex.RLock()
	defer r.memoryMutex.RUnlock()

	var results []model.Address

	// Filtrar por ciudad si se especifica
	for _, addr := range r.memoryStore {
		if city == "" || strings.EqualFold(addr.City, city) {
			results = append(results, *addr)
		}
	}

	// Aplicar paginación
	start := offset
	if start > len(results) {
		return []model.Address{}, nil
	}

	end := start + limit
	if end > len(results) {
		end = len(results)
	}

	r.lg.Debug("listed addresses from memory",
		zap.Int("total", len(results)),
		zap.Int("returned", end-start))

	return results[start:end], nil
}

// Update actualiza una dirección (retorna error porque las direcciones son inmutables)
func (r *AdaptiveAddressRepo) Update(*model.Address) error {
	return &rerrors.ForbiddenError{Msg: "addresses are immutable; create a new record"}
}

// Delete elimina una dirección
func (r *AdaptiveAddressRepo) Delete(id string) error {
	if r.IsConnected() {
		return r.deleteFromDB(id)
	}
	return r.deleteFromMemory(id)
}

// deleteFromDB elimina la dirección de la base de datos
func (r *AdaptiveAddressRepo) deleteFromDB(id string) error {
	r.mu.RLock()
	db := r.db
	r.mu.RUnlock()

	if db == nil {
		return r.deleteFromMemory(id) // Fallback a memoria
	}

	r.mu.Lock()
	r.dbOps++
	r.mu.Unlock()

	if err := db.Delete(&model.Address{}, "id = ?", id).Error; err != nil {
		r.lg.Error("delete address from db failed", zap.String("id", id), zap.Error(err))
		// También intentar eliminar de memoria como fallback
		return r.deleteFromMemory(id)
	}

	r.lg.Info("address deleted from database", zap.String("id", id))
	return nil
}

// deleteFromMemory elimina la dirección de la memoria
func (r *AdaptiveAddressRepo) deleteFromMemory(id string) error {
	r.memoryMutex.Lock()
	defer r.memoryMutex.Unlock()

	if _, exists := r.memoryStore[id]; !exists {
		return &rerrors.NotFoundError{Msg: "address not found"}
	}

	delete(r.memoryStore, id)
	r.lg.Info("address deleted from memory", zap.String("id", id))
	return nil
}

// normalize helper function para normalización de datos
func (r *AdaptiveAddressRepo) normalize(a *model.Address) (floor, unit, street, city, state, zip, country string) {
	tl := func(s string) string { return strings.ToLower(strings.TrimSpace(s)) }
	floor = tl(a.Floor)
	unit = tl(a.Unit)
	street = tl(a.Street)
	city = tl(a.City)
	state = tl(a.State)
	zip = strings.TrimSpace(a.Zip)
	country = strings.ToUpper(strings.TrimSpace(a.Country))
	return
}
