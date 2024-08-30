package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// BookingRequest struct to capture user inputs
type BookingRequest struct {
	CarID             uint      `json:"car_id" validate:"required"`
	DriverID          *uint     `json:"driver_id" validate:"required"`
	PackageID         *uint     `json:"package_id" validate:"required"`
	RentalDate        time.Time `json:"rental_date" validate:"required"`
	ReturnDate        time.Time `json:"return_date" validate:"required"`
	PickupLocation    string    `json:"pickup_location"`  // Optional
	DropoffLocation   string    `json:"dropoff_location"` // Optional
	RentalDuration    string    `json:"rental_duration" validate:"required,oneof=daily weekly monthly"`
	AirportTransfer   bool      `json:"airport_transfer"`   // Optional
	ConciergeServices bool      `json:"concierge_services"` // Optional
}

// @Summary Book a car
// @Description Book a car
// @Tags Role User
// @Accept json
// @Produce json
// @Param bookingReq body BookingRequest true "Booking request body containing car ID and other booking details"
// @Success 200 {object} map[string]interface{} "Success message and details of the car booking"
// @Failure 400 {object} map[string]string "Invalid request format, validation error, or insufficient stock"
// @Failure 404 {object} map[string]string "Car, driver, or package not found"
// @Failure 500 {object} map[string]string "Failed to create rental history, send notifications, or process booking"
// @Router /users/booking [post]
// @Security BearerAuth
func BookCar(c echo.Context) error {
	var bookingReq BookingRequest

	// Bind and validate the request body to BookingRequest struct
	if err := c.Bind(&bookingReq); err != nil {
		return jsonResponse(c, http.StatusBadRequest, "Invalid request format", err.Error())
	}
	if err := c.Validate(&bookingReq); err != nil {
		return jsonResponse(c, http.StatusBadRequest, "Validation error", err.Error())
	}

	// Extract user ID from JWT token
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwt.MapClaims)
	userID := uint((*claims)["user_id"].(float64))

	// Check if the selected car is available
	car, err := getCarByID(bookingReq.CarID)
	if err != nil {
		return jsonResponse(c, http.StatusNotFound, "Car not found", err.Error())
	}

	// Check if there is enough stock to reduce
	if car.StockAvailability < 1 {
		return jsonResponse(c, http.StatusBadRequest, "Insufficient stock", fmt.Sprintf("Current stock: %d", car.StockAvailability))
	}

	// checks for driver and package
	driver, eventPackage, err := getDriverandPackage(bookingReq.DriverID, bookingReq.PackageID)
	if err != nil {
		return jsonResponse(c, http.StatusNotFound, "driver and eventPackage not found", err.Error())
	}

	// Calculate the total cost
	totalCost := calculateTotalCost(bookingReq, car, eventPackage)

	// Apply discount based on membership level
	totalCost = applyMembershipDiscount(userID, totalCost)

	// Prepare data rental history entry
	rentalHistory := createRentalHistoryEntry(bookingReq, userID, totalCost)

	// Save the rental history to the database
	if err := database.DB.Create(&rentalHistory).Error; err != nil {
		return jsonResponse(c, http.StatusInternalServerError, "Failed to create rental history", err.Error())
	}

	// // // Send confirmation notification
	if err := sendConfirmationNotification(userID, rentalHistory, car, driver, eventPackage); err != nil {

		return jsonResponse(c, http.StatusInternalServerError, "Failed to send WhatsApp notification", err.Error())
	}

	// // Send Invoice notification
	if err := CreateInvoiceAndSendWhatsApp(userID, bookingReq, rentalHistory, car, driver, eventPackage); err != nil {
		return jsonResponse(c, http.StatusInternalServerError, "Failed to send Invoice notification", err.Error())
	}

	// Return success response
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Car booking successfully created",
		"data":    rentalHistory,
	})
}

func jsonResponse(c echo.Context, statusCode int, message, errorDetail string) error {
	return c.JSON(statusCode, echo.Map{
		"message": message,
		"error":   errorDetail,
	})
}

func getCarByID(carID uint) (models.Car, error) {
	var car models.Car
	if err := database.DB.First(&car, carID).Error; err != nil {
		return car, err
	}
	return car, nil
}

func getDriverandPackage(driverID, packageID *uint) (models.Driver, models.EventPackage, error) {
	var driver models.Driver
	var eventPackage models.EventPackage

	if driverID != nil {
		if err := database.DB.First(&driver, *driverID).Error; err != nil {
			return driver, eventPackage, err
		}
	}
	if packageID != nil {
		if err := database.DB.First(&eventPackage, *packageID).Error; err != nil {
			return driver, eventPackage, err
		}
	}

	return driver, eventPackage, nil
}

