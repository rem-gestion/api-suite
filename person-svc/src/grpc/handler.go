package grpcserver

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/person/src/dto"
	model "github.com/rem-gestion/api-suite/person/src/models"
	"github.com/rem-gestion/api-suite/person/src/services"
	personpb "github.com/rem-gestion/rem-common/protos/person/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

/* ─────────────────── helpers de conversión ─────────────────── */

// ------------ enums ------------
func enumPersonType(t model.PersonType) personpb.PersonType {
	switch t {
	case model.PersonIndividual:
		return personpb.PersonType_INDIVIDUAL
	case model.PersonCompany:
		return personpb.PersonType_COMPANY
	default:
		return personpb.PersonType_PERSON_TYPE_UNSPECIFIED
	}
}
func enumSexo(s *model.Sexo) personpb.Sexo {
	if s == nil {
		return personpb.Sexo_SEXO_UNSPECIFIED
	}
	if *s == model.SexoMasculino {
		return personpb.Sexo_MASCULINO
	}
	if *s == model.SexoFemenino {
		return personpb.Sexo_FEMENINO
	}
	return personpb.Sexo_SEXO_UNSPECIFIED
}
func enumContactType(t string) personpb.ContactType {
	switch t {
	case "email":
		return personpb.ContactType_EMAIL
	case "phone":
		return personpb.ContactType_PHONE
	case "whatsapp":
		return personpb.ContactType_WHATSAPP
	default:
		return personpb.ContactType_CONTACT_TYPE_UNSPECIFIED
	}
}

