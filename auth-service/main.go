package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Emoto13/photo-viewer-rest/auth-service/src/service"
	"github.com/Emoto13/photo-viewer-rest/auth-service/src/token"
	"github.com/Emoto13/photo-viewer-rest/auth-service/src/user"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"go.etcd.io/etcd/clientv3"

	_ "github.com/lib/pq"
)

var (
	host                = os.Getenv("POSTGRE_HOST")
	port                = os.Getenv("POSTGRE_PORT")
	dbuser              = os.Getenv("POSTGRE_USER")
	password            = os.Getenv("POSTGRE_PASSWORD")
	dbname              = os.Getenv("POSTGRE_DB_NAME")
	serverPort          = os.Getenv("AUTH_SERVICE_PORT")
	fullHostname        = os.Getenv("FULL_HOSTNAME") // + ":10000"
	redisDatabaseNumber = os.Getenv("REDIS_TOKEN_DATABASE_NUMBER")
	redisAddress        = os.Getenv("REDIS_ADDRESS")
	redisPassword       = os.Getenv("REDIS_PASSWORD")
	etcdAddress         = os.Getenv("ETCD_ADDRESS")
	etcdUsername        = os.Getenv("ETCD_USERNAME")
	etcdPassword        = os.Getenv("ETCD_PASSWORD")
)

func main() {
	router := getRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	registerService()
	handler := c.Handler(router)
	fmt.Println("Starting to listen at port", serverPort)
	http.ListenAndServe(serverPort, handler)
}

func OpenDatabaseConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, dbuser, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	fmt.Println("connected to: ", host, port)
	return db, nil
}
func getRouter() *mux.Router {
	redisClient := createRedisClient()
	db, err := OpenDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to open database connection: %v", err)
	}
	state := user.NewState(db)

	tokenManager := token.NewTokenManager(redisClient)
	authServer := service.New(state, tokenManager)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/auth-service/login", authServer.Login).Methods("POST")
	router.HandleFunc("/auth-service/logout", authServer.Logout).Methods("POST")
	router.HandleFunc("/auth-service/authenticate", authServer.Authenticate).Methods("GET")
	router.HandleFunc("/auth-service/health-check", authServer.HealthCheck).Methods("GET")
	return router
}

func createRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: redisPassword,
		DB:       getRedisDatabaseNumber(),
	})

	return rdb
}

func getRedisDatabaseNumber() int {
	val, err := strconv.Atoi(redisDatabaseNumber)
	if err != nil {
		return 1
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
		fmt.Println("error connecting to etcd:", err)
		return
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err = client.Put(ctx, "auth-service", fullHostname)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}

	fmt.Println("write to etcd successful")
}
