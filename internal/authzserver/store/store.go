// Package store defines the Store interface.
package store

import "context"

//go:generate mockgen -self_package=iam-authz/internal/authzserver/store -destination mock_store.go -package store iam-authz/internal/authzserver/store Store,SecretStore,PolicyStore

var client Store

// Store defines the iam platform storage interface.
type Store interface {
	Policies() PolicyStore
	Secrets() SecretStore
	Close(ctx context.Context) error
}

// Client returns the store client.
func Client() Store {
	return client
}

// SetClient sets the store client.
func SetClient(store Store) {
	client = store
}
