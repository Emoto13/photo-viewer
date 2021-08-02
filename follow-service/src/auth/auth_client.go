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
	req, err := http.NewRequest("GET", c.address+"/auth-service/authenticate", nil)
	if err != nil {
		fmt.Println("something went wrong with creation authenticate request: ", err.Error())
		return "", err
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := c.client.Do(req)
	if err != nil {
		fmt.Println("something went wrong with sending authenticate request: ", err.Error())
		return "", err
	}

	var responseMap map[string]string
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		fmt.Println("failed to parse json response body: ", err.Error())
		return "", err
	}

	return responseMap["username"], nil
}
