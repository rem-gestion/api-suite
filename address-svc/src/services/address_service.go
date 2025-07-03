package services

import (
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/address/src/dto"
	model "github.com/rem-gestion/api-suite/address/src/models"
	"github.com/rem-gestion/api-suite/address/src/repository"
	"go.uber.org/zap"
)

type AddressService struct {
	repo *repository.AddressRepo
	lg   *zap.Logger
}

func New(r *repository.AddressRepo, lg *zap.Logger) *AddressService {
	return &AddressService{repo: r, lg: lg.Named("service")}
}

/* ---------- helpers ---------- */

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

/* ---------- CRUD ---------- */

func (s *AddressService) Create(in dto.AddressCreate) (*model.Address, error) {
	s.lg.Debug("create request", zap.Any("payload", in))

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

	out, err := s.repo.Create(&a)
	if err != nil {
		s.lg.Error("create failed", zap.Error(err))
		return nil, err
	}

	s.lg.Info("create ok", zap.String("id", out.ID))
	return out, nil
}

func (s *AddressService) Get(id string) (*model.Address, error) {
	a, err := s.repo.Get(id)
	if err != nil {
		s.lg.Warn("get failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return a, nil
}

func (s *AddressService) Update(id string, in dto.AddressUpdate) (*model.Address, error) {
	s.lg.Debug("update request", zap.String("id", id), zap.Any("payload", in))

	a, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	copyUpdate(a, in)

	if err := s.repo.Update(a); err != nil {
		s.lg.Warn("update rejected (immutable)", zap.String("id", id))
		return nil, err
	}

	s.lg.Info("update noop (immutable, new record required)", zap.String("id", id))
	return a, nil
}

func (s *AddressService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.lg.Info("delete ok", zap.String("id", id))
	return nil
}
