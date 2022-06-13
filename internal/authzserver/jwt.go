package authzserver

import (
	"github.com/che-kwas/iam-kit/middleware"
	"github.com/che-kwas/iam-kit/middleware/auth"

	"iam-authz/internal/authzserver/cache"
)

func newJWTExAuth() middleware.AuthStrategy {
	return auth.NewJWTExStrategy(secretGetter)
}

func secretGetter(kid string) (auth.Secret, error) {
	secret, err := cache.CacheIns().GetSecret(kid)
	if err != nil {
		return auth.Secret{}, err
	}

	return auth.Secret{
		Username: secret.Username,
		ID:       secret.SecretId,
		Key:      secret.SecretKey,
		Expires:  secret.Expires,
	}, nil
}
