package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rem-gestion/api-suite/person/src/controllers"
)

func Setup(r *gin.Engine, ctrl *controller.Ctrl) {
	// API Gateway enrutará /api/persons/ a este servicio
	// Por lo tanto, usamos "/" como base para evitar duplicación

	// Personas
	r.POST("/", ctrl.Create)
	r.GET("/", ctrl.List)
	r.GET("/:id", ctrl.Get)
	r.GET("/:id/full", ctrl.GetFull) // NUEVO: persona con dirección expandida
	r.PUT("/:id", ctrl.Update)
	r.DELETE("/:id", ctrl.Delete)

	// Bulk operations
	r.POST("/bulk", ctrl.BulkCreate) // NUEVO: creación masiva

	// Contactos (usando la ruta /api/contacts/ del gateway)
	r.POST("/:id/contacts", ctrl.AddContact)
	r.GET("/:id/contacts", ctrl.ListContacts)

	// contactos por id (sin necesidad del id de persona)
	// Estas rutas serán manejadas por /api/contacts/ en el gateway
	contacts := r.Group("/contacts")
	{
		contacts.DELETE("/:contactID", ctrl.DeleteContact)
		contacts.PUT("/:contactID", ctrl.UpdateContact)
		contacts.GET("/primary/:personId", ctrl.GetPrimaryContact) // NUEVO: contacto primario por persona
	}
}
