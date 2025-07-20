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

type PropertyController struct {
	service *services.PropertyService
}

func NewPropertyController(service *services.PropertyService) *PropertyController {
	return &PropertyController{service: service}
}

func (c *PropertyController) Create(ctx *gin.Context) {
	var req dto.PropertyCreate
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

func (c *PropertyController) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	includeRelations := ctx.Query("include_relations") == "true"

	response, err := c.service.GetByID(id, includeRelations)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyController) GetByInternalCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		ctx.Error(&errors.BadRequestError{Msg: "internal code is required"})
		return
	}

	response, err := c.service.GetByInternalCode(code)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyController) List(ctx *gin.Context) {
	var req dto.PropertyListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	response, err := c.service.List(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyController) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	var req dto.PropertyUpdate
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

func (c *PropertyController) SoftDelete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	// Get user ID from context (set by auth middleware)
	var deletedBy uuid.UUID
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(uuid.UUID); ok {
			deletedBy = uid
		} else {
			ctx.Error(&errors.UnauthorizedError{Msg: "user not authenticated"})
			return
		}
	} else {
		ctx.Error(&errors.UnauthorizedError{Msg: "user not authenticated"})
		return
	}

	err = c.service.SoftDelete(id, deletedBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PropertyController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *PropertyController) Search(ctx *gin.Context) {
	var req dto.PropertySearchRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	// Set defaults
	if req.Page == 0 {
		req.Page = 1
	}
	if req.Limit == 0 {
		req.Limit = 20
	}

	response, err := c.service.Search(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyController) BulkCreate(ctx *gin.Context) {
	var req dto.BulkPropertyCreateRequest
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

	response, err := c.service.BulkCreate(req, createdBy)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// ListPropertyAmenities - GET /properties/:id/amenities
func (c *PropertyController) ListPropertyAmenities(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property UUID format"})
		return
	}

	amenities, err := c.service.ListPropertyAmenities(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, amenities)
}

// AddAmenityToProperty - POST /properties/:id/amenities
func (c *PropertyController) AddAmenityToProperty(ctx *gin.Context) {
	idStr := ctx.Param("id")
	propertyID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property UUID format"})
		return
	}

	var req struct {
		AmenityID int32  `json:"amenity_id" form:"amenity_id" binding:"required"`
		Note      string `json:"note" form:"note"`
	}

	// Try to bind from JSON first, then from query/form params
	if err := ctx.ShouldBindJSON(&req); err != nil {
		// If JSON binding fails, try form/query binding
		if bindErr := ctx.ShouldBind(&req); bindErr != nil {
			ctx.Error(&errors.BadRequestError{Msg: "amenity_id is required (can be sent in JSON body or as query parameter)"})
			return
		}
	}

	err = c.service.AddAmenityToProperty(propertyID, req.AmenityID, req.Note)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":     "amenity added to property successfully",
		"property_id": propertyID,
		"amenity_id":  req.AmenityID,
	})
}

// RemoveAmenityFromProperty - DELETE /properties/:id/amenities/:amenity_id
func (c *PropertyController) RemoveAmenityFromProperty(ctx *gin.Context) {
	idStr := ctx.Param("id")
	propertyID, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property UUID format"})
		return
	}

	amenityID, err := parseInt32Param(ctx, "amenity_id")
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.RemoveAmenityFromProperty(propertyID, amenityID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":     "amenity removed from property successfully",
		"property_id": propertyID,
		"amenity_id":  amenityID,
	})
}

// Utility function to parse int32 from string parameter
func parseInt32Param(ctx *gin.Context, paramName string) (int32, error) {
	paramStr := ctx.Param(paramName)
	param, err := strconv.ParseInt(paramStr, 10, 32)
	if err != nil {
		return 0, &errors.BadRequestError{Msg: "invalid " + paramName + " format"}
	}
	return int32(param), nil
}
