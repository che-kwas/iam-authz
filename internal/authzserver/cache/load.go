package cache

import (
	"context"

	"github.com/che-kwas/iam-kit/logger"

	"iam-authz/internal/pkg/redis"
)

// Redis pub/sub events.
const (
	channel            = "iam.notifications"
	eventPolicyChanged = "PolicyChanged"
	eventSecretChanged = "SecretChanged"
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
	go l.startEventLoop()

	l.reloadAll()
}

func (l *Loader) startEventLoop() {
	pubsub := redis.Client().Subscribe(l.ctx, channel)
	defer pubsub.Close()

	ch := pubsub.Channel()
	for {
		select {
		case <-l.ctx.Done():
			return
		case msg := <-ch:
			l.reload(msg.Payload)
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
