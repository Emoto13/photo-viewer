package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Emoto13/photo-viewer-rest/feed-service/src/feed/models"
	postModels "github.com/Emoto13/photo-viewer-rest/feed-service/src/post/models"
)

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload map[string]string) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func getFollowingsRequestBody(requestBody io.ReadCloser) (map[string][]*models.Following, error) {
	bodyBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}

	var result map[string][]*models.Following
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getPostRequestBody(requestBody io.ReadCloser) (map[string]*postModels.Post, error) {
	bodyBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}

	var result map[string]*postModels.Post
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func getAuthToken(authHeader string) (string, error) {
	values := strings.Split(authHeader, "Bearer ")
	if len(values) < 2 {
		return "", fmt.Errorf("Authorization token not set")
	}

	return values[1], nil
}
