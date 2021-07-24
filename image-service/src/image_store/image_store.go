package image_store

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store/image_data"
)

type ImageStore interface {
	UploadImage(imageData *image_data.UploadImage) error
}

type imageStore struct {
	connector S3Connector
	db        *sql.DB
	mu        sync.RWMutex
}

func NewImageStore(connector S3Connector, db *sql.DB) ImageStore {
	return &imageStore{
		connector: connector,
		db:        db,
		mu:        sync.RWMutex{},
	}
}

func (store *imageStore) UploadImage(imageData *image_data.UploadImage) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	fileUrl, err := store.connector.UploadFile(imageData)
	if err != nil {
		fmt.Println("could not upload to s3")
		return err
	}

	_, err = store.db.Exec(addImageToDatabase, imageData.Owner, fileUrl, imageData.Name)
	if err != nil {
		return err
	}

	return nil
}
