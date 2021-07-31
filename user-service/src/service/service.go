package service

import (
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/user-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/user-service/src/store"
)

type userService struct {
	store        store.UserStore
	followClient follow.FollowClient
}

func NewUserService(store store.UserStore, followClient follow.FollowClient) *userService {
	return &userService{
		store:        store,
		followClient: followClient,
	}
}

func (us *userService) CreateUser(w http.ResponseWriter, r *http.Request) {
	body, err := getRequestBody(r.Body)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := store.NewUser(body["username"], body["password"], body["role"], body["email"])
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = us.store.SaveUser(user)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = us.followClient.CreateUser(body["username"])
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("User created successfully")
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "success"})
	return
}

func (us *userService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
