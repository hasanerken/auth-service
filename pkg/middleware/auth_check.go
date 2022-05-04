package mw

import (
	"github.com/gofiber/fiber/v2"
	"weight-tracker/pkg/api"
)

func TokenAuthMiddleware(ts api.TokenService) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// extract the token from the request header
		tokenString, err := ts.ExtractToken(c.Request())
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "token can't be extracted",
				"data":    err,
			})
		}

		if ok, accessToken, err := ts.ValidateToken(tokenString); !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "unauthorized user, request is blocked",
				"data":    err.Error(),
			})
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "users access token refreshed",
				"data":    accessToken,
			})

		}

		return c.Next()
	}
}
