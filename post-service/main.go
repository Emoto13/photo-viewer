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
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/post_store/cache_store"
	"github.com/Emoto13/photo-viewer-rest/post-service/src/service"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
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
	serverPort   = os.Getenv("POST_SERVICE_PORT")
	redisPort    = os.Getenv("REDIS_PORT")
	fullHostname = os.Getenv("FULL_HOSTNAME")
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
	authClient := auth.NewAuthClient(&http.Client{}, getAuthServiceAddress())
	postStore := post_store.NewPostStore(db)
	redisCache := cache_store.NewPostCacheStore(newCacheDatabase())
	postService := service.New(authClient, postStore, redisCache)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/post-service/get-following-posts", postService.GetFollowingPosts).Methods("GET")
	router.HandleFunc("/post-service/search-posts/{query}", postService.SearchPosts).Methods("GET")
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

func newCacheDatabase() *cache.Cache {
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"server1": redisPort,
		},
	})

	redisCache := cache.New(&cache.Options{
		Redis:      ring,
		LocalCache: cache.NewTinyLFU(1000, 15*time.Minute),
	})

	return redisCache
}

func registerService() {
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
	}

	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, "post-service", fullHostname+serverPort)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
}

func getAuthServiceAddress() string {
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"},
		DialTimeout: 5 * time.Second,
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
