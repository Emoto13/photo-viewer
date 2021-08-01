package post

import (
	"encoding/json"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
)

type PostClient interface {
	GetUserPosts(username string) ([]*models.Post, error)
}

type postClient struct {
	client  *http.Client
	address string
}

func NewPostClient(client *http.Client, address string) PostClient {
	return &postClient{client: client, address: address}
}

func (c *postClient) GetUserPosts(username string) ([]*models.Post, error) {
	req, err := http.NewRequest("GET", c.address+"/post-service/posts/"+username, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var responseMap map[string][]*models.Post
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		return nil, err
	}

	return responseMap["posts"], nil
}
