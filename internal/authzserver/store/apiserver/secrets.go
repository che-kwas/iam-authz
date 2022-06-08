package apiserver

import (
	"context"

	"github.com/avast/retry-go"
	"github.com/che-kwas/iam-kit/logger"
	"github.com/marmotedu/errors"

	pb "iam-authz/api/apiserver/proto/v1"
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
		s.log.Errorf("failed to list secrets: %s", err.Error())
		return nil, errors.Wrap(err, "failed to list secrets")
	}

	total := len(resp.Items)
	s.log.Debugf("secrets found %d total", total)

	secrets := make(map[string]*pb.SecretInfo, total)
	for _, v := range resp.Items {
		s.log.Debugf(" - %s:%s", v.Username, v.SecretId)
		secrets[v.SecretId] = v
	}

	return secrets, nil
}
