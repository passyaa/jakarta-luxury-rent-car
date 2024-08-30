package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/labstack/echo/v4"
)

// @Summary Get drivers
// @Description Get drivers
// @Tags Public
// @Accept json
// @Produce json
// @Success 200 {array} models.Driver "List of available drivers"
// @Failure 500 {object} map[string]string "Error get drivers"
// @Router /drivers [get]
func GetDriver(c echo.Context) error {
	var drivers []models.Driver

	// Query the database for available drivers
	if err := database.DB.Find(&drivers).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": "Error get drivers"})
	}

	return c.JSON(http.StatusOK, drivers)
}
