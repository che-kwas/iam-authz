// Package queue defines the Queue interface.
package queue

import "context"

//go:generate mockgen -self_package=iam-authz/internal/authzserver/queue -destination mock_queue.go -package queue iam-authz/internal/authzserver/queue Queue

var que Queue

// Queue defines the behavior of a queue.
type Queue interface {
	PushMany(ctx context.Context, key string, values [][]byte) error
	Close() error
}

// Que returns the queue instance.
func Que() Queue {
	return que
}

// SetQue sets the queue instance.
func SetQue(queue Queue) {
	que = queue
}
