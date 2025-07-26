package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/services"
	"github.com/rem-gestion/rem-common/errors"
)

type PropertyListingController struct {
	service *services.PropertyListingService
}

func NewPropertyListingController(service *services.PropertyListingService) *PropertyListingController {
	return &PropertyListingController{service: service}
}

// CreateListing creates a new property listing
func (c *PropertyListingController) CreateListing(ctx *gin.Context) {
	var req dto.PropertyListingCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	var createdBy *uuid.UUID
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			createdBy = &uid
		}
	}

	response, err := c.service.Create(req, createdBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// GetListingByID retrieves a property listing by ID
func (c *PropertyListingController) GetListingByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	response, err := c.service.GetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetListingsByPropertyID retrieves all listings for a specific property
func (c *PropertyListingController) GetListingsByPropertyID(ctx *gin.Context) {
	propertyIDStr := ctx.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property ID format"})
		return
	}

	response, err := c.service.GetByPropertyID(propertyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListListings retrieves property listings with optional filters
func (c *PropertyListingController) ListListings(ctx *gin.Context) {
	// Parse query parameters for pagination
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Parse optional filters
	filters := dto.PropertyListingFilter{}

	if operationType := ctx.Query("operation_type"); operationType != "" {
		opType := models.OperationType(operationType)
		filters.OperationType = &opType
	}

	if status := ctx.Query("listing_status"); status != "" {
		listingStatus := models.ListingStatus(status)
		filters.ListingStatus = &listingStatus
	}

	if minPriceStr := ctx.Query("min_price"); minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			filters.MinPrice = &minPrice
		}
	}

	if maxPriceStr := ctx.Query("max_price"); maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			filters.MaxPrice = &maxPrice
		}
	}

	if currency := ctx.Query("currency"); currency != "" {
		filters.Currency = &currency
	}

	if highlightedStr := ctx.Query("highlighted"); highlightedStr == "true" {
		highlighted := true
		filters.Highlighted = &highlighted
	}

	response, total, err := c.service.List(filters, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.JSON(http.StatusOK, response)
}

// UpdateListing updates an existing property listing
func (c *PropertyListingController) UpdateListing(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	var req dto.PropertyListingUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	var updatedBy *uuid.UUID
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			updatedBy = &uid
		}
	}

	response, err := c.service.Update(id, req, updatedBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// DeleteListing deletes a property listing
func (c *PropertyListingController) DeleteListing(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	// Get user ID from context (set by auth middleware)
	var deletedBy *uuid.UUID
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			deletedBy = &uid
		}
	}

	err = c.service.Delete(id, deletedBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}

// SearchListings performs a text search on property listings
func (c *PropertyListingController) SearchListings(ctx *gin.Context) {
	searchQuery := ctx.Query("q")
	if searchQuery == "" {
		ctx.Error(&errors.BadRequestError{Msg: "search query is required"})
		return
	}

	// Parse query parameters for pagination
	pageStr := ctx.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		page = 1
	}

	limitStr := ctx.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	// Use empty filters for search
	filters := dto.PropertyListingFilter{}

	response, total, err := c.service.Search(searchQuery, filters, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.JSON(http.StatusOK, response)
}

// IncrementViews increments the view count for a listing
func (c *PropertyListingController) IncrementViews(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	err = c.service.IncrementViews(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Views incremented successfully"})
}

// IncrementFavorites increments the favorite count for a listing
func (c *PropertyListingController) IncrementFavorites(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	err = c.service.IncrementFavorites(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Favorites incremented successfully"})
}

// DecrementFavorites decrements the favorite count for a listing
func (c *PropertyListingController) DecrementFavorites(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	err = c.service.DecrementFavorites(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Favorites decremented successfully"})
}
