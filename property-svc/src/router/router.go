package router

import (
	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/api-suite/property/src/controllers"
	"github.com/rem-gestion/rem-common/middleware"
)

type Router struct {
	propertyController           *controllers.PropertyController
	propertyTypeController       *controllers.PropertyTypeController
	managerTypeController        *controllers.ManagerTypeController
	amenityController            *controllers.AmenityController
	propertyManagementController *controllers.PropertyManagementController
}

func New(
	propertyController *controllers.PropertyController,
	propertyTypeController *controllers.PropertyTypeController,
	managerTypeController *controllers.ManagerTypeController,
	amenityController *controllers.AmenityController,
	propertyManagementController *controllers.PropertyManagementController,
) *Router {
	return &Router{
		propertyController:           propertyController,
		propertyTypeController:       propertyTypeController,
		managerTypeController:        managerTypeController,
		amenityController:            amenityController,
		propertyManagementController: propertyManagementController,
	}
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	// Middleware setup
	engine.Use(middleware.ErrorHandler())

	// Property routes
	properties := engine.Group("/properties")
	{
		properties.POST("", r.propertyController.Create)
		properties.GET("", r.propertyController.List)
		properties.GET("/search", r.propertyController.Search)
		properties.POST("/bulk", r.propertyController.BulkCreate)
		properties.GET("/:id", r.propertyController.GetByID)
		properties.PUT("/:id", r.propertyController.Update)
		properties.DELETE("/:id", r.propertyController.SoftDelete)
		properties.DELETE("/:id/permanent", r.propertyController.Delete)
		properties.GET("/internal-code/:code", r.propertyController.GetByInternalCode)

		// Property-Amenity relations
		properties.GET("/:id/amenities", r.propertyController.ListPropertyAmenities)
		properties.POST("/:id/amenities", r.propertyController.AddAmenityToProperty)
		properties.DELETE("/:id/amenities/:amenity_id", r.propertyController.RemoveAmenityFromProperty)
	}

	// Property Type routes (catalog)
	propertyTypes := engine.Group("/property-types")
	{
		propertyTypes.POST("", r.propertyTypeController.Create)
		propertyTypes.GET("", r.propertyTypeController.List)
		propertyTypes.GET("/:id", r.propertyTypeController.GetByID)
		propertyTypes.PUT("/:id", r.propertyTypeController.Update)
		propertyTypes.DELETE("/:id", r.propertyTypeController.Delete)
		propertyTypes.GET("/code/:code", r.propertyTypeController.GetByCode)
	}

	// Manager Type routes (catalog)
	managerTypes := engine.Group("/manager-types")
	{
		managerTypes.POST("", r.managerTypeController.Create)
		managerTypes.GET("", r.managerTypeController.List)
		managerTypes.GET("/:id", r.managerTypeController.GetByID)
		managerTypes.PUT("/:id", r.managerTypeController.Update)
		managerTypes.DELETE("/:id", r.managerTypeController.Delete)
	}

	// Amenity routes (catalog)
	amenities := engine.Group("/amenities")
	{
		amenities.POST("", r.amenityController.Create)
		amenities.GET("", r.amenityController.List)
		amenities.GET("/:id", r.amenityController.GetByID)
		amenities.PUT("/:id", r.amenityController.Update)
		amenities.DELETE("/:id", r.amenityController.Delete)
	}

	// Property Management routes
	propertyManagements := engine.Group("/property-managements")
	{
		propertyManagements.POST("", r.propertyManagementController.Create)
		propertyManagements.GET("", r.propertyManagementController.List)
		propertyManagements.GET("/:id", r.propertyManagementController.GetByID)
		propertyManagements.PUT("/:id", r.propertyManagementController.Update)
		propertyManagements.DELETE("/:id", r.propertyManagementController.Delete)
	}
}
