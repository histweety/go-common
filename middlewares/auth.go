package auth

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/histweety/go-common/types"
)

type Config struct {
	Secret string
}

func New(config Config) fiber.Handler {
	var accessTokenSecret = []byte(config.Secret)

	return func(c *fiber.Ctx) error {
		token := c.Get("Authorization")
		if token != "" {
			claims := &types.Claims{}
			token = strings.Replace(token, "Bearer ", "", 1)
			jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return accessTokenSecret, nil
			})

			c.Locals("UserID", claims.UserID)
		}

		return c.Next()
	}
}
