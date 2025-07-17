package controllers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rem-gestion/property-svc/src/models"
	"github.com/rem-gestion/property-svc/src/services"
	"gorm.io/gorm"
)

var db *gorm.DB

func SetDB(database *gorm.DB) {
	db = database
}

func CreateProperty(c *gin.Context) {
	var input models.Property
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "JSON inválido: " + err.Error()})
		return
	}
	// Validaciones básicas
	if input.Type == "" || input.Title == "" || input.AddressID == 0 || input.OwnerID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Faltan campos obligatorios: type, title, address_id, owner_id"})
		return
	}
	if input.Area < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El área no puede ser negativa"})
		return
	}
	if input.Rooms < 0 || input.Bathrooms < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rooms y bathrooms no pueden ser negativos"})
		return
	}
	// Validación real con address-svc y person-svc
	addressSvcURL := os.Getenv("ADDRESS_SVC_URL")
	if addressSvcURL == "" {
		addressSvcURL = "http://localhost:4000"
	}
	personSvcURL := os.Getenv("PERSON_SVC_URL")
	if personSvcURL == "" {
		personSvcURL = "http://localhost:4001"
	}

	addrResp, err := http.Get(fmt.Sprintf("%s/addresses/%d", addressSvcURL, input.AddressID))
	if err != nil || addrResp.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("address_id %d no existe", input.AddressID)})
		return
	}
	ownerResp, err := http.Get(fmt.Sprintf("%s/persons/%d", personSvcURL, input.OwnerID))
	if err != nil || ownerResp.StatusCode != 200 {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("owner_id %d no existe", input.OwnerID)})
		return
	}

	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB not initialized"})
		return
	}
	err = services.CreateProperty(db, &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error al guardar en base de datos: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}
