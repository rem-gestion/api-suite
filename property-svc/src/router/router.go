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
	propertyListingController    *controllers.PropertyListingController
	propertyMediaController      *controllers.PropertyMediaController
	propertyValuationController  *controllers.PropertyValuationController
}

func New(
	propertyController *controllers.PropertyController,
	propertyTypeController *controllers.PropertyTypeController,
	managerTypeController *controllers.ManagerTypeController,
	amenityController *controllers.AmenityController,
	propertyManagementController *controllers.PropertyManagementController,
	propertyListingController *controllers.PropertyListingController,
	propertyMediaController *controllers.PropertyMediaController,
	propertyValuationController *controllers.PropertyValuationController,
) *Router {
	return &Router{
		propertyController:           propertyController,
		propertyTypeController:       propertyTypeController,
		managerTypeController:        managerTypeController,
		amenityController:            amenityController,
		propertyManagementController: propertyManagementController,
		propertyListingController:    propertyListingController,
		propertyMediaController:      propertyMediaController,
		propertyValuationController:  propertyValuationController,
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

		// Property Listing routes by property (using :id instead of :property_id)
		properties.GET("/:id/listings", r.propertyListingController.GetListingsByPropertyID)

		// Property Media routes by property (using :id instead of :property_id)
		properties.GET("/:id/media", r.propertyMediaController.GetMediaByPropertyID)

		// Property Valuation routes by property (using :id instead of :property_id)
		properties.GET("/:id/valuations", r.propertyValuationController.GetValuationsByPropertyID)
		properties.GET("/:id/valuations/latest", r.propertyValuationController.GetLatestValuation)
		properties.GET("/:id/valuations/type/:type", r.propertyValuationController.GetValuationsByType)

		// Property Listing routes under /properties/listings (as expected by Postman)
		listings := properties.Group("/listings")
		{
			listings.POST("", r.propertyListingController.CreateListing)
			listings.GET("", r.propertyListingController.ListListings)
			listings.GET("/search", r.propertyListingController.SearchListings)
			listings.GET("/:id", r.propertyListingController.GetListingByID)
			listings.PUT("/:id", r.propertyListingController.UpdateListing)
			listings.DELETE("/:id", r.propertyListingController.DeleteListing)
		}

		// Property Media routes under /properties/media (as expected by Postman)
		media := properties.Group("/media")
		{
			media.POST("", r.propertyMediaController.CreateMedia)
			media.GET("", r.propertyMediaController.ListMedia)
			media.GET("/:id", r.propertyMediaController.GetMediaByID)
			media.PUT("/:id", r.propertyMediaController.UpdateMedia)
			media.DELETE("/:id", r.propertyMediaController.DeleteMedia)
			media.PUT("/:id/set-main", r.propertyMediaController.SetMainImage)
			media.PUT("/display-order", r.propertyMediaController.UpdateDisplayOrder)
		}

		// Property Valuation routes under /properties/valuations (as expected by Postman)
		valuations := properties.Group("/valuations")
		{
			valuations.POST("", r.propertyValuationController.CreateValuation)
			valuations.GET("", r.propertyValuationController.ListValuations)
			valuations.GET("/:id", r.propertyValuationController.GetValuationByID)
			valuations.PUT("/:id", r.propertyValuationController.UpdateValuation)
			valuations.DELETE("/:id", r.propertyValuationController.DeleteValuation)
		}
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

	// Property Listing routes
	propertyListings := engine.Group("/property-listings")
	{
		propertyListings.POST("", r.propertyListingController.CreateListing)
		propertyListings.GET("", r.propertyListingController.ListListings)
		propertyListings.GET("/search", r.propertyListingController.SearchListings)
		propertyListings.GET("/:id", r.propertyListingController.GetListingByID)
		propertyListings.PUT("/:id", r.propertyListingController.UpdateListing)
		propertyListings.DELETE("/:id", r.propertyListingController.DeleteListing)
		propertyListings.POST("/:id/views", r.propertyListingController.IncrementViews)
		propertyListings.POST("/:id/favorites/increment", r.propertyListingController.IncrementFavorites)
		propertyListings.POST("/:id/favorites/decrement", r.propertyListingController.DecrementFavorites)
	}

	// Property Media routes
	propertyMedia := engine.Group("/property-media")
	{
		propertyMedia.POST("", r.propertyMediaController.CreateMedia)
		propertyMedia.GET("", r.propertyMediaController.ListMedia)
		propertyMedia.GET("/:id", r.propertyMediaController.GetMediaByID)
		propertyMedia.PUT("/:id", r.propertyMediaController.UpdateMedia)
		propertyMedia.DELETE("/:id", r.propertyMediaController.DeleteMedia)
		propertyMedia.PUT("/:id/set-main", r.propertyMediaController.SetMainImage)
		propertyMedia.PUT("/display-order", r.propertyMediaController.UpdateDisplayOrder)
	}

	// Property Valuation routes
	propertyValuations := engine.Group("/property-valuations")
	{
		propertyValuations.POST("", r.propertyValuationController.CreateValuation)
		propertyValuations.GET("", r.propertyValuationController.ListValuations)
		propertyValuations.GET("/:id", r.propertyValuationController.GetValuationByID)
		propertyValuations.PUT("/:id", r.propertyValuationController.UpdateValuation)
		propertyValuations.DELETE("/:id", r.propertyValuationController.DeleteValuation)
	}
}
