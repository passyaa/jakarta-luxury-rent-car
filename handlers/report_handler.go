package handlers

import (
	"fmt"
	"net/http"
	"time"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// RentalReportResponse struct to structure the report output
type RentalReportResponse struct {
	Email              string    `json:"email"`
	PhoneNumber        string    `json:"phone_number"`
	Address            string    `json:"address"`
	CarName            string    `json:"car_name"`
	CarCategory        string    `json:"car_category"`
	CarMake            string    `json:"car_make"`
	CarModel           string    `json:"car_model"`
	CarTransmission    string    `json:"car_transmission"`
	CarYear            int       `json:"car_year"`
	CarFuelType        string    `json:"car_fuel_type"`
	CarClass           string    `json:"car_class"`
	DriverName         string    `json:"driver_name,omitempty"`
	DriverPhoneNumber  string    `json:"driver_phone_number,omitempty"`
	PackageName        string    `json:"package_name,omitempty"`
	PackageDescription string    `json:"package_description,omitempty"`
	RentalDate         time.Time `json:"rental_date"`
	ReturnDate         time.Time `json:"return_date"`
	PickupLocation     string    `json:"pickup_location"`
	DropoffLocation    string    `json:"dropoff_location"`
	Duration           string    `json:"duration"`
	CostDetails        string    `json:"cost_details"`
	Status             string    `json:"status"`
	ConciergeServices  bool      `json:"concierge_services"`
	AirportTransfer    bool      `json:"airport_transfer"`
}

// @Summary Generate rental reports
// @Description Generate detailed rental history reports for the owner
// @Tags Role Owner
// @Accept json
// @Produce json
// @Success 200 {array} RentalReportResponse "List of rental reports"
// @Failure 403 {object} map[string]interface{} "Permission denied. Only owners can approve or reject bookings."
// @Failure 404 {object} map[string]interface{} "User or related entity not found"
// @Failure 500 {object} map[string]interface{} "Failed to fetch rental histories or related data"
// @Router /owner/report [get]
// @Security BearerAuth

func Report(c echo.Context) error {
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

	var rentalHistories []models.RentalHistory
	var reports []RentalReportResponse

	if err := database.DB.Find(&rentalHistories).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to fetch rental histories",
			"error":   err.Error(),
		})
	}

	// Loop through each rental history to build detailed report
	for _, rental := range rentalHistories {
		var user models.User
		var car models.Car
		var driver models.Driver
		var eventPackage models.EventPackage

		if err := database.DB.First(&user, rental.UserID).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "User not found for rental history",
				"error":   err.Error(),
			})
		}

		if err := database.DB.First(&car, rental.CarID).Error; err != nil {
			return c.JSON(http.StatusNotFound, echo.Map{
				"message": "Car not found for rental history",
				"error":   err.Error(),
			})
		}

		driverName := ""
		driverPhoneNumber := ""
		if rental.DriverID != nil {
			if err := database.DB.First(&driver, *rental.DriverID).Error; err == nil {
				driverName = driver.Name
				driverPhoneNumber = driver.PhoneNumber
			}
		}

		packageName := ""
		packageDescription := ""
		if rental.PackageID != nil {
			if err := database.DB.First(&eventPackage, *rental.PackageID).Error; err == nil {
				packageName = eventPackage.PackageName
				packageDescription = eventPackage.Description
			}
		}

		durationDays := rental.ReturnDate.Sub(rental.RentalDate).Hours() / 24
		duration := fmt.Sprintf("%.0f days", durationDays)

		costDetails := fmt.Sprintf("Total Cost for %s: %.2f", duration, car.RentalCosts*durationDays)
		if rental.AirportTransfer {
			costDetails += " + Airport Transfer: 50"
		}
		if rental.ConciergeServices {
			costDetails += " + Concierge Services: 100"
		}
		if rental.DriverID != nil {
			costDetails += fmt.Sprintf(" + Driver: %.2f", durationDays*100)
		}
		if rental.PackageID != nil {
			costDetails += fmt.Sprintf(" + Package: %.2f", eventPackage.Cost)
		}
		costDetails += fmt.Sprintf(" = Total: %.2f", rental.TotalCost)

		// Build the report entry
		report := RentalReportResponse{
			Email:              user.Email,
			PhoneNumber:        user.PhoneNumber,
			Address:            user.Address,
			CarName:            car.Name,
			CarCategory:        car.Category,
			CarMake:            car.Make,
			CarModel:           car.Model,
			CarTransmission:    car.Transmission,
			CarYear:            car.Year,
			CarFuelType:        car.FuelType,
			CarClass:           car.Class,
			DriverName:         driverName,
			DriverPhoneNumber:  driverPhoneNumber,
			PackageName:        packageName,
			PackageDescription: packageDescription,
			RentalDate:         rental.RentalDate,
			ReturnDate:         *rental.ReturnDate,
			PickupLocation:     rental.PickupLocation,
			DropoffLocation:    rental.DropoffLocation,
			Duration:           duration,
			CostDetails:        costDetails,
			Status:             rental.Status,
			ConciergeServices:  rental.ConciergeServices,
			AirportTransfer:    rental.AirportTransfer,
		}

		// Add report entry to the reports slice
		reports = append(reports, report)
	}

	return c.JSON(http.StatusOK, reports)
}
