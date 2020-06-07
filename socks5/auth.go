package socks5

import (
	"errors"
	"fmt"
	"io"
)

const (
	NoAuth          = uint8(0)
	UserPassAuth    = uint8(2)
	noAcceptable    = uint8(255)
	Socks5Version 	= uint8(5)

	userAuthVersion = uint8(1)

	authSuccess     = uint8(0)
	authFailure     = uint8(1)
)

type BytesGenerator func() []byte

type Authenticator interface {
	Authenticate(reader io.Reader, writer io.Writer) error
	GetCode() uint8
}

type NoAuthAuthenticator struct{}

func (a NoAuthAuthenticator) GetCode() uint8 {
	return NoAuth
}

func (a NoAuthAuthenticator) Authenticate(_ io.Reader, writer io.Writer) error {
	_, err := writer.Write([]byte{Socks5Version, NoAuth})
	return err
}

type UserPassAuthenticator struct {
	Accounts Accounts
}

func (a UserPassAuthenticator) GetCode() uint8 {
	return UserPassAuth
}

func (a UserPassAuthenticator) Authenticate(reader io.Reader, writer io.Writer) error {
	if _, err := writer.Write([]byte{Socks5Version, UserPassAuth}); err != nil {
		return err
	}
	req := &UserPassAuthRequest{}
	if err := req.Read(reader); err != nil {
		return err
	}
	if req.Ver != userAuthVersion {
		return fmt.Errorf("Unsupported auth version: %v ", req.Ver)
	}
	if a.Accounts.contains(string(req.Uname), string(req.Passwd)) {
		if _, err := writer.Write([]byte{userAuthVersion, authSuccess}); err != nil {
			return err
		}
		return nil
	}
	if _, err := writer.Write([]byte{userAuthVersion, authFailure}); err != nil {
		return err
	}
	return errors.New("User authentication failed ")
}

func NoAcceptableAuth(conn io.Writer) error {
	_, _ = conn.Write([]byte{Socks5Version, noAcceptable})
	return errors.New("No supported authentication mechanism ")
}

type MemoryUser map[string]string

type Accounts struct {
	MemoryUser MemoryUser
}

func (a *Accounts) contains(username, password string) bool {
	if pass, ok := a.MemoryUser[username]; !ok || pass != password {
		return false
	}
	return true
}