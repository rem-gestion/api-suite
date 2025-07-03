package repository

import (
	"strings"

	model "github.com/rem-gestion/api-suite/address/src/models"
	rerrors "github.com/rem-gestion/rem-common/errors"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AddressRepo struct {
	db *gorm.DB
	lg *zap.Logger
}

func New(db *gorm.DB, lg *zap.Logger) *AddressRepo {
	return &AddressRepo{db: db, lg: lg.Named("repo")}
}

/* ---------- helpers ---------- */

func normalize(a *model.Address) (floor, unit, street, city, state, zip, country string) {
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

/* ---------- CRUD ---------- */

func (r *AddressRepo) Create(a *model.Address) (*model.Address, error) {
	// 1) intento de inserción
	if err := r.db.Create(a).Error; err != nil {
		r.lg.Warn("insert failed, trying dedup", zap.Error(err))

		f, u, s, c, st, z, co := normalize(a)
		var existing model.Address

		findErr := r.db.
			Where(`lower(trim(floor))   = ? AND
				   lower(trim(unit))    = ? AND
				   lower(trim(street))  = ? AND
				   number              = ? AND
				   lower(trim(city))    = ? AND
				   lower(trim(state))   = ? AND
				   trim(zip)            = ? AND
				   upper(trim(country)) = ?`,
				f, u, s, a.Number, c, st, z, co).
			First(&existing).Error

		if findErr == nil {
			r.lg.Info("deduplicated → returning existing address",
				zap.String("id", existing.ID))
			return &existing, nil
		}

		r.lg.Error("dedup query failed", zap.Error(findErr))
		return nil, err // devolvemos el error real de inserción
	}

	r.lg.Info("address created", zap.String("id", a.ID))
	return a, nil
}

func (r *AddressRepo) Get(id string) (*model.Address, error) {
	var a model.Address
	if err := r.db.First(&a, "id = ?", id).Error; err != nil {
		r.lg.Warn("get address failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return &a, nil
}

func (r *AddressRepo) List(city string, limit, offset int) ([]model.Address, error) {
	var res []model.Address
	q := r.db.Limit(limit).Offset(offset)
	if city != "" {
		q = q.Where("city = ?", city)
	}
	if err := q.Find(&res).Error; err != nil {
		r.lg.Warn("list addresses failed", zap.Error(err))
		return nil, err
	}
	return res, nil
}

// direcciones inmutables → 403
func (r *AddressRepo) Update(*model.Address) error {
	return &rerrors.ForbiddenError{Msg: "addresses are immutable; create a new record"}
}

func (r *AddressRepo) Delete(id string) error {
	if err := r.db.Delete(&model.Address{}, "id = ?", id).Error; err != nil {
		r.lg.Error("delete address failed", zap.String("id", id), zap.Error(err))
		return err
	}
	r.lg.Info("address deleted", zap.String("id", id))
	return nil
}
