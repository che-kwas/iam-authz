// Package store defines the storage interface for authserver.
package store

//go:generate mockgen -self_package=iam-auth/internal/authserver/store -destination mock_store.go -package store iam-auth/internal/authserver/store Factory,SecretStore,PolicyStore

var client Store

// Store defines the iam platform storage interface.
type Store interface {
	Policies() PolicyStore
	Secrets() SecretStore
	Close() error
}

// Client returns the store client.
func Client() Store {
	return client
}

// SetClient sets the store client.
func SetClient(store Store) {
	client = store
}
