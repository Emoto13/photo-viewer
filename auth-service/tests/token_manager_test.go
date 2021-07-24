package tests

import (
	"testing"

	"github.com/Emoto13/photo-viewer-rest/auth-service/src/token"
	"github.com/Emoto13/photo-viewer-rest/auth-service/tests/setup"
)

func assertEquals(a, b string, message string) {
	if a != b {
		panic(message)
	}
}

func TestGenerateToken(t *testing.T) {
	tokenManager := setup.NewMockTokenManager()

	var tests = []struct {
		name     string
		input    string
		expected string
		message  string
	}{
		{"Test GenerateToken", "username", "token", "Token should be generated successfully"},
	}

	for _, test := range tests {
		res, _ := tokenManager.GenerateToken(test.input)
		assertEquals(res, test.expected, test.message)
	}
}

func TestSaveGetRemoveTokenManager(t *testing.T) {
	tokenManager := setup.NewMockTokenManager()

	var tests = []struct {
		name     string
		input    token.Token
		expected string
		message  string
	}{
		{"Test SaveToken", token.NewToken("token", "username"), "username", "Token manager should save and retrieve token successfully"},
	}

	for _, test := range tests {
		tokenManager.SaveToken(test.input)
		res, _ := tokenManager.GetUsernameFromToken(test.input.GetValue())
		assertEquals(res, test.expected, test.message)
		tokenManager.RemoveToken(test.input.GetValue())
	}
}
