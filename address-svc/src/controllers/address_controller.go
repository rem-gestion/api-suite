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
