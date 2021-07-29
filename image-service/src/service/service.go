package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store/image_data"
)

type imageService struct {
	authClient auth.AuthClient
	imageStore image_store.ImageStore
}

func New(authClient auth.AuthClient, imageStore image_store.ImageStore) *imageService {
	return &imageService{
		authClient: authClient,
		imageStore: imageStore,
	}
}

func (s *imageService) UploadImage(w http.ResponseWriter, r *http.Request) {
	username, err := s.authClient.Authenticate(r.Header.Get("Authorization"))
	if err != nil {
		fmt.Println("Failed to authenticate")
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	imageName := r.FormValue("imageName")
	fileName := r.FormValue("fileName")
	image, _, _ := r.FormFile("image")
	defer image.Close()
	fmt.Println(imageName, fileName)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, image); err != nil {
		fmt.Println(err)
		return
	}

	err = s.imageStore.UploadImage(&image_data.UploadImage{
		Name:     imageName,
		FileName: fileName,
		Data:     buf.Bytes(),
		Owner:    username,
	})

	if err != nil {
		fmt.Println("failed to upload: ", err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	fmt.Println("Image ", imageName, " was successfully uploaded")
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "success"})
	return
}

func (s *imageService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
