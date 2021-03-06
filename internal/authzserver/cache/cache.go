package cache

import (
	"sync"
	"unsafe"

	"github.com/dgraph-io/ristretto"
	"github.com/marmotedu/errors"
	"github.com/ory/ladon"

	pb "iam-authz/api/apiserver/proto/v1"
	"iam-authz/internal/authzserver/store"
)

// Cache stores the secrets and policies.
type Cache struct {
	lock     *sync.RWMutex
	secrets  *ristretto.Cache
	policies *ristretto.Cache
}

var _ Loadable = &Cache{}

const (
	secretSize = unsafe.Sizeof(pb.SecretInfo{})
	policySize = unsafe.Sizeof(pb.PolicyInfo{})
)

var (
	// ErrSecretNotFound defines secret not found error.
	ErrSecretNotFound = errors.New("secret not found")
	// ErrPolicyNotFound defines policy not found error.
	ErrPolicyNotFound = errors.New("policy not found")
)

var cacheIns *Cache

// InitCacheIns initilizes the cache instance and returns it.
func InitCacheIns() (*Cache, error) {
	var (
		err         error
		secretCache *ristretto.Cache
		policyCache *ristretto.Cache
	)

	c := &ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Cost:        nil,
	}

	if secretCache, err = ristretto.NewCache(c); err != nil {
		return nil, err
	}
	if policyCache, err = ristretto.NewCache(c); err != nil {
		return nil, err
	}

	cacheIns = &Cache{
		lock:     new(sync.RWMutex),
		secrets:  secretCache,
		policies: policyCache,
	}

	return cacheIns, err
}

// CacheIns returns the global cache instance.
func CacheIns() *Cache {
	return cacheIns
}

// GetSecret returns secret detail for the given key.
func (c *Cache) GetSecret(key string) (*pb.SecretInfo, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.secrets.Get(key)
	if !ok {
		return nil, ErrSecretNotFound
	}

	return value.(*pb.SecretInfo), nil
}

// ListPolicies returns user's ladon policies for the given user.
func (c *Cache) ListPolicies(username string) ([]*ladon.DefaultPolicy, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	value, ok := c.policies.Get(username)
	if !ok {
		return nil, ErrPolicyNotFound
	}

	return value.([]*ladon.DefaultPolicy), nil
}

// ReloadSecrets reloads secrets.
func (c *Cache) ReloadSecrets() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	secrets, err := store.Client().Secrets().List()
	if err != nil {
		return err
	}

	c.secrets.Clear()
	for key, val := range secrets {
		c.secrets.Set(key, val, int64(secretSize))
	}

	return nil
}

// ReloadPolicies reloads policies.
func (c *Cache) ReloadPolicies() error {
	c.lock.Lock()
	defer c.lock.Unlock()

	policies, err := store.Client().Policies().List()
	if err != nil {
		return err
	}

	c.policies.Clear()
	for key, val := range policies {
		c.policies.Set(key, val, int64(policySize))
	}

	return nil
}
