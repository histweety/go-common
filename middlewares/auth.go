package middlewares

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/histweety/go-common/types"
)

type ConfigAuth struct {
	Secret string
}

func NewAuth(config ConfigAuth) fiber.Handler {
	var accessTokenSecret = []byte(config.Secret)

	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString != "" {
			claims := &types.Claims{}
			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return accessTokenSecret, nil
			})

			if err == nil && token.Valid {
				c.Locals("UserID", claims.UserID)
			}
		}

		return c.Next()
	}
}
