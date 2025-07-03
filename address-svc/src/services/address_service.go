package services

import (
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/address/src/dto"
	model "github.com/rem-gestion/api-suite/address/src/models"
	"github.com/rem-gestion/api-suite/address/src/repository"
)

type AddressService struct{ repo *repository.AddressRepo }

func New(r *repository.AddressRepo) *AddressService { return &AddressService{r} }

func copyUpdate(a *model.Address, in dto.AddressUpdate) {
	if in.Floor != nil {
		a.Floor = *in.Floor
	}
	if in.Unit != nil {
		a.Unit = *in.Unit
	}
	if in.Street != nil {
		a.Street = *in.Street
	}
	if in.Number != nil {
		a.Number = *in.Number
	}
	if in.City != nil {
		a.City = *in.City
	}
	if in.State != nil {
		a.State = *in.State
	}
	if in.Zip != nil {
		a.Zip = *in.Zip
	}
	if in.Country != nil {
		a.Country = *in.Country
	}
}

func (s *AddressService) Create(in dto.AddressCreate) (*model.Address, error) {
	a := model.Address{
		ID:      uuid.NewString(),
		Floor:   in.Floor,
		Unit:    in.Unit,
		Street:  in.Street,
		Number:  in.Number,
		City:    in.City,
		State:   in.State,
		Zip:     in.Zip,
		Country: in.Country,
	}
	if err := s.repo.Create(&a); err != nil {
		return nil, err
	}
	return &a, nil
}

func (s *AddressService) Get(id string) (*model.Address, error) {
	return s.repo.Get(id)
}

func (s *AddressService) Update(id string, in dto.AddressUpdate) (*model.Address, error) {
	a, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	copyUpdate(a, in)

	if err := s.repo.Update(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *AddressService) Delete(id string) error {
	return s.repo.Delete(id)
}
