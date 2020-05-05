package JadeSocks

import (
	"errors"
	"fmt"
	"io"
)

const (
	NoAuth          = uint8(0)
	UserPassAuth    = uint8(2)
	noAcceptable    = uint8(255)

	userAuthVersion = uint8(1)

	authSuccess     = uint8(0)
	authFailure     = uint8(1)
)

var (
	UserAuthFailed  = errors.New("User authentication failed ")
	NoSupportedAuth = errors.New("No supported authentication mechanism ")
)

type AuthContext struct {
	Method uint8
	Payload map[string]string
}

type Authenticator interface {
	Authenticate(reader io.Reader, writer io.Writer) (*AuthContext, error)
	GetCode() uint8
}

type NoAuthAuthenticator struct{}

func (a NoAuthAuthenticator) GetCode() uint8 {
	return NoAuth
}

func (a NoAuthAuthenticator) Authenticate(_ io.Reader, writer io.Writer) (*AuthContext, error) {
	_, err := writer.Write([]byte{Socks5Version, NoAuth})
	return &AuthContext{NoAuth, nil}, err
}

type UserPassAuthenticator struct {
	Credentials Credential
}

func (a UserPassAuthenticator) GetCode() uint8 {
	return UserPassAuth
}

func (a UserPassAuthenticator) Authenticate(reader io.Reader, writer io.Writer) (*AuthContext, error) {
	if _, err := writer.Write([]byte{Socks5Version, UserPassAuth}); err != nil {
		return nil, err
	}

	request := &UserPassNegotiationRequest{}
	if err := request.ReadFrom(reader); err != nil {
		return nil, err
	}

	if request.Ver != userAuthVersion {
		return nil, fmt.Errorf("Unsupported auth version: %v ", request.Ver)
	}

	if a.Credentials.Valid(string(request.Uname), string(request.Passwd)) {
		if _, err := writer.Write([]byte{userAuthVersion, authSuccess}); err != nil {
			return nil, err
		}
	} else {
		if _, err := writer.Write([]byte{userAuthVersion, authFailure}); err != nil {
			return nil, err
		}
		return nil, UserAuthFailed
	}

	return &AuthContext{UserPassAuth, map[string]string{"user": string(request.Uname)}}, nil
}

func NoAcceptableAuth(conn io.Writer) error {
	_, _ = conn.Write([]byte{Socks5Version, noAcceptable})
	return NoSupportedAuth
}