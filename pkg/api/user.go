package api

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type UserService interface {
	New(user NewUserRequest) error
	GetUserByUsername(username string) User
}

type UserRepository interface {
	CreateUser(NewUserRequest) error
	FindUserByUsername(username string) (User, error)
}

type userService struct {
	storage UserRepository
}

func NewUserService(userRepo UserRepository) UserService {
	return &userService{storage: userRepo}
}

func (u *userService) New(user NewUserRequest) error {
	// basic validations
	if user.Email == "" {
		return errors.New("user service:: - email required")
	}
	if user.Name == "" {
		return errors.New("user service - name required")
	}
	if user.WeightGoal == "" {
		return errors.New("user service - weight goal required")
	}
	// basic normalisation
	user.Name = strings.ToLower(user.Name)
	user.Email = strings.TrimSpace(user.Email)

	if err := u.storage.CreateUser(user); err != nil {
		fmt.Println("in service", err)
		return err
	}
	return nil
}

func (u *userService) GetUserByUsername(username string) User {
	user, err := u.storage.FindUserByUsername(username)
	if err != nil {
		log.Warning("Can not find the user", err)
	}
	return user
}
