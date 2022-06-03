// Package load subscribes to PolicyChanged/SecretChanged events
// and reloads data when changes occur.
package load

import (
	"context"
	"time"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/spf13/viper"

	"iam-auth/internal/pkg/redis"
)

// Redis pub/sub events.
const (
	channel            = "iam.notifications"
	eventPolicyChanged = "PolicyChanged"
	eventSecretChanged = "SecretChanged"

	defaultReloadInterval = time.Duration(30 * time.Second)
)

// Loadable defines the behavior of a loader.
type Loadable interface {
	ReloadSecrets() error
	ReloadPolicies() error
}

// Loader is used to do reload.
type Loader struct {
	ctx    context.Context
	loader Loadable
	log    *logger.Logger
}

// NewLoader creates a loader with a loaderImpl.
func NewLoader(ctx context.Context, loaderImpl Loadable) *Loader {
	return &Loader{
		ctx:    ctx,
		loader: loaderImpl,
		log:    logger.L(),
	}
}

// Start starts a reloading loop.
func (l *Loader) Start() {
	// reload when event is received
	go l.startEventLoop()
	// reload when timer ticks
	l.startTimerLoop()
}

func (l *Loader) startEventLoop() {
	pubsub := redis.Client().Subscribe(l.ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		for msg := range ch {
			l.reload(msg.Payload)
		}
	}
}

func (l *Loader) startTimerLoop() {
	var interval time.Duration
	if err := viper.UnmarshalKey("main.reload-interval", &interval); err != nil {
		interval = defaultReloadInterval
	}
	ticker := time.NewTicker(interval)

	for {
		select {
		case <-l.ctx.Done():
			return
		case <-ticker.C:
			l.reloadAll()
		}
	}
}

func (l *Loader) reloadAll() {
	if err := l.loader.ReloadSecrets(); err != nil {
		l.log.Errorf("faild to reload secrets: %s", err.Error())
	}
	if err := l.loader.ReloadPolicies(); err != nil {
		l.log.Errorf("faild to reload policies: %s", err.Error())
	}
}

func (l *Loader) reload(event string) {
	switch event {
	case eventSecretChanged:
		if err := l.loader.ReloadSecrets(); err != nil {
			l.log.Errorf("faild to reload secrets: %s", err.Error())
		}
	case eventPolicyChanged:
		if err := l.loader.ReloadPolicies(); err != nil {
			l.log.Errorf("faild to reload policies: %s", err.Error())
		}
	}
}
