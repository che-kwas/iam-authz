package apiserver

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/che-kwas/iam-kit/logger"
	"github.com/marmotedu/errors"

	pb "iam-auth/api/apiserver/proto/v1"
)

type secrets struct {
	cli pb.CacheClient
	log *logger.Logger
}

func newSecrets(ds *datastore) *secrets {
	return &secrets{cli: ds.cli, log: logger.L()}
}

// List returns all secrets.
func (s *secrets) List() (map[string]*pb.SecretInfo, error) {
	s.log.Info("loading secrets")

	req := &pb.ListSecretsRequest{Offset: 0, Limit: -1}
	var resp *pb.ListSecretsResponse
	err := retry.Do(
		func() error {
			var listErr error
			if resp, listErr = s.cli.ListSecrets(context.Background(), req); listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		return nil, errors.Wrap(err, "list secrets failed")
	}

	total := len(resp.Items)
	s.log.Infof("secrets found %d total:", total)

	secrets := make(map[string]*pb.SecretInfo, total)
	for _, v := range resp.Items {
		s.log.Infof(" - %s:%s", v.Username, v.SecretId)
		secrets[v.SecretId] = v
	}

	return secrets, nil
}
