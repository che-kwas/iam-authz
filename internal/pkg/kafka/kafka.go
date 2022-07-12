// Package kafka is the kafka producer builder.
package kafka // import "iam-authz/internal/pkg/kafka"

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/spf13/viper"

	"github.com/che-kwas/iam-kit/logger"
)

const (
	confKey = "kafka"

	defaultFlushFrequency = time.Second
	defaultFlushMessages  = 1000
)

// KafkaOptions defines options for building a kafka producer.
type KafkaOptions struct {
	Brokers        []string
	FlushFrequency time.Duration `mapstructure:"flush-frequency"`
	FlushMessages  int           `mapstructure:"flush-messages"`
}

// NewKafkaProducer creates a kafka producer.
func NewKafkaProducer() (sarama.AsyncProducer, error) {
	opts, err := getKafkaOpts()
	if err != nil {
		return nil, err
	}
	logger.L().Debugf("new kafka instance with options: %+v", opts)

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal     // Only wait for the leader to ack
	config.Producer.Compression = sarama.CompressionSnappy // Compress messages
	config.Producer.Flush.Frequency = opts.FlushFrequency
	config.Producer.Flush.Messages = opts.FlushMessages

	producer, err := sarama.NewAsyncProducer(opts.Brokers, config)
	if err != nil {
		return nil, err
	}

	// We will just log to STDOUT if we're not able to produce messages.
	// Note: messages will only be returned here after all retry attempts are exhausted.
	go func() {
		for err := range producer.Errors() {
			logger.L().Errorw("kafka producer", "error", err)
		}
	}()

	return producer, nil
}

func getKafkaOpts() (*KafkaOptions, error) {
	opts := &KafkaOptions{
		FlushFrequency: defaultFlushFrequency,
		FlushMessages:  defaultFlushMessages,
	}

	if err := viper.UnmarshalKey(confKey, opts); err != nil {
		return nil, err
	}
	return opts, nil
}
