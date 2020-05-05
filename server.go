package JadeSocks

import (
	"context"
	"errors"
	"github.com/phachon/go-logger"
	"net"
)

type Config struct {
	AuthMethods []Authenticator
	Resolver    NameResolver
	Rules       RuleSet
	BindIP      net.IP
	Logger      *go_logger.Logger
	Dial        func(ctx context.Context, network, addr string) (net.Conn, error)
}

type Server struct {
	Config *Config
}

func New(conf *Config) (*Server, error) {
	if len(conf.AuthMethods) == 0 {
		return nil, errors.New("Ensure we have at least one authentication method enabled ")
	}

	if conf.Resolver == nil {
		conf.Resolver = DNSResolver{}
	}

	if conf.Rules == nil {
		conf.Rules = PermitAll()
	}

	server := &Server{
		Config: conf,
	}
	return server, nil
}

func (s *Server) ListenAndServe(network, addr string) error {
	l, err := net.Listen(network, addr)
	if err != nil {
		s.Config.Logger.Errorf("Failed listen to %s:%s %v", network, addr, err)
		return err
	}
	s.Config.Logger.Infof("Successfully listen to %s:%s", network, addr)
	return s.serve(l)
}

func (s *Server) serve(l net.Listener) error {
	for {
		conn, err := l.Accept()
		s.Config.Logger.Infof("TCP connection established successfully, %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
		if err != nil {
			s.Config.Logger.Errorf("TCP connection established failed, %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
			return err
		}
		go func() {
			err = s.serveConn(conn)
			if err != nil {
				s.Config.Logger.Errorf("Processing connection failure %v", err)
			}
		}()
	}
}

func (s *Server) serveConn(conn net.Conn) error {
	return (&Proxy{
		Server: s,
	}).HandleConn(conn)
}
