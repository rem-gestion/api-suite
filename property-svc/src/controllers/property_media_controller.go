package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/services"
	"github.com/rem-gestion/rem-common/errors"
)

type PropertyMediaController struct {
	service *services.PropertyMediaService
}

func NewPropertyMediaController(service *services.PropertyMediaService) *PropertyMediaController {
	return &PropertyMediaController{service: service}
}

// CreateMedia creates a new property media
func (c *PropertyMediaController) CreateMedia(ctx *gin.Context) {
	var req dto.PropertyMediaCreate
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

// GetMediaByID retrieves a property media by ID
func (c *PropertyMediaController) GetMediaByID(ctx *gin.Context) {
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

// GetMediaByPropertyID retrieves all media for a specific property
func (c *PropertyMediaController) GetMediaByPropertyID(ctx *gin.Context) {
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

// ListMedia retrieves property media with optional filters
func (c *PropertyMediaController) ListMedia(ctx *gin.Context) {
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

	// Parse optional property_id filter
	var filters dto.PropertyMediaFilter
	propertyIDStr := ctx.Query("property_id")
	if propertyIDStr != "" {
		propertyID, err := uuid.Parse(propertyIDStr)
		if err != nil {
			ctx.Error(&errors.BadRequestError{Msg: "invalid property_id format"})
			return
		}
		filters.PropertyID = &propertyID
	}

	if mediaType := ctx.Query("media_type"); mediaType != "" {
		// Convert string to MediaType if needed
		filters.MediaType = nil // Set based on mediaType
	}

	if isMainStr := ctx.Query("is_main"); isMainStr == "true" {
		isMain := true
		filters.IsMain = &isMain
	}

	if isPublicStr := ctx.Query("is_public"); isPublicStr == "true" {
		isPublic := true
		filters.IsPublic = &isPublic
	}

	response, total, err := c.service.List(filters, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.JSON(http.StatusOK, response)
}

// UpdateMedia updates an existing property media
func (c *PropertyMediaController) UpdateMedia(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	var req dto.PropertyMediaUpdate
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

// DeleteMedia deletes a property media
func (c *PropertyMediaController) DeleteMedia(ctx *gin.Context) {
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

// SetMainImage sets a media as the main image for a property
func (c *PropertyMediaController) SetMainImage(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	// Get property_id from request body or query
	propertyIDStr := ctx.Query("property_id")
	if propertyIDStr == "" {
		ctx.Error(&errors.BadRequestError{Msg: "property_id is required"})
		return
	}

	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property_id format"})
		return
	}

	// Get user ID from context (set by auth middleware)
	var updatedBy uuid.UUID
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			updatedBy = uid
		}
	}

	err = c.service.SetAsMain(id, propertyID, updatedBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Media set as main image successfully"})
}

// UpdateDisplayOrder updates the display order of media items
func (c *PropertyMediaController) UpdateDisplayOrder(ctx *gin.Context) {
	var req struct {
		PropertyID uuid.UUID   `json:"property_id" binding:"required"`
		MediaIDs   []uuid.UUID `json:"media_ids" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	// Get user ID from context (set by auth middleware)
	var updatedBy uuid.UUID
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			updatedBy = uid
		}
	}

	err := c.service.UpdateDisplayOrder(req.PropertyID, req.MediaIDs, updatedBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "Display order updated successfully"})
}
