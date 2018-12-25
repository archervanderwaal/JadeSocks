package JadeSocks

import (
	"golang.org/x/net/context"
	"net"
)

// NameResolver is used to implement custom name resolution.
// It contains the Resolve method for name resolution to IP.
type NameResolver interface {
	Resolve(ctx context.Context, name string) (context.Context, net.IP, error)
}

// DNSResolver uses the System DNS to resolve host names.
type DNSResolver struct {}

func (resolver DNSResolver) Resolve(ctx context.Context, name string) (context.Context, net.IP, error) {
	addr, err := net.ResolveIPAddr("ip", name)
	if err != nil {
		return ctx, nil, err
	}
	return ctx, addr.IP, err
}

