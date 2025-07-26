package server

import (
	"context"
	"fmt"
	"net"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/services"
	propertypb "github.com/rem-gestion/rem-common/protos/property/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PropertyGRPCServer struct {
	propertypb.UnimplementedPropertyServiceServer
	server          *grpc.Server
	propertyService *services.PropertyService
	logger          *zap.Logger
}

func NewPropertyGRPCServer(propertyService *services.PropertyService, logger *zap.Logger) *PropertyGRPCServer {
	s := &PropertyGRPCServer{
		server:          grpc.NewServer(),
		propertyService: propertyService,
		logger:          logger.Named("property-grpc-server"),
	}

	// Register the service
	propertypb.RegisterPropertyServiceServer(s.server, s)

	return s
}

// Implement PropertyServiceServer interface
func (s *PropertyGRPCServer) GetProperty(ctx context.Context, req *propertypb.GetPropertyRequest) (*propertypb.PropertyResponse, error) {
	s.logger.Debug("get property request", zap.String("id", req.Id))

	propertyID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid property id: %v", err)
	}

	property, err := s.propertyService.GetByID(propertyID, false)
	if err != nil {
		s.logger.Error("failed to get property", zap.Error(err))
		return nil, status.Errorf(codes.NotFound, "property not found: %v", err)
	}

	return &propertypb.PropertyResponse{
		Property: s.convertToProtoProperty(property),
	}, nil
}

func (s *PropertyGRPCServer) ListProperties(ctx context.Context, req *propertypb.ListPropertiesRequest) (*propertypb.ListPropertiesResponse, error) {
	s.logger.Debug("list properties request", zap.Any("request", req))

	// Convert gRPC request to DTO
	listReq := dto.PropertyListRequest{
		Page:           int(req.Page),
		Limit:          int(req.Limit),
		InternalCode:   &req.InternalCode,
		YearBuilt:      &req.YearBuilt,
		Bedrooms:       &req.Bedrooms,
		MinTotalArea:   &req.MinTotalArea,
		MaxTotalArea:   &req.MaxTotalArea,
		City:           &req.City,
		State:          &req.State,
		IncludeDeleted: req.IncludeDeleted,
	}

	if req.PropertyTypeId > 0 {
		listReq.PropertyTypeID = &req.PropertyTypeId
	}

	if req.OwnerPersonId != "" {
		ownerID, err := uuid.Parse(req.OwnerPersonId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner_person_id: %v", err)
		}
		listReq.OwnerPersonID = &ownerID
	}

	result, err := s.propertyService.List(listReq)
	if err != nil {
		s.logger.Error("failed to list properties", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to list properties: %v", err)
	}

	s.logger.Debug("list properties result",
		zap.Int64("total", result.Total),
		zap.Any("data_type", fmt.Sprintf("%T", result.Data)))

	// Convert to proto response - handle different possible data types
	var properties []dto.PropertyResponse

	switch data := result.Data.(type) {
	case []dto.PropertyResponse:
		properties = data
	case []interface{}:
		// Convert from []interface{} to []dto.PropertyResponse
		properties = make([]dto.PropertyResponse, len(data))
		for i, item := range data {
			if prop, ok := item.(dto.PropertyResponse); ok {
				properties[i] = prop
			} else {
				s.logger.Error("invalid property item type", zap.Any("type", fmt.Sprintf("%T", item)))
				return nil, status.Errorf(codes.Internal, "invalid property item type")
			}
		}
	default:
		s.logger.Error("invalid data type in list response", zap.Any("type", fmt.Sprintf("%T", result.Data)))
		return nil, status.Errorf(codes.Internal, "invalid data type in list response: %T", result.Data)
	}

	s.logger.Debug("converted properties", zap.Int("count", len(properties)))

	protoProperties := make([]*propertypb.Property, len(properties))
	for i, p := range properties {
		protoProperties[i] = s.convertToProtoProperty(&p)
	}

	return &propertypb.ListPropertiesResponse{
		Properties: protoProperties,
		Total:      int32(result.Total),
		Page:       int32(result.Page),
		Limit:      int32(result.Limit),
	}, nil
}

func (s *PropertyGRPCServer) CreateProperty(ctx context.Context, req *propertypb.CreatePropertyRequest) (*propertypb.PropertyResponse, error) {
	s.logger.Debug("create property request", zap.Any("request", req))

	ownerID, err := uuid.Parse(req.OwnerPersonId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid owner_person_id: %v", err)
	}

	addressID, err := uuid.Parse(req.AddressId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid address_id: %v", err)
	}

	createReq := dto.PropertyCreate{
		OwnerPersonID:  ownerID,
		AddressID:      addressID,
		PropertyTypeID: req.PropertyTypeId,
	}

	if req.InternalCode != "" {
		createReq.InternalCode = &req.InternalCode
	}
	if req.YearBuilt > 0 {
		createReq.YearBuilt = &req.YearBuilt
	}
	if req.Bedrooms > 0 {
		createReq.Bedrooms = &req.Bedrooms
	}
	if req.Bathrooms > 0 {
		bathrooms := float32(req.Bathrooms)
		createReq.Bathrooms = &bathrooms
	}
	if req.TotalAreaSqm > 0 {
		createReq.TotalAreaSqm = &req.TotalAreaSqm
	}
	if req.CoveredAreaSqm > 0 {
		createReq.CoveredAreaSqm = &req.CoveredAreaSqm
	}
	if req.Description != "" {
		createReq.Description = &req.Description
	}

	property, err := s.propertyService.Create(createReq, nil)
	if err != nil {
		s.logger.Error("failed to create property", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to create property: %v", err)
	}

	return &propertypb.PropertyResponse{
		Property: s.convertToProtoProperty(property),
	}, nil
}

// Helper function to convert DTO to proto
func (s *PropertyGRPCServer) convertToProtoProperty(p *dto.PropertyResponse) *propertypb.Property {
	property := &propertypb.Property{
		Id:             p.ID.String(),
		OwnerPersonId:  p.OwnerPersonID.String(),
		AddressId:      p.AddressID.String(),
		PropertyTypeId: p.PropertyTypeID,
	}

	if p.InternalCode != nil {
		property.InternalCode = *p.InternalCode
	}
	if p.YearBuilt != nil {
		property.YearBuilt = *p.YearBuilt
	}
	if p.Bedrooms != nil {
		property.Bedrooms = *p.Bedrooms
	}
	if p.Bathrooms != nil {
		property.Bathrooms = float64(*p.Bathrooms)
	}
	if p.TotalAreaSqm != nil {
		property.TotalAreaSqm = *p.TotalAreaSqm
	}
	if p.CoveredAreaSqm != nil {
		property.CoveredAreaSqm = *p.CoveredAreaSqm
	}
	if p.Description != nil {
		property.Description = *p.Description
	}
	if p.UpdatedBy != nil {
		property.UpdatedBy = p.UpdatedBy.String()
	}

	if !p.CreatedAt.IsZero() {
		property.CreatedAt = timestamppb.New(p.CreatedAt)
	}
	if p.UpdatedAt != nil {
		property.UpdatedAt = timestamppb.New(*p.UpdatedAt)
	}

	return property
}

// Placeholder implementations for other methods
func (s *PropertyGRPCServer) UpdateProperty(ctx context.Context, req *propertypb.UpdatePropertyRequest) (*propertypb.PropertyResponse, error) {
	s.logger.Debug("update property request", zap.Any("request", req))

	propertyID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid property id: %v", err)
	}

	updateReq := dto.PropertyUpdate{}

	// Only update fields that are provided
	if req.OwnerPersonId != "" {
		ownerID, err := uuid.Parse(req.OwnerPersonId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid owner_person_id: %v", err)
		}
		updateReq.OwnerPersonID = &ownerID
	}

	if req.AddressId != "" {
		addressID, err := uuid.Parse(req.AddressId)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "invalid address_id: %v", err)
		}
		updateReq.AddressID = &addressID
	}

	if req.PropertyTypeId > 0 {
		updateReq.PropertyTypeID = &req.PropertyTypeId
	}

	if req.InternalCode != "" {
		updateReq.InternalCode = &req.InternalCode
	}

	if req.YearBuilt > 0 {
		updateReq.YearBuilt = &req.YearBuilt
	}

	if req.Bedrooms > 0 {
		updateReq.Bedrooms = &req.Bedrooms
	}

	if req.Bathrooms > 0 {
		bathrooms := float32(req.Bathrooms)
		updateReq.Bathrooms = &bathrooms
	}

	if req.TotalAreaSqm > 0 {
		updateReq.TotalAreaSqm = &req.TotalAreaSqm
	}

	if req.CoveredAreaSqm > 0 {
		updateReq.CoveredAreaSqm = &req.CoveredAreaSqm
	}

	if req.Description != "" {
		updateReq.Description = &req.Description
	}

	property, err := s.propertyService.Update(propertyID, updateReq, nil)
	if err != nil {
		s.logger.Error("failed to update property", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to update property: %v", err)
	}

	return &propertypb.PropertyResponse{
		Property: s.convertToProtoProperty(property),
	}, nil
}

