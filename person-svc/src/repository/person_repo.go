package repository

import (
	"fmt"
	"strings"

	"github.com/rem-gestion/api-suite/person/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

/* ───────────────── struct & ctor ────────────────── */

type PersonRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func NewPersonRepo(db *gorm.DB, lg *zap.Logger) *PersonRepo {
	return &PersonRepo{db: db, lg: lg}
}

/* ───────────────── Personas CRUD ─────────────────── */

func (r *PersonRepo) Create(p *models.Person) (*models.Person, error) {
	err := r.withTx(func(tx *gorm.DB) error {
		return tx.Create(p).Error
	})
	return p, err
}

func (r *PersonRepo) Get(id string) (*models.Person, error) {
	var p models.Person
	err := r.preloads(r.db).
		First(&p, "id = ?", id).Error
	return &p, err
}

func (r *PersonRepo) List(search string, pType *models.PersonType, limit, offset int) ([]models.Person, int64, error) {
	var list []models.Person
	var total int64

	q := r.preloads(r.db.Model(&models.Person{})).
		Limit(limit).Offset(offset)

	if pType != nil {
		q = q.Where("type = ?", *pType)
	}

	if s := strings.TrimSpace(search); s != "" {
		sLike := fmt.Sprintf("%%%s%%", s)
		q = q.Joins(`
			LEFT JOIN individual ON individual.person_id = person.id
			LEFT JOIN company    ON company.person_id    = person.id`).
			Where(`
				lower(individual.first_name) LIKE lower(?) OR
				lower(individual.last_name)  LIKE lower(?) OR
				individual.dni               = ?          OR
				lower(company.legal_name)    LIKE lower(?) OR
				company.cuit                 = ?`,
				sLike, sLike, s, sLike, s)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := q.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *PersonRepo) Update(p *models.Person) error {
	return r.withTx(func(tx *gorm.DB) error {
		return tx.Session(&gorm.Session{FullSaveAssociations: true}).Updates(p).Error
	})
}

func (r *PersonRepo) Delete(id string) error {
	return r.db.Delete(&models.Person{}, "id = ?", id).Error
}

/* ───────────────── Contactos CRUD ────────────────── */

func (r *PersonRepo) AddContact(c *models.Contacto) (*models.Contacto, error) {
	err := r.withTx(func(tx *gorm.DB) error {
		if c.IsPrimary {
			if err := tx.Model(&models.Contacto{}).
				Where("persona_id = ?", c.PersonaID).
				Update("is_primary", false).Error; err != nil {
				return err
			}
		}
		return tx.Create(c).Error
	})
	return c, err
}

func (r *PersonRepo) UpdateContact(id string, up map[string]any) (*models.Contacto, error) {
	var contact models.Contacto
	err := r.withTx(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&contact, "id = ?", id).Error; err != nil {
			return err
		}
		if v, ok := up["is_primary"]; ok && v.(bool) {
			if err := tx.Model(&models.Contacto{}).
				Where("persona_id = ?", contact.PersonaID).
				Update("is_primary", false).Error; err != nil {
				return err
			}
		}
		return tx.Model(&contact).Updates(up).Error
	})
	return &contact, err
}

func (r *PersonRepo) DeleteContact(id string) error {
	return r.db.Delete(&models.Contacto{}, "id = ?", id).Error
}

func (r *PersonRepo) ListContacts(personID string) ([]models.Contacto, error) {
	var list []models.Contacto
	err := r.db.Where("persona_id = ?", personID).
		Order("is_primary DESC, created_at DESC").
		Find(&list).Error
	return list, err
}

/* ───────────────── helpers internos ───────────────── */

func (r *PersonRepo) preloads(q *gorm.DB) *gorm.DB {
	return q.
		Preload("Individual").
		Preload("Company").
		Preload("Contactos")
}

// withTx envuelve una fn en BEGIN / COMMIT / ROLLBACK.
// Devuelve el error que retorne la fn o de Commit.
func (r *PersonRepo) withTx(fn func(*gorm.DB) error) error {
	tx := r.db.Begin()
	defer func() {
		if rec := recover(); rec != nil {
			tx.Rollback()
		}
	}()
	if err := fn(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return &rerrors.InternalServerError{Msg: err.Error()}
	}
	return nil
}
