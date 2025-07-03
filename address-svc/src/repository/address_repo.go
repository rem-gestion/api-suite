// repository/address_repo.go
package repository

import (
	model "github.com/rem-gestion/api-suite/address/src/models"
	"gorm.io/gorm"
)

type AddressRepo struct{ db *gorm.DB }

func New(db *gorm.DB) *AddressRepo { return &AddressRepo{db} }

func (r *AddressRepo) Create(a *model.Address) error { return r.db.Create(a).Error }
func (r *AddressRepo) Get(id string) (*model.Address, error) {
	var a model.Address
	if err := r.db.First(&a, "id = ?", id).Error; err != nil {
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
	return res, q.Find(&res).Error
}
func (r *AddressRepo) Update(a *model.Address) error { return r.db.Save(a).Error }
func (r *AddressRepo) Delete(id string) error {
	return r.db.Delete(&model.Address{}, "id = ?", id).Error
}
