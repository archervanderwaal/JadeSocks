package JadeSocks

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

const (
	ConnectCommand = uint8(1)
	BindCommand = uint8(2)
	AssociateCommand = uint8(3)
)

const (
	succeeded uint8 = iota
	serverFailure
	notAllowRuleset
	networkUnreachable
	hostUnreachable
	connectionRefused
	ttlExpired
	commandNotSupported
	addrTypeNotSupported
)

var (
	unrecognizedAddrType = fmt.Errorf("Unrecognized address type ")
)

type AddrSpec struct {
	Domain   string
	AddrType byte
	IP       net.IP
	Port     uint16
}

func (a *AddrSpec) String() string {
	if a.Domain != "" {
		return fmt.Sprintf("%s (%s):%d", a.Domain, a.IP, a.Port)
	}
	return fmt.Sprintf("%s:%d", a.IP, a.Port)
}

func (a AddrSpec) Address() string {
	if 0 != len(a.IP) {
		return net.JoinHostPort(a.IP.String(), strconv.Itoa(int(a.Port)))
	}
	return net.JoinHostPort(a.Domain, strconv.Itoa(int(a.Port)))
}

type Request struct {
	Version uint8
	Command uint8
	AuthContext *AuthContext
	RemoteAddr *AddrSpec
	DestAddr *AddrSpec
	bufConn      io.Reader
}

type conn interface {
	Write([]byte) (int, error)
	RemoteAddr() net.Addr
}

func NewRequest(bufConn io.Reader) (*Request, error) {
	header := []byte{0, 0, 0}
	if _, err := io.ReadAtLeast(bufConn, header, 3); err != nil {
		return nil, fmt.Errorf("Failed to get command version: %v ", err)
	}

	if header[0] != Socks5Version {
		return nil, fmt.Errorf("Unsupported command version: %v ", header[0])
	}

	dest, err := parseAddrSpec(bufConn)
	if err != nil {
		return nil, err
	}

	request := &Request{
		Version:  Socks5Version,
		Command:  header[1],
		DestAddr: dest,
		bufConn:  bufConn,
	}

	return request, nil
}

func (req *Request) handleRequest(proxy *Proxy, conn conn) error {
	ctx := context.Background()

	dest := req.DestAddr
	if dest.Domain != "" {
		ctx_, addr, err := proxy.Server.Config.Resolver.Resolve(ctx, dest.Domain)
		if err != nil {
			if err := SendSocks5Reply(conn, hostUnreachable, nil); err != nil {
				return fmt.Errorf("Failed to send reply: %v ", err)
			}
			return fmt.Errorf("Failed to resolve destination '%v': %v ", dest.Domain, err)
		}
		ctx = ctx_
		dest.IP = addr
	}

	switch req.Command {
	case ConnectCommand:
		return req.handleConnect(ctx, conn, proxy)
	case BindCommand:
		return req.handleBind(ctx, conn, proxy)
	case AssociateCommand:
		return req.handleAssociate(ctx, conn, proxy)
	default:
		if err := SendSocks5Reply(conn, commandNotSupported, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v ", err)
		}
		return fmt.Errorf("Unsupported command: %v ", req.Command)
	}
}

func (req *Request) handleConnect(ctx context.Context, conn conn, proxy *Proxy) error {
	if ctx_, ok := proxy.Server.Config.Rules.Allow(ctx, req); !ok {
		if err := SendSocks5Reply(conn, notAllowRuleset, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v ", err)
		}
		return fmt.Errorf("Connect to %v blocked by rules ", req.DestAddr)
	} else {
		ctx = ctx_
	}

	dial := proxy.Server.Config.Dial
	if dial == nil {
		dial = func(ctx context.Context, net_, addr string) (net.Conn, error) {
			return net.Dial(net_, addr)
		}
	}
	target, err := dial(ctx, "tcp", req.DestAddr.Address())
	if err != nil {
		msg := err.Error()
		resp := hostUnreachable
		if strings.Contains(msg, "refused") {
			resp = connectionRefused
		} else if strings.Contains(msg, "network is unreachable") {
			resp = networkUnreachable
		}
		if err := SendSocks5Reply(conn, resp, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v ", err)
		}
		return fmt.Errorf("Connect to %v failed: %v ", req.DestAddr, err)
	}
	defer target.Close()

	local := target.LocalAddr().(*net.TCPAddr)
	bind := AddrSpec{IP: local.IP, Port: uint16(local.Port), AddrType: IPV4Address}
	if err := SendSocks5Reply(conn, succeeded, &bind); err != nil {
		return fmt.Errorf("Failed to send reply: %v ", err)
	}

	errCh := make(chan error, 2)
	go copyData(target, req.bufConn, errCh)
	go copyData(conn, target, errCh)

	for i := 0; i < 2; i++ {
		e := <-errCh
		if e != nil {
			return e
		}
	}
	return nil
}

func (req *Request) handleBind(ctx context.Context, conn conn, proxy *Proxy) error {
	if ctx_, ok := proxy.Server.Config.Rules.Allow(ctx, req); !ok {
		if err := SendSocks5Reply(conn, notAllowRuleset, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v ", err)
		}
		return fmt.Errorf("Bind to %v blocked by rules ", req.DestAddr)
	} else {
		ctx = ctx_
	}

	// TODO: Support bind
	if err := SendSocks5Reply(conn, commandNotSupported, nil); err != nil {
		return fmt.Errorf("Failed to send reply: %v ", err)
	}
	return nil
}

func (req *Request) handleAssociate(ctx context.Context, conn conn, proxy *Proxy) error {
	if ctx_, ok := proxy.Server.Config.Rules.Allow(ctx, req); !ok {
		if err := SendSocks5Reply(conn, notAllowRuleset, nil); err != nil {
			return fmt.Errorf("Failed to send reply: %v ", err)
		}
		return fmt.Errorf("Associate to %v blocked by rules ", req.DestAddr)
	} else {
		ctx = ctx_
	}

	// TODO: Support associate
	if err := SendSocks5Reply(conn, commandNotSupported, nil); err != nil {
		return fmt.Errorf("Failed to send reply: %v ", err)
	}
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
		return nil, unrecognizedAddrType
	}
	port := []byte{0, 0}
	if _, err := io.ReadAtLeast(r, port, 2); err != nil {
		return nil, err
	}
	addrSpec.Port = uint16((int(port[0]) << 8) | int(port[1]))

	return addrSpec, nil
}

type closeWriter interface {
	CloseWrite() error
}

func copyData(dst io.Writer, src io.Reader, errCh chan error) {
	_, err := io.Copy(dst, src)
	if tcpConn, ok := dst.(closeWriter); ok {
		_ = tcpConn.CloseWrite()
	}
	errCh <- err
}
