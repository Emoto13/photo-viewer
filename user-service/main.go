package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Emoto13/photo-viewer-rest/user-service/src/service"
	"github.com/Emoto13/photo-viewer-rest/user-service/src/store"
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
	serverPort   = os.Getenv("USER_SERVICE_PORT")
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
	fmt.Println("Starting to listen at port", serverPort)
	registerService()
	http.ListenAndServe(serverPort, handler)
}

func getRouter() *mux.Router {
	db, err := OpenDatabaseConnection()
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)

	}

	userStore := store.NewUserStore(db)
	userService := service.NewUserService(userStore)

	router := mux.NewRouter()
	router.StrictSlash(true)
	router.HandleFunc("/user-service/create-user", userService.CreateUser).Methods("POST")
	router.HandleFunc("/user-service/health-check", userService.HealthCheck).Methods("POST")

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
	_, err = client.Put(ctx, "user-service", fullHostname+serverPort)
	cancel()
	if err != nil {
		fmt.Println("put failed, err:", err)
		return
	}
}
