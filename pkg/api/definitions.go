package api

import (
	"github.com/google/uuid"
	"time"
)

type NewUserRequest struct {
	Name          string `json:"name"`
	Age           int    `json:"age"`
	Height        int    `json:"height"`
	Sex           string `json:"sex"`
	ActivityLevel int    `json:"activity_level"`
	WeightGoal    string `json:"weight_goal"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	Username      string `json:"username"`
}

type User struct {
	ID            uuid.UUID `json:"id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Name          string    `json:"name"`
	Age           int       `json:"age"`
	Height        int       `json:"height"`
	Sex           string    `json:"sex"`
	ActivityLevel int       `json:"activity_level"`
	WeightGoal    string    `json:"weight_goal"`
	Email         string    `json:"email"`
	Username      string    `json:"username"`
	Password      []byte    `json:"-"`
}

type UsersList []*User

type Weight struct {
	Weight             int       `json:"weight"`
	UserID             uuid.UUID `json:"user_id"`
	BMR                int       `json:"bmr"`
	DailyCaloricIntake int       `json:"daily_caloric_intake"`
}

type NewWeightRequest struct {
	Weight int       `json:"weight"`
	UserID uuid.UUID `json:"user_id"`
}

type TokenResponse struct {
	Type         string `json:"type"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type AccessToken struct {
	Type         string `json:"type"`
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshToken struct {
	RefreshToken string    `json:"refresh_token"`
	UserID       uuid.UUID `json:"user_id"`
	Expired      int64     `json:"expired"`
}
