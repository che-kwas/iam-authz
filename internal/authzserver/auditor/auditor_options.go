package auditor

import (
	"time"

	"github.com/spf13/viper"
)

const (
	confKey = "audit"

	defaultEnable        = false
	defaultPoolSize      = 50
	defaultBufferSize    = 2000
	defaultFlushInterval = time.Duration(time.Second)
)

// AuditorOptions defines options for building an auditor.
type AuditorOptions struct {
	Enable   bool
	PoolSize int `mapstructure:"pool-size"`

	// BufferSize defines the buffer size for quantitative delivery.
	BufferSize int `mapstructure:"buffer-size"`

	// BufferSize defines the interval for scheduled delivery.
	FlushInterval time.Duration `mapstructure:"flush-interval"`
}

func NewAuditorOptions() *AuditorOptions {
	opts := &AuditorOptions{
		Enable:        defaultEnable,
		PoolSize:      defaultPoolSize,
		BufferSize:    defaultBufferSize,
		FlushInterval: defaultFlushInterval,
	}

	_ = viper.UnmarshalKey(confKey, opts)
	return opts
}
