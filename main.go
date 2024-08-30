package main

import (
	"log"
	"os"

	_ "jakarta-luxury-rent-car/docs"

	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/handlers"
	"jakarta-luxury-rent-car/middlewares"
)

// struct untuk validator custom
type CustomValidator struct {
	validator *validator.Validate
}

// memvalidasi input menggunakan validator custom
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// @title Jakarta Luxury Rent Car API
// @version 1.0
// @description This is Jakarta Luxury Rent Car service API documentation.

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @type apiKey
// @in header
// @name Authorization

// @schemes http https
func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	// Initialize the database
	database.InitDB()

	// Create a new Echo instance
	e := echo.New()

	// Set custom validator
	e.Validator = &CustomValidator{validator: validator.New()} // Register custom validator

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/register", handlers.RegisterUser)
	e.POST("/login", handlers.LoginUser)
	e.GET("/cars", handlers.GetLuxuryCars)
	e.GET("/drivers", handlers.GetDriver)
	e.GET("/packages", handlers.GetEventPackage)

	// Secure Routes
	r := e.Group("")
	r.Use(middlewares.JWTMiddleware())
	r.POST("/users/register-membership", handlers.RegisterMembership)
	r.GET("/users/get-membership", handlers.GetMembership)
	r.GET("/users/get-deposit", handlers.GetDepositAmount)
	r.POST("/users/topup", handlers.TopUp)
	r.POST("/users/booking", handlers.BookCar)
	r.POST("/users/making-payment", handlers.MakingPayment)
	r.POST("/users/call-assistance", handlers.CallAssistance)

	r.POST("/owner/approve-booking", handlers.ApprovalBooking)
	r.GET("/owner/report", handlers.Report)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified
	}
	e.Logger.Fatal(e.Start(":" + port))
}
