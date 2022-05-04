package repository

import (
	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"log"
	"weight-tracker/pkg/api"
	encryption "weight-tracker/pkg/utils"
)

type Storage interface {
	CreateUser(request api.NewUserRequest) error
	CreateWeightEntry(request api.Weight) error
	FindUser(userID uuid.UUID) (api.User, error)
	FindUserByUsername(username string) (api.User, error)
}

type storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) Storage {
	return &storage{db: db}
}

func (s *storage) CreateUser(request api.NewUserRequest) error {
	newUserStatement := `
			INSERT INTO users (id, name, age, height, sex, activity_level, weight_goal, email, password, username) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id;
			`
	// now := time.Now().Local()
	var ID uuid.UUID
	hashedPwd := encryption.GenerateHashedPwd(request.Password)
	err := s.db.QueryRow(newUserStatement, uuid.New(), request.Name, request.Age, request.Height, request.Sex, request.ActivityLevel, request.WeightGoal, request.Email, hashedPwd, request.Username).Scan(&ID)
	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return err
	}

	return nil
}

func (s *storage) CreateWeightEntry(request api.Weight) error {
	newWeightStatement := `
			INSERT INTO "weight" (weight, user_id, bmr, daily_caloric_intake)
			VALUES ($1, $2, $3, $4)
			RETURNING id;
			`
	var ID int
	err := s.db.QueryRow(newWeightStatement, request.Weight, request.UserID, request.BMR, request.DailyCaloricIntake).Scan(&ID)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return err
	}
	return nil
}

func (s *storage) FindUser(userID uuid.UUID) (api.User, error) {
	getUserStatement := `
					SELECT id, name, age, height, sex, activity_level, email, weight_goal FROM users
					WHERE id=$1;
					`

	var user api.User
	err := s.db.QueryRow(getUserStatement, userID).Scan(&user.ID, &user.Name, &user.Age, &user.Height, &user.Sex, &user.ActivityLevel, &user.Email, &user.WeightGoal)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return api.User{}, err
	}
	return user, nil
}
func (s *storage) FindUserByUsername(username string) (api.User, error) {
	getUserStatement := `
					SELECT id, name, username, password FROM users
					WHERE username=$1;
					`

	var user api.User
	err := s.db.QueryRow(getUserStatement, username).Scan(&user.ID, &user.Name, &user.Username, &user.Password)

	if err != nil {
		log.Printf("this was the error: %v", err.Error())
		return api.User{}, err
	}
	return user, nil
}
