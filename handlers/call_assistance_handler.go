package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// CallAssistanceRequest struct to get request body
type CallAssistanceRequest struct {
	RentalID    uint   `json:"rental_id" validate:"required"`
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
}

// @Summary Call assistance request
// @Description Call assistance request
// @Tags Role User
// @Accept json
// @Produce json
// @Param assistanceReq body CallAssistanceRequest true "Call assistance request body containing rental ID, location, and description"
// @Success 200 {object} map[string]interface{} "Success message and details of the call assistance request"
// @Failure 400 {object} map[string]interface{} "Invalid request format or validation error"
// @Failure 404 {object} map[string]interface{} "User, rental history, or car not found"
// @Failure 500 {object} map[string]interface{} "Failed to create call assistance record or send WhatsApp notification"
// @Router /users/call-assistance [post]
// @Security BearerAuth
func CallAssistance(c echo.Context) error {
	var assistanceReq CallAssistanceRequest

	// Bind the request body to CallAssistanceRequest struct
	if err := c.Bind(&assistanceReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid request format",
			"error":   err.Error(),
		})
	}

	// Validate the request
	if err := c.Validate(&assistanceReq); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Validation error",
			"error":   err.Error(),
		})
	}

	// Extract user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userID := uint((*claims)["user_id"].(float64))

	// Retrieve User model based on userID
	var userModel models.User
	if err := database.DB.Where("user_id = ?", userID).First(&userModel).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	// Retrieve RentalHistory model based on userID and RentailID
	var rentalHistory models.RentalHistory
	if err := database.DB.Where("rental_id = ? AND user_id = ?", assistanceReq.RentalID, userID).First(&rentalHistory).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Rental history not found or does not belong to the user",
			"error":   err.Error(),
		})
	}

	// Retrieve Car model based on RentailID.CarID
	var car models.Car
	if err := database.DB.First(&car, rentalHistory.CarID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Car not found for rental history",
			"error":   err.Error(),
		})
	}

	callAssistance := models.CallAssistance{
		RentalID:           assistanceReq.RentalID,
		UserID:             userID,
		CallAssistanceDate: time.Now(),
		Description:        assistanceReq.Description,
		Location:           assistanceReq.Location,
	}

	// Insert db callAssistance
	if err := database.DB.Create(&callAssistance).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to create call assistance record",
			"error":   err.Error(),
		})
	}

	toPhoneNumber := fmt.Sprintf("whatsapp:+%s", userModel.PhoneNumber)

	gmapsLink := "https://www.google.com/maps/search/?api=1&query=" + url.QueryEscape(callAssistance.Location)

	messageBody := "Subject: Call Assistance Request - [Rental ID: " + strconv.Itoa(int(rentalHistory.RentalID)) + "]\n\n" +
		"Dear" + userModel.Email + " - " + userModel.Role + ",\n\n" +
		"You received call assistance. Below are the details of user request:\n\n" +
		"User Details:\n" +
		"  - Email: " + userModel.Email + "\n" +
		"  - Phone Number: " + userModel.PhoneNumber + "\n\n" +
		"Rental Details:\n" +
		"  - Rental ID: " + strconv.Itoa(int(rentalHistory.RentalID)) + "\n" +
		"  - Car Name: " + car.Name + "\n" +
		"  - Car Category: " + car.Category + "\n" +
		"  - Car Brand: " + car.Make + "\n" +
		"  - Car Model: " + car.Model + "\n" +
		"  - Car Transmission: " + car.Transmission + "\n" +
		"  - Car Year: " + strconv.Itoa(car.Year) + "\n" +
		"  - Car Fuel Type: " + car.FuelType + "\n" +
		"  - Car Class: " + car.Class + "\n\n" +
		"Assistance Request Details:\n" +
		"  - Date: " + callAssistance.CallAssistanceDate.Format("02 January 2006 15:04:05") + "\n" +
		"  - Location: " + callAssistance.Location + "\n" +
		"  - Link to Location: " + gmapsLink + "\n" +
		"  - Description: " + callAssistance.Description

	// fmt.Println("toPhoneNumber", toPhoneNumber)
	// fmt.Println("gmapsLink", gmapsLink)
	// fmt.Println("messageBody", messageBody)

	// sendWhatsAppNotification
	err := sendWhatsAppNotification(toPhoneNumber, messageBody)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to send WhatsApp notification",
			"error":   err.Error(),
		})
	}

	// Return success response with car and call assistance details
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Call assistance request created successfully",
		"data": echo.Map{
			"assistance_id":       callAssistance.CallAssistanceID,
			"user_email":          userModel.Email,
			"user_phone_number":   userModel.PhoneNumber,
			"rental_id":           rentalHistory.RentalID,
			"car_name":            car.Name,
			"car_category":        car.Category,
			"car_make":            car.Make,
			"car_model":           car.Model,
			"car_transmission":    car.Transmission,
			"car_year":            car.Year,
			"car_fuel_type":       car.FuelType,
			"car_class":           car.Class,
			"callassistance_date": callAssistance.CallAssistanceDate,
			"location":            callAssistance.Location,
			"link_location":       gmapsLink,
			"description":         callAssistance.Description,
		},
	})
}
