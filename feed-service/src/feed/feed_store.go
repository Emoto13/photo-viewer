package feed

import (
	"fmt"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/feed-service/src/follow/models"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post"
	postModels "github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
	"github.com/gocql/gocql"
)

type FeedStore interface {
	GetFeed(username string) ([]*postModels.Post, error)
	UpdateFeed(username string, followings []*models.Following) error
	AddToFeed(username string, post *postModels.Post) error
}

type feedStore struct {
	mu         sync.RWMutex
	session    *gocql.Session
	postClient post.PostClient
}

func NewFeedStore(postClient post.PostClient, session *gocql.Session) FeedStore {
	return &feedStore{
		session: session,
		mu:      sync.RWMutex{},
	}
}

func (s *feedStore) GetFeed(username string) ([]*postModels.Post, error) {
	res := []*postModels.Post{}
	err := s.session.Query(`SELECT feed FROM feed WHERE username=?`, username).Scan(&res)
	if err != nil && err != gocql.ErrNotFound {
		fmt.Println("error retrieving feed posts: ", err.Error())
		return nil, err
	}

	return res, nil
}

func (s *feedStore) UpdateFeed(username string, followings []*models.Following) error {
	feed := []*postModels.Post{}
	for _, following := range followings {
		temp := []*postModels.Post{}
		err := s.session.Query(`SELECT posts FROM posts WHERE username=?`, following.Username).Scan(&temp)
		if err != nil {
			fmt.Println("could not execute select query", err.Error())
			return err
		}

		feed = append(feed, temp...)
	}

	sortPosts(feed)
	err := s.session.Query(`INSERT INTO feed(username, feed) VALUES (?, ?);`, username, feed).Exec()
	if err != nil {
		fmt.Println("could not execute insert query", err.Error())
		return err
	}

	return nil
}

func (s *feedStore) AddToFeed(username string, post *postModels.Post) error {
	feed := []*postModels.Post{}
	err := s.session.Query(`SELECT feed FROM feed WHERE username=?`, username).Scan(&feed)
	if err != nil && err != gocql.ErrNotFound {
		return err
	}

	feed = append([]*postModels.Post{post}, feed...)
	sortPosts(feed)
	err = s.session.Query(`INSERT INTO feed(username, feed) VALUES (?, ?);`, username, feed).Exec()
	if err != nil {
		return err
	}

	return nil
}
