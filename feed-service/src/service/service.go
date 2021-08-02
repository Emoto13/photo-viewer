package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/feed-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/feed"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post/cache_store"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
)

type feedService struct {
	authClient     auth.AuthClient
	followClient   follow.FollowClient
	postClient     post.PostClient
	postCacheStore cache_store.PostCacheStore
	feedStore      feed.FeedStore
}

func NewFeedService(authClient auth.AuthClient,
	followClient follow.FollowClient,
	postClient post.PostClient,
	postCacheStore cache_store.PostCacheStore,
	feedStore feed.FeedStore) *feedService {
	return &feedService{
		authClient:     authClient,
		followClient:   followClient,
		postClient:     postClient,
		postCacheStore: postCacheStore,
		feedStore:      feedStore,
	}
}

func (s *feedService) GetFollowingPosts(w http.ResponseWriter, r *http.Request) {
	username, err := s.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("Couldn't authenticate", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cachedResult, err := s.postCacheStore.Get(context.Background(), username)
	if err == nil {
		response, _ := json.Marshal(map[string][]*models.Post{"feed": cachedResult})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

	feed, err := s.feedStore.GetFeed(username)
	if err != nil {
		fmt.Println("Couldn't get feed", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.postCacheStore.Set(context.Background(), username, feed)
	response, _ := json.Marshal(map[string][]*models.Post{"feed": feed})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (s *feedService) UpdateFeed(w http.ResponseWriter, r *http.Request) {
	username, err := s.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("couldn't authenticate", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	followings, err := s.followClient.GetFollowing(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("couldnt get followings", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(followings)
	err = s.feedStore.UpdateFeed(username, followings)
	if err != nil {
		fmt.Println("couldn't get feed", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
}

func (s *feedService) AddToFeed(w http.ResponseWriter, r *http.Request) {
	username, err := s.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("couldn't authenticate", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := getPostRequestBody(r.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = s.feedStore.AddToFeed(username, body["post"])
	if err != nil {
		fmt.Println("error adding to feed", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Successfully added to feed")
	respondWithJSON(w, http.StatusOK, map[string]string{"Message": "Post added to feed succssfully"})
}

func (s *feedService) AddToFollowersFeed(w http.ResponseWriter, r *http.Request) {
	body, err := getPostRequestBody(r.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	followers, err := s.followClient.GetFollowers(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	for _, follower := range followers {
		err = s.feedStore.AddToFeed(follower.Username, body["post"])
		if err != nil {
			fmt.Println("error adding to feed", err.Error())
			respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
	}

	fmt.Println("Post added to followers feed successfully")
	respondWithJSON(w, http.StatusOK, map[string]string{"Message": "Post added to followers feed succssfully"})
}

func (s *feedService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
}
