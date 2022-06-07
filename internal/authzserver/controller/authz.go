// Package controller defines the controllers.
package controller

import (
	basecode "github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/che-kwas/iam-kit/logger"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"
	"github.com/ory/ladon"

	"iam-authz/internal/authzserver/service"
)

// AuthController handles requests for authorization.
type AuthController struct {
	srv service.Authorizer
	log *logger.Logger
}

// NewAuthController creates a auth controller.
func NewAuthController() *AuthController {
	return &AuthController{
		srv: *service.NewAuthorizer(),
		log: logger.L(),
	}
}

// Authorize returns whether a request is allow or deny to access a resource
// and do some action under specified condition.
func (a *AuthController) Authorize(c *gin.Context) {
	var req ladon.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.WriteResponse(c, errors.WithCode(basecode.ErrBadParams, err.Error()), nil)
		return
	}
	a.log.X(c).Infow("authorize params", "request", req)

	if req.Context == nil {
		req.Context = ladon.Context{}
	}
	req.Context["username"] = c.GetString("username")
	resp := a.srv.Authorize(&req)

	httputil.WriteResponse(c, nil, resp)
}
