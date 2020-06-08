package socks5

import (
	"net"
)

type NameResolver interface {
	Resolve(name string) (net.IP, error)
}

type DNSResolver struct{}

func (d DNSResolver) Resolve(name string) (net.IP, error) {
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return nil, err
	}
	return addr.IP, err
}
