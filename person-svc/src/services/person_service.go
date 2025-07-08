package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/person/src/dto"
	"github.com/rem-gestion/api-suite/person/src/models"
	"github.com/rem-gestion/api-suite/person/src/repository"
	rerrors "github.com/rem-gestion/rem-common/errors"
	addresspb "github.com/rem-gestion/rem-common/protos/address/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

/* ───────────────────── struct & ctor ─────────────────────── */

type PersonService struct {
	repo    *repository.PersonRepo
	addrCli addresspb.AddressServiceClient
	lg      *zap.Logger
	timeout time.Duration
}

func New(repo *repository.PersonRepo, addrConn *grpc.ClientConn, lg *zap.Logger) *PersonService {
	return &PersonService{
		repo:    repo,
		addrCli: addresspb.NewAddressServiceClient(addrConn),
		lg:      lg.Named("service"),
		timeout: 3 * time.Second,
	}
}

/* ─────────────────────── helpers ─────────────────────────── */

func (s *PersonService) validateCreate(in dto.CreatePersonDTO) error {
	if in.AddressID != nil && in.Address != nil {
		return &rerrors.BadRequestError{Msg: "address_id y address_payload son mutuamente excluyentes"}
	}
	switch in.Type {
	case "individual":
		if in.FirstName == nil || in.LastName == nil {
			return &rerrors.ValidationError{Msg: "falta first_name o last_name"}
		}
	case "company":
		if in.LegalName == nil {
			return &rerrors.ValidationError{Msg: "falta legal_name"}
		}
	default:
		return &rerrors.BadRequestError{Msg: "type debe ser individual o company"}
	}
	return nil
}

func (s *PersonService) createRemoteAddress(ctx context.Context, a *dto.AddressPayloadDTO) (*uuid.UUID, error) {
	req := &addresspb.CreateAddressRequest{
		Address: &addresspb.Address{
			Street:  a.Street,
			Number:  int32(a.Number),
			City:    a.City,
			State:   a.State,
			Zip:     a.Zip,
			Country: a.Country,
			Floor:   valueOrEmpty(a.Floor),
			Unit:    valueOrEmpty(a.Unit),
		},
	}

	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.addrCli.Create(ctx, req) //  ←  método Create
	if err != nil {
		st, _ := status.FromError(err)
		return nil, &rerrors.InternalServerError{Msg: "address-svc: " + st.Message()}
	}
	id := uuid.MustParse(res.Address.Id)
	return &id, nil
}

