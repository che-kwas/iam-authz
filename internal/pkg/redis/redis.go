// Package redis defines the global redis client.
package redis

import (
	"sync"

	"github.com/che-kwas/iam-kit/redis"
	redisv8 "github.com/go-redis/redis/v8"
)

var (
	rdb  redisv8.UniversalClient
	once sync.Once
)

func Client() redisv8.UniversalClient {
	if rdb != nil {
		return rdb
	}

	once.Do(func() {
		rdb, _ = redis.NewRedisIns()
	})

	return rdb
}
