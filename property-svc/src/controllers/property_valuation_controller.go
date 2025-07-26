package controllers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/services"
	"github.com/rem-gestion/rem-common/errors"
)

type PropertyValuationController struct {
	service *services.PropertyValuationService
}

func NewPropertyValuationController(service *services.PropertyValuationService) *PropertyValuationController {
	return &PropertyValuationController{service: service}
}

// CreateValuation creates a new property valuation
func (c *PropertyValuationController) CreateValuation(ctx *gin.Context) {
	var req dto.PropertyValuationCreate
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

// GetValuationByID retrieves a property valuation by ID
func (c *PropertyValuationController) GetValuationByID(ctx *gin.Context) {
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

// GetValuationsByPropertyID retrieves all valuations for a specific property
func (c *PropertyValuationController) GetValuationsByPropertyID(ctx *gin.Context) {
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

// GetLatestValuation retrieves the latest valuation for a specific property
func (c *PropertyValuationController) GetLatestValuation(ctx *gin.Context) {
	propertyIDStr := ctx.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property ID format"})
		return
	}

	response, err := c.service.GetLatestByPropertyID(propertyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// GetValuationsByType retrieves valuations for a property filtered by type
func (c *PropertyValuationController) GetValuationsByType(ctx *gin.Context) {
	propertyIDStr := ctx.Param("id")
	propertyID, err := uuid.Parse(propertyIDStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid property ID format"})
		return
	}

	valuationType := ctx.Param("type")
	if valuationType == "" {
		ctx.Error(&errors.BadRequestError{Msg: "valuation type is required"})
		return
	}

	response, err := c.service.GetValuationsByType(propertyID, valuationType)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

// ListValuations retrieves property valuations with optional filters
func (c *PropertyValuationController) ListValuations(ctx *gin.Context) {
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
	var filters dto.PropertyValuationFilter
	propertyIDStr := ctx.Query("property_id")
	if propertyIDStr != "" {
		propertyID, err := uuid.Parse(propertyIDStr)
		if err != nil {
			ctx.Error(&errors.BadRequestError{Msg: "invalid property_id format"})
			return
		}
		filters.PropertyID = &propertyID
	}

	// Parse optional filters
	// filters := dto.PropertyValuationFilter{
	// 	PropertyID: propertyID,
	// }

	if valuationType := ctx.Query("valuation_type"); valuationType != "" {
		filters.ValuationType = &valuationType
	}

	if minValueStr := ctx.Query("min_value"); minValueStr != "" {
		if minValue, err := strconv.ParseFloat(minValueStr, 64); err == nil {
			filters.MinValue = &minValue
		}
	}

	if maxValueStr := ctx.Query("max_value"); maxValueStr != "" {
		if maxValue, err := strconv.ParseFloat(maxValueStr, 64); err == nil {
			filters.MaxValue = &maxValue
		}
	}

	if dateFromStr := ctx.Query("date_from"); dateFromStr != "" {
		if dateFrom, err := time.Parse("2006-01-02", dateFromStr); err == nil {
			filters.DateFrom = &dateFrom
		}
	}

	if dateToStr := ctx.Query("date_to"); dateToStr != "" {
		if dateTo, err := time.Parse("2006-01-02", dateToStr); err == nil {
			filters.DateTo = &dateTo
		}
	}

	response, total, err := c.service.List(filters, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Header("X-Total-Count", strconv.FormatInt(total, 10))
	ctx.JSON(http.StatusOK, response)
}

// UpdateValuation updates an existing property valuation
func (c *PropertyValuationController) UpdateValuation(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		ctx.Error(&errors.BadRequestError{Msg: "invalid UUID format"})
		return
	}

	var req dto.PropertyValuationUpdate
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

// DeleteValuation deletes a property valuation
func (c *PropertyValuationController) DeleteValuation(ctx *gin.Context) {
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
