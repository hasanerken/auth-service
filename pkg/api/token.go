package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
	"weight-tracker/config"
)

type TokenService interface {
	CreateAccessToken(claims jwt.MapClaims) AccessToken
	ValidateToken(tokenString string) (bool, *AccessToken, error)
	InvalidateToken(tokenString string) (bool, error)
	ExtractToken(request *fiber.Request) (string, error)
}

type TokenRepository interface {
	SaveToken(userID uuid.UUID, refreshToken []byte) error
	FindRefreshToken(tokenID string) (uuid.UUID, error)
	DeleteRefreshToken(tokenID string) bool
}

type tokenService struct {
	redisStore TokenRepository
}

func NewTokenService(tokenStore TokenRepository) TokenService {
	return &tokenService{redisStore: tokenStore}
}

func (ts *tokenService) CreateAccessToken(claims jwt.MapClaims) AccessToken {
	now := time.Now()
	accessDuration, err := time.ParseDuration(config.Get("ACCESS_EXPIRED"))
	tokenExpired := now.Add(accessDuration).Unix()

	if claims["user_id"] == nil {
		return AccessToken{}
	}
	userID := claims["user_id"].(uuid.UUID)

	token := jwt.New(jwt.SigningMethodHS256)
	var _, checkExp = claims["exp"]
	var _, checkIat = claims["iat"]

	if !checkExp {
		claims["exp"] = tokenExpired
	}
	if !checkIat {
		claims["iat"] = now.Unix()
	}
	claims["token_type"] = "access_token"
	claims["token_id"], _ = uuid.NewRandom()
	token.Claims = claims

	authToken := new(AccessToken)
	tokenString, err := token.SignedString([]byte(config.Get("SECRET_KEY")))

	if err != nil {
		fmt.Println(err)
		return AccessToken{}
	}
	authToken.Token = tokenString
	authToken.Type = "Bearer"

	// create refresh token
	refreshDuration, err := time.ParseDuration(config.Get("REFRESH_EXPIRED"))
	if err != nil {
		log.Warning(err)
	}
	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenExpired := now.Add(refreshDuration).Unix()

	claims["exp"] = refreshTokenExpired
	claims["token_type"] = "refresh_token"
	refreshToken.Claims = claims

	refreshTokenString, err := refreshToken.SignedString([]byte(config.Get("SECRET_KEY")))
	if err != nil {
		return AccessToken{}
	}

	authToken.RefreshToken = refreshTokenString
	fmt.Println("userid", userID)
	tokenID := claims["token_id"].(uuid.UUID)
	rt := RefreshToken{
		RefreshToken: refreshTokenString,
		Expired:      refreshTokenExpired,
		UserID:       userID,
	}
	rtJson, err := json.Marshal(rt)
	if err != nil {
		log.Warning(err)
	}

	go ts.saveTokens(tokenID, rtJson)

	return AccessToken{
		Type:         "Bearer",
		Token:        authToken.Token,
		RefreshToken: authToken.RefreshToken,
	}
}

func (ts *tokenService) saveTokens(tokenID uuid.UUID, rtJson []byte) {
	if err := ts.redisStore.SaveToken(tokenID, rtJson); err != nil {
		log.Warning("Redis can not save")
	}
}

func (ts *tokenService) ValidateToken(tokenString string) (bool, *AccessToken, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Get("SECRET_KEY")), nil
	})

	if token.Valid {
		fmt.Println("You look nice today")
		return true, nil, nil
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			fmt.Println("That's not even a token")
			return false, nil, errors.New("token is fake")
		} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
			// Token is either expired or not active yet
			fmt.Println("Timing is everything")
			// We need to handle whether refresh token is still in use
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				tokenID := claims["token_id"].(string)
				aToken, err := ts.tokenRefreshRequest(tokenID)
				if err != nil {
					return false, nil, errors.New("refresh token has expired")
				}
				return true, aToken, nil
			}

			return false, nil, errors.New("token has expired")
		} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
			fmt.Println("You need to wait to use this token")
			return false, nil, errors.New("token not valid yet")
		} else {
			fmt.Println("Couldn't handle this token:", err)
			return false, nil, errors.New("problem with token")
		}
	} else {
		fmt.Println("Couldn't handle this token:", err)
		return false, nil, errors.New("problem with token")
	}
}

func (ts *tokenService) tokenRefreshRequest(tokenID string) (*AccessToken, error) {
	fmt.Println("token refresh request called")
	userID, err := ts.redisStore.FindRefreshToken(tokenID)
	if err != nil {
		// ts.CreateAccessToken()
		fmt.Println("mmm- token refrehsed error occurred")
		return nil, err
	} else {
		fmt.Println("user id", userID)
		claims := jwt.MapClaims{
			"user_id": userID,
		}
		accessToken := ts.CreateAccessToken(claims)

		return &accessToken, nil
	}

}

func (ts *tokenService) InvalidateToken(tokenString string) (bool, error) {
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Get("SECRET_KEY")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		tokenID := claims["token_id"].(string)
		if ok := ts.redisStore.DeleteRefreshToken(tokenID); ok {
			fmt.Println("refresh p token deleted successfully")
		}

		return true, nil
	}

	return true, nil
}

func (ts *tokenService) ExtractToken(request *fiber.Request) (string, error) {
	AuthorizationHeader := string(request.Header.Peek("Authorization"))
	tokenPart := strings.Split(AuthorizationHeader, "Bearer ")
	if len(tokenPart) != 2 {
		// Error: Bearer token not in proper format
		fmt.Println("Not properly formed token")
		return "", fiber.NewError(400, "Not properly formed token")

	}
	tokenString := tokenPart[1]
	return tokenString, nil
}
