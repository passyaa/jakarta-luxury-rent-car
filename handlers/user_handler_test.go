package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"jakarta-luxury-rent-car/models"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func setupTestingInitDB() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		"localhost",
		"postgres",
		"Jakarta1!",
		"jakarta-luxury-rent-car",
		"5432")

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	if DB == nil {
		log.Fatal("Database connection is nil after initialization")
	} else {
		fmt.Println("Database initialized successfully.")
	}

	err = DB.AutoMigrate(
		&models.User{},
		&models.Car{},
		&models.Driver{},
		&models.EventPackage{},
		&models.RentalHistory{},
		&models.CallAssistance{},
		&models.Membership{},
	)

	if err != nil {
		log.Fatal("Failed to auto migrate: ", err)
	}

	fmt.Println("Database migrations completed successfully.")
}

// CustomValidator is a custom validator for Echo using go-playground/validator
type CustomValidator struct {
	validator *validator.Validate
}

func (cv *CustomValidator) Validate(i interface{}) error {
	if cv.validator == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Validator is not initialized")
	}
	if err := cv.validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}

func TestRegisterUser_Success(t *testing.T) {
	setupTestingInitDB()

	if DB == nil {
		t.Fatalf("Database is not initialized")
	}

	e := echo.New()

	validatorInstance := validator.New()
	if validatorInstance == nil {
		t.Fatalf("Validator is not initialized")
	}
	e.Validator = &CustomValidator{validator: validatorInstance}

	reqBody := `{"email":"test10@example.com","password":"testtest","phone_number":"1234567890","address":"Jalan 123", "role":"user"}`

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewReader([]byte(reqBody)))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	fmt.Println("Echo context initialized successfully.")

	err := RegisterUser(c)
	if err != nil {
		t.Fatalf("Handler returned error: %v", err)
	}

	if rec.Code != http.StatusCreated {
		t.Fatalf("Expected status code %d, got %d", http.StatusCreated, rec.Code)
	}

	var userResponse UserResponse

	err = json.Unmarshal(rec.Body.Bytes(), &userResponse)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	assert.Equal(t, "test10@example.com", userResponse.Email)
	assert.Equal(t, "1234567890", userResponse.PhoneNumber)
	assert.Equal(t, "Jalan 123", userResponse.Address)
	assert.Equal(t, "user", userResponse.Role)

	// Check that the user ID is not empty
	assert.NotEmpty(t, userResponse.ID)
}
