package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/property-svc/src/controllers"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/properties", controllers.CreateProperty)
	r.GET("/properties", controllers.GetAllProperties)
	r.GET("/properties/:id", controllers.GetPropertyByID)
	r.PUT("/properties/:id", controllers.UpdateProperty)
	r.DELETE("/properties/:id", controllers.DeleteProperty)
	return r
}
