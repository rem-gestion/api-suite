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
	g.PUT("/:id", ctrl.Update)
	g.DELETE("/:id", ctrl.Delete)

	// Contactos
	g.POST("/:id/contacts", ctrl.AddContact)
	g.GET("/:id/contacts", ctrl.ListContacts)

	// contactos por id (sin necesidad del id de persona)
	r.DELETE("/contacts/:contactID", ctrl.DeleteContact)
	r.PUT("/contacts/:contactID", ctrl.UpdateContact)
}
