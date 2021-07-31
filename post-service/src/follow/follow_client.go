package follow

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/follow/models"
)

type FollowClient interface {
	GetFollowing(authHeader string) ([]models.Following, error)
}

type followClient struct {
	client  *http.Client
	address string
}

func NewFollowClient(client *http.Client, address string) FollowClient {
	return &followClient{client: client, address: address}
}

func (c *followClient) GetFollowing(authHeader string) ([]models.Following, error) {
	req, err := http.NewRequest("GET", c.address+"/follow-service/get-following", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", authHeader)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Println("here")
	var responseMap map[string][]models.Following
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return nil, err
	}

	return responseMap["following"], nil
}
