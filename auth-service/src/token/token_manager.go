package token

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	"golang.org/x/crypto/bcrypt"
)

type TokenManager interface {
	GenerateToken(username string) (string, error)
	GetUsernameFromToken(token string) (string, error)
	SaveToken(token Token) error
	RemoveToken(token string) error
}

type tokenManager struct {
	redisClient *redis.Client
	mutex       sync.RWMutex
}

func NewTokenManager(redisClient *redis.Client) TokenManager {
	return &tokenManager{
		redisClient: redisClient,
		mutex:       sync.RWMutex{},
	}
}

func (m *tokenManager) GenerateToken(username string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println("Token: ", hash)
	return string(hash), nil
}

func (m *tokenManager) GetUsernameFromToken(token string) (string, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	username, err := m.redisClient.Get(context.Background(), token).Result()
	if err != nil {
		return "", err
	}

	return username, nil
}

func (m *tokenManager) SaveToken(token Token) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := m.redisClient.Set(context.Background(), token.GetValue(), token.GetOwner(), 0).Err()
	if err != nil {
		return err
	}

	return nil
}

func (m *tokenManager) RemoveToken(token string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	err := m.redisClient.Del(context.Background(), token).Err()
	if err != nil {
		return err
	}

	return nil
}
