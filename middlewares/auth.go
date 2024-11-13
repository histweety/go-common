package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

type claims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

type iAuthService interface {
	ValidateToken(token string) (*claims, error)
}

func AuthRequired(authService iAuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Missing or malformed JWT",
				"data":       nil,
			})
		}

		claims, err := authService.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Invalid or expired JWT",
				"data":       nil,
			})
		}

		// Store user ID in locals for use in handlers
		c.Locals("userID", claims.UserID)
		return c.Next()
	}
}
