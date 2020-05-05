package JadeSocks

import (
	"bytes"
	"fmt"
	"io"
)

const (
	Socks5Version = uint8(5)
	IPV4Address   = uint8(1)
	DomainAddress = uint8(3)
	IPV6Address   = uint8(4)
)

func SendSocks5Reply(w io.Writer, resp uint8, addr *AddrSpec) error {
	reply := &Socks5Reply{
		Ver:     Socks5Version,
		Rep:     resp,
		Rsv:     uint8(0),
		Atyp:    0,
		BndAddr: nil,
		BndPort: nil,
	}

	if addr == nil {
		reply.Atyp = IPV4Address
		reply.BndAddr = []byte{0, 0, 0, 0}
		reply.BndPort = []byte{0, 0}
		return reply.WriteTo(w)
	}
	reply.Atyp = addr.AddrType
	reply.BndPort = []byte{byte(addr.Port >> 8), byte(addr.Port & 0xff)}
	switch addr.AddrType {
	case IPV4Address:
		reply.BndAddr = addr.IP.To4()
	case DomainAddress:
		reply.BndAddr = bytesCombine([]byte{byte(len(addr.Domain))}, []byte(addr.Domain))
	case IPV6Address:
		reply.BndAddr = addr.IP.To16()
	default:
		return fmt.Errorf("Failed to format address: %v ", addr)
	}
	return reply.WriteTo(w)
}

type NegotiationRequest struct {
	Ver      byte
	NMethods byte
	Methods  []byte
}

func (request *NegotiationRequest) WriteTo(w io.Writer) error {
	_, err := w.Write(bytesCombine([]byte{request.Ver, request.NMethods}, request.Methods))
	return err
}

func (request *NegotiationRequest) ReadFrom(r io.Reader) error {
	version := []byte{0}
	if _, err := r.Read(version); err != nil {
		return fmt.Errorf("Failed to get version byte: %v ", err)
	}
	request.Ver = version[0]
	nMethods := []byte{0}
	if _, err := r.Read(nMethods); err != nil {
		return err
	}
	request.NMethods = nMethods[0]
	methods := make([]byte, int(request.NMethods))
	_, err := io.ReadAtLeast(r, methods, int(request.NMethods))
	request.Methods = methods
	return err
}

type NegotiationReply struct {
	Ver    byte
	Method byte
}

type UserPassNegotiationRequest struct {
	Ver    byte
	Ulen   byte
	Uname  []byte
	Plen   byte
	Passwd []byte
}

func (request *UserPassNegotiationRequest) ReadFrom(r io.Reader) error {
	header := []byte{0, 0}
	if _, err := io.ReadAtLeast(r, header, 2); err != nil {
		return err
	}
	request.Ver = header[0]
	request.Ulen = header[1]
	userLen := int(request.Ulen)
	user := make([]byte, userLen)
	if _, err := io.ReadAtLeast(r, user, userLen); err != nil {
		return err
	}
	request.Uname = user
	if _, err := r.Read(header[:1]); err != nil {
		return err
	}
	request.Plen = header[0]
	passLen := int(request.Plen)
	pass := make([]byte, passLen)
	if _, err := io.ReadAtLeast(r, pass, passLen); err != nil {
		return err
	}
	request.Passwd = pass
	return nil
}

func (request *UserPassNegotiationRequest) WriteTo(w io.Writer) error {
	return nil
}

type UserPassNegotiationReply struct {
	Ver    byte
	Status byte
}

type Socks5Request struct {
	Ver     byte
	Cmd     byte
	Rsv     byte // 0x00
	Atyp    byte
	DstAddr []byte
	DstPort []byte
}

type Socks5Reply struct {
	Ver  byte
	Rep  byte
	Rsv  byte // 0x00
	Atyp byte
	BndAddr []byte
	BndPort []byte
}

func (reply *Socks5Reply) WriteTo(w io.Writer) error {
	_, err := w.Write(bytesCombine([]byte{reply.Ver, reply.Rep, reply.Rsv, reply.Atyp}, reply.BndAddr, reply.BndPort))
	return err
}

func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

type Datagram struct {
	Rsv     []byte
	Frag    byte
	Atyp    byte
	DstAddr []byte
	DstPort []byte
	Data    []byte
}