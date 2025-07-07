package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rem-gestion/api-suite/person/src/controllers"
)

func Setup(r *gin.Engine, ctrl *controller.Ctrl) {
	g := r.Group("/persons")

	// Personas
	g.POST("/", ctrl.Create)
	g.GET("/", ctrl.List)
	g.GET("/:id", ctrl.Get)
	g.GET("/:id/full", ctrl.GetFull) // NUEVO: persona con dirección expandida
	g.PUT("/:id", ctrl.Update)
	g.DELETE("/:id", ctrl.Delete)

	// Bulk operations
	g.POST("/bulk", ctrl.BulkCreate) // NUEVO: creación masiva

	// Contactos
	g.POST("/:id/contacts", ctrl.AddContact)
	g.GET("/:id/contacts", ctrl.ListContacts)

	// contactos por id (sin necesidad del id de persona)
	r.DELETE("/contacts/:contactID", ctrl.DeleteContact)
	r.PUT("/contacts/:contactID", ctrl.UpdateContact)

	// NUEVO: contacto primario por persona
	r.GET("/contacts/primary/:personId", ctrl.GetPrimaryContact)
}
