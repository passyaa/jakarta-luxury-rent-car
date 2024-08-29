package handlers

import (
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

type TopUpRequest struct {
	Amount float64 `json:"deposit_amount"`
}

func TopUp(c echo.Context) error {
	var topUpReq TopUpRequest
	if err := c.Bind(&topUpReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Extract user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userID := uint((*claims)["user_id"].(float64))

	var userModel models.User
	if err := database.DB.First(&userModel, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	userModel.DepositAmount += topUpReq.Amount

	if err := database.DB.Save(&userModel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update deposit amount",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"Email":          userModel.Email,
		"deposit_amount": userModel.DepositAmount,
	})
}

func GetDepositAmount(c echo.Context) error {
	// Extract user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userID := uint((*claims)["user_id"].(float64))

	var userModel models.User
	if err := database.DB.First(&userModel, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"email":          userModel.Email,
		"deposit_amount": userModel.DepositAmount,
	})
}
