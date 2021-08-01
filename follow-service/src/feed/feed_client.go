package feed

import (
	"fmt"
	"net/http"
)

type FeedClient interface {
	UpdateFeed(authHeader string) error
}

type feedClient struct {
	client  *http.Client
	address string
}

func NewFeedClient(client *http.Client, address string) FeedClient {
	return &feedClient{client: client, address: address}
}

func (c *feedClient) UpdateFeed(authHeader string) error {
	req, err := http.NewRequest("PATCH", c.address+"/feed-service/update-feed", nil)
	if err != nil {
		fmt.Println("could not create update feed request: ", err.Error())
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
