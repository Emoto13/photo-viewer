package service

import (
	"fmt"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/auth-service/src/token"
	"github.com/Emoto13/photo-viewer-rest/auth-service/src/user"
)

type AuthServer struct {
	state        user.State
	tokenManager token.TokenManager
}

func New(state user.State, tokenManager token.TokenManager) *AuthServer {
	return &AuthServer{state, tokenManager}
}

func (service *AuthServer) Login(w http.ResponseWriter, r *http.Request) {
	body, err := getRequestBody(r.Body)
	if err != nil {
		fmt.Println("bad login request", err)
		respondWithError(w, http.StatusBadRequest, "Bad login request")
		return
	}

	user, err := service.state.RetrieveUser(body["username"])
	if err != nil {
		fmt.Println("didn't retrieve user", err)
		respondWithError(w, http.StatusBadRequest, "Wrong username/password")
		return
	}

	if user == nil || !user.IsCorrectPassword(body["password"]) {
		fmt.Println("wrong username/password", err)
		respondWithError(w, http.StatusBadRequest, "Wrong username/password")
		return
	}

	authToken, err := service.tokenManager.GenerateToken(body["username"])
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	err = service.tokenManager.SaveToken(token.NewToken(authToken, body["username"]))
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to save token")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"token": authToken})
	return
}

func (service *AuthServer) Logout(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		fmt.Println("Invalid authorization header")
		respondWithError(w, http.StatusBadRequest, "Invalid authorization header")
		return
	}
	fmt.Println("here")

	token, err := getAuthToken(authHeader)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Println("here2")

	err = service.tokenManager.RemoveToken(token)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't logout")
		return
	}

	fmt.Println("here3")
	respondWithJSON(w, http.StatusOK, map[string]string{"Message": "Successfully logged out"})
	return
}

func (service *AuthServer) Authenticate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Authenticating", r)
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		respondWithError(w, http.StatusBadRequest, "Invalid authorization header")
		return
	}

	token, err := getAuthToken(authHeader)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	username, err := service.tokenManager.GetUsernameFromToken(token)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid token")
		return
	}

	fmt.Println(username)
	respondWithJSON(w, http.StatusOK, map[string]string{"username": username})
	return
}

func (service *AuthServer) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
