package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rem-gestion/api-suite/person/src/controllers"
)

func Setup(r *gin.Engine, ctrl *controller.Ctrl) {
	g := r.Group("/persons")
	g.POST("/", ctrl.Create)
	g.GET("/:id", ctrl.Get)
	g.PUT("/:id", ctrl.Update)
	g.DELETE("/:id", ctrl.Delete)
}
