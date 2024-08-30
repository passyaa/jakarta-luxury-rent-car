package handlers

import (
	"fmt"
	"net/http"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// ApprovalRequest struct to capture approval or rejection
type ApprovalRequest struct {
	RentalID uint   `json:"rental_id" validate:"required"`
	Action   string `json:"action" validate:"required,oneof=approve reject"`
}

// @Summary Approve or reject a car booking
// @Description Approve or reject a car booking
// @Tags Role Owner
// @Accept json
// @Produce json
// @Param approvalReq body ApprovalRequest true "Approval request body containing rental ID and action (approve/reject)"
// @Success 200 {object} map[string]interface{} "Success message indicating the booking has been approved or rejected"
// @Failure 400 {object} map[string]interface{} "Invalid request format, validation error, car out of stock, or invalid action"
// @Failure 403 {object} map[string]interface{} "Permission denied. Only owners can approve or reject bookings."
// @Failure 404 {object} map[string]interface{} "User, car, or rental history not found"
// @Failure 500 {object} map[string]interface{} "Failed to update car stock, rental history, or send notifications"
// @Router /owner/approve-booking [post]
// @Security BearerAuth
func ApprovalBooking(c echo.Context) error {
	var approvalReq ApprovalRequest

	// Bind the request body to ApprovalRequest struct
	if err := c.Bind(&approvalReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	// Validate the request
	if err := c.Validate(&approvalReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Validation error",
			"error":   err.Error(),
		})
	}

	// Extract user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userID := uint((*claims)["user_id"].(float64))

	// Fetch user role from the database
	var userModel models.User
	if err := database.DB.First(&userModel, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	// Check if the user has the role of "owner"
	if userModel.Role != "owner" {
		return c.JSON(http.StatusForbidden, echo.Map{
			"message": "Permission denied. Only owners can approve or reject bookings.",
		})
	}

	// Fetch the rental history record
	var rentalHistory models.RentalHistory
	if err := database.DB.First(&rentalHistory, approvalReq.RentalID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Rental history not found",
			"error":   err.Error(),
		})
	}

	var carModel models.Car
	if err := database.DB.First(&carModel, rentalHistory.CarID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Car not found",
			"error":   err.Error(),
		})
	}

	var userModelWithRoleUser models.User
	if err := database.DB.First(&userModelWithRoleUser, rentalHistory.UserID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Car not found",
			"error":   err.Error(),
		})
	}

	// Approve or Reject booking based on action
	if approvalReq.Action == "approve" {
		// Approve the booking: set status to "Rent"
		rentalHistory.Status = "Rent"

		// Update the car stock availability
		if carModel.StockAvailability <= 0 {
			return c.JSON(http.StatusBadRequest, echo.Map{
				"message": "Car is out of stock",
			})
		}

		carModel.StockAvailability -= 1

		// Update the car and rental history records in the database
		if err := database.DB.Save(&carModel).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Failed to update car stock",
				"error":   err.Error(),
			})
		}

		if err := database.DB.Save(&rentalHistory).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Failed to approve booking",
				"error":   err.Error(),
			})
		}

		// Send To User
		toPhoneNumberUser := fmt.Sprintf("whatsapp:+%s", userModelWithRoleUser.PhoneNumber)

		messageBodyUser := fmt.Sprintf(
			"Dear %s - %s,\n\n"+
				"Your booking for the car '%s' is confirmed and ready for use. Enjoy your ride!\n\n"+
				"We hope you have a great experience! Please don't forget to leave us a 5-star review!\n\n"+
				"Best regards,\nJakarta Luxury Rent Car",
			userModelWithRoleUser.Email,
			userModelWithRoleUser.Role,
			carModel.Name,
		)

		sendWhatsAppNotification(toPhoneNumberUser, messageBodyUser)

		// Send To Owner
		toPhoneNumberOwner := fmt.Sprintf("whatsapp:+%s", userModel.PhoneNumber)

		messageBodyOwner := fmt.Sprintf(
			"Dear %s - %s,\n\n"+
				"The car '%s' has been successfully booked by '%s' and is ready for use. Please ensure it is in perfect condition for the customer\n",
			userModel.Email,
			userModel.Role,
			carModel.Name,
			userModelWithRoleUser.Email,
		)

		sendWhatsAppNotification(toPhoneNumberOwner, messageBodyOwner)

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Booking approved and car stock updated",
		})

	} else if approvalReq.Action == "reject" {
		// Reject the booking: set status to "Cancel"
		rentalHistory.Status = "Cancel"

		// Update the rental history record in the database
		if err := database.DB.Save(&rentalHistory).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"message": "Failed to reject booking",
				"error":   err.Error(),
			})
		}

		return c.JSON(http.StatusOK, echo.Map{
			"message": "Booking rejected",
		})
	}

	// Return error if action is not recognized
	return c.JSON(http.StatusBadRequest, echo.Map{
		"message": "Invalid action. Only 'approve' or 'reject' are allowed.",
	})
}
