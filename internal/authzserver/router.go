package authzserver

import (
	"github.com/che-kwas/iam-kit/code"
	"github.com/che-kwas/iam-kit/httputil"
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/errors"

	"iam-authz/internal/authzserver/controller"
)

func initRouter(g *gin.Engine) {
	auth := newJWTExAuth()
	g.NoRoute(auth.AuthFunc(), notFound())

	v1 := g.Group("/v1")
	{
		authController := controller.NewAuthController()
		v1.POST("/authz", authController.Authorize)
	}
}

func notFound() func(c *gin.Context) {
	return func(c *gin.Context) {
		httputil.WriteResponse(c, errors.WithCode(code.ErrNotFound, "Not found."), nil)
	}
}
