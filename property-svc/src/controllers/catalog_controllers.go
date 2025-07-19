package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/models"
	"github.com/rem-gestion/api-suite/property/src/services"
	"github.com/rem-gestion/rem-common/errors"
)

// PropertyTypeController handles property type HTTP requests
type PropertyTypeController struct {
	service *services.PropertyTypeService
}

func NewPropertyTypeController(service *services.PropertyTypeService) *PropertyTypeController {
	return &PropertyTypeController{service: service}
}

func (c *PropertyTypeController) Create(ctx *gin.Context) {
	var req dto.PropertyTypeCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.service.Create(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *PropertyTypeController) GetByID(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.GetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyTypeController) GetByCode(ctx *gin.Context) {
	code := ctx.Param("code")
	if code == "" {
		ctx.Error(&errors.BadRequestError{Msg: "code is required"})
		return
	}

	response, err := c.service.GetByCode(code)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyTypeController) List(ctx *gin.Context) {
	var isActive *bool
	if activeStr := ctx.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			isActive = &active
		}
	}

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	response, err := c.service.List(isActive, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyTypeController) Update(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dto.PropertyTypeUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.service.Update(id, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PropertyTypeController) Delete(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// ManagerTypeController handles manager type HTTP requests
type ManagerTypeController struct {
	service *services.ManagerTypeService
}

func NewManagerTypeController(service *services.ManagerTypeService) *ManagerTypeController {
	return &ManagerTypeController{service: service}
}

func (c *ManagerTypeController) Create(ctx *gin.Context) {
	var req dto.ManagerTypeCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.service.Create(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *ManagerTypeController) GetByID(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.GetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ManagerTypeController) List(ctx *gin.Context) {
	var isActive *bool
	if activeStr := ctx.Query("is_active"); activeStr != "" {
		if active, err := strconv.ParseBool(activeStr); err == nil {
			isActive = &active
		}
	}

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	response, err := c.service.List(isActive, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ManagerTypeController) Update(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dto.ManagerTypeUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.service.Update(id, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *ManagerTypeController) Delete(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

// AmenityController handles amenity HTTP requests
type AmenityController struct {
	service *services.AmenityService
}

func NewAmenityController(service *services.AmenityService) *AmenityController {
	return &AmenityController{service: service}
}

func (c *AmenityController) Create(ctx *gin.Context) {
	var req dto.AmenityCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.service.Create(req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *AmenityController) GetByID(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.GetByID(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AmenityController) List(ctx *gin.Context) {
	var category *models.AmenityCategory
	if catStr := ctx.Query("category"); catStr != "" {
		cat := models.AmenityCategory(catStr)
		// Validate category
		if cat == models.AmenityCategoryGeneral ||
			cat == models.AmenityCategoryServices ||
			cat == models.AmenityCategoryEnvironments ||
			cat == models.AmenityCategorySecurity ||
			cat == models.AmenityCategoryComfort {
			category = &cat
		}
	}

	page := 1
	if pageStr := ctx.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	limit := 20
	if limitStr := ctx.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	response, err := c.service.List(category, page, limit)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AmenityController) Update(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	var req dto.AmenityUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}

	response, err := c.service.Update(id, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *AmenityController) Delete(ctx *gin.Context) {
	id, err := parseInt32Param(ctx, "id")
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Delete(id)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
