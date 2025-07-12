package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rem-gestion/api-suite/address/src/controllers"
)

func Setup(r *gin.Engine, ctrl *controller.Ctrl) {
	// API Gateway enrutará /api/addresses/ a este servicio
	// Por lo tanto, usamos "/" como base para evitar duplicación

	r.GET("/", ctrl.List) // Lista todas las direcciones
	r.POST("/", ctrl.Create)
	r.GET("/:id", ctrl.Get)
	r.PUT("/:id", ctrl.Update)
	r.DELETE("/:id", ctrl.Delete)
}
