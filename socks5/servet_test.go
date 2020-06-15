package socks5

import (
	"testing"
)

func TestServer_WithUserPassAuthenticator(t *testing.T) {
	accounts := Accounts{MemoryUser: map[string]string{"root": "123456"}}
	serverConf := &ServerConfig{
		AuthMethods: []Authenticator{UserPassAuthenticator{Accounts: accounts}},
		ListenAddr:  ":7890",
	}
	socks5Server, err := New(serverConf)
	if err != nil {
		t.Fail()
		return
	}
	if err = socks5Server.ListenAndServe(); err != nil {
		t.Fail()
		return
	}
}

func TestServer_WithNoAuthAuthenticator(t *testing.T) {
	serverConf := &ServerConfig{
		AuthMethods: []Authenticator{NoAuthAuthenticator{}},
		ListenAddr:  ":7890",
	}
	socks5Server, err := New(serverConf)
	if err != nil {
		t.Fail()
		return
	}
	if err = socks5Server.ListenAndServe(); err != nil {
		t.Fail()
		return
	}
}