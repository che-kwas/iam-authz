package apiserver

import (
	"context"
	"encoding/json"

	"github.com/avast/retry-go"
	"github.com/che-kwas/iam-kit/logger"
	"github.com/marmotedu/errors"
	"github.com/ory/ladon"

	pb "iam-auth/api/apiserver/proto/v1"
)

type policies struct {
	cli pb.CacheClient
	log *logger.Logger
}

func newPolicies(ds *datastore) *policies {
	return &policies{cli: ds.cli, log: logger.L()}
}

// List returns all policies.
func (p *policies) List() (map[string][]*ladon.DefaultPolicy, error) {
	p.log.Info("list policies")

	req := &pb.ListPoliciesRequest{Offset: 0, Limit: -1}
	var resp *pb.ListPoliciesResponse
	err := retry.Do(
		func() error {
			var listErr error
			resp, listErr = p.cli.ListPolicies(context.Background(), req)
			if listErr != nil {
				return listErr
			}

			return nil
		}, retry.Attempts(3),
	)
	if err != nil {
		p.log.Errorf("failed to list policies: %s", err.Error())
		return nil, errors.Wrap(err, "failed to list policies")
	}

	total := len(resp.Items)
	p.log.Infof("policies found %d total", total)

	pols := make(map[string][]*ladon.DefaultPolicy)
	for _, v := range resp.Items {
		p.log.Infof(" - %s:%s", v.Username, v.Name)
		var policy ladon.DefaultPolicy

		if err := json.Unmarshal([]byte(v.PolicyShadow), &policy); err != nil {
			p.log.Warnf("failed to load policy for %s:%s, error: %s", v.Username, v.Name, err.Error())
			continue
		}

		pols[v.Username] = append(pols[v.Username], &policy)
	}

	return pols, nil
}
