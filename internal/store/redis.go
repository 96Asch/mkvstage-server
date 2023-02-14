package store

import (
	"context"
	"fmt"
	"log"

	"github.com/go-redis/redis/v8"
)

func GetRedis(host, port string) (*redis.Client, error) {
	log.Printf("Connecting to Redis: (%s:%s)\n", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}
	log.Println("Connection to Redis succeeded!")
	return rdb, nil
}
