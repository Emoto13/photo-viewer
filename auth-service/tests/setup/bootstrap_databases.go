package setup

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/go-redis/redis/v8"

	_ "github.com/lib/pq"
)

var (
	posgre_host         = os.Getenv("POSTGRE_HOST")
	port                = os.Getenv("POSTGRE_PORT")
	dbuser              = os.Getenv("POSTGRE_USER")
	password            = os.Getenv("POSTGRE_PASSWORD")
	test_dbname         = os.Getenv("POSTGRE_TEST_DB_NAME")
	redisDatabaseNumber = os.Getenv("REDIS_TOKEN_TEST_DATABASE")
	redisAddress        = os.Getenv("HOSTNAME") + os.Getenv("REDIS_PORT")
)

func OpenPostgresDatabaseConnection() (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s",
		posgre_host, port, dbuser, password, test_dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func CreateRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddress,
		Password: "",
		DB:       getRedisDatabaseNumber(),
	})
	return rdb
}

func getRedisDatabaseNumber() int {
	val, err := strconv.Atoi(redisDatabaseNumber)
	if err != nil {
		return 11
	}
	return val
}
