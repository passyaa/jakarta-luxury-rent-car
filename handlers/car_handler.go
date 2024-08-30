package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/labstack/echo/v4"
)

// @Summary Get available luxury cars
// @Description Get available luxury cars
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {array} models.Car "List of available luxury cars"
// @Failure 500 {object} map[string]interface{} "Error retrieving cars from database"
// @Router /cars [get]
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
