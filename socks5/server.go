package socks5

import (
	"bufio"
	"errors"
	"fmt"
	go_logger "github.com/phachon/go-logger"
	"io"
	"net"
)

type ServerConfig struct {
	AuthMethods []Authenticator
	Resolver    NameResolver
	Network 	string
	ListenAddr  string
	Logger      *go_logger.Logger
	Dial        func(network, addr string) (net.Conn, error)
}

type Server struct {
	Config *ServerConfig
}

func New(conf *ServerConfig) (*Server, error) {
	if len(conf.AuthMethods) == 0 {
		return nil, errors.New("Ensure we have at least one authentication method enabled ")
	}
	if conf.Resolver == nil {
		conf.Resolver = DNSResolver{}
	}
	server := &Server{
		Config: conf,
	}
	return server, nil
}

func (server *Server) ListenAndServe() error {
	listener, err := net.Listen(server.Config.Network, server.Config.ListenAddr)
	if err != nil {
		server.Config.Logger.Errorf("Failed listen to %s:%s %v", server.Config.Network, server.Config.ListenAddr, err)
		return err
	}
	server.Config.Logger.Infof("Successfully listen to %s:%s", server.Config.Network, server.Config.ListenAddr)
	return server.serve(listener)
}

func (s *Server) serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		s.Config.Logger.Infof("TCP connection established successfully, %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
		if err != nil {
			s.Config.Logger.Errorf("TCP connection established failed, %s -> %s", conn.RemoteAddr(), conn.LocalAddr())
			return err
		}
		go func() {
			// 2. 处理连接
			_ = s.handleConn(conn)
		}()
	}
}

func (server *Server) handleConn(conn net.Conn) error {
	server.Config.Logger.Infof("Start handle connection, remoteAddr: %s:%s", conn.RemoteAddr().Network(), conn.RemoteAddr().String())
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	negotiationRequest := &NegotiationRequest{}
	err := negotiationRequest.Read(conn)
	if err != nil {
		server.Config.Logger.Errorf("Failed to parse socks5 negotiation request: %v", err)
		return err
	}

	if negotiationRequest.Ver != Socks5Version {
		err := fmt.Errorf("Unsupported SOCKS version: %v ", negotiationRequest.Ver)
		server.Config.Logger.Errorf("%v", err)
		return err
	}

	if err := server.authenticate(conn, bufConn, negotiationRequest); err != nil {
		return err
	}

	request, err := NewRequest(bufConn)
	if err != nil {
		if err == UnrecognizedAddrType {
			if err = sendResponse(conn, addrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("Failed to send response: %v ", err)
			}
		}
		return fmt.Errorf("Failed to read destination address: %v ", err)
	}

	if client, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: uint16(client.Port)}
	}

	// process client request
	if err := server.process(request, conn); err != nil {
		err = fmt.Errorf("Failed to handle request: %v ", err)
		server.Config.Logger.Errorf("Failed to handle request: %v ", err)
		return err
	}
	return nil
}

func (s *Server) authenticate(conn io.Writer, reader io.Reader, request *NegotiationRequest) error {
	for _, method := range request.Methods {
		for _, authenticator := range s.Config.AuthMethods {
			if authenticator.GetCode() == method {
				if err := authenticator.Authenticate(reader, conn); err != nil {
					s.Config.Logger.Errorf("Use the %d method of authentication failed", authenticator.GetCode())
					return err
				}
				s.Config.Logger.Infof("Use the %d method of authentication success", authenticator.GetCode())
				return nil
			}
		}
	}
	return NoAcceptableAuth(conn)
}