func calculateTotalCost(bookingReq BookingRequest, car models.Car, eventPackage models.EventPackage) float64 {
	// Calculate rental duration
	rentalDays := bookingReq.ReturnDate.Sub(bookingReq.RentalDate).Hours() / 24

	totalCost := rentalDays * car.RentalCosts

	if bookingReq.DriverID != nil {
		totalCost += rentalDays * 100
	}

	if bookingReq.PackageID != nil {
		totalCost += eventPackage.Cost
	}

	// Add costs for optional services
	if bookingReq.PickupLocation != "" && bookingReq.DropoffLocation != "" {
		totalCost += 100
	}

	if bookingReq.AirportTransfer {
		totalCost += 50
	}

	if bookingReq.ConciergeServices {
		totalCost += 100
	}

	return totalCost
}

func applyMembershipDiscount(userID uint, totalCost float64) float64 {
	var membership models.Membership
	if err := database.DB.Where("user_id = ?", userID).First(&membership).Error; err == nil {
		// Apply discount based on discount level
		switch membership.DiscountLevel {
		case "Silver":
			totalCost *= 0.90 // 10% discount
		case "Gold":
			totalCost *= 0.80 // 20% discount
		case "Platinum":
			totalCost *= 0.70 // 30% discount
		}
	}

	return totalCost
}

func createRentalHistoryEntry(bookingReq BookingRequest, userID uint, totalCost float64) models.RentalHistory {
	return models.RentalHistory{
		UserID:            userID,
		CarID:             bookingReq.CarID,
		DriverID:          bookingReq.DriverID,
		PackageID:         bookingReq.PackageID,
		RentalDate:        bookingReq.RentalDate,
		ReturnDate:        &bookingReq.ReturnDate,
		PickupLocation:    bookingReq.PickupLocation,
		DropoffLocation:   bookingReq.DropoffLocation,
		TotalCost:         totalCost,
		Status:            "Book",
		AirportTransfer:   bookingReq.AirportTransfer,
		ConciergeServices: bookingReq.ConciergeServices,
	}
}

func sendConfirmationNotification(userID uint, rentalHistory models.RentalHistory, car models.Car, driver models.Driver, eventPackage models.EventPackage) error {
	var userModel models.User
	if err := database.DB.Where("user_id = ?", userID).First(&userModel).Error; err != nil {
		return err
	}

	toPhoneNumber := fmt.Sprintf("whatsapp:+%s", userModel.PhoneNumber)

	// Format the message body
	messageBody := fmt.Sprintf(
		"Subject: Booking Confirmation - [Rental ID: %d]\n\n"+
			"Dear %s - %s,\n\n"+
			"Congratulations! Your booking has been successfully confirmed. Below are the details of your booking:\n\n"+
			"User Details:\n"+
			"  - Email: %s\n"+
			"  - Phone Number: %s\n\n"+
			"Rental Details:\n"+
			"  - Rental ID: %d\n"+
			"  - Car Name: %s\n"+
			"  - Car Category: %s\n"+
			"  - Car Brand: %s\n"+
			"  - Car Model: %s\n"+
			"  - Car Transmission: %s\n"+
			"  - Car Year: %d\n"+
			"  - Car Fuel Type: %s\n"+
			"  - Car Class: %s\n\n"+
			"Booking Details:\n"+
			"  - Rental Date: %s\n"+
			"  - Return Date: %s\n"+
			"  - Pickup Location: %s\n"+
			"  - Dropoff Location: %s\n"+
			"  - Total Cost: %.2f\n"+
			"  - Airport Transfer: %t\n"+
			"  - Concierge Services: %t\n\n"+
			"Driver Details:\n"+
			"  - Driver Name: %s\n"+
			"  - Driver Contact: %s\n\n"+
			"Package Details:\n"+
			"  - Package Name: %s\n"+
			"  - Package Description: %s\n\n"+
			"Thank you for choosing our service! We look forward to serving you.\n\n"+
			"Best regards,\nJakarta Luxury Rent Car",
		rentalHistory.RentalID,
		userModel.Email,
		userModel.Role,
		userModel.Email,
		userModel.PhoneNumber,
		rentalHistory.RentalID,
		car.Name,
		car.Category,
		car.Make,
		car.Model,
		car.Transmission,
		car.Year,
		car.FuelType,
		car.Class,
		rentalHistory.RentalDate.Format("02 January 2006 15:04"),
		rentalHistory.ReturnDate.Format("02 January 2006 15:04"),
		rentalHistory.PickupLocation,
		rentalHistory.DropoffLocation,
		rentalHistory.TotalCost,
		rentalHistory.AirportTransfer,
		rentalHistory.ConciergeServices,
		driver.Name,
		driver.PhoneNumber,
		eventPackage.PackageName,
		eventPackage.Description,
	)

	return sendWhatsAppNotification(toPhoneNumber, messageBody)
}

