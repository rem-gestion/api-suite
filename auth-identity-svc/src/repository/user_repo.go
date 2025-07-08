package repository

import (
	"errors"
	"time"

	"github.com/rem-gestion/api-suite/auth-identity/src/models"
	"gorm.io/gorm"
)

// UserRepository define la interfaz para operaciones de User
type UserRepository interface {
	Create(user *models.User) error
	GetByID(id string) (*models.User, error)
	GetByAccountID(accountID string) (*models.User, error)
	GetByPersonID(personID string) (*models.User, error)
	Update(user *models.User) error
	UpdateLastLogin(id string) error
	UpdateOnboardStatus(id string, status models.OnboardStatus) error
	UpdatePersonID(id string, personID string) error
	Delete(id string) error
	List(page, perPage int, filter *UserFilter) ([]*models.User, int64, error)
}

// UserFilter define los filtros para listar usuarios
type UserFilter struct {
	OnboardStatus *models.OnboardStatus
	AccountStatus *models.AccountStatus
	HasPerson     *bool
	Search        string // busca en email
}

// userRepository implementa UserRepository
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository crea una nueva instancia del repositorio
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create crea un nuevo usuario
func (r *userRepository) Create(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	return r.db.Create(user).Error
}

// GetByID obtiene un usuario por ID
func (r *userRepository) GetByID(id string) (*models.User, error) {
	if id == "" {
		return nil, errors.New("id cannot be empty")
	}

	var user models.User
	err := r.db.Where("id = ? AND deleted_at IS NULL", id).
		Preload("Account").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByAccountID obtiene un usuario por Account ID
func (r *userRepository) GetByAccountID(accountID string) (*models.User, error) {
	if accountID == "" {
		return nil, errors.New("accountID cannot be empty")
	}

	var user models.User
	err := r.db.Where("account_id = ? AND deleted_at IS NULL", accountID).
		Preload("Account").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByPersonID obtiene un usuario por Person ID
func (r *userRepository) GetByPersonID(personID string) (*models.User, error) {
	if personID == "" {
		return nil, errors.New("personID cannot be empty")
	}

	var user models.User
	err := r.db.Where("person_id = ? AND deleted_at IS NULL", personID).
		Preload("Account").
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Update actualiza un usuario
func (r *userRepository) Update(user *models.User) error {
	if user == nil {
		return errors.New("user cannot be nil")
	}

	now := time.Now()
	user.UpdatedAt = &now

	return r.db.Save(user).Error
}

// UpdateLastLogin actualiza la fecha de último login
func (r *userRepository) UpdateLastLogin(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"last_login": now,
			"updated_at": now,
		}).Error
}

// UpdateOnboardStatus actualiza el estado de onboarding
func (r *userRepository) UpdateOnboardStatus(id string, status models.OnboardStatus) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(map[string]interface{}{
			"onboard_status": status,
			"updated_at":     now,
		}).Error
}

// UpdatePersonID actualiza el ID de persona
func (r *userRepository) UpdatePersonID(id string, personID string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"updated_at": now,
	}

	if personID == "" {
		updates["person_id"] = nil
	} else {
		updates["person_id"] = personID
	}

	return r.db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Updates(updates).Error
}

// Delete realiza soft delete de un usuario
func (r *userRepository) Delete(id string) error {
	if id == "" {
		return errors.New("id cannot be empty")
	}

	now := time.Now()
	return r.db.Model(&models.User{}).
		Where("id = ? AND deleted_at IS NULL", id).
		Update("deleted_at", now).Error
}

// List obtiene una lista paginada de usuarios con filtros
func (r *userRepository) List(page, perPage int, filter *UserFilter) ([]*models.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}

	query := r.db.Model(&models.User{}).
		Where("users.deleted_at IS NULL").
		Preload("Account")

	// Aplicar filtros
	if filter != nil {
		if filter.OnboardStatus != nil {
			query = query.Where("users.onboard_status = ?", *filter.OnboardStatus)
		}

		if filter.HasPerson != nil {
			if *filter.HasPerson {
				query = query.Where("users.person_id IS NOT NULL")
			} else {
				query = query.Where("users.person_id IS NULL")
			}
		}

		if filter.AccountStatus != nil {
			query = query.Joins("JOIN accounts ON accounts.id = users.account_id").
				Where("accounts.status = ? AND accounts.deleted_at IS NULL", *filter.AccountStatus)
		}

		if filter.Search != "" {
			query = query.Joins("JOIN accounts ON accounts.id = users.account_id").
				Where("accounts.email ILIKE ? AND accounts.deleted_at IS NULL", "%"+filter.Search+"%")
		}
	}

	// Contar total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Obtener datos paginados
	var users []*models.User
	offset := (page - 1) * perPage
	err := query.Offset(offset).Limit(perPage).
		Order("users.created_at DESC").
		Find(&users).Error

	return users, total, err
}
