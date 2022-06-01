// Package apiserver is the proxy of iam-apiserver
// which implements `iam-auth/internal/authserver/store.Store` interface.
package apiserver

import (
	"sync"

	pb "iam-auth/api/apiserver/proto/v1"

	"github.com/marmotedu/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"iam-auth/internal/authserver/store"
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
func APIServerStore(address string) (store.Store, error) {
	if apiserverStore != nil {
		return apiserverStore, nil
	}

	var err error
	once.Do(func() {
		var conn *grpc.ClientConn
		conn, err = grpc.Dial(address, grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			apiserverStore = &datastore{
				cli:  pb.NewCacheClient(conn),
				conn: conn,
			}
		}
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to apiserver")
	}

	return apiserverStore, nil
}
