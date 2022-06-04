package service

import (
	"github.com/che-kwas/iam-kit/logger"
	"github.com/ory/ladon"
)

// AuditLogger tracks denied and granted authorizations.
type AuditLogger struct {
	log *logger.Logger
}

var _ ladon.AuditLogger = &AuditLogger{}

// NewAuditLogger creates a AuditLogger.
func NewAuditLogger() *AuditLogger {
	return &AuditLogger{
		log: logger.L(),
	}
}

// LogRejectedAccessRequest write rejected subject access to log.
func (a *AuditLogger) LogRejectedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	a.log.Debug("TODO")
}

// LogGrantedAccessRequest write granted subject access to log.
func (a *AuditLogger) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	a.log.Debug("TODO")
}
