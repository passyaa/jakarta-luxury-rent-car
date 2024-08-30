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

// @Summary MakingPayment updates status from "Book" to "Paid"
// @Description MakingPayment updates status from "Book" to "Paid"
// @Tags Role User
// @Accept json
// @Produce json
// @Param paymentReq body PaymentRequest true "Payment request body containing rental ID and payment details"
// @Success 200 {object} map[string]string "Rental status updated to 'Paid' successfully"
// @Failure 400 {object} map[string]string "Invalid request body or rental status is not 'Book'"
// @Failure 404 {object} map[string]string "Rental history or user not found"
// @Failure 500 {object} map[string]string "Failed to update rental status or deposit amount"
// @Router /users/making-payment [post]
// @Security BearerAuth
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

	var userModel models.User
	if err := database.DB.Where("user_id = ?", userID).First(&userModel).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	// update data models.User terhadap DepositAmount
	currentAmount := userModel.DepositAmount - rentalHistory.TotalCost
	userModel.DepositAmount = currentAmount

	// Save the updated user model back to the database
	if err := database.DB.Save(&userModel).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update deposit amount",
			"error":   err.Error(),
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
