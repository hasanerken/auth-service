package app

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"time"
	encryption "weight-tracker/pkg/utils"
)

func (s *Server) Login() fiber.Handler {
	return func(c *fiber.Ctx) error {
		type Credentials struct {
			Identity string `json:"identity"`
			Password string `json:"password"`
		}

		var credentials Credentials
		if err := c.BodyParser(&credentials); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "error on login request", "data": err})
		}

		user := s.userService.GetUserByUsername(credentials.Identity)
		if user.Username == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{"message": "Hatalı şifre veya kullanıcı adı"})
		}
		fmt.Println("ss", user.Password)
		auth, err := encryption.ComparePwd(user.Password, credentials.Password)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{"message": "Hatalı şifre veya kullanıcı adı"})
		}
		if !auth {
			return c.Status(fiber.StatusUnauthorized).JSON(
				fiber.Map{"message": "Hatalı şifre veya kullanıcı adı"})
		}

		// create auth token cookie
		claims := jwt.MapClaims{
			"user_id": user.ID,
		}
		accessToken := s.tokenService.CreateAccessToken(claims)
		cookie := fiber.Cookie{
			Name:     "accessToken",
			Value:    accessToken.Token,
			Expires:  time.Now().Add(time.Minute * 15),
			HTTPOnly: true,
		}

		c.Cookie(&cookie)

		return c.JSON(fiber.Map{
			"status": "successful login", "data": accessToken,
		})
	}
}

func (s *Server) Logout() fiber.Handler {

	return func(c *fiber.Ctx) error {
		// extract token for other steps
		tokenString, err := s.tokenService.ExtractToken(c.Request())
		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"message": "token can't be extracted",
				"data":    err,
			})
		}
		fmt.Println("token string", tokenString)

		// delete refresh token
		s.tokenService.InvalidateToken(tokenString)

		// todo: delete auth cookies
		cookie := fiber.Cookie{
			Name:     "accessToken",
			Value:    "",
			Path:     "",
			Expires:  time.Now().Add(-time.Hour),
			HTTPOnly: false,
		}
		c.Cookie(&cookie)

		// redirect to main page
		return c.Status(fiber.StatusOK).Redirect("/")
	}
}
