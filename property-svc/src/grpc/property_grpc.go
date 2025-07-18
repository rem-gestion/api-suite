package grpc

import (
	"github.com/rem-gestion/api-suite/property/src/services"
	"go.uber.org/zap"
	// TODO: Importar el proto cuando esté generado
	// propertypb "github.com/rem-gestion/rem-common/protos/property/v1"
)

// PropertyGRPCHandler maneja las peticiones gRPC
type PropertyGRPCHandler struct {
	svc *services.PropertyService
	lg  *zap.Logger
	// TODO: Implementar UnimplementedPropertyServiceServer cuando se genere el proto
}

func New(svc *services.PropertyService, lg *zap.Logger) *PropertyGRPCHandler {
	return &PropertyGRPCHandler{
		svc: svc,
		lg:  lg.Named("grpc"),
	}
}

// TODO: Implementar métodos gRPC cuando se generen los protos
// Ejemplo de estructura que tendríamos:

/*
func (h *PropertyGRPCHandler) GetProperty(ctx context.Context, req *propertypb.GetPropertyRequest) (*propertypb.PropertyResponse, error) {
	property, err := h.svc.GetProperty(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	// Convertir DTO a proto response
	return &propertypb.PropertyResponse{
		Id: property.ID.String(),
		OwnerPersonId: property.OwnerPersonID.String(),
		AddressId: property.AddressID.String(),
		PropertyType: string(property.PropertyType),
		// ... otros campos
	}, nil
}

func (h *PropertyGRPCHandler) ListProperties(ctx context.Context, req *propertypb.ListPropertiesRequest) (*propertypb.ListPropertiesResponse, error) {
	// Implementar listado con paginación
}

func (h *PropertyGRPCHandler) CreateProperty(ctx context.Context, req *propertypb.CreatePropertyRequest) (*propertypb.PropertyResponse, error) {
	// Implementar creación
}
*/
