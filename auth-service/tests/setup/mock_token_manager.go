package setup

import (
	"github.com/Emoto13/photo-viewer-rest/auth-service/src/token"
)

type mockTokenManager struct {
}

func NewMockTokenManager() token.TokenManager {
	return &mockTokenManager{}
}

func (m *mockTokenManager) GenerateToken(username string) (string, error) {
	return "token", nil
}

func (m *mockTokenManager) GetUsernameFromToken(token string) (string, error) {
	return "username", nil
}

func (m *mockTokenManager) SaveToken(token token.Token) error {
	return nil
}

func (m *mockTokenManager) RemoveToken(token string) error {
	return nil
}
