package authzserver

import (
	"context"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
	"github.com/che-kwas/iam-kit/shutdown"
	"github.com/spf13/viper"

	"iam-authz/internal/authzserver/auditor"
	"iam-authz/internal/authzserver/cache"
	"iam-authz/internal/authzserver/store"
	"iam-authz/internal/authzserver/store/apiserver"
	"iam-authz/internal/authzserver/subscriber"
	"iam-authz/internal/authzserver/subscriber/redis"
)

type authServer struct {
	*server.Server
	name        string
	enableAudit bool
	ctx         context.Context
	cancel      context.CancelFunc
	log         *logger.Logger

	err error
}

// NewServer builds a new authServer.
func NewServer(name string) *authServer {
	ctx, cancel := context.WithCancel(context.Background())

	s := &authServer{
		name:   name,
		ctx:    ctx,
		cancel: cancel,
		log:    logger.L(),
	}

	return s.initStore().
		initCache().
		initAudit().
		newServer().
		setupHTTP()
}

// Run runs the authServer.
func (s *authServer) Run() {
	if s.err != nil {
		s.log.Fatal(s.err)
	}

	defer s.cancel()
	defer s.log.Sync()
	defer store.Client().Close()
	defer subscriber.Sub().Close()

	if err := s.Server.Run(); err != nil {
		s.log.Fatal(err)
	}
}

func (s *authServer) initStore() *authServer {
	var addr string
	if s.err = viper.UnmarshalKey("apiserver.addr", &addr); s.err != nil {
		return s
	}

	var cli store.Store
	opts := apiserver.NewAPIServerOptions()
	if cli, s.err = apiserver.NewAPIServerStore(opts); s.err != nil {
		return s
	}
	store.SetClient(cli)

	return s
}

func (s *authServer) initCache() *authServer {
	if s.err != nil {
		return s
	}

	var sub subscriber.Subscriber
	if sub, s.err = redis.NewRedisSub(); s.err != nil {
		return s
	}

	var loaderImpl cache.Loadable
	if loaderImpl, s.err = cache.InitCacheIns(); s.err != nil {
		return s
	}

	cache.NewLoader(s.ctx, sub, loaderImpl).Start()

	return s
}

func (s *authServer) initAudit() *authServer {
	if s.err != nil {
		return s
	}

	opts := auditor.NewAuditorOptions()
	s.enableAudit = opts.Enable

	if opts.Enable {
		auditor.InitAuditor(s.ctx, opts).Start()
	}

	return s
}

func (s *authServer) newServer() *authServer {
	if s.err != nil {
		return s
	}

	opts := []server.Option{}
	if s.enableAudit {
		sd := shutdown.ShutdownFunc(auditor.GetAuditor().Stop)
		opts = append(opts, server.WithShutdown(sd))
	}

	s.Server, s.err = server.NewServer(s.name, opts...)
	return s
}

func (s *authServer) setupHTTP() *authServer {
	if s.err != nil {
		return s
	}

	initRouter(s.Server.HTTPServer.Engine)
	return s
}