func (s *PropertyGRPCServer) DeleteProperty(ctx context.Context, req *propertypb.DeletePropertyRequest) (*emptypb.Empty, error) {
	s.logger.Debug("delete property request", zap.String("id", req.Id))

	propertyID, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid property id: %v", err)
	}

	err = s.propertyService.Delete(propertyID)
	if err != nil {
		s.logger.Error("failed to delete property", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "failed to delete property: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *PropertyGRPCServer) GetPropertyType(ctx context.Context, req *propertypb.GetPropertyTypeRequest) (*propertypb.PropertyTypeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "GetPropertyType not implemented yet")
}

func (s *PropertyGRPCServer) ListPropertyTypes(ctx context.Context, req *propertypb.ListPropertyTypesRequest) (*propertypb.ListPropertyTypesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListPropertyTypes not implemented")
}

func (s *PropertyGRPCServer) ListPropertyAmenities(ctx context.Context, req *propertypb.ListPropertyAmenitiesRequest) (*propertypb.ListPropertyAmenitiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "ListPropertyAmenities not implemented")
}

func (s *PropertyGRPCServer) AddPropertyAmenity(ctx context.Context, req *propertypb.AddPropertyAmenityRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "AddPropertyAmenity not implemented")
}

func (s *PropertyGRPCServer) RemovePropertyAmenity(ctx context.Context, req *propertypb.RemovePropertyAmenityRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "RemovePropertyAmenity not implemented")
}

func (s *PropertyGRPCServer) Start(port int) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	s.logger.Info("starting gRPC server", zap.Int("port", port))
	return s.server.Serve(lis)
}

func (s *PropertyGRPCServer) Stop() {
	s.logger.Info("stopping gRPC server")
	s.server.GracefulStop()
}
