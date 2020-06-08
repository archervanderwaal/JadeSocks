package socks5

import (
	"bytes"
	"fmt"
	"io"
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