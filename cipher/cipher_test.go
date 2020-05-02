package cipher

import (
	"fmt"
	"github.com/archervanderwaal/JadeSocks/utils"
	"testing"
)

const (
	key = "Archer"
	keyLength = 16
)

var (
	encryptDataResult string
	decryptDataResult string
	cipherTest        = &Cipher{
		Algorithm: &AesCryptoAlgorithm{
			Key: utils.NewKey(keyLength, key),
		},
	}
	originalText = "Archer"
)

func TestCipher_EncryptData(t *testing.T) {
	resultByte, err := cipherTest.EncryptData([]byte(originalText))
	if err != nil {
		t.Error()
	}
	encryptResult = string(resultByte)
}

func TestCipher_DecryptData(t *testing.T) {
	resultByte, err := cipherTest.DecryptData([]byte(encryptDataResult))
	if err != nil {
		t.Error()
	}
	decryptDataResult = string(resultByte)
	if decryptDataResult != originalText {
		t.Error()
	}
}

func TestCipher_DecryptData2(t *testing.T) {
	resultByte, _ := cipherTest.EncryptData([]byte(originalText))
	resultString := string(resultByte)
	fmt.Println("加密结果: " + resultString)

	cipher2 := &Cipher{
		Algorithm: &AesCryptoAlgorithm{
			Key: utils.NewKey(16, "Archer"),
		},
	}

	decryptResultByte, err := cipher2.DecryptData(resultByte)
	fmt.Println(string(decryptResultByte))
	if err != nil {
		fmt.Println("error : " + err.Error())
	}
}