func valueOrEmpty(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

/* ───────────────────── Personas CRUD ─────────────────────── */

func (s *PersonService) Create(in dto.CreatePersonDTO) (*models.Person, error) {
	if err := s.validateCreate(in); err != nil {
		return nil, err
	}

	// --- Address ------------------------------------------------
	var addrID *uuid.UUID
	if in.AddressID != nil {
		addrID = in.AddressID
	} else if in.Address != nil {
		id, err := s.createRemoteAddress(context.Background(), in.Address)
		if err != nil {
			return nil, err
		}
		addrID = id
	}

	// --- map DTO -> model --------------------------------------
	p := &models.Person{
		ID:        uuid.New(),
		Type:      models.PersonType(in.Type),
		AddressID: addrID,
		AvatarURL: in.AvatarURL,
		CreatedAt: time.Now(),
	}
	if in.Sexo != nil {
		sexo := models.Sexo(*in.Sexo)
		p.Sexo = &sexo
	}

	switch in.Type {
	case "individual":
		p.Individual = &models.Individual{
			PersonID:  p.ID,
			FirstName: *in.FirstName,
			LastName:  *in.LastName,
			DNI:       valueOrEmpty(in.DNI),
			CreatedAt: time.Now(),
		}
	case "company":
		p.Company = &models.Company{
			PersonID:    p.ID,
			LegalName:   *in.LegalName,
			CUIT:        valueOrEmpty(in.CUIT),
			SocietyType: valueOrEmpty(in.SocietyType),
			CreatedAt:   time.Now(),
		}
	}

	// --- persist persona ---------------------------------------
	out, err := s.repo.Create(p)
	if err != nil {
		return nil, err
	}

	// --- contactos iniciales -----------------------------------
	for _, c := range in.Contacts {
		id := out.ID
		c.PersonaID = &id
		_, err := s.AddContact(c)
		if err != nil {
			s.lg.Warn("contact init failed", zap.Error(err))
		}
	}

	return out, nil
}

func (s *PersonService) Get(id string) (*models.Person, error) {
	return s.repo.Get(id)
}

func (s *PersonService) List(page, per int, search string, filter *models.PersonType) ([]models.Person, int64, error) {
	if per <= 0 {
		per = 20
	}
	offset := (page - 1) * per
	return s.repo.List(search, filter, per, offset)
}

func (s *PersonService) Update(id string, in dto.UpdatePersonDTO) (*models.Person, error) {
	p, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	// address logic
	if in.AddressID != nil {
		p.AddressID = in.AddressID
	} else if in.Address != nil {
		id, err := s.createRemoteAddress(context.Background(), in.Address)
		if err != nil {
			return nil, err
		}
		p.AddressID = id
	}

	if in.AvatarURL != nil {
		p.AvatarURL = in.AvatarURL
	}
	if in.Sexo != nil {
		sexo := models.Sexo(*in.Sexo)
		p.Sexo = &sexo
	}

	// subtype updates
	if p.Type == models.PersonIndividual && p.Individual != nil {
		setIf := func(src *string, dst *string) {
			if src != nil {
				*dst = *src
			}
		}
		setIf(in.FirstName, &p.Individual.FirstName)
		setIf(in.LastName, &p.Individual.LastName)
		setIf(in.DNI, &p.Individual.DNI)
	} else if p.Type == models.PersonCompany && p.Company != nil {
		setIf := func(src *string, dst *string) {
			if src != nil {
				*dst = *src
			}
		}
		setIf(in.LegalName, &p.Company.LegalName)
		setIf(in.CUIT, &p.Company.CUIT)
		setIf(in.SocietyType, &p.Company.SocietyType)
	}

	now := time.Now()
	p.UpdatedAt = &now

	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PersonService) Delete(id string) error {
	return s.repo.Delete(id)
}

/* ───────────────────── Contactos CRUD ─────────────────────── */

func (s *PersonService) AddContact(in dto.CreateContactoDTO) (*models.Contacto, error) {
	if in.Tipo != "email" && in.Tipo != "phone" && in.Tipo != "whatsapp" {
		return nil, &rerrors.ValidationError{Msg: "tipo inválido"}
	}
	c := &models.Contacto{
		ID:        uuid.New(),
		PersonaID: *in.PersonaID,
		Tipo:      in.Tipo,
		Dato:      in.Dato,
		IsPrimary: in.IsPrimary,
		CreatedAt: time.Now(),
	}
	return s.repo.AddContact(c)
}

func (s *PersonService) UpdateContact(id string, in dto.UpdateContactoDTO) (*models.Contacto, error) {
	changes := map[string]any{}
	if in.Dato != nil {
		changes["dato"] = *in.Dato
	}
	if in.IsPrimary != nil {
		changes["is_primary"] = *in.IsPrimary
	}
	if len(changes) == 0 {
		return nil, &rerrors.BadRequestError{Msg: "no hay campos para actualizar"}
	}
	return s.repo.UpdateContact(id, changes)
}

func (s *PersonService) DeleteContact(id string) error {
	return s.repo.DeleteContact(id)
}

func (s *PersonService) ListContacts(personID string) ([]models.Contacto, error) {
	return s.repo.ListContacts(personID)
}

/* ─────────────────────── NUEVOS ENDPOINTS ─────────────────────── */

// GetFullPerson obtiene una persona con dirección completa expandida y contactos
func (s *PersonService) GetFullPerson(id string) (*dto.FullPersonResponse, error) {
	// Obtener la persona con sus relaciones
	person, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	response := &dto.FullPersonResponse{
		Person: person,
	}

	// Si tiene address_id, obtener la dirección completa del address-svc
	if person.AddressID != nil {
		ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
		defer cancel()

		req := &addresspb.GetAddressRequest{
			Id: person.AddressID.String(),
		}

		res, err := s.addrCli.Get(ctx, req)
		if err != nil {
			// Log error pero no fallar completamente
			s.lg.Warn("failed to fetch address",
				zap.String("address_id", person.AddressID.String()),
				zap.Error(err))
		} else {
			response.Address = &dto.AddressResponse{
				ID:      res.Address.Id,
				Floor:   res.Address.Floor,
				Unit:    res.Address.Unit,
				Street:  res.Address.Street,
				Number:  res.Address.Number,
				City:    res.Address.City,
				State:   res.Address.State,
				Zip:     res.Address.Zip,
				Country: res.Address.Country,
			}
		}
	}

	return response, nil
}

// GetPrimaryContact obtiene solo el contacto primario de una persona
func (s *PersonService) GetPrimaryContact(personID string) (*models.Contacto, error) {
	contacts, err := s.repo.ListContacts(personID)
	if err != nil {
		return nil, err
	}

	// Buscar el contacto primario
	for _, contact := range contacts {
		if contact.IsPrimary {
			return &contact, nil
		}
	}

	// Si no hay contacto primario, devolver error específico
	return nil, &rerrors.NotFoundError{Msg: "no se encontró contacto primario para esta persona"}
}

// BulkCreate crea múltiples personas en lote
func (s *PersonService) BulkCreate(requests []dto.CreatePersonDTO) (*dto.BulkCreateResponse, error) {
	response := &dto.BulkCreateResponse{
		Success: make([]dto.BulkPersonResult, 0),
		Errors:  make([]dto.BulkErrorResult, 0),
		Total:   len(requests),
	}

	for i, req := range requests {
		person, err := s.Create(req)
		if err != nil {
			response.Errors = append(response.Errors, dto.BulkErrorResult{
				Index: i,
				Error: err.Error(),
			})
		} else {
			response.Success = append(response.Success, dto.BulkPersonResult{
				Index:  i,
				Person: person,
			})
		}
	}

	s.lg.Info("bulk create completed",
		zap.Int("total", response.Total),
		zap.Int("success", len(response.Success)),
		zap.Int("errors", len(response.Errors)))

	return response, nil
}
