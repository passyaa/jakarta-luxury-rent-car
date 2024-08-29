package middlewares

import (
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
	jwtMiddleware "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware() echo.MiddlewareFunc {
	signingKey := os.Getenv("JWT_SECRET")
	if signingKey == "" {
		log.Fatal("JWT_SECRET environment variable is not set or is empty")
	}

	return jwtMiddleware.WithConfig(jwtMiddleware.Config{
		SigningKey: []byte(signingKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwt.MapClaims)
		},
	})
}
