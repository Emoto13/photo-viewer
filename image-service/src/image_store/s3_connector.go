package image_store

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store/image_data"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	AWS_S3_REGION = os.Getenv("AWS_S3_REGION")
	AWS_S3_BUCKET = os.Getenv("AWS_S3_BUCKET")
)

type S3Connector interface {
	UploadFile(imageData *image_data.UploadImage) (string, error)
	DownloadFile(name string) ([]byte, error)
}

type s3Connector struct {
	s3Service  *s3.S3
	awsSession *session.Session
}

func NewS3Connector(s3Service *s3.S3, awsSession *session.Session) S3Connector {
	return &s3Connector{
		s3Service:  s3Service,
		awsSession: awsSession,
	}
}

func (conn *s3Connector) UploadFile(imageData *image_data.UploadImage) (string, error) {
	uploader := s3manager.NewUploaderWithClient(conn.s3Service)
	hashedFilename := sha256.Sum256([]byte(imageData.Owner + imageData.Name + imageData.FileName))
	hash := fmt.Sprintf("%x", hashedFilename)
	extension := filepath.Ext(imageData.FileName)
	key := hash + extension

	upParams := &s3manager.UploadInput{
		Bucket: aws.String(AWS_S3_BUCKET),
		Key:    aws.String(key),
		Body:   bytes.NewReader(imageData.Data),
	}

	_, err := uploader.Upload(upParams)
	if err != nil {
		return "", err
	}

	fileUrl := fmt.Sprintf("https://%s.s3-%s.amazonaws.com/%s", AWS_S3_BUCKET, AWS_S3_REGION, key)
	return fileUrl, nil
}

func (conn *s3Connector) DownloadFile(name string) ([]byte, error) {
	downloader := s3manager.NewDownloader(conn.awsSession)
	file, err := os.Create(name)
	defer file.Close()

	if err != nil {
		return nil, err
	}

	_, err = downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(AWS_S3_BUCKET),
			Key:    aws.String(name),
		})

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return data, nil
}