// ------------ util ------------
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
func ts(t time.Time) *timestamppb.Timestamp { return timestamppb.New(t) }
func tsPtr(t *time.Time) *timestamppb.Timestamp {
	if t == nil {
		return nil
	}
	return timestamppb.New(*t)
}
func uuidPtr(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

// ------------ proto → DTO ------------
func toCreateDTO(r *personpb.CreatePersonRequest) dto.CreatePersonDTO {
	out := dto.CreatePersonDTO{
		AvatarURL: strPtr(r.AvatarUrl),
	}

	switch r.Type {
	case personpb.PersonType_INDIVIDUAL:
		out.Type = "individual"
	case personpb.PersonType_COMPANY:
		out.Type = "company"
	default:
		out.Type = ""
	}
	if r.Sexo != personpb.Sexo_SEXO_UNSPECIFIED {
		s := r.Sexo.String()
		out.Sexo = &s
	}
	// address oneof
	switch a := r.Address.(type) {
	case *personpb.CreatePersonRequest_AddressId:
		if id, err := uuid.Parse(a.AddressId); err == nil {
			out.AddressID = &id
		}
	case *personpb.CreatePersonRequest_AddressPayload:
		ap := a.AddressPayload
		out.Address = &dto.AddressPayloadDTO{
			Floor:   strPtr(ap.Floor),
			Unit:    strPtr(ap.Unit),
			Street:  ap.Street,
			Number:  int(ap.Number),
			City:    ap.City,
			State:   ap.State,
			Zip:     ap.Zip,
			Country: ap.Country,
		}
	}
	// details oneof
	switch d := r.Details.(type) {
	case *personpb.CreatePersonRequest_Individual:
		i := d.Individual
		out.FirstName = strPtr(i.FirstName)
		out.LastName = strPtr(i.LastName)
		out.DNI = strPtr(i.Dni)
	case *personpb.CreatePersonRequest_Company:
		c := d.Company
		out.LegalName = strPtr(c.LegalName)
		out.CUIT = strPtr(c.Cuit)
		out.SocietyType = strPtr(c.SocietyType)
	}
	return out
}

func toUpdateDTO(r *personpb.UpdatePersonRequest) dto.UpdatePersonDTO {
	out := dto.UpdatePersonDTO{
		AvatarURL: strPtr(r.AvatarUrl),
	}
	if r.Sexo != personpb.Sexo_SEXO_UNSPECIFIED {
		s := r.Sexo.String()
		out.Sexo = &s
	}
	switch a := r.Address.(type) {
	case *personpb.UpdatePersonRequest_AddressId:
		if id, err := uuid.Parse(a.AddressId); err == nil {
			out.AddressID = &id
		}
	case *personpb.UpdatePersonRequest_AddressPayload:
		ap := a.AddressPayload
		out.Address = &dto.AddressPayloadDTO{
			Floor:   strPtr(ap.Floor),
			Unit:    strPtr(ap.Unit),
			Street:  ap.Street,
			Number:  int(ap.Number),
			City:    ap.City,
			State:   ap.State,
			Zip:     ap.Zip,
			Country: ap.Country,
		}
	}
	switch d := r.Details.(type) {
	case *personpb.UpdatePersonRequest_Individual:
		i := d.Individual
		out.FirstName = strPtr(i.FirstName)
		out.LastName = strPtr(i.LastName)
		out.DNI = strPtr(i.Dni)
	case *personpb.UpdatePersonRequest_Company:
		c := d.Company
		out.LegalName = strPtr(c.LegalName)
		out.CUIT = strPtr(c.Cuit)
		out.SocietyType = strPtr(c.SocietyType)
	}
	return out
}

// ------------ model → proto ------------
func contactToProto(c *model.Contacto) *personpb.Contact {
	return &personpb.Contact{
		Id:        c.ID.String(),
		Tipo:      enumContactType(c.Tipo),
		Dato:      c.Dato,
		IsPrimary: c.IsPrimary,
		CreatedAt: ts(c.CreatedAt),
	}
}

func personToProto(m *model.Person, withContacts bool) *personpb.PersonResponse {
	p := &personpb.PersonResponse{
		Id:        m.ID.String(),
		Type:      enumPersonType(m.Type),
		AvatarUrl: strVal(m.AvatarURL),
		Sexo:      enumSexo(m.Sexo),
		CreatedAt: ts(m.CreatedAt),
		UpdatedAt: tsPtr(m.UpdatedAt),
		AddressId: uuidPtr(m.AddressID),
	}
	if m.Individual != nil {
		p.Details = &personpb.PersonResponse_Individual{
			Individual: &personpb.Individual{
				FirstName: m.Individual.FirstName,
				LastName:  m.Individual.LastName,
				Dni:       m.Individual.DNI,
			},
		}
	}
	if m.Company != nil {
		p.Details = &personpb.PersonResponse_Company{
			Company: &personpb.Company{
				LegalName:   m.Company.LegalName,
				Cuit:        m.Company.CUIT,
				SocietyType: m.Company.SocietyType,
			},
		}
	}
	if withContacts {
		for _, c := range m.Contactos {
			p.Contacts = append(p.Contacts, contactToProto(&c))
		}
	}
	return p
}
func strVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

/* ───────────────────── server struct ─────────────────────── */

type Handler struct {
	personpb.UnimplementedPersonServiceServer
	svc *services.PersonService
}

func New(s *services.PersonService) *Handler { return &Handler{svc: s} }

/* ───────────────────── Personas CRUD ─────────────────────── */

func (h *Handler) CreatePerson(ctx context.Context, r *personpb.CreatePersonRequest) (*personpb.PersonResponse, error) {
	p, err := h.svc.Create(toCreateDTO(r))
	if err != nil {
		return nil, err
	}
	return personToProto(p, true), nil
}

func (h *Handler) GetPerson(ctx context.Context, r *personpb.GetPersonRequest) (*personpb.PersonResponse, error) {
	p, err := h.svc.Get(r.GetId())
	if err != nil {
		return nil, err
	}
	return personToProto(p, true), nil
}

func (h *Handler) ListPersons(ctx context.Context, r *personpb.ListPersonsRequest) (*personpb.ListPersonsResponse, error) {
	var filter *model.PersonType
	if r.TypeFilter != personpb.PersonType_PERSON_TYPE_UNSPECIFIED {
		ft := model.PersonType(r.TypeFilter.String())
		filter = &ft
	}

	data, total, err := h.svc.List(int(r.Page), int(r.PerPage), r.Search, filter)
	if err != nil {
		return nil, err
	}

	res := &personpb.ListPersonsResponse{Total: uint32(total)}
	for _, m := range data {
		res.Data = append(res.Data, personToProto(&m, r.WithContacts))
	}
	return res, nil
}

func (h *Handler) UpdatePerson(ctx context.Context, r *personpb.UpdatePersonRequest) (*personpb.PersonResponse, error) {
	p, err := h.svc.Update(r.GetId(), toUpdateDTO(r))
	if err != nil {
		return nil, err
	}
	return personToProto(p, true), nil
}

func (h *Handler) DeletePerson(ctx context.Context, r *personpb.DeletePersonRequest) (*emptypb.Empty, error) {
	if err := h.svc.Delete(r.GetId()); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

/* ───────────────────── Contactos CRUD ─────────────────────── */

func (h *Handler) AddContact(ctx context.Context, r *personpb.AddContactRequest) (*personpb.ContactResponse, error) {
	id, err := uuid.Parse(r.PersonId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "person_id inválido")
	}
	in := dto.CreateContactoDTO{
		PersonaID: &id,
		Tipo:      r.Payload.Tipo.String(),
		Dato:      r.Payload.Dato,
		IsPrimary: r.Payload.IsPrimary,
	}
	c, err := h.svc.AddContact(in)
	if err != nil {
		return nil, err
	}
	return &personpb.ContactResponse{
		Id:        c.ID.String(),
		PersonId:  c.PersonaID.String(),
		Tipo:      enumContactType(c.Tipo),
		Dato:      c.Dato,
		IsPrimary: c.IsPrimary,
		CreatedAt: ts(c.CreatedAt),
	}, nil
}

func (h *Handler) UpdateContact(ctx context.Context, r *personpb.UpdateContactRequest) (*personpb.ContactResponse, error) {
	upd := dto.UpdateContactoDTO{
		Dato:      strPtr(r.Dato),
		IsPrimary: &r.IsPrimary,
	}
	c, err := h.svc.UpdateContact(r.Id, upd)
	if err != nil {
		return nil, err
	}
	return &personpb.ContactResponse{
		Id:        c.ID.String(),
		PersonId:  c.PersonaID.String(),
		Tipo:      enumContactType(c.Tipo),
		Dato:      c.Dato,
		IsPrimary: c.IsPrimary,
		CreatedAt: ts(c.CreatedAt),
	}, nil
}

func (h *Handler) DeleteContact(ctx context.Context, r *personpb.DeleteContactRequest) (*emptypb.Empty, error) {
	if err := h.svc.DeleteContact(r.Id); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) ListContacts(ctx context.Context, r *personpb.ListContactsRequest) (*personpb.ListContactsResponse, error) {
	contacts, err := h.svc.ListContacts(r.PersonId)
	if err != nil {
		return nil, err
	}
	out := &personpb.ListContactsResponse{}
	for _, c := range contacts {
		out.Data = append(out.Data, contactToProto(&c))
	}
	return out, nil
}
