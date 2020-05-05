package JadeSocks
//
//import (
//	"bytes"
//	"github.com/archervanderwaal/JadeSocks"
//	"testing"
//)
//
//func TestNoAuthAuthenticator(t *testing.T) {
//	req := bytes.NewBuffer(nil)
//	req.Write([]byte{1, NoAuth})
//	var resp bytes.Buffer
//
//	s, _ := JadeSocks.New(&JadeSocks.Config{
//		AuthMethods: []Authenticator{NoAuthAuthenticator{}},
//	})
//	ctx, err := s.authenticate(&resp, req)
//
//	if err != nil {
//		t.Error()
//	}
//
//	if ctx == nil || ctx.Method != NoAuth {
//		t.Error()
//	}
//
//	out := resp.Bytes()
//	if !bytes.Equal(out, []byte{JadeSocks.socks5Version, NoAuth}) {
//		t.Error()
//	}
//}
//
//func TestUserPassAuthenticator(t *testing.T) {
//	req := bytes.NewBuffer(nil)
//	req.Write([]byte{2, NoAuth, UserPassAuth})
//	req.Write([]byte{1, 3, 'f', 'o', 'o', 3, 'b', 'a', 'r'})
//	var resp bytes.Buffer
//
//	cred := StaticCredentials{
//		"foo": "bar",
//	}
//
//	authenticator := UserPassAuthenticator{Credentials: cred}
//
//	s, _ := JadeSocks.New(&JadeSocks.Config{AuthMethods: []Authenticator{authenticator}})
//
//	ctx, err := s.authenticate(&resp, req)
//	if err != nil {
//		t.Error()
//	}
//
//	if ctx == nil {
//		t.Error()
//		return
//	}
//
//	if ctx.Method != UserPassAuth {
//		t.Error()
//	}
//
//	val, ok := ctx.Payload["user"]
//	if !ok {
//		t.Error()
//	}
//
//	if val != "foo" {
//		t.Error()
//	}
//
//	out := resp.Bytes()
//	if !bytes.Equal(out, []byte{JadeSocks.socks5Version, UserPassAuth, 1, authSuccess}) {
//		t.Error()
//	}
//}
//
//func TestPasswordAuth_Invalid(t *testing.T) {
//	req := bytes.NewBuffer(nil)
//	req.Write([]byte{2, NoAuth, UserPassAuth})
//	req.Write([]byte{1, 3, 'f', 'o', 'o', 3, 'b', 'a', 'z'})
//	var resp bytes.Buffer
//
//	cred := StaticCredentials{
//		"foo": "bar",
//	}
//	cator := UserPassAuthenticator{Credentials: cred}
//	s, _ := JadeSocks.New(&JadeSocks.Config{AuthMethods: []Authenticator{cator}})
//
//	ctx, err := s.authenticate(&resp, req)
//	if err != UserAuthFailed {
//		t.Error()
//	}
//
//	if ctx != nil {
//		t.Error()
//	}
//
//	out := resp.Bytes()
//	if !bytes.Equal(out, []byte{JadeSocks.socks5Version, UserPassAuth, 1, authFailure}) {
//		t.Error()
//	}
//}
//
//func TestNoSupportedAuth(t *testing.T) {
//	req := bytes.NewBuffer(nil)
//	req.Write([]byte{1, NoAuth})
//	var resp bytes.Buffer
//
//	cred := StaticCredentials{
//		"foo": "bar",
//	}
//	authenticator := UserPassAuthenticator{Credentials: cred}
//
//	s, _ := JadeSocks.New(&JadeSocks.Config{AuthMethods: []Authenticator{authenticator}})
//
//	ctx, err := s.authenticate(&resp, req)
//	if err != NoSupportedAuth {
//		t.Error()
//	}
//
//	if ctx != nil {
//		t.Error()
//	}
//
//	out := resp.Bytes()
//	if !bytes.Equal(out, []byte{JadeSocks.socks5Version, noAcceptable}) {
//		t.Error()
//	}
//}
