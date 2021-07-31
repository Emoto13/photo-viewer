package follow

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type FollowClient interface {
	CreateUser(username string) error
}

type followClient struct {
	client  *http.Client
	address string
}

func NewFollowClient(client *http.Client, address string) FollowClient {
	return &followClient{client: client, address: address}
}

func (c *followClient) CreateUser(username string) error {
	values := map[string]string{"username": username}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", c.address+"/follow-service/create-user", bytes.NewBuffer(jsonValue))
	if err != nil {
		return err
	}

	_, err = c.client.Do(req)
	if err != nil {
		return err
	}

	fmt.Println("User created successfully as follower/ following")
	return nil
}
