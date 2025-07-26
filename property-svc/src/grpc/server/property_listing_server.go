package server

import (
	"context"
	"strconv"

	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/grpc/pb"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/services"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PropertyListingServer struct {
	pb.UnimplementedPropertyListingServiceServer
	propertyListingService *services.PropertyListingService
}

func NewPropertyListingServer(propertyListingService *services.PropertyListingService) *PropertyListingServer {
	return &PropertyListingServer{
		propertyListingService: propertyListingService,
	}
}

func (s *PropertyListingServer) CreatePropertyListing(ctx context.Context, req *pb.CreatePropertyListingRequest) (*pb.PropertyListingResponse, error) {
	propertyID, err := uuid.Parse(strconv.FormatInt(req.PropertyId, 10))
	if err != nil {
		return &pb.PropertyListingResponse{
			Success: false,
			Message: "Invalid property ID",
		}, status.Error(codes.InvalidArgument, "Invalid property ID")
	}

	var salePrice *float64
	if req.SalePrice != "" {
		if price, err := strconv.ParseFloat(req.SalePrice, 64); err == nil {
			salePrice = &price
		}
	}

	var currency *string
	if req.Currency != "" {
		currency = &req.Currency
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	createDTO := dto.PropertyListingCreate{
		PropertyID:    propertyID,
		OperationType: models.OperationType(req.ListingType),
		Price:         salePrice,
		Currency:      currency,
		Title:         req.SeoTitle,
		Description:   description,
		PublishOnWeb:  &req.IsPublished,
	}

	listing, err := s.propertyListingService.Create(createDTO, nil)
	if err != nil {
		return &pb.PropertyListingResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.PropertyListingResponse{
		Listing: convertToPropertyListingPB(listing),
		Success: true,
		Message: "Property listing created successfully",
	}, nil
}

func (s *PropertyListingServer) GetPropertyListing(ctx context.Context, req *pb.GetPropertyListingRequest) (*pb.PropertyListingResponse, error) {
	id, err := uuid.Parse(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return &pb.PropertyListingResponse{
			Success: false,
			Message: "Invalid listing ID",
		}, status.Error(codes.InvalidArgument, "Invalid listing ID")
	}

	listing, err := s.propertyListingService.GetByID(id)
	if err != nil {
		return &pb.PropertyListingResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.NotFound, err.Error())
	}

	return &pb.PropertyListingResponse{
		Listing: convertToPropertyListingPB(listing),
		Success: true,
		Message: "Property listing retrieved successfully",
	}, nil
}

func (s *PropertyListingServer) UpdatePropertyListing(ctx context.Context, req *pb.UpdatePropertyListingRequest) (*pb.PropertyListingResponse, error) {
	id, err := uuid.Parse(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return &pb.PropertyListingResponse{
			Success: false,
			Message: "Invalid listing ID",
		}, status.Error(codes.InvalidArgument, "Invalid listing ID")
	}

	var price *float64
	if req.SalePrice != "" {
		if salePrice, err := strconv.ParseFloat(req.SalePrice, 64); err == nil {
			price = &salePrice
		}
	}

	var operationType *models.OperationType
	if req.ListingType != "" {
		opType := models.OperationType(req.ListingType)
		operationType = &opType
	}

	var listingStatus *models.ListingStatus
	if req.Status != "" {
		status := models.ListingStatus(req.Status)
		listingStatus = &status
	}

	var currency *string
	if req.Currency != "" {
		currency = &req.Currency
	}

	var description *string
	if req.Description != "" {
		description = &req.Description
	}

	updateDTO := dto.PropertyListingUpdate{
		OperationType: operationType,
		ListingStatus: listingStatus,
		Price:         price,
		Currency:      currency,
		Description:   description,
	}

	listing, err := s.propertyListingService.Update(id, updateDTO, nil)
	if err != nil {
		return &pb.PropertyListingResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.PropertyListingResponse{
		Listing: convertToPropertyListingPB(listing),
		Success: true,
		Message: "Property listing updated successfully",
	}, nil
}

func (s *PropertyListingServer) DeletePropertyListing(ctx context.Context, req *pb.DeletePropertyListingRequest) (*pb.DeletePropertyListingResponse, error) {
	id, err := uuid.Parse(strconv.FormatInt(req.Id, 10))
	if err != nil {
		return &pb.DeletePropertyListingResponse{
			Success: false,
			Message: "Invalid listing ID",
		}, status.Error(codes.InvalidArgument, "Invalid listing ID")
	}

	err = s.propertyListingService.Delete(id, nil)
	if err != nil {
		return &pb.DeletePropertyListingResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeletePropertyListingResponse{
		Success: true,
		Message: "Property listing deleted successfully",
	}, nil
}

func (s *PropertyListingServer) ListPropertyListings(ctx context.Context, req *pb.ListPropertyListingsRequest) (*pb.ListPropertyListingsResponse, error) {
	var operationType *models.OperationType
	if req.ListingType != "" {
		opType := models.OperationType(req.ListingType)
		operationType = &opType
	}

	var listingStatus *models.ListingStatus
	if req.Status != "" {
		status := models.ListingStatus(req.Status)
		listingStatus = &status
	}

	filter := dto.PropertyListingFilter{
		OperationType: operationType,
		ListingStatus: listingStatus,
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}

	listings, total, err := s.propertyListingService.List(filter, page, limit)
	if err != nil {
		return &pb.ListPropertyListingsResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	pbListings := make([]*pb.PropertyListing, len(listings))
	for i, listing := range listings {
		pbListings[i] = convertToPropertyListingPB(&listing)
	}

	return &pb.ListPropertyListingsResponse{
		Listings: pbListings,
		Total:    int32(total),
		Page:     req.Page,
		Limit:    req.Limit,
		Success:  true,
		Message:  "Property listings retrieved successfully",
	}, nil
}

func (s *PropertyListingServer) SearchPropertyListings(ctx context.Context, req *pb.SearchPropertyListingsRequest) (*pb.ListPropertyListingsResponse, error) {
	var operationType *models.OperationType
	if req.ListingType != "" {
		opType := models.OperationType(req.ListingType)
		operationType = &opType
	}

	var minPrice *float64
	if req.MinPrice != "" {
		if price, err := strconv.ParseFloat(req.MinPrice, 64); err == nil {
			minPrice = &price
		}
	}

	var maxPrice *float64
	if req.MaxPrice != "" {
		if price, err := strconv.ParseFloat(req.MaxPrice, 64); err == nil {
			maxPrice = &price
		}
	}

	var currency *string
	if req.Currency != "" {
		currency = &req.Currency
	}

	filter := dto.PropertyListingFilter{
		OperationType: operationType,
		MinPrice:      minPrice,
		MaxPrice:      maxPrice,
		Currency:      currency,
	}

	page := int(req.Page)
	if page <= 0 {
		page = 1
	}
	limit := int(req.Limit)
	if limit <= 0 {
		limit = 10
	}

	listings, total, err := s.propertyListingService.Search(req.Query, filter, page, limit)
	if err != nil {
		return &pb.ListPropertyListingsResponse{
			Success: false,
			Message: err.Error(),
		}, status.Error(codes.Internal, err.Error())
	}

	pbListings := make([]*pb.PropertyListing, len(listings))
	for i, listing := range listings {
		pbListings[i] = convertToPropertyListingPB(&listing)
	}

	return &pb.ListPropertyListingsResponse{
		Listings: pbListings,
		Total:    int32(total),
		Page:     req.Page,
		Limit:    req.Limit,
		Success:  true,
		Message:  "Property listings search completed successfully",
	}, nil
}

func convertToPropertyListingPB(listing *dto.PropertyListingResponse) *pb.PropertyListing {
	// Convert UUID to int64 by using a simple hash of the first 8 bytes
	// This is for protobuf compatibility, not for reconstruction
	idBytes := listing.ID[:]
	var idInt64 int64
	for i := 0; i < 8 && i < len(idBytes); i++ {
		idInt64 = (idInt64 << 8) | int64(idBytes[i])
	}

	propertyIDBytes := listing.PropertyID[:]
	var propertyIDInt64 int64
	for i := 0; i < 8 && i < len(propertyIDBytes); i++ {
		propertyIDInt64 = (propertyIDInt64 << 8) | int64(propertyIDBytes[i])
	}

	pbListing := &pb.PropertyListing{
		Id:             idInt64,
		PropertyId:     propertyIDInt64,
		ListingType:    string(listing.OperationType),
		Status:         string(listing.ListingStatus),
		Currency:       listing.Currency,
		ViewsCount:     listing.ViewsCount,
		InquiriesCount: listing.InquiriesCount,
		CreatedAt:      listing.CreatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if listing.Price != nil {
		pbListing.SalePrice = strconv.FormatFloat(*listing.Price, 'f', 2, 64)
	}

	if listing.Description != nil {
		pbListing.Description = *listing.Description
	}

	if listing.UpdatedAt != nil {
		pbListing.UpdatedAt = listing.UpdatedAt.Format("2006-01-02T15:04:05Z")
	}

	return pbListing
}
