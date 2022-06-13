// Package apiserver is the proxy of iam-apiserver
// which implements `iam-authz/internal/authzserver/store.Store` interface.
package apiserver

import (
	"errors"

	pb "iam-authz/api/apiserver/proto/v1"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"iam-authz/internal/authzserver/store"
)

type datastore struct {
	cli  pb.CacheClient
	conn *grpc.ClientConn
}

func (ds *datastore) Secrets() store.SecretStore {
	return newSecrets(ds)
}

func (ds *datastore) Policies() store.PolicyStore {
	return newPolicies(ds)
}

func (ds *datastore) Close() error {
	return ds.conn.Close()
}

// NewAPIServerStore returns a apiserver store instance.
func NewAPIServerStore(opts *APIServerOptions) (store.Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	dialOpts := []grpc.DialOption{
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	conn, err := grpc.DialContext(ctx, opts.Addr, dialOpts...)
	if err != nil {
		if err == ctx.Err() {
			err = errors.New("connect to apiserver timeout")
		}
		return nil, err
	}

	return &datastore{
		cli:  pb.NewCacheClient(conn),
		conn: conn,
	}, nil
}
