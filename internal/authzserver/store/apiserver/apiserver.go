// Package apiserver is the proxy of iam-apiserver
// which implements `iam-authz/internal/authzserver/store.Store` interface.
package apiserver

import (
	"sync"

	pb "iam-authz/api/apiserver/proto/v1"

	"github.com/marmotedu/errors"
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

var (
	apiserverStore store.Store
	once           sync.Once
)

// APIServerStore returns a apiserver store instance.
func APIServerStore(opts *APIServerOptions) (store.Store, error) {
	if apiserverStore != nil {
		return apiserverStore, nil
	}

	var conn *grpc.ClientConn
	var err error
	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		dialOpts := []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		}

		conn, err = grpc.DialContext(ctx, opts.Addr, dialOpts...)
		if err == ctx.Err() {
			err = errors.Wrap(err, "connect to apiserver timeout")
		}
	})

	if err == nil {
		apiserverStore = &datastore{
			cli:  pb.NewCacheClient(conn),
			conn: conn,
		}
	}

	return apiserverStore, err
}
