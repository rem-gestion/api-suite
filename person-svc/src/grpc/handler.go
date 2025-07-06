// src/grpc/handler.go
package grpcserver

import (
	"context"

	pb "github.com/rem-gestion/api-suite/person/internal/pb"
	"github.com/rem-gestion/api-suite/person/src/dto"
	model "github.com/rem-gestion/api-suite/person/src/models"
	"github.com/rem-gestion/api-suite/person/src/services"
)

/* ─── helpers de conversión ─────────────────────────── */

func dtoToCreate(p *pb.Person) dto.CreatePersonDTO {
	dto := dto.CreatePersonDTO{
		Type:      p.Type,
		AvatarURL: stringPtr(p.AvatarUrl),
	}

	if p.Sexo != "" {
		dto.Sexo = &p.Sexo
	}

	// Set fields based on type
	if p.Type == "individual" {
		dto.FirstName = stringPtr(p.FirstName)
		dto.LastName = stringPtr(p.LastName)
		dto.DNI = stringPtr(p.Dni)
	} else if p.Type == "company" {
		dto.LegalName = stringPtr(p.LegalName)
		dto.CUIT = stringPtr(p.Cuit)
		dto.SocietyType = stringPtr(p.SocietyType)
	}

	return dto
}

func modelToProto(m *model.Person) *pb.Person {
	p := &pb.Person{
		Id:   m.ID.String(),
		Type: string(m.Type),
	}

	if m.AvatarURL != nil {
		p.AvatarUrl = *m.AvatarURL
	}
	if m.Sexo != nil {
		p.Sexo = string(*m.Sexo)
	}
	if m.AddressID != nil {
		p.AddressId = m.AddressID.String()
	}

	// Set sub-entity fields
	if m.Individual != nil {
		p.FirstName = m.Individual.FirstName
		p.LastName = m.Individual.LastName
		p.Dni = m.Individual.DNI
	}
	if m.Company != nil {
		p.LegalName = m.Company.LegalName
		p.Cuit = m.Company.CUIT
		p.SocietyType = m.Company.SocietyType
	}

	return p
}

func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

/* ─── implementación del servicio gRPC ───────────────── */

type Handler struct {
	pb.UnimplementedPersonServiceServer
	svc *services.PersonService
}

func New(svc *services.PersonService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(ctx context.Context, req *pb.CreatePersonRequest) (*pb.CreatePersonResponse, error) {
	person, err := h.svc.Create(dtoToCreate(req.GetPerson()))
	if err != nil {
		return nil, err
	}
	return &pb.CreatePersonResponse{Person: modelToProto(person)}, nil
}

func (h *Handler) Get(ctx context.Context, req *pb.GetPersonRequest) (*pb.GetPersonResponse, error) {
	person, err := h.svc.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	return &pb.GetPersonResponse{Person: modelToProto(person)}, nil
}
