package utils

import (
	"crypto/rand"
	"io"
)

func NewKey(length int, password string) string {
	//return randChar(length, []byte(password))
	return password
}

func randChar(length int, chars []byte) string {
	password := make([]byte, length)
	data := make([]byte, length+(length/4))
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, data); err != nil {
			panic(err)
		}
		for _, c := range data {
			if c >= maxrb {
				continue
			}
			password[i] = chars[c%clen]
			i++
			if i == length {
				return string(password)
			}
		}
	}
}