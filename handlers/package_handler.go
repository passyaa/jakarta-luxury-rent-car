package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/labstack/echo/v4"
)

// @Summary Get available event packages
// @Description Retrieve a list of all available event packages
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {array} models.EventPackage "List of available event packages"
// @Failure 500 {object} map[string]string "Error get event packages"
// @Router /packages [get]
func GetEventPackage(c echo.Context) error {
	var eventPackages []models.EventPackage

	// Query the database for available event packages
	if err := database.DB.Find(&eventPackages).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error get event packages"})
	}

	// Return the list of available event packages in JSON format
	return c.JSON(http.StatusOK, eventPackages)
}
