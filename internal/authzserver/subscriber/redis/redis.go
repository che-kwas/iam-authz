// Package redis implements `iam-authz/internal/authzserver/subscriber.Subscriber` interface.
package redis

import (
	"context"

	rdb "github.com/che-kwas/iam-kit/redis"
	"github.com/go-redis/redis/v8"

	"iam-authz/internal/authzserver/subscriber"
)

type redisSub struct {
	cli redis.UniversalClient
}

var _ subscriber.Subscriber = &redisSub{}

func (r *redisSub) PubSubLoop(ctx context.Context, channel string, handleFunc func(string)) {
	pubsub := r.cli.Subscribe(ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			handleFunc(msg.Payload)
		}
	}
}

func (r *redisSub) Close() error {
	return r.cli.Close()
}

// NewRedisSub returns a redis subscriber.
func NewRedisSub() (subscriber.Subscriber, error) {
	cli, err := rdb.NewRedisIns()
	if err != nil {
		return nil, err
	}

	return &redisSub{cli}, nil
}
