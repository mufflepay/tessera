package db

import (
	"time"

	"github.com/go-redis/redis"
)

type IRedis interface {
	GetRedis() *redis.Client
}

type redisCache struct {
	host string
	db   int
	exp  time.Duration
}

func NewRedisCache(host string, db int, exp time.Duration) IRedis {
	return &redisCache{
		host: host,
		db:   db,
		exp:  exp,
	}
}

func (c *redisCache) GetRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     c.host,
		Password: "",
		DB:       c.db,
	})
}
