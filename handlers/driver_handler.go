package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/labstack/echo/v4"
)

func GetDriver(c echo.Context) error {
	var drivers []models.Driver

	// Query the database for available drivers
	if err := database.DB.Find(&drivers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error get drivers"})
	}

	return c.JSON(http.StatusOK, drivers)
}
