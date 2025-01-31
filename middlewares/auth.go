package middlewares

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/histweety/go-common/types"
)

type ConfigAuth struct {
	Secret    string
	WhiteList []string
}

func NewAuth(cfg ConfigAuth) fiber.Handler {
	var accessTokenSecret = []byte(cfg.Secret)

	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		isWhiteList := isContain(c.Path(), cfg.WhiteList)

		if tokenString == "" && isWhiteList {
			return c.Next()
		}

		claims := &types.Claims{}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return accessTokenSecret, nil
		})

		if err != nil && isWhiteList {
			return c.Next()
		}
		if !token.Valid && isWhiteList {
			return c.Next()
		}

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Unauthorized",
				"data":       "Something went wrong",
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    "Unauthorized",
				"data":       "Token is invalid",
			})
		}

		c.Locals("UserID", claims.UserID)

		return c.Next()
	}
}

func isContain(path string, whiteList []string) bool {
	for _, p := range whiteList {
		if p == path {
			return true
		}
	}

	return false
}
