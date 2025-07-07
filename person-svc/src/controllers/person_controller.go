package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/rem-gestion/api-suite/person/src/dto"
	"github.com/rem-gestion/api-suite/person/src/models"
	"github.com/rem-gestion/api-suite/person/src/services"
	rerrors "github.com/rem-gestion/rem-common/errors"
)

type Ctrl struct{ svc *services.PersonService }

func New(s *services.PersonService) *Ctrl { return &Ctrl{s} }

/* ───────────────────── Personas ────────────────────── */

// POST /persons
func (c *Ctrl) Create(ctx *gin.Context) {
	var in dto.CreatePersonDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.Create(in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// GET /persons/:id
func (c *Ctrl) Get(ctx *gin.Context) {
	p, err := c.svc.Get(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, p)
}

// GET /persons?page=&per_page=&search=&type=&with_contacts=
func (c *Ctrl) List(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(ctx.DefaultQuery("per_page", "20"))
	search := ctx.Query("search")
	withCt := ctx.DefaultQuery("with_contacts", "false") == "true"

	// type filter
	var filter *models.PersonType
	if t := ctx.Query("type"); t == "individual" || t == "company" {
		ft := models.PersonType(t)
		filter = &ft
	}

	data, total, err := c.svc.List(page, per, search, filter)
	if err != nil {
		ctx.Error(err)
		return
	}

	// si el caller NO quiere contactos → limpiamos la slice
	if !withCt {
		for i := range data {
			data[i].Contactos = nil
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"total": total,
		"data":  data,
		"page":  page,
		"per":   per,
	})
}

// PUT /persons/:id
func (c *Ctrl) Update(ctx *gin.Context) {
	var in dto.UpdatePersonDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	p, err := c.svc.Update(ctx.Param("id"), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, p)
}

// DELETE /persons/:id
func (c *Ctrl) Delete(ctx *gin.Context) {
	if err := c.svc.Delete(ctx.Param("id")); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

/* ───────────────────── Contactos ───────────────────── */

// POST /persons/:id/contacts
func (c *Ctrl) AddContact(ctx *gin.Context) {
	var in dto.CreateContactoDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	in.PersonaID = uuidFromParam(ctx.Param("id"), ctx)
	if ctx.IsAborted() {
		return
	}
	out, err := c.svc.AddContact(in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusCreated, out)
}

// GET /persons/:id/contacts
func (c *Ctrl) ListContacts(ctx *gin.Context) {
	personID := ctx.Param("id")
	list, err := c.svc.ListContacts(personID)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, list)
}

// PUT /contacts/:contactID
func (c *Ctrl) UpdateContact(ctx *gin.Context) {
	var in dto.UpdateContactoDTO
	if err := ctx.ShouldBindJSON(&in); err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: err.Error()})
		return
	}
	out, err := c.svc.UpdateContact(ctx.Param("contactID"), in)
	if err != nil {
		ctx.Error(err)
		return
	}
	ctx.JSON(http.StatusOK, out)
}

// DELETE /contacts/:contactID
func (c *Ctrl) DeleteContact(ctx *gin.Context) {
	if err := c.svc.DeleteContact(ctx.Param("contactID")); err != nil {
		ctx.Error(err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

/* ───────────── helpers ───────────── */

func uuidFromParam(id string, ctx *gin.Context) uuid.UUID {
	u, err := uuid.Parse(id)
	if err != nil {
		ctx.Error(&rerrors.BadRequestError{Msg: "id inválido"})
		ctx.Abort()
	}
	return u
}
