package auditor

import (
	"time"

	"github.com/spf13/viper"
)

const (
	confKey = "audit"

	defaultEnable        = false
	defaultPoolSize      = 100
	defaultBufferSize    = 100
	defaultFlushInterval = time.Duration(time.Second)
	defaultOmitDetails   = true
)

// AuditorOptions defines options for building an auditor.
type AuditorOptions struct {
	Enable   bool
	PoolSize int `mapstructure:"pool-size"`

	// BufferSize defines the channel buffer size for receiving AuditRecord.
	BufferSize int `mapstructure:"buffer-size"`

	// OmitDetails defines whether to omit the details in AuditRecord.
	OmitDetails bool `mapstructure:"omit-details"`
}

func NewAuditorOptions() *AuditorOptions {
	opts := &AuditorOptions{
		Enable:      defaultEnable,
		PoolSize:    defaultPoolSize,
		BufferSize:  defaultBufferSize,
		OmitDetails: defaultOmitDetails,
	}

	_ = viper.UnmarshalKey(confKey, opts)
	return opts
}
