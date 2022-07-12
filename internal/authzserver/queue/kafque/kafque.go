// Package kafque implements `iam-authz/internal/authzserver/queue.Queue` interface.
package kafque

import (
	"context"

	"iam-authz/internal/authzserver/queue"
	"iam-authz/internal/pkg/kafka"

	"github.com/Shopify/sarama"
)

const topic = "iam"

type kafkaQue struct {
	producer sarama.AsyncProducer
}

var _ queue.Queue = &kafkaQue{}

func (k *kafkaQue) Push(ctx context.Context, key string, value []byte) error {
	k.producer.Input() <- &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(key),
		Value: sarama.ByteEncoder(value),
	}

	return nil
}

func (k *kafkaQue) Close(ctx context.Context) error {
	return k.producer.Close()
}

// NewKafkaQue returns a kafka queue.
func NewKafkaQue() (queue.Queue, error) {
	producer, err := kafka.NewKafkaProducer()
	if err != nil {
		return nil, err
	}

	return &kafkaQue{producer}, nil
}
