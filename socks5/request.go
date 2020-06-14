package socks5

import (
	"fmt"
	"io"
	"net"
	"strings"
)

const (
	connectCommand   = uint8(1)
	bindCommand      = uint8(2)
	associateCommand = uint8(3)
)

const (
	// 0x00 succeeded
	succeeded uint8 = iota
	// 0x01 general SOCKS server failure
	serverFailure
	// 0x02 connection not allowed by ruleset
	ruleNotAllowed
	// 0x03 network unreachable
	networkUnreachable
	// 0x04 host unreachable
	hostUnreachable
	// 0x05 connection refused
	connectionRefused
	// 0x06 TTL expired
	ttlExpired
	// 0x07 command not supported
	commandNotSupported
	// 0x08 address type not supported
	addrTypeNotSupported
)

type Request struct {
	Version    uint8
	Command    uint8
	RemoteAddr *AddrSpec
	DestAddr   *AddrSpec
	reader     io.Reader
}

type AddrSpec struct {
	Domain   string
	AddrType byte
	IP       net.IP
	Port     uint16
}

func NewRequest(reader io.Reader) (*Request, error) {
	header := []byte{0, 0, 0}
	if _, err := io.ReadAtLeast(reader, header, 3); err != nil {
		return nil, fmt.Errorf("Failed to get command version: %v ", err)
	}
	if header[0] != Socks5Version {
		return nil, fmt.Errorf("Unsupported command version: %v ", header[0])
	}
	dest, err := parseAddrSpec(reader)
	if err != nil {
		return nil, err
	}
	return &Request{
		Version:  Socks5Version,
		Command:  header[1],
		DestAddr: dest,
		reader:   reader,
	}, nil
}

func (server *Server) process(req *Request, conn net.Conn) error {
	dest := req.DestAddr
	if dest.Domain != "" {
		addr, err := server.Config.Resolver.Resolve(dest.Domain)
		if err != nil {
			if err = sendResponse(conn, hostUnreachable, nil); err != nil {
				server.Config.Logger.Errorf("Failed to send response %v ", err)
				return fmt.Errorf("Failed to send response %v ", err)
			}
			server.Config.Logger.Errorf("Failed to resolve destination '%v': %v ", dest.Domain, err)
			return fmt.Errorf("Failed to resolve destination '%v': %v ", dest.Domain, err)
		}
		dest.IP = addr
	}
	switch req.Command {
	case connectCommand:
		return server.handleConnect(req, conn)
	case bindCommand:
		return server.handleBind(req, conn)
	case associateCommand:
		return server.handleAssociate(req, conn)
	default:
		if err := sendResponse(conn, commandNotSupported, nil); err != nil {
			server.Config.Logger.Errorf("Failed to send response: %v", err)
			return err
		}
		server.Config.Logger.Errorf("Unsupported command: %v ", req.Command)
		return fmt.Errorf("Unsupported command: %v ", req.Command)
	}
}

func (server *Server) handleConnect(req *Request, conn net.Conn) error {
	if err := server.checkRules(req, conn); err != nil {
		return err
	}
	dial := server.Config.Dial
	if dial == nil {
		dial = func(network string, addr AddrSpec) (net.Conn, error) {
			return net.DialTCP(network, nil, &net.TCPAddr{IP: addr.IP, Port:int(addr.Port)})
		}
	}
	target, err := dial("tcp", *req.DestAddr)
	if err != nil {
		msg := err.Error()
		resp := hostUnreachable
		if strings.Contains(msg, "refused") {
			resp = connectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			resp = networkUnreachable
		}
		if err := sendResponse(conn, resp, nil); err != nil {
			server.Config.Logger.Errorf("Failed to send response: %v ", err)
			return err
		}
		server.Config.Logger.Errorf("Connect to %v failed: %s", req.DestAddr, msg)
		return err
	}
	defer target.Close()

	local := target.LocalAddr().(*net.TCPAddr)
	bind := AddrSpec{
		IP:       local.IP,
		Port:     uint16(local.Port),
		AddrType: IPV4Address,
	}
	if err := sendResponse(conn, succeeded, &bind); err != nil {
		server.Config.Logger.Errorf("Failed to send response: %v ", err)
		return err
	}

	errCh := make(chan error, 2)
	go copyData(target, req.reader, errCh)
	go copyData(conn, target, errCh)

	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			return e
		}
	}
	return nil
}

func (server *Server) handleBind(req *Request, conn net.Conn) error {
	if err := server.checkRules(req, conn); err != nil {
		return err
	}
	if err := sendResponse(conn, commandNotSupported, nil); err != nil {
		server.Config.Logger.Errorf("Failed to send response: %v ", err)
		return err
	}
	server.Config.Logger.Infof("Bind command is temporarily not supported")
	return nil
}

