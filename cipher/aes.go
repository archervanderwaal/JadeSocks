package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
)

type AesCryptoAlgorithm struct {
	Key string
}

func (aesAlgorithm *AesCryptoAlgorithm) Encrypt(bs []byte) ([]byte, error) {
	return aesEncrypt([]byte(bs), []byte(aesAlgorithm.Key))
}

func (aesAlgorithm *AesCryptoAlgorithm) Decrypt(bs []byte) ([]byte, error) {
	return aesDecrypt(bs, []byte(aesAlgorithm.Key))
}

func aesEncrypt(origData, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = pkcs7Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	result := make([]byte, len(origData))
	blockMode.CryptBlocks(result, origData)
	return result, nil
}

func aesDecrypt(encryptResult, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encryptResult))
	blockMode.CryptBlocks(origData, encryptResult)
	origData = pkcs7UnPadding(origData)
	return origData, nil
}

func pkcs7Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	paddingText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, paddingText...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}
