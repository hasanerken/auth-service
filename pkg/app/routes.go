package app

import (
	"github.com/gofiber/fiber/v2"
	"weight-tracker/pkg/middleware"
)

func (s *Server) Routes() *fiber.App {
	router := s.router

	api := router.Group("/api")
	api.Get("/status", mw.TokenAuthMiddleware(s.tokenService), s.ApiStatus())
	api.Post("/users", s.CreateUser())

	auth := router.Group("/auth")
	auth.Post("/login", s.Login())
	auth.Post("/logout", s.Logout())

	return router
}
