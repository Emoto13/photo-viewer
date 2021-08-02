package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/follow-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/follow-service/src/feed"
	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow/models"
)

type followService struct {
	authClient  auth.AuthClient
	followStore follow.FollowStore
	feedClient  feed.FeedClient
	mu          sync.RWMutex
}

func NewFollowService(authClient auth.AuthClient, feedClient feed.FeedClient, followStore follow.FollowStore) *followService {
	return &followService{
		authClient:  authClient,
		feedClient:  feedClient,
		followStore: followStore,
		mu:          sync.RWMutex{},
	}
}

func (fs *followService) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := getRequestBody(r.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = fs.followStore.CreateUser(body["username"])
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"Message": fmt.Sprintf("User %s created successfully", body["username"])})
	return
}

func (fs *followService) Follow(w http.ResponseWriter, r *http.Request) {
	username, err := fs.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := getRequestBody(r.Body)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = fs.followStore.SaveFollow(models.NewFollow(username, body["follow"]))
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = fs.feedClient.UpdateFeed(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("couldnt update feed", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(username, "followed", body["follow"])
	respondWithJSON(w, http.StatusOK, map[string]string{"Message": fmt.Sprintf("You are now following %s", body["follow"])})
	return
}

func (fs *followService) Unfollow(w http.ResponseWriter, r *http.Request) {
	username, err := fs.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("failed to authenticate: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	body, err := getRequestBody(r.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = fs.followStore.RemoveFollow(models.NewFollow(username, body["unfollow"]))
	if err != nil {
		fmt.Println("failed to remove follow: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = fs.feedClient.UpdateFeed(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("couldnt update feed: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println(username, "unfollowed", body["unfollow"])
	respondWithJSON(w, http.StatusOK, map[string]string{"Message": fmt.Sprintf("%s is unfollowed", body["unfollow"])})
	return
}

func (fs *followService) GetFollowers(w http.ResponseWriter, r *http.Request) {
	username, err := fs.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("failed to authenticate: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	followers, err := fs.followStore.GetFollowers(username)
	if err != nil {
		fmt.Println("couldnt get followers: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, _ := json.Marshal(map[string][]*models.Follower{"followers": followers})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (fs *followService) GetFollowing(w http.ResponseWriter, r *http.Request) {
	username, err := fs.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("failed to authenticate: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	following, err := fs.followStore.GetFollowing(username)
	if err != nil {
		fmt.Println("couldnt get following: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, _ := json.Marshal(map[string][]*models.Following{"following": following})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (fs *followService) GetSuggestions(w http.ResponseWriter, r *http.Request) {
	username, err := fs.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("failed to authenticate: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	suggestions, err := fs.followStore.GetSuggestions(username)
	if err != nil {
		fmt.Println("couldnt get suggestions: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, _ := json.Marshal(map[string][]*models.Suggestion{"suggestions": suggestions})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (fs *followService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
