package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/histweety/go-common/types"
	"github.com/histweety/go-common/utils"
)

type ConfigAuth struct {
	Secret    string
	WhiteList []string
}

func NewAuth(cfg ConfigAuth) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		isWhiteList := isContain(c.Path(), cfg.WhiteList)

		if tokenString == "" && isWhiteList {
			return c.Next()
		}

		claims := &types.Claims{}
		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		claims, err := utils.ParseToken(tokenString, false)
		if err != nil {
			if isWhiteList {
				return c.Next()
			}

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"statusCode": 401,
				"message":    err.Error(),
				"data":       nil,
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
