package cipher

import (
	"github.com/archervanderwaal/JadeSocks/utils"
	"testing"
)

var (
	aesCryptoAlgorithm = &AesCryptoAlgorithm{
		Key: utils.NewKey(32, "Archer"),
	}
	originalData = "111"
	encryptResult string
)

func TestAesCryptoAlgorithm_Encrypt(t *testing.T) {
	resultByte, err := aesCryptoAlgorithm.Encrypt([]byte(originalData))
	if err != nil {
		t.Error()
	}
	encryptResult = string(resultByte)
}

func TestAesCryptoAlgorithm_Decrypt(t *testing.T) {
	resultByte, err := aesCryptoAlgorithm.Decrypt([]byte(encryptResult))
	if err != nil {
		t.Error()
	}
	if originalData != string(resultByte) {
		t.Error()
	}
}
