package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"

	"jakarta-luxury-rent-car/database"
	"jakarta-luxury-rent-car/models"
)

// Struct untuk mengirimkan response
type UserResponse struct {
	ID          uint   `json:"id"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
	Token       string `json:"token,omitempty"`
}

// Struct untuk input registrasi
type RegisterRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=3"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Address     string `json:"address" validate:"required"`
	Role        string `json:"role"` // Optional
}

// Struct untuk input login
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// @Summary Register new user
// @Description Register new user
// @Tags Register and Login
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "User registration request body"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 409 {object} map[string]interface{} "User already exists with this email"
// @Failure 500 {object} map[string]interface{} "Failed to hash password or register user"
// @Router /register [post]
func RegisterUser(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Check if the user already exists by email
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		// User already exists
		return c.JSON(http.StatusConflict, echo.Map{"error": "User already exists with this email"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to hash password"})
	}

	// Create user model
	user := models.User{
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		Address:     req.Address,
		Role:        req.Role,
	}

	// Save user to database
	if err := database.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, UserResponse{
		ID:          user.UserID,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
	})
}

// @Summary User login
// @Description Login user and return JWT token
// @Tags Register and Login
// @Accept json
// @Produce json
// @Param request body LoginRequest true "User login request body"
// @Success 200 {object} UserResponse
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Invalid email or password"
// @Failure 500 {object} map[string]interface{} "Failed to generate token"
// @Router /login [post]
func LoginUser(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid input"})
	}

	if err := c.Validate(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{"error": "Invalid email or password"})
	}

	// Generate JWT token
	token, err := generateJWT(user)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, UserResponse{
		ID:          user.UserID,
		Email:       user.Email,
		Address:     user.Address,
		PhoneNumber: user.PhoneNumber,
		Role:        user.Role,
		Token:       token,
	})
}

// Fungsi untuk menghasilkan JWT
func generateJWT(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id": user.UserID,
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	return token.SignedString([]byte(secret))
}
