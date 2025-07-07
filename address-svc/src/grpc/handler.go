package grpcserver

import (
	"context"

	"github.com/rem-gestion/api-suite/address/src/dto"
	model "github.com/rem-gestion/api-suite/address/src/models"
	"github.com/rem-gestion/api-suite/address/src/services"
	addresspb "github.com/rem-gestion/rem-common/protos/address/v1"
)

/* ─── helpers de conversión ─────────────────────────── */

func dtoToCreate(a *addresspb.Address) dto.AddressCreate {
	return dto.AddressCreate{
		Floor:   a.Floor,
		Unit:    a.Unit,
		Street:  a.Street,
		Number:  int(a.Number),
		City:    a.City,
		State:   a.State,
		Zip:     a.Zip,
		Country: a.Country,
	}
}

func modelToProto(m *model.Address) *addresspb.Address {
	return &addresspb.Address{
		Id:      m.ID,
		Floor:   m.Floor,
		Unit:    m.Unit,
		Street:  m.Street,
		Number:  int32(m.Number),
		City:    m.City,
		State:   m.State,
		Zip:     m.Zip,
		Country: m.Country,
	}
}

/* ─── implementación del servicio gRPC ───────────────── */

type Handler struct {
	addresspb.UnimplementedAddressServiceServer
	svc *services.AddressService
}

func New(svc *services.AddressService) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(ctx context.Context, req *addresspb.CreateAddressRequest) (*addresspb.CreateAddressResponse, error) {
	addr, err := h.svc.Create(dtoToCreate(req.GetAddress()))
	if err != nil {
		return nil, err
	}
	return &addresspb.CreateAddressResponse{Address: modelToProto(addr)}, nil
}

func (h *Handler) Get(ctx context.Context, req *addresspb.GetAddressRequest) (*addresspb.GetAddressResponse, error) {
	addr, err := h.svc.Get(req.GetId())
	if err != nil {
		return nil, err
	}
	return &addresspb.GetAddressResponse{Address: modelToProto(addr)}, nil
}
