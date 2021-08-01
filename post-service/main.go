package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Emoto13/photo-viewer-rest/post-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/service"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"go.etcd.io/etcd/clientv3"
)

var (
	host              = os.Getenv("POSTGRE_HOST")
	port              = os.Getenv("POSTGRE_PORT")
	dbuser            = os.Getenv("POSTGRE_USER")
	password          = os.Getenv("POSTGRE_PASSWORD")
	dbname            = os.Getenv("POSTGRE_DB_NAME")
	serverPort        = os.Getenv("POST_SERVICE_PORT")
	fullHostname      = os.Getenv("FULL_HOSTNAME")
	etcdAddress       = os.Getenv("ETCD_ADDRESS")
	etcdUsername      = os.Getenv("ETCD_USERNAME")
	etcdPassword      = os.Getenv("ETCD_PASSWORD")
	cassandraAddress  = os.Getenv("CASSANDRA_ADDRESS")
	cassandraUsername = os.Getenv("CASSANDRA_USERNAME")
	cassandraPassword = os.Getenv("CASSANDRA_PASSWORD")
	cassandraKeyspace = os.Getenv("CASSANDRA_KEYSPACE")
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
	fmt.Println("Starting to listen at port:", serverPort)
	http.ListenAndServe(serverPort, handler)
}

func getRouter() *mux.Router {
	db, err := OpenDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}

	session, err := newCassandraSession()
	if err != nil {
		fmt.Println("cannot open Cassandra session:", err.Error())
	}

	authClient := auth.NewAuthClient(&http.Client{}, getAuthServiceAddress())
	followClient := follow.NewFollowClient(&http.Client{}, getFollowServiceAddress())

	postStore := post_store.NewPostStore(db, session)
	postService := service.New(authClient, postStore, followClient)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/post-service/search-posts/{query}", postService.SearchPosts).Methods("GET")
	router.HandleFunc("/post-service/create-post", postService.CreatePost).Methods("POST")
	router.HandleFunc("/post-service/posts/{username}", postService.GetUserPosts).Methods("GET")
	router.HandleFunc("/post-service/health-check", postService.HealthCheck).Methods("GET")

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

func newCassandraSession() (*gocql.Session, error) {
	cluster := gocql.NewCluster(cassandraAddress)
	cluster.ProtoVersion = 4
	cluster.ConnectTimeout = time.Second * 15
	cluster.Keyspace = cassandraKeyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{Username: cassandraUsername, Password: cassandraPassword}

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, err
	}

	return session, nil
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
	_, err = client.Put(ctx, "post-service", fullHostname)
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

func getFollowServiceAddress() string {
	config := clientv3.Config{
		Endpoints:   []string{etcdAddress},
		DialTimeout: 15 * time.Second,
		Username:    etcdUsername,
		Password:    etcdPassword,
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return fullHostname + ":10001"
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	resp, err := client.Get(ctx, "follow-service")
	cancel()
	if err != nil {
		fmt.Println("get failed err:", err)
		return fullHostname + ":10001"
	}

	if len(resp.Kvs) == 0 || resp.Kvs[0].Value == nil {
		return fullHostname + ":10001"
	}

	fmt.Println(string(resp.Kvs[0].Value))
	return string(resp.Kvs[0].Value)
}
