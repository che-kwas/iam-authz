// Package redis implements `iam-authz/internal/authzserver/subscriber.Subscriber` interface.
package redis

import (
	"context"

	rdb "github.com/che-kwas/iam-kit/redis"
	"github.com/go-redis/redis/v8"

	"iam-authz/internal/authzserver/subscriber"
)

type redisSub struct {
	cli    redis.UniversalClient
	pubsub *redis.PubSub
}

var _ subscriber.Subscriber = &redisSub{}

func (r *redisSub) PubSubLoop(ctx context.Context, channel string, handleFunc func(string)) {
	r.pubsub = r.cli.Subscribe(ctx, channel)

	ch := r.pubsub.Channel()
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-ch:
			handleFunc(msg.Payload)
		}
	}
}

func (r *redisSub) Close(ctx context.Context) error {
	r.pubsub.Close()
	return r.cli.Close()
}

// NewRedisSub returns a redis subscriber.
func NewRedisSub() (subscriber.Subscriber, error) {
	cli, err := rdb.NewRedisIns()
	if err != nil {
		return nil, err
	}

	return &redisSub{cli: cli}, nil
}
