package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Emoto13/photo-viewer-rest/follow-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/follow-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/follow-service/src/service"
	"go.etcd.io/etcd/clientv3"

	"github.com/gorilla/mux"

	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

var (
	host         = os.Getenv("POSTGRE_HOST")
	port         = os.Getenv("POSTGRE_PORT")
	dbuser       = os.Getenv("POSTGRE_USER")
	password     = os.Getenv("POSTGRE_PASSWORD")
	dbname       = os.Getenv("POSTGRE_DB_NAME")
	serverPort   = os.Getenv("FOLLOW_SERVICE_PORT")
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
	registerService()
	fmt.Println("Starting to listen at port", serverPort)
	err := http.ListenAndServe(serverPort, handler)
	fmt.Println(err)
}

func getRouter() *mux.Router {
	db, err := OpenDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	followStore := follow.NewFollowStore(db)
	authClient := auth.NewAuthClient(&http.Client{}, getAuthServiceAddress())
	followService := service.NewFollowService(authClient, followStore)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/follow-service/follow", followService.Follow).Methods("POST")
	router.HandleFunc("/follow-service/unfollow", followService.Unfollow).Methods("POST")
	router.HandleFunc("/follow-service/get-followers", followService.GetFollowers).Methods("GET")
	router.HandleFunc("/follow-service/get-following", followService.GetFollowing).Methods("GET")
	router.HandleFunc("/follow-service/get-suggestions", followService.GetSuggestions).Methods("GET")
	router.HandleFunc("/follow-service/health-check", followService.HealthCheck).Methods("GET")

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
	_, err = client.Put(ctx, "follow-service", fullHostname)
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
