package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/address/src/dto"
	"github.com/rem-gestion/api-suite/address/src/services"
	"github.com/rem-gestion/rem-common/errors"
)

type Ctrl struct{ svc *services.AddressService }

func New(s *services.AddressService) *Ctrl { return &Ctrl{s} }

func (c *Ctrl) Create(ctx *gin.Context) {
	var in dto.AddressCreate
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.Create(in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

func (c *Ctrl) Get(ctx *gin.Context) {
	a, err := c.svc.Get(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, a)
}

func (c *Ctrl) List(ctx *gin.Context) {
	// Parámetros de consulta opcionales
	city := ctx.DefaultQuery("city", "")
	limit := 20 // por defecto
	offset := 0 // por defecto

	// Aquí podrías agregar lógica para parsear limit y offset desde query params si es necesario

	addresses, err := c.svc.List(city, limit, offset)
	if err != nil {
		ctx.Error(err)
		return
	}

	// Respuesta con formato similar al person service
	response := map[string]interface{}{
		"data":  addresses,
		"page":  1,
		"per":   limit,
		"total": len(addresses),
	}

	ctx.JSON(200, response)
}

func (c *Ctrl) Update(ctx *gin.Context) {
	var in dto.AddressUpdate
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&errors.BadRequestError{Msg: err.Error()})
		return
	}
	a, err := c.svc.Update(ctx.Param("id"), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(200, a)
}

func (c *Ctrl) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Param("id")); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