func (server *Server) handleAssociate(req *Request, conn net.Conn) error {
	if err := server.checkRules(req, conn); err != nil {
		return err
	}
	//// Check if this is allowed
	//	_ctx, ok := s.config.Rules.Allow(ctx, req)
	//	if !ok {
	//		if err := sendReply(conn, ReplyRuleFailure, nil); err != nil {
	//			return fmt.Errorf("failed to send reply: %v", err)
	//		}
	//		return fmt.Errorf("associate to %v blocked by rules", req.DestAddr)
	//	}
	//	ctx = _ctx
	//
	//	// check bindIP 1st
	//	if len(s.config.BindIP) == 0 || s.config.BindIP.IsUnspecified() {
	//		s.config.BindIP = net.ParseIP("127.0.0.1")
	//	}
	//
	//	bindAddr := AddrSpec{IP: s.config.BindIP, Port: s.config.BindPort}
	//
	//	if err := sendReply(conn, ReplySucceeded, &bindAddr); err != nil {
	//		return fmt.Errorf("failed to send reply: %v", err)
	//	}
	//
	//	// wait here till the client close the connection
	//	// check every 10 secs
	//	tmp := []byte{}
	//	var neverTimeout time.Time
	//	for {
	//		conn.SetReadDeadline(time.Now())
	//		if _, err := conn.Read(tmp); err == io.EOF {
	//			break
	//		} else {
	//			conn.SetReadDeadline(neverTimeout)
	//		}
	//		time.Sleep(10 * time.Second)
	//	}
	//
	//	return nil
	if err := sendResponse(conn, commandNotSupported, nil); err != nil {
		server.Config.Logger.Errorf("Failed to send response: %v ", err)
		return err
	}
	server.Config.Logger.Infof("Associate command is temporarily not supported")
	return nil
}

func parseAddrSpec(r io.Reader) (*AddrSpec, error) {
	addrSpec := &AddrSpec{}

	addrType := []byte{0}
	if _, err := r.Read(addrType); err != nil {
		return nil, err
	}

	switch addrType[0] {
	case IPV4Address:
		addr := make([]byte, 4)
		if _, err := io.ReadAtLeast(r, addr, len(addr)); err != nil {
			return nil, err
		}
		addrSpec.IP = addr
		addrSpec.AddrType = IPV4Address

	case IPV6Address:
		addr := make([]byte, 16)
		if _, err := io.ReadAtLeast(r, addr, len(addr)); err != nil {
			return nil, err
		}
		addrSpec.IP = addr
		addrSpec.AddrType = IPV6Address

	case DomainAddress:
		if _, err := r.Read(addrType); err != nil {
			return nil, err
		}
		addrLen := int(addrType[0])
		domain := make([]byte, addrLen)
		if _, err := io.ReadAtLeast(r, domain, addrLen); err != nil {
			return nil, err
		}
		addrSpec.Domain = string(domain)
		addrSpec.AddrType = DomainAddress

	default:
		return nil, UnrecognizedAddrType
	}
	port := []byte{0, 0}
	if _, err := io.ReadAtLeast(r, port, 2); err != nil {
		return nil, err
	}
	addrSpec.Port = uint16((int(port[0]) << 8) | int(port[1]))

	return addrSpec, nil
}

func sendResponse(writer io.Writer, resp uint8, addr *AddrSpec) error {
	rep := &Response{}
	rep.Ver = Socks5Version
	rep.Rep = resp
	rep.Rsv = uint8(0)
	rep.Atyp = 0
	if addr == nil {
		rep.Atyp = IPV4Address
		rep.BndAddr = []byte{0, 0, 0, 0}
		rep.BndPort = []byte{0, 0}
		return nil
	}
	rep.Atyp = addr.AddrType
	rep.BndPort = []byte{byte(addr.Port >> 8), byte(addr.Port & 0xff)}
	switch addr.AddrType {
	case IPV4Address:
		rep.BndAddr = addr.IP.To4()
	case DomainAddress:
		rep.BndAddr = bytesCombine([]byte{byte(len(addr.Domain))}, []byte(addr.Domain))
	case IPV6Address:
		rep.BndAddr = addr.IP.To16()
	default:
		return fmt.Errorf("Failed to format address: %v ", addr)
	}
	_, err := writer.Write(bytesCombine([]byte{rep.Ver, rep.Rep, rep.Rsv, rep.Atyp}, rep.BndAddr, rep.BndPort))
	return err
}

type closeWriter interface {
	CloseWrite() error
}

// checkRules is used to check request command is allowed
func (server *Server) checkRules(req* Request, conn net.Conn) error {
	if ok := server.Config.Rules.Allow(req); !ok {
		if err := sendResponse(conn, ruleNotAllowed, nil); err != nil {
			server.Config.Logger.Errorf("Failed to send response: %v ", err)
			return err
		}
		return fmt.Errorf("command %d to %v blocked by rules ", req.Command, req.DestAddr)
	}
	return nil
}

func copyData(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	if tcpConn, ok := dst.(closeWriter); ok {
		_ = tcpConn.CloseWrite()
	}
	errCh <- err
}
