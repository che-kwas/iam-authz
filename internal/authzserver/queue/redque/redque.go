// Package redque implements `iam-authz/internal/authzserver/queue.Queue` interface.
package redque

import (
	"context"
	"time"

	rdb "github.com/che-kwas/iam-kit/redis"
	"github.com/go-redis/redis/v8"

	"iam-authz/internal/authzserver/queue"
)

const queueExpiration = time.Duration(24 * time.Hour)

type redisQueue struct {
	cli redis.UniversalClient
}

var _ queue.Queue = &redisQueue{}

func (r *redisQueue) Push(ctx context.Context, key string, values ...interface{}) error {
	pipe := r.cli.Pipeline()
	for _, record := range values {
		pipe.RPush(ctx, key, record)
	}
	pipe.Expire(ctx, key, queueExpiration)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *redisQueue) Close() error {
	return r.cli.Close()
}

// NewRedisQue returns a redis queue.
func NewRedisQue() (queue.Queue, error) {
	cli, err := rdb.NewRedisIns()
	if err != nil {
		return nil, err
	}

	return &redisQueue{cli}, nil
}
