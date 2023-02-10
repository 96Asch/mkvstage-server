package repository

import "github.com/go-redis/redis/v8"

type redisTokenRepo struct {
	R *redis.Client
}

func NewRedisTokenRepository(client *redis.Client) *redisTokenRepo {
	return &redisTokenRepo{
		R: client,
	}
}
