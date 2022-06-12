package apiserver

import (
	"time"

	"github.com/spf13/viper"
)

const (
	confKey = "apiserver"

	defaultTimeout = time.Duration(5 * time.Second)
)

// APIServerOptions defines options for building an apiserver.
type APIServerOptions struct {
	Addr    string
	Timeout time.Duration
}

func NewAPIServerOptions() *APIServerOptions {
	opts := &APIServerOptions{
		Timeout: defaultTimeout,
	}

	_ = viper.UnmarshalKey(confKey, opts)
	return opts
}
