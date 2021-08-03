package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Emoto13/photo-viewer-rest/image-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/image_store"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/post"
	"github.com/Emoto13/photo-viewer-rest/image-service/src/service"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"go.etcd.io/etcd/clientv3"
)

var (
	host         = os.Getenv("POSTGRE_HOST")
	port         = os.Getenv("POSTGRE_PORT")
	dbuser       = os.Getenv("POSTGRE_USER")
	password     = os.Getenv("POSTGRE_PASSWORD")
	dbname       = os.Getenv("POSTGRE_DB_NAME")
	serverPort   = os.Getenv("IMAGE_SERVICE_PORT")
	awsS3Region  = os.Getenv("AWS_S3_REGION")
	fullHostname = os.Getenv("FULL_HOSTNAME")
	etcdAddress  = os.Getenv("ETCD_ADDRESS")
	etcdUsername = os.Getenv("ETCD_USERNAME")
	etcdPassword = os.Getenv("ETCD_PASSWORD")
)

func main() {
	router := getRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	fmt.Println("Starting to listen at port:", serverPort)
	http.ListenAndServe(serverPort, handler)
}

func getRouter() *mux.Router {
	db, err := OpenDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(awsS3Region),
	}))

	s3Service := s3.New(sess)
	s3Connector := image_store.NewS3Connector(s3Service, sess)

	authClient := auth.NewAuthClient(&http.Client{}, getAuthServiceAddress())
	postClient := post.NewPostClient(&http.Client{}, getPostServiceAddress())
	imageStore := image_store.NewImageStore(s3Connector, db)

	imageService := service.New(authClient, postClient, imageStore)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/image-service/upload-image", imageService.UploadImage).Methods("POST", "OPTIONS")
	router.HandleFunc("/image-service/health-check", imageService.HealthCheck).Methods("GET")

	return router
}

func OpenDatabaseConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, dbuser, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func registerService() {
	config := clientv3.Config{
		Endpoints:   []string{etcdAddress},
		DialTimeout: 15 * time.Second,
		Username:    etcdUsername,
		Password:    etcdPassword,
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, "image-service", fullHostname)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
}

func getAuthServiceAddress() string {
	config := clientv3.Config{
		Endpoints:   []string{etcdAddress},
		DialTimeout: 15 * time.Second,
		Username:    etcdUsername,
		Password:    etcdPassword,
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return fullHostname + ":10000"
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, "auth-service")
	cancel()
	if err != nil {
		fmt.Println("get failed err:", err)
		return fullHostname + ":10000"
	}

	if len(resp.Kvs) == 0 || resp.Kvs[0].Value == nil {
		return fullHostname + ":10000"
	}

	fmt.Println(string(resp.Kvs[0].Value))
	return string(resp.Kvs[0].Value)
}

func getPostServiceAddress() string {
	config := clientv3.Config{
		Endpoints:   []string{etcdAddress},
		DialTimeout: 15 * time.Second,
		Username:    etcdUsername,
		Password:    etcdPassword,
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return fullHostname + ":10005"
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, "post-service")
	cancel()
	if err != nil {
		fmt.Println("get failed err:", err)
		return fullHostname + ":10005"
	}

	if len(resp.Kvs) == 0 || resp.Kvs[0].Value == nil {
		fmt.Println("get failed err 2:", err)
		return fullHostname + ":10005"
	}

	fmt.Println(string(resp.Kvs[0].Value))
	return string(resp.Kvs[0].Value)
}
