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
	defaultOmitDetails   = true
)

// AuditorOptions defines options for building an auditor.
type AuditorOptions struct {
	Enable   bool
	PoolSize int `mapstructure:"pool-size"`

	// BufferSize defines the buffer size for quantitative delivery.
	BufferSize int `mapstructure:"buffer-size"`

	// FlushInterval defines the interval for scheduled delivery.
	FlushInterval time.Duration `mapstructure:"flush-interval"`

	// OmitDetails defines whether to omit the details in AuditRecord.
	OmitDetails bool `mapstructure:"omit-details"`
}

func NewAuditorOptions() *AuditorOptions {
	opts := &AuditorOptions{
		Enable:        defaultEnable,
		PoolSize:      defaultPoolSize,
		BufferSize:    defaultBufferSize,
		FlushInterval: defaultFlushInterval,
		OmitDetails:   defaultOmitDetails,
	}

	_ = viper.UnmarshalKey(confKey, opts)
	return opts
}
