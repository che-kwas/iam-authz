package cache

import (
	"context"

	"github.com/che-kwas/iam-kit/logger"

	"iam-authz/internal/authzserver/subscriber"
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
	sub    subscriber.Subscriber
	loader Loadable
	ctx    context.Context
	log    *logger.Logger
}

// NewLoader creates a loader by a subscriber and a loaderImpl.
func NewLoader(ctx context.Context, sub subscriber.Subscriber, loaderImpl Loadable) *Loader {
	return &Loader{
		sub:    sub,
		loader: loaderImpl,
		ctx:    ctx,
		log:    logger.L(),
	}
}

// Start starts a reloading loop.
func (l *Loader) Start() {
	go l.sub.PubSubLoop(l.ctx, channel, l.reload)

	l.reloadAll()
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
