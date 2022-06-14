// Package subscriber defines the Subscriber interface.
package subscriber

import "context"

//go:generate mockgen -self_package=iam-authz/internal/authzserver/subscriber -destination mock_subscriber.go -package subscriber iam-authz/internal/authzserver/subscriber Subscriber

var sub Subscriber

// Subscriber defines the behavior of a subscriber.
type Subscriber interface {
	PubSubLoop(ctx context.Context, channel string, handleFunc func(string))
	Close() error
}

// Sub returns the subscriber instance.
func Sub() Subscriber {
	return sub
}

// SetSub sets the subscriber instance.
func SetSub(subscriber Subscriber) {
	sub = subscriber
}
