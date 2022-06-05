package authserver

import (
	"context"

	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
	"github.com/spf13/viper"

	"iam-auth/internal/authserver/cache"
	"iam-auth/internal/authserver/store"
	"iam-auth/internal/authserver/store/apiserver"
)

type authServer struct {
	*server.Server
	name   string
	ctx    context.Context
	cancel context.CancelFunc
	log    *logger.Logger

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

	return s.initStore().initCache().newServer().setupHTTP()
}

// Run runs the authServer.
func (s *authServer) Run() {
	defer s.log.Sync()
	defer store.Client().Close()
	defer s.cancel()

	if s.err != nil {
		s.log.Fatal("failed to build the server: ", s.err)
	}

	if err := s.Server.Run(); err != nil {
		s.log.Fatal("server stopped unexpectedly: ", err)
	}
}

func (s *authServer) initStore() *authServer {
	var addr string
	if s.err = viper.UnmarshalKey("main.apiserver", &addr); s.err != nil {
		return s
	}

	var storeIns store.Store
	if storeIns, s.err = apiserver.APIServerStore(addr); s.err != nil {
		return s
	}
	store.SetClient(storeIns)

	return s
}

func (s *authServer) initCache() *authServer {
	if s.err != nil {
		return s
	}

	var cacheIns cache.Loadable
	if cacheIns, s.err = cache.CacheIns(); s.err != nil {
		return s
	}
	cache.NewLoader(s.ctx, cacheIns).Start()

	return s
}

func (s *authServer) newServer() *authServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(s.name)
	return s
}

func (s *authServer) setupHTTP() *authServer {
	if s.err != nil {
		return s
	}

	initRouter(s.Server.HTTPServer.Engine)
	return s
}
