package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Emoto13/photo-viewer-rest/feed-service/src/auth"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/feed"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/follow"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/post/cache_store"
	"github.com/Emoto13/photo-viewer-rest/feed-service/src/service"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.etcd.io/etcd/clientv3"
)

var (
	serverPort          = os.Getenv("FEED_SERVICE_PORT")
	redisPort           = os.Getenv("REDIS_PORT")
	fullHostname        = os.Getenv("FULL_HOSTNAME")
	redisDatabaseNumber = os.Getenv("REDIS_POST_DATABASE_NUMBER")
	etcdAddress         = os.Getenv("ETCD_ADDRESS")
	etcdUsername        = os.Getenv("ETCD_USERNAME")
	etcdPassword        = os.Getenv("ETCD_PASSWORD")
	cassandraAddress    = os.Getenv("CASSANDRA_ADDRESS")
	cassandraUsername   = os.Getenv("CASSANDRA_USERNAME")
	cassandraPassword   = os.Getenv("CASSANDRA_PASSWORD")
	cassandraKeyspace   = os.Getenv("CASSANDRA_KEYSPACE")
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
	authClient := auth.NewAuthClient(&http.Client{}, getAuthServiceAddress())
	followClient := follow.NewFollowClient(&http.Client{}, getFollowServiceAddress())
	postClient := post.NewPostClient(&http.Client{}, getPostServiceAddress())
	postCacheStore := cache_store.NewPostCacheStore(newCacheDatabase())

	session, err := newCassandraSession()
	if err != nil {
		log.Fatal(err.Error())
	}
	feedStore := feed.NewFeedStore(postClient, session)
	feedService := service.NewFeedService(authClient, followClient, postClient, postCacheStore, feedStore)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/feed-service/get-feed", feedService.GetFollowingPosts).Methods("GET")
	router.HandleFunc("/feed-service/update-feed", feedService.UpdateFeed).Methods("PATCH")
	router.HandleFunc("/feed-service/add-to-feed", feedService.AddToFeed).Methods("PATCH")
	router.HandleFunc("/feed-service/add-to-followers-feed", feedService.AddToFollowersFeed).Methods("PATCH")

	router.HandleFunc("/feed-service/health-check", feedService.HealthCheck).Methods("GET")
	return router
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

func getRedisDatabaseNumber() int {
	val, err := strconv.Atoi(redisDatabaseNumber)
	if err != nil {
		return 0
	}
	return val
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
	_, err = client.Put(ctx, "feed-service", fullHostname)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	fmt.Println("feed-service registered at: ", fullHostname)
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
		return fullHostname + ":10005"
	}

	fmt.Println(string(resp.Kvs[0].Value))
	return string(resp.Kvs[0].Value)
}
