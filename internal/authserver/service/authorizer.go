// Package service defines the core business logic.
package service

import (
	"github.com/ory/ladon"

	v1 "iam-auth/api/authserver/v1"
)

// Authorizer implements the authorize interface.
type Authorizer struct {
	warden ladon.Warden
}

// NewAuthorizer creates a local repository authorizer.
func NewAuthorizer() *Authorizer {
	return &Authorizer{
		warden: &ladon.Ladon{
			Manager:     NewPolicyManager(),
			AuditLogger: NewAuditLogger(),
		},
	}
}

// Authorize authorizes the subject access.
func (a *Authorizer) Authorize(request *ladon.Request) *v1.Response {
	if err := a.warden.IsAllowed(request); err != nil {
		return &v1.Response{
			Denied: true,
			Reason: err.Error(),
		}
	}

	return &v1.Response{Allowed: true}
}
