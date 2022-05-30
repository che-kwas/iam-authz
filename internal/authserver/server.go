package authserver

import (
	"github.com/che-kwas/iam-kit/logger"
	"github.com/che-kwas/iam-kit/server"
)

type authServer struct {
	*server.Server
	name string
	log  *logger.Logger

	err error
}

// NewServer builds a new authServer.
func NewServer(name string) *authServer {
	s := &authServer{
		name: name,
		log:  logger.L(),
	}

	return s.initStore().newServer().registerRouter()
}

// Run runs the authServer.
func (s *authServer) Run() {
	s.log.Sync()

	if s.err != nil {
		s.log.Fatal("failed to build the server: ", s.err)
	}

	if err := s.Server.Run(); err != nil {
		s.log.Fatal("server stopped unexpectedly: ", err)
	}
}

func (s *authServer) initStore() *authServer {
	return s
}

func (s *authServer) newServer() *authServer {
	if s.err != nil {
		return s
	}

	s.Server, s.err = server.NewServer(s.name)
	return s
}

func (s *authServer) registerRouter() *authServer {
	if s.err != nil {
		return s
	}

	return s
}
