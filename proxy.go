package JadeSocks

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type Proxy struct {
	Server *Server
}

func (proxy *Proxy) HandleConn(conn net.Conn) error {
	defer conn.Close()
	bufConn := bufio.NewReader(conn)

	negotiationRequest := &NegotiationRequest{}
	err := negotiationRequest.ReadFrom(conn)
	if err != nil {
		proxy.Server.Config.Logger.Errorf("Failed to parse socks5 negotiation request: %v", err)
		return err
	}

	if negotiationRequest.Ver != Socks5Version {
		err := fmt.Errorf("Unsupported SOCKS version: %v ", negotiationRequest.Ver)
		proxy.Server.Config.Logger.Errorf("%v", err)
		return err
	}

	authContext, err := proxy.authenticate(conn, bufConn, negotiationRequest)
	if err != nil {
		err = fmt.Errorf("Failed to authenticate: %v ", err)
		proxy.Server.Config.Logger.Errorf("socks: %v", err)
		return err
	}
	proxy.Server.Config.Logger.Infof("Authentication success, use method %d", authContext.Method)

	request, err := NewRequest(bufConn)
	if err != nil {
		if err == unrecognizedAddrType {
			if err := SendSocks5Reply(conn, addrTypeNotSupported, nil); err != nil {
				return fmt.Errorf("Failed to send reply: %v ", err)
			}
		}
		return fmt.Errorf("Failed to read destination address: %v ", err)
	}
	request.AuthContext = authContext
	if client, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
		request.RemoteAddr = &AddrSpec{IP: client.IP, Port: uint16(client.Port), AddrType: IPV4Address}
	}

	if err := request.handleRequest(proxy, conn); err != nil {
		err = fmt.Errorf("Failed to handle request: %v ", err)
		proxy.Server.Config.Logger.Errorf("socks: %v", err)
		return err
	}
	return nil
}

func (proxy *Proxy) authenticate(conn io.Writer, reader io.Reader, request *NegotiationRequest) (*AuthContext, error) {
	for _, method := range request.Methods {
		for _, authenticator := range proxy.Server.Config.AuthMethods {
			if authenticator.GetCode() == method {
				return authenticator.Authenticate(reader, conn)
			}
		}
	}
	return nil, NoAcceptableAuth(conn)
}