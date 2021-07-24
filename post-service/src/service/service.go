package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/cache_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/post_data"
	"github.com/gorilla/mux"
)

type postService struct {
	authClient auth.AuthClient
	postStore  post_store.PostStore
	postCache  cache_store.PostCacheStore
}

func New(authClient auth.AuthClient, postStore post_store.PostStore, postCache cache_store.PostCacheStore) *postService {
	return &postService{
		authClient: authClient,
		postStore:  postStore,
		postCache:  postCache,
	}
}

func (s *postService) GetFollowingPosts(w http.ResponseWriter, r *http.Request) {
	username, err := s.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("Couldn't authenticate")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	cachedResult, err := s.postCache.Get(context.Background(), username)
	if err == nil {
		response, _ := json.Marshal(map[string][]*post_data.PostData{"posts": cachedResult})
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

	posts, err := s.postStore.RetrieveFollowingPosts(username)
	if err != nil {
		fmt.Println("Couldnt retrieve following posts")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	s.postCache.Set(context.Background(), username, posts)

	response, _ := json.Marshal(map[string][]*post_data.PostData{"posts": posts})
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

	response, _ := json.Marshal(map[string][]*post_data.PostData{"posts": posts})
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
	return
}

func (s *postService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
