package repository

import (
	"errors"
	"time"

	"github.com/rem-gestion/api-suite/auth-identity/src/models"
	"gorm.io/gorm"
)

// AccountRepository define la interfaz para operaciones de Account
type AccountRepository interface {
	Create(account *models.Account) error
	GetByEmail(email string) (*models.Account, error)
	GetByID(id string) (*models.Account, error)
	Update(account *models.Account) error
	UpdateStatus(id string, status models.AccountStatus) error
	Delete(id string) error
	Exists(email string) (bool, error)
}

// accountRepository implementa AccountRepository
type accountRepository struct {
	db *gorm.DB
}

// NewAccountRepository crea una nueva instancia del repositorio
func NewAccountRepository(db *gorm.DB) AccountRepository {
	return &accountRepository{db: db}
}

// Create crea una nueva cuenta
func (r *accountRepository) Create(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}

	return r.db.Create(account).Error
}

// GetByEmail obtiene una cuenta por email
func (r *accountRepository) GetByEmail(email string) (*models.Account, error) {
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}

	var account models.Account
	err := r.db.Where("email = ? AND deleted_at IS NULL", email).
		Preload("User").
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // No encontrado, pero no es error
		}
		return nil, err
	}

	return &account, nil
}

// GetByID obtiene una cuenta por ID
func (r *accountRepository) GetByID(id string) (*models.Account, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var account models.Account
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("User").
		First(&account).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &account, nil
}

// Update actualiza una cuenta
func (r *accountRepository) Update(account *models.Account) error {
	if account == nil {
		return errors.New("account cannot be nil")
	}

	now := time.Now()
	account.UpdatedAt = &now

	return r.db.Save(account).Error
}

// UpdateStatus actualiza solo el estado de la cuenta
func (r *accountRepository) UpdateStatus(id string, status models.AccountStatus) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	now := time.Now()
	return r.db.Model(&models.Account{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"status":     status,
			"updated_at": now,
		}).Error
}

// Delete realiza soft delete de una cuenta
func (r *accountRepository) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	now := time.Now()
	return r.db.Model(&models.Account{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

// Exists verifica si existe una cuenta con el email dado
func (r *accountRepository) Exists(email string) (bool, error) {
	if email == "" {
		return false, errors.New("email cannot be empty")
	}

	var count int64
	err := r.db.Model(&models.Account{}).
		Where("email = ? AND deleted_at IS NULL", email).
		Count(&count).Error

	return count > 0, err
}
