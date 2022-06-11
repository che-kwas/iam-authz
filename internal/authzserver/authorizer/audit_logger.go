package authorizer

import (
	"encoding/json"
	"fmt"
	"iam-authz/internal/authzserver/auditor"
	"strings"
	"time"

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
	var conclusion string
	if len(d) > 1 {
		allowed := joinPoliciesNames(d[0 : len(d)-1])
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policies [%s] allowed access, but policy [%s] forcefully denied it", allowed, denied)
	} else if len(d) == 1 {
		denied := d[len(d)-1].GetID()
		conclusion = fmt.Sprintf("policy [%s] forcefully denied the access", denied)
	} else {
		conclusion = "no policy allowed access"
	}
	a.log.Debugw("access request rejected", "conclusion", conclusion)

	rstring, pstring, dstring := convertToString(r, p, d)
	record := auditor.AuditRecord{
		Timestamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.DenyAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}

	auditor.GetAuditor().RecordHit(&record)
}

// LogGrantedAccessRequest write granted subject access to log.
func (a *AuditLogger) LogGrantedAccessRequest(r *ladon.Request, p ladon.Policies, d ladon.Policies) {
	conclusion := fmt.Sprintf("policies [%s] allowed access", joinPoliciesNames(d))
	a.log.Debugw("access request granted", "conclusion", conclusion)

	rstring, pstring, dstring := convertToString(r, p, d)
	record := auditor.AuditRecord{
		Timestamp:  time.Now().Unix(),
		Username:   r.Context["username"].(string),
		Effect:     ladon.AllowAccess,
		Conclusion: conclusion,
		Request:    rstring,
		Policies:   pstring,
		Deciders:   dstring,
	}

	auditor.GetAuditor().RecordHit(&record)
}

func joinPoliciesNames(policies ladon.Policies) string {
	names := []string{}
	for _, policy := range policies {
		names = append(names, policy.GetID())
	}

	return strings.Join(names, ", ")
}

func convertToString(r *ladon.Request, p ladon.Policies, d ladon.Policies) (string, string, string) {
	rbytes, _ := json.Marshal(r)
	pbytes, _ := json.Marshal(p)
	dbytes, _ := json.Marshal(d)

	return string(rbytes), string(pbytes), string(dbytes)
}
