package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthClient interface {
	Authenticate(token string) (string, error)
}

type authClient struct {
	client  *http.Client
	address string
}

func NewAuthClient(client *http.Client, address string) AuthClient {
	return &authClient{client: client, address: address}
}

func (c *authClient) Authenticate(authHeader string) (string, error) {
	req, _ := http.NewRequest("GET", c.address+"/auth-service/authenticate", nil)
	req.Header.Set("Authorization", authHeader)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return responseMap["username"], nil
}
