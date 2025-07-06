package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/person/src/dto"
	model "github.com/rem-gestion/api-suite/person/src/models"
	"github.com/rem-gestion/api-suite/person/src/repository"
	"go.uber.org/zap"
)

type PersonService struct {
	repo *repository.PersonRepo
	lg   *zap.Logger
}

func New(r *repository.PersonRepo, lg *zap.Logger) *PersonService {
	return &PersonService{repo: r, lg: lg.Named("service")}
}

/* ---------- helpers ---------- */

func copyUpdate(p *model.Person, in dto.UpdatePersonDTO) {
	if in.AvatarURL != nil {
		p.AvatarURL = in.AvatarURL
	}
	if in.Sexo != nil {
		sexo := model.Sexo(*in.Sexo)
		p.Sexo = &sexo
	}

	// Update sub-entity based on type
	if p.Type == model.PersonIndividual && p.Individual != nil {
		if in.FirstName != nil {
			p.Individual.FirstName = *in.FirstName
		}
		if in.LastName != nil {
			p.Individual.LastName = *in.LastName
		}
		if in.DNI != nil {
			p.Individual.DNI = *in.DNI
		}
	} else if p.Type == model.PersonCompany && p.Company != nil {
		if in.LegalName != nil {
			p.Company.LegalName = *in.LegalName
		}
		if in.CUIT != nil {
			p.Company.CUIT = *in.CUIT
		}
		if in.SocietyType != nil {
			p.Company.SocietyType = *in.SocietyType
		}
	}
}

/* ---------- CRUD ---------- */

func (s *PersonService) Create(in dto.CreatePersonDTO) (*model.Person, error) {
	s.lg.Debug("create request", zap.Any("payload", in))

	p := model.Person{
		ID:        uuid.New(),
		Type:      model.PersonType(in.Type),
		AddressID: in.AddressID,
		AvatarURL: in.AvatarURL,
		CreatedAt: time.Now(),
	}

	if in.Sexo != nil {
		sexo := model.Sexo(*in.Sexo)
		p.Sexo = &sexo
	}

	// Create sub-entity based on type
	if in.Type == "individual" {
		p.Individual = &model.Individual{
			FirstName: *in.FirstName,
			LastName:  *in.LastName,
			DNI:       *in.DNI,
			CreatedAt: time.Now(),
		}
	} else if in.Type == "company" {
		p.Company = &model.Company{
			LegalName:   *in.LegalName,
			CUIT:        *in.CUIT,
			SocietyType: *in.SocietyType,
			CreatedAt:   time.Now(),
		}
	}

	out, err := s.repo.Create(&p)
	if err != nil {
		s.lg.Error("create failed", zap.Error(err))
		return nil, err
	}

	s.lg.Info("create ok", zap.String("id", out.ID.String()))
	return out, nil
}

func (s *PersonService) Get(id string) (*model.Person, error) {
	p, err := s.repo.Get(id)
	if err != nil {
		s.lg.Warn("get failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return p, nil
}

func (s *PersonService) Update(id string, in dto.UpdatePersonDTO) (*model.Person, error) {
	s.lg.Debug("update request", zap.String("id", id), zap.Any("payload", in))

	p, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	copyUpdate(p, in)
	now := time.Now()
	p.UpdatedAt = &now

	if err := s.repo.Update(p); err != nil {
		s.lg.Error("update failed", zap.String("id", id), zap.Error(err))
		return nil, err
	}

	s.lg.Info("update ok", zap.String("id", id))
	return p, nil
}

func (s *PersonService) Delete(id string) error {
	if err := s.repo.Delete(id); err != nil {
		return err
	}
	s.lg.Info("delete ok", zap.String("id", id))
	return nil
}
