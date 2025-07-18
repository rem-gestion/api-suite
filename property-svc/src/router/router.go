package router

import (
	"github.com/gin-gonic/gin"
	controller "github.com/rem-gestion/api-suite/property/src/controllers"
)

func Setup(r *gin.Engine, ctrl *controller.Ctrl) {
	// API Gateway enrutará /api/properties/ a este servicio
	// Por lo tanto, usamos "/" como base para evitar duplicación

	// Properties
	r.POST("/", ctrl.CreateProperty)
	r.GET("/", ctrl.ListProperties)
	r.GET("/:id", ctrl.GetProperty)
	r.PUT("/:id", ctrl.UpdateProperty)
	r.DELETE("/:id", ctrl.DeleteProperty)

	// Property Types (master data)
	r.GET("/property-types", ctrl.GetPropertyTypes)

	// Manager Types (master data)
	r.GET("/manager-types", ctrl.GetManagerTypes)

	// Property-specific endpoints (use different pattern to avoid conflicts)
	properties := r.Group("/properties")
	{
		// Property-Amenity Relations
		properties.GET("/:property_id/amenities", ctrl.GetPropertyAmenities)
	}

	// Property-Amenity management - using /manage prefix to avoid route conflicts
	manage := r.Group("/manage")
	{
		manage.POST("/:property_id/amenities/:amenity_id", ctrl.AddAmenityToProperty)
		manage.DELETE("/:property_id/amenities/:amenity_id", ctrl.RemoveAmenityFromProperty)
	}

	// Amenities management (using /api/amenities/ route from gateway)
	amenities := r.Group("/amenities")
	{
		amenities.POST("/", ctrl.CreateAmenity)
		amenities.GET("/", ctrl.ListAmenities)
		amenities.GET("/:id", ctrl.GetAmenity)
		amenities.PUT("/:id", ctrl.UpdateAmenity)
		amenities.DELETE("/:id", ctrl.DeleteAmenity)
	}
}
