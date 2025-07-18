package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/rem-gestion/api-suite/property/src/dto"
	"github.com/rem-gestion/api-suite/property/src/services"
	rerrors "github.com/rem-gestion/rem-common/errors"
)

type Ctrl struct{ svc *services.PropertyService }

func New(s *services.PropertyService) *Ctrl { return &Ctrl{s} }

/* ───────────────────── Properties ────────────────────── */

// POST /properties
func (c *Ctrl) CreateProperty(ctx *gin.Context) {
	var in dto.CreatePropertyDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.CreateProperty(ctx.Request.Context(), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// GET /properties/:id
func (c *Ctrl) GetProperty(ctx *gin.Context) {
	p, err := c.svc.GetProperty(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, p)
}

// GET /properties?page=&per_page=&search=&type=&owner_person_id=
func (c *Ctrl) ListProperties(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if per < 1 || per > 100 {
		per = 20
	}

	search := ctx.Query("search")
	propertyType := ctx.Query("type")
	ownerPersonID := ctx.Query("owner_person_id")

	var pType *string
	if propertyType != "" {
		pType = &propertyType
	}
	var ownerID *string
	if ownerPersonID != "" {
		ownerID = &ownerPersonID
	}

	limit := per
	offset := (page - 1) * per

	list, total, err := c.svc.ListProperties(ctx.Request.Context(), search, pType, ownerID, limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":        list,
		"total":       total,
		"page":        page,
		"per_page":    per,
		"total_pages": (int(total) + per - 1) / per,
	})
}

// PUT /properties/:id
func (c *Ctrl) UpdateProperty(ctx *gin.Context) {
	var in dto.UpdatePropertyDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.UpdateProperty(ctx.Request.Context(), ctx.Param("id"), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// DELETE /properties/:id
func (c *Ctrl) DeleteProperty(ctx *gin.Context) {
	err := c.svc.DeleteProperty(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

/* ───────────────────── Amenities ────────────────────── */

// POST /amenities
func (c *Ctrl) CreateAmenity(ctx *gin.Context) {
	var in dto.CreateAmenityDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.CreateAmenity(ctx.Request.Context(), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// GET /amenities/:id
func (c *Ctrl) GetAmenity(ctx *gin.Context) {
	a, err := c.svc.GetAmenity(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, a)
}

// GET /amenities?page=&per_page=&search=&category=
func (c *Ctrl) ListAmenities(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))

	if page < 1 {
		page = 1
	}
	if per < 1 || per > 100 {
		per = 20
	}

	search := ctx.Query("search")
	category := ctx.Query("category")

	var cat *string
	if category != "" {
		cat = &category
	}

	limit := per
	offset := (page - 1) * per

	list, total, err := c.svc.ListAmenities(ctx.Request.Context(), search, cat, limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":        list,
		"total":       total,
		"page":        page,
		"per_page":    per,
		"total_pages": (int(total) + per - 1) / per,
	})
}

// PUT /amenities/:id
func (c *Ctrl) UpdateAmenity(ctx *gin.Context) {
	var in dto.UpdateAmenityDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.UpdateAmenity(ctx.Request.Context(), ctx.Param("id"), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// DELETE /amenities/:id
func (c *Ctrl) DeleteAmenity(ctx *gin.Context) {
	err := c.svc.DeleteAmenity(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}

/* ───────────────────── Property Amenities Management ────────────────────── */

// POST /properties/:property_id/amenities/:amenity_id
func (c *Ctrl) AddAmenityToProperty(ctx *gin.Context) {
	propertyID := ctx.Param("property_id")
	amenityID := ctx.Param("amenity_id")

	err := c.svc.AddAmenityToProperty(ctx.Request.Context(), propertyID, amenityID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"message": "amenity added to property"})
}

// DELETE /properties/:property_id/amenities/:amenity_id
func (c *Ctrl) RemoveAmenityFromProperty(ctx *gin.Context) {
	propertyID := ctx.Param("property_id")
	amenityID := ctx.Param("amenity_id")

	err := c.svc.RemoveAmenityFromProperty(ctx.Request.Context(), propertyID, amenityID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusNoContent, nil)
}
