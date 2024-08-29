package handlers

import (
	"fmt"
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// PaymentRequest is the struct for the incoming JSON request body
type PaymentRequest struct {
	RentalID uint `json:"rental_id"`
}

// MakingPayment updates the status from "Book" to "Paid"
func MakingPayment(c echo.Context) error {
	// Get the user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userID := uint((*claims)["user_id"].(float64))

	// Bind the incoming JSON to the PaymentRequest struct
	var paymentReq PaymentRequest
	if err := c.Bind(&paymentReq); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	// Find the rental history entry by RentalID and UserID
	var rentalHistory models.RentalHistory
	if err := database.DB.Where("rental_id = ? AND user_id = ?", paymentReq.RentalID, userID).First(&rentalHistory).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Rental history not found",
		})
	}

	// Check if the current status is "Book"
	if rentalHistory.Status != "Book" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Rental status is not 'Book', cannot proceed to 'Paid'",
		})
	}

	// Update the status to "Paid"
	rentalHistory.Status = "Paid"
	if err := database.DB.Save(&rentalHistory).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to update rental status",
		})
	}

	var userOwnerModel models.User
	if err := database.DB.Where("role = ?", "owner").First(&userOwnerModel).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Owner not found",
			"error":   err.Error(),
		})
	}

	// Send To Owner
	toPhoneNumberOwner := fmt.Sprintf("whatsapp:+%s", userOwnerModel.PhoneNumber)

	messageBodyOwner := fmt.Sprintf(
		"Dear %s - %s,\n\nPayment for Rental ID - %d has been successfully completed!, please approve the process",
		userOwnerModel.Email,
		userOwnerModel.Role,
		rentalHistory.RentalID,
	)

	sendWhatsAppNotification(toPhoneNumberOwner, messageBodyOwner)

	// Return a success response
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Rental status updated to 'Paid' successfully",
	})
}