func CreateInvoiceAndSendWhatsApp(userID uint, bookingReq BookingRequest, rentalHistory models.RentalHistory, car models.Car, driver models.Driver, eventPackage models.EventPackage) error {

	var costDriver, costPackage, costPickup, costTranferAirport, costConcierge, discountAmount float64

	var userModel models.User
	if err := database.DB.Where("user_id = ?", userID).First(&userModel).Error; err != nil {
		return err
	}

	var membership models.Membership
	if err := database.DB.Where("user_id = ?", userID).First(&membership).Error; err == nil {
		// Apply discount based on discount level
		switch membership.DiscountLevel {
		case "Silver":
			originalCost := rentalHistory.TotalCost / 0.90 // 10% discount
			discountAmount = originalCost * 0.10
			fmt.Println("Go originalCost", originalCost)
			fmt.Println("Go discountAmount", discountAmount)
		case "Gold":
			originalCost := rentalHistory.TotalCost / 0.20 // 20% discount
			discountAmount = originalCost * 0.20
			fmt.Println("Go originalCost", originalCost)
			fmt.Println("Go discountAmount", discountAmount)
		case "Platinum":
			originalCost := rentalHistory.TotalCost / 0.30 // 30% discount
			discountAmount = originalCost * 0.30
			fmt.Println("Go originalCost", originalCost)
			fmt.Println("Go discountAmount", discountAmount)
		}
	}

	rentalDays := bookingReq.ReturnDate.Sub(bookingReq.RentalDate).Hours() / 24

	if bookingReq.DriverID != nil {
		costDriver = 100
	}
	if bookingReq.PackageID != nil {
		costPackage = eventPackage.Cost
	}

	// Add costs for optional services
	if bookingReq.PickupLocation != "" && bookingReq.DropoffLocation != "" {
		costPickup = 100
	}

	if bookingReq.AirportTransfer {
		costTranferAirport = 50
	}

	if bookingReq.ConciergeServices {
		costConcierge = 100
	}

	// Xendit API URL
	url := "https://api.xendit.co/v2/invoices"

	// Create request body
	requestBody := map[string]interface{}{
		"external_id":      "Invoice Jakarta Luxury Car For : " + car.Name,
		"amount":           rentalHistory.TotalCost,
		"description":      "Invoice Jakarta Luxury Car To : " + userModel.Email + " - " + userModel.PhoneNumber,
		"invoice_duration": 86400,
		"currency":         "IDR",
		"customer": map[string]interface{}{
			"email":         userModel.Email,
			"mobile_number": userModel.PhoneNumber,
		},
		"items": []map[string]interface{}{
			{
				"name":     car.Name,
				"quantity": rentalDays,
				"price":    car.RentalCosts,
				"category": car.Category,
			},
			{
				"name":     "Driver",
				"quantity": rentalDays,
				"price":    costDriver,
			},
			{
				"name":     "Event Package",
				"quantity": 1,
				"price":    costPackage,
			},
			{
				"name":     "Pickup and Dropoff",
				"quantity": 1,
				"price":    costPickup,
			},
			{
				"name":     "Transfer Airport",
				"quantity": 1,
				"price":    costTranferAirport,
			},
			{
				"name":     "Concierge Service",
				"quantity": 1,
				"price":    costConcierge,
			},
		},
		"fees": []map[string]interface{}{
			{
				"type":  "Discount Level - " + membership.DiscountLevel,
				"value": discountAmount,
			},
		},
	}

	// Convert request body to JSON
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set Content-Type header
	req.Header.Set("Content-Type", "application/json")

	// Set basic authentication header
	username := os.Getenv("API_KEY_XENDIT")
	req.SetBasicAuth(username, "")

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send HTTP request: %v", err)
	}
	defer resp.Body.Close()

	// Check if the response status is 200 OK
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("received non-200 response: %s - %s", resp.Status, string(bodyBytes))
	}

	// Read the response body
	var responseBody map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("failed to decode response body: %v", err)
	}

	// Extract invoice URL or ID for WhatsApp message
	invoiceURL, ok := responseBody["invoice_url"].(string)
	if !ok {
		return fmt.Errorf("invoice URL not found in response")
	}

	toPhoneNumber := fmt.Sprintf("whatsapp:+%s", userModel.PhoneNumber)

	messageBody := fmt.Sprintf(
		"Dear %s - %s,\n\nThank you for using Jakarta Luxury Car Rental. Please find your invoice at the following link:\n%s\n\nKindly complete the payment within the next 24 hours. If you have any questions, feel free to contact us.\n\nBest regards,\nJakarta Luxury Car Rental",
		userModel.Email,
		userModel.Role,
		invoiceURL,
	)

	return sendWhatsAppNotification(toPhoneNumber, messageBody)
}
