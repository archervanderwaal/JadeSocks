package JadeSocks

import (
	"testing"
	"golang.org/x/net/context"
	"fmt"
)

func TestDNSResolver_Resolve(t *testing.T) {
	resolver := DNSResolver{}
	ctx := context.Background()
	_, addr, err := resolver.Resolve(ctx, "archervanderwaal.com")
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	fmt.Println(addr)
}
