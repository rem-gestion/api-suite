package repository

import (
	model "github.com/rem-gestion/api-suite/person/src/models"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type PersonRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func New(db *gorm.DB, lg *zap.Logger) *PersonRepo {
	return &PersonRepo{db: db, lg: lg.Named("repo")}
}

/* ---------- CRUD ---------- */

func (r *PersonRepo) Create(p *model.Person) (*model.Person, error) {
	// Begin transaction for person creation with sub-entities
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create the main Person record
	if err := tx.Create(p).Error; err != nil {
		tx.Rollback()
		r.lg.Error("person creation failed", zap.Error(err))
		return nil, err
	}

	// Create sub-entity based on type
	if p.Type == model.PersonIndividual && p.Individual != nil {
		p.Individual.PersonID = p.ID
		if err := tx.Create(p.Individual).Error; err != nil {
			tx.Rollback()
			r.lg.Error("individual creation failed", zap.Error(err))
			return nil, err
		}
	} else if p.Type == model.PersonCompany && p.Company != nil {
		p.Company.PersonID = p.ID
		if err := tx.Create(p.Company).Error; err != nil {
			tx.Rollback()
			r.lg.Error("company creation failed", zap.Error(err))
			return nil, err
		}
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		r.lg.Error("transaction commit failed", zap.Error(err))
		return nil, err
	}

	r.lg.Info("person created", zap.String("id", p.ID.String()))
	return p, nil
}

func (r *PersonRepo) Get(id string) (*model.Person, error) {
	var p model.Person

	// Load person with sub-entities
	if err := r.db.Preload("Individual").Preload("Company").Preload("Contactos").
		First(&p, "id = ?", id).Error; err != nil {
		r.lg.Warn("get person failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return &p, nil
}

func (r *PersonRepo) List(personType string, limit, offset int) ([]model.Person, error) {
	var res []model.Person
	q := r.db.Preload("Individual").Preload("Company").Preload("Contactos").
		Limit(limit).Offset(offset)

	if personType != "" {
		q = q.Where("type = ?", personType)
	}

	if err := q.Find(&res).Error; err != nil {
		r.lg.Warn("list persons failed", zap.Error(err))
		return nil, err
	}
	return res, nil
}

func (r *PersonRepo) Update(p *model.Person) error {
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update main Person record
	if err := tx.Save(p).Error; err != nil {
		tx.Rollback()
		r.lg.Error("person update failed", zap.Error(err))
		return err
	}

	// Update sub-entity based on type
	if p.Type == model.PersonIndividual && p.Individual != nil {
		if err := tx.Save(p.Individual).Error; err != nil {
			tx.Rollback()
			r.lg.Error("individual update failed", zap.Error(err))
			return err
		}
	} else if p.Type == model.PersonCompany && p.Company != nil {
		if err := tx.Save(p.Company).Error; err != nil {
			tx.Rollback()
			r.lg.Error("company update failed", zap.Error(err))
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		r.lg.Error("transaction commit failed", zap.Error(err))
		return err
	}

	r.lg.Info("person updated", zap.String("id", p.ID.String()))
	return nil
}

func (r *PersonRepo) Delete(id string) error {
	// Soft delete will cascade to sub-entities due to GORM constraints
	if err := r.db.Delete(&model.Person{}, "id = ?", id).Error; err != nil {
		r.lg.Error("delete person failed", zap.String("id", id), zap.Error(err))
		return err
	}
	r.lg.Info("person deleted", zap.String("id", id))
	return nil
}
