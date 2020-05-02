package cipher

type Cipher struct {
	Algorithm CryptoAlgorithm
}

type CryptoAlgorithm interface {
	Encrypt(bs []byte) ([]byte, error)
	Decrypt(bs []byte) ([]byte, error)
}

func (cipher *Cipher) EncryptData(originalData []byte) ([]byte, error) {
	return cipher.Algorithm.Encrypt(originalData)
}

func (cipher *Cipher) DecryptData(encryptResult []byte) ([]byte, error) {
	return cipher.Algorithm.Decrypt(encryptResult)
}