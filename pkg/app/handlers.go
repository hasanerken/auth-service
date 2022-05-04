package app

import (
	"github.com/gofiber/fiber/v2"
	"weight-tracker/pkg/api"
)

func (s *Server) ApiStatus() fiber.Handler {
	return func(c *fiber.Ctx) error {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   "weight tracker API running smoothly",
		})
	}
}

func (s *Server) CreateUser() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var newUser api.NewUserRequest

		err := c.BodyParser(&newUser)
		if err != nil {
			return err
		}
		err = s.userService.New(newUser)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status": "can not write on db", "data": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "success",
			"data":   newUser,
		})
	}
}
