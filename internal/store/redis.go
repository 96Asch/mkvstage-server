package store

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
)

func GetRedis(host, port string) (*redis.Client, error) {
	fmt.Printf("Connecting to Redis: (%s:%s)\n", host, port)
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
