package repository

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"time"
	"weight-tracker/pkg/api"
)

type RedisStore interface {
	Ping() error
	SaveToken(userID uuid.UUID, refreshToken []byte) error
	FindRefreshToken(tokenID string) (uuid.UUID, error)
	DeleteRefreshToken(tokenID string) bool
}

type redisStore struct {
	client *redis.Client
}

func NewRedisStore(client *redis.Client) RedisStore {
	return &redisStore{client: client}
}

func (rs *redisStore) Ping() error {
	panic("implement me")
	return nil
}

func (rs *redisStore) SaveToken(tokenID uuid.UUID, refreshToken []byte) error {
	m := new(api.RefreshToken)
	if err := json.Unmarshal(refreshToken, &m); err != nil {
		fmt.Println(err)
	}
	now := time.Now()
	exp := time.Unix(m.Expired, 0)

	if err := rs.client.Set(tokenID.String(), refreshToken, exp.Sub(now)).Err(); err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (rs *redisStore) FindRefreshToken(tokenID string) (uuid.UUID, error) {
	fmt.Println("find refresh called")
	m := new(api.RefreshToken)
	value, err := rs.client.Get(tokenID).Result()
	if err != nil {
		fmt.Println("redis get errors:", err)
		return uuid.Nil, err
	}
	if err := json.Unmarshal([]byte(value), &m); err != nil {
		fmt.Println("json unmarshalling errors:", err)
		return uuid.Nil, err
	}

	return m.UserID, nil
}

func (rs *redisStore) DeleteRefreshToken(tokenID string) bool {
	
	deleted, err := rs.client.Del(tokenID).Result()
	if err != nil {
		fmt.Println("could not delete the token")
		return false
	}
	fmt.Println("value :::", deleted)
	return true
}
