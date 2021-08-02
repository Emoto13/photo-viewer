package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/cache_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/models"
	"github.com/gorilla/mux"
)

type postService struct {
	authClient   auth.AuthClient
	postStore    post_store.PostStore
	postCache    cache_store.PostCacheStore
	followClient follow.FollowClient
}

func New(authClient auth.AuthClient, postStore post_store.PostStore, followClient follow.FollowClient) *postService {
	return &postService{
		authClient:   authClient,
		followClient: followClient,
		postStore:    postStore,
	}
}

func (s *postService) CreatePost(w http.ResponseWriter, r *http.Request) {
	username, err := s.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("couldn't authenticate", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	post, err := getRequestBody(r.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("Post recieved in CreatePost:", post)

	err = s.postStore.CreatePost(username, post)
	if err != nil {
		fmt.Println("couldnt write post to cassandra:", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Post created successfully")
	respondWithJSON(w, http.StatusOK, map[string]string{"Message": "Post created successfully"})
}

func (s *postService) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["username"]
	posts, err := s.postStore.GetUserPosts(name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, _ := json.Marshal(map[string][]*models.Post{"posts": posts})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return

}

func (s *postService) SearchPosts(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := params["query"]
	posts, err := s.postStore.SearchPosts(name)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, _ := json.Marshal(map[string][]*models.Post{"posts": posts})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (s *postService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
