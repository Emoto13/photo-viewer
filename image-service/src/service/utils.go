package service

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
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

func getRequestBody(requestBody io.ReadCloser) (map[string]interface{}, error) {
	bodyBytes, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
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
