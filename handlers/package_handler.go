package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/labstack/echo/v4"
)

func GetEventPackage(c echo.Context) error {
	var eventPackages []models.EventPackage

	// Query the database for available event packages
	if err := database.DB.Find(&eventPackages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error get event packages"})
	}

	// Return the list of available event packages in JSON format
	return c.JSON(http.StatusOK, eventPackages)
}
