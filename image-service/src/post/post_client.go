package post

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/post/models"
)

type PostClient interface {
	CreatePost(authHeader string, post *models.Post) error
}

type postClient struct {
	client  *http.Client
	address string
}

func NewPostClient(client *http.Client, address string) PostClient {
	return &postClient{client: client, address: address}
}

func (c *postClient) CreatePost(authHeader string, post *models.Post) error {
	values := map[string]*models.Post{"post": post}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("POST", c.address+"/post-service/create-post", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("failed to create create-post request:", err.Error())
		return err
	}
	req.Header.Set("Authorization", authHeader)

	_, err = c.client.Do(req)
	if err != nil {
		fmt.Println("error create-post response:", err.Error())
		return err
	}

	return nil
}
