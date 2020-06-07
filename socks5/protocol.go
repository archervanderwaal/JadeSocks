package socks5

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"strconv"
)

const (
	IPV4Address   = uint8(1)
	DomainAddress = uint8(3)
	IPV6Address   = uint8(4)
)

var UnrecognizedAddrType = fmt.Errorf("Unrecognized address type ")

type RequestReader interface {
	Read(reader io.Reader) error
}

type ResponseWriter interface {
	Write(writer io.Writer) error
}

type UserPassAuthRequest struct {
	Ver    	byte
	Ulen   	byte
	Uname  	[]byte
	Plen   	byte
	Passwd 	[]byte
}

type NegotiationRequest struct {
	Ver      byte
	NMethods byte
	Methods  []byte
}

//type Request struct {
//	Version 		uint8
//	Command 		uint8
//	RemoteAddr 		*AddrSpec
//	DestAddr 		*AddrSpec
//	reader      	io.Reader
//}

type AddrSpec struct {
	Domain   string
	AddrType byte
	IP       net.IP
	Port     uint16
}

func (a AddrSpec) Address() string {
	if 0 != len(a.IP) {
		return net.JoinHostPort(a.IP.String(), strconv.Itoa(int(a.Port)))
	}
	return net.JoinHostPort(a.Domain, strconv.Itoa(int(a.Port)))
}

type Response struct {
	Ver  byte
	Rep  byte
	Rsv  byte // 0x00
	Atyp byte
	BndAddr []byte
	BndPort []byte
}

func (req *UserPassAuthRequest) Read(reader io.Reader) error {
	header := []byte{0, 0}
	if _, err := io.ReadAtLeast(reader, header, 2); err != nil {
		return err
	}
	req.Ver = header[0]
	req.Ulen = header[1]
	userLen := int(req.Ulen)
	user := make([]byte, userLen)
	if _, err := io.ReadAtLeast(reader, user, userLen); err != nil {
		return err
	}
	req.Uname = user
	if _, err := reader.Read(header[:1]); err != nil {
		return err
	}
	req.Plen = header[0]
	passLen := int(req.Plen)
	pass := make([]byte, passLen)
	if _, err := io.ReadAtLeast(reader, pass, passLen); err != nil {
		return err
	}
	req.Passwd = pass
	return nil
}

func (req *NegotiationRequest) Read(r io.Reader) error {
	version := []byte{0}
	if _, err := r.Read(version); err != nil {
		return fmt.Errorf("Failed to get version byte: %v ", err)
	}
	req.Ver = version[0]
	nMethods := []byte{0}
	if _, err := r.Read(nMethods); err != nil {
		return err
	}
	req.NMethods = nMethods[0]
	methods := make([]byte, int(req.NMethods))
	_, err := io.ReadAtLeast(r, methods, int(req.NMethods))
	req.Methods = methods
	return err
}

func (req *NegotiationRequest) Write(w io.Writer) error {
	_, err := w.Write(bytesCombine([]byte{req.Ver, req.NMethods}, req.Methods))
	return err
}

func bytesCombine(pBytes ...[]byte) []byte {
	return bytes.Join(pBytes, []byte(""))
}

//func (req *Request) Read(reader io.Reader) error {
//	header := []byte{0, 0, 0}
//	if _, err := io.ReadAtLeast(reader, header, 3); err != nil {
//		return fmt.Errorf("Failed to get command version: %v ", err)
//	}
//	if header[0] != Socks5Version {
//		return fmt.Errorf("Unsupported command version: %v ", header[0])
//	}
//	dest, err := parseAddrSpec(reader)
//	if err != nil {
//		return err
//	}
//	req.Version = Socks5Version
//	req.Command = header[1]
//	req.DestAddr = dest
//	req.reader = reader
//	return nil
//}


//func parseAddrSpec(r io.Reader) (*AddrSpec, error) {
//	addrSpec := &AddrSpec{}
//
//	addrType := []byte{0}
//	if _, err := r.Read(addrType); err != nil {
//		return nil, err
//	}
//
//	switch addrType[0] {
//	case IPV4Address:
//		addr := make([]byte, 4)
//		if _, err := io.ReadAtLeast(r, addr, len(addr)); err != nil {
//			return nil, err
//		}
//		addrSpec.IP = addr
//		addrSpec.AddrType = IPV4Address
//
//	case IPV6Address:
//		addr := make([]byte, 16)
//		if _, err := io.ReadAtLeast(r, addr, len(addr)); err != nil {
//			return nil, err
//		}
//		addrSpec.IP = addr
//		addrSpec.AddrType = IPV6Address
//
//	case DomainAddress:
//		if _, err := r.Read(addrType); err != nil {
//			return nil, err
//		}
//		addrLen := int(addrType[0])
//		domain := make([]byte, addrLen)
//		if _, err := io.ReadAtLeast(r, domain, addrLen); err != nil {
//			return nil, err
//		}
//		addrSpec.Domain = string(domain)
//		addrSpec.AddrType = DomainAddress
//
//	default:
//		return nil, UnrecognizedAddrType
//	}
//	port := []byte{0, 0}
//	if _, err := io.ReadAtLeast(r, port, 2); err != nil {
//		return nil, err
//	}
//	addrSpec.Port = uint16((int(port[0]) << 8) | int(port[1]))
//
//	return addrSpec, nil
//}

//func (rep *Response) Write(writer io.Writer) error {
//	if _, err := writer.Write(bytesCombine([]byte{rep.Ver, rep.Rep, rep.Rsv, rep.Atyp}, rep.BndAddr, rep.BndPort)); err != nil {
//		return err
//	}
//	return nil
//}

//func (rep *Response) Generate(resp uint8, addr *AddrSpec) error {
//	rep.Ver = Socks5Version
//	rep.Rep = resp
//	rep.Rsv = uint8(0)
//	rep.Atyp = 0
//	if addr == nil {
//		rep.Atyp = IPV4Address
//		rep.BndAddr = []byte{0, 0, 0, 0}
//		rep.BndPort = []byte{0, 0}
//		return nil
//	}
//	rep.Atyp = addr.AddrType
//	rep.BndPort = []byte{byte(addr.Port >> 8), byte(addr.Port & 0xff)}
//	switch addr.AddrType {
//	case IPV4Address:
//		rep.BndAddr = addr.IP.To4()
//	case DomainAddress:
//		rep.BndAddr = bytesCombine([]byte{byte(len(addr.Domain))}, []byte(addr.Domain))
//	case IPV6Address:
//		rep.BndAddr = addr.IP.To16()
//	default:
//		return fmt.Errorf("Failed to format address: %v ", addr)
//	}
//	return nil
//}