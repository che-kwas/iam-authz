package authorizer

import (
	"github.com/ory/ladon"

	"iam-authz/internal/authzserver/cache"
)

// PolicyManager implements the ladon.Manager interface.
type PolicyManager struct {
	cache *cache.Cache
}

var _ ladon.Manager = &PolicyManager{}

// NewPolicyManager creates a ladon.Manager with the cache storage.
func NewPolicyManager() ladon.Manager {
	cacheIns, _ := cache.CacheIns()
	return &PolicyManager{cache: cacheIns}
}

// Create persists the policy.
// Does nothing because we use apiserver to manage the policies.
func (*PolicyManager) Create(policy ladon.Policy) error {
	return nil
}

// Update updates an existing policy.
// Does nothing because we use apiserver to manage the policies.
func (*PolicyManager) Update(policy ladon.Policy) error {
	return nil
}

// Get retrieves a policy.
// Does nothing because we use apiserver to manage the policies.
func (*PolicyManager) Get(id string) (ladon.Policy, error) {
	return nil, nil
}

// Delete removes a policy.
// Does nothing because we use apiserver to manage the policies.
func (*PolicyManager) Delete(id string) error {
	return nil
}

// GetAll retrieves all policies.
// Does nothing because we use apiserver to manage the policies.
func (*PolicyManager) GetAll(limit, offset int64) (ladon.Policies, error) {
	return nil, nil
}

// FindRequestCandidates returns candidates that could match the request object. It either returns
// a set that exactly matches the request, or a superset of it. If an error occurs, it returns nil and
// the error.
func (m *PolicyManager) FindRequestCandidates(r *ladon.Request) (ladon.Policies, error) {
	username := ""

	if user, ok := r.Context["username"].(string); ok {
		username = user
	}

	policies, err := m.cache.ListPolicies(username)
	if err != nil {
		return nil, err
	}

	ret := make([]ladon.Policy, 0, len(policies))
	for _, policy := range policies {
		ret = append(ret, policy)
	}

	return ret, nil
}

// FindPoliciesForSubject returns policies that could match the subject. It either returns
// a set of policies that applies to the subject, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *PolicyManager) FindPoliciesForSubject(subject string) (ladon.Policies, error) {
	return nil, nil
}

// FindPoliciesForResource returns policies that could match the resource. It either returns
// a set of policies that apply to the resource, or a superset of it.
// If an error occurs, it returns nil and the error.
func (m *PolicyManager) FindPoliciesForResource(resource string) (ladon.Policies, error) {
	return nil, nil
}
