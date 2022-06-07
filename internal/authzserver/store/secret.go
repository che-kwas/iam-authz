package store

import pb "iam-authz/api/apiserver/proto/v1"

// SecretStore defines the secret storage interface.
type SecretStore interface {
	List() (map[string]*pb.SecretInfo, error)
}
