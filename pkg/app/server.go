package app

import (
	"github.com/gofiber/fiber/v2"
	"log"
	"weight-tracker/pkg/api"
)

type Server struct {
	router        *fiber.App
	userService   api.UserService
	weightService api.WeightService
	tokenService  api.TokenService
}

func NewServer(router *fiber.App, userService api.UserService, weightService api.WeightService, tokenService api.TokenService,
) *Server {
	return &Server{
		router:        router,
		userService:   userService,
		weightService: weightService,
		tokenService:  tokenService,
	}
}

func (s *Server) Run() error {
	r := s.Routes()

	err := r.Listen(":3003")
	if err != nil {
		log.Printf("Server - there was an error calling Run on router: %v/n", err)
		return err
	}
	return nil
}
