package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/labstack/echo/v4"
)

func GetLuxuryCars(c echo.Context) error {
	var cars []models.Car

	// Query the database for available cars
	if err := database.DB.Where("stock_availability > ?", 0).Find(&cars).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Error retrieving cars from database",
			"error":   err.Error(),
		})
	}

	// Return the list of available luxury cars in JSON format
	return c.JSON(http.StatusOK, cars)
}
