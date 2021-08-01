package feed

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/post/models"
)

type FeedClient interface {
	AddToFollowersFeed(authHeader string, post *models.Post) error
}

type feedClient struct {
	client  *http.Client
	address string
}

func NewFeedClient(client *http.Client, address string) FeedClient {
	return &feedClient{client: client, address: address}
}

func (c *feedClient) AddToFollowersFeed(authHeader string, post *models.Post) error {
	values := map[string]*models.Post{"post": post}
	jsonValue, _ := json.Marshal(values)

	req, err := http.NewRequest("PATCH", c.address+"/feed-service/add-to-followers-feed", bytes.NewBuffer(jsonValue))
	if err != nil {
		fmt.Println("could not create add to feed request: ", err.Error())
		return err
	}
	req.Header.Set("Authorization", authHeader)

	_, err = c.client.Do(req)
	if err != nil {
		fmt.Println("error sending update feed request: ", err.Error())
		return err
	}

	return nil
}
