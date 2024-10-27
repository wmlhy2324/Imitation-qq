package core

import (
	"context"
	"github.com/go-redis/redis"
	"time"
)

var rdb *redis.Client

func InitRedis(addr string, pwd string, db int) (rdb *redis.Client) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       db,
		PoolSize: 100,
	})
	_, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	_, err := rdb.Ping().Result()
	if err != nil {
		panic(err)

	}
	return rdb
}
