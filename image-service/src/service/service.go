package service

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/feed"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store/image_data"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/post"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/post/models"
)

type imageService struct {
	authClient auth.AuthClient
	postClient post.PostClient
	feedClient feed.FeedClient
	imageStore image_store.ImageStore
}

func New(authClient auth.AuthClient, postClient post.PostClient, feedClient feed.FeedClient, imageStore image_store.ImageStore) *imageService {
	return &imageService{
		authClient: authClient,
		imageStore: imageStore,
		feedClient: feedClient,
		postClient: postClient,
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

	path, err := s.imageStore.UploadImage(&image_data.UploadImage{
		Name:     imageName,
		FileName: fileName,
		Data:     buf.Bytes(),
		Owner:    username,
	})

	if err != nil {
		fmt.Println("failed to upload: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	post := &models.Post{
		Username:  username,
		Name:      imageName,
		Path:      path,
		CreatedOn: time.Now(),
	}
	fmt.Println("Post to create", post)

	err = s.feedClient.AddToFollowersFeed(r.Header.Get("Authorization"), post)
	if err != nil {
		fmt.Println("failed to add to followers feed: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	err = s.postClient.CreatePost(r.Header.Get("Authorization"), post)
	if err != nil {
		fmt.Println("failed to create post: ", err.Error())
		respondWithError(w, http.StatusBadRequest, err.Error())
	}

	fmt.Println("Image ", imageName, " was successfully uploaded")
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "success"})
	return
}

func (s *imageService) HealthCheck(w http.ResponseWriter, r *http.Request) {
	respondWithJSON(w, http.StatusOK, map[string]string{"Status": "OK"})
	return
}
