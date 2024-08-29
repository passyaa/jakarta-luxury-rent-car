package handlers

import (
	"math/rand"
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func RegisterMembership(c echo.Context) error {
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

	// Check if the user already has a membership
	var existingMembership models.Membership
	if err := database.DB.Where("user_id = ?", userID).First(&existingMembership).Error; err == nil {
		// Membership exists, return an error response
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "User already has a membership registered",
		})
	}

	discountLevels := []string{"Silver", "Gold", "Platinum"}
	randomIndex := rand.Intn(len(discountLevels))
	randomDiscountLevel := discountLevels[randomIndex]

	newMembership := models.Membership{
		UserID:        userID,
		DiscountLevel: randomDiscountLevel,
	}

	// Save the membership 
	if err := database.DB.Create(&newMembership).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to register membership",
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"membership_id":  newMembership.MembershipID,
		"user_id":        newMembership.UserID,
		"email":          userModel.Email,
		"discount_level": newMembership.DiscountLevel,
	})
}


func GetMembership(c echo.Context) error {
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

	var membership models.Membership
	if err := database.DB.Where("user_id = ?", userID).First(&membership).Error; err != nil {
		// Membership does not exist, ask user to register
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "No membership found. Please register for a membership first.",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"membership_id":  membership.MembershipID,
		"user_id":        membership.UserID,
		"email":          userModel.Email,
		"discount_level": membership.DiscountLevel,
	})
}
