package auth

import (
	"encoding/json"
	"net/http"
)

type AuthClient interface {
	Authenticate(authHeader string) (string, error)
}

type authClient struct {
	client  *http.Client
	address string
}

func NewAuthClient(client *http.Client, address string) AuthClient {
	return &authClient{client: client, address: address}
}

func (c *authClient) Authenticate(authHeader string) (string, error) {
	req, err := http.NewRequest("GET", c.address+"/auth-service/authenticate", nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := c.client.Do(req)
	if err != nil {
		return "", err
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return "", err
	}

	return responseMap["username"], nil
}
