package lsd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
)

// func (lsd *LSD) encryptNotification(privateKey []byte, message string) ([]byte, error) {
// 	block, _ := pem.Decode(privateKey)
// 	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	encryptedMessage, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, []byte(message))
// 	// encryptedMessage, err := rsa.EncryptOAEP(
// 	// 	lsd.hash,
// 	// 	lsd.random,
// 	// 	pub,
// 	// 	[]byte(message),
// 	// 	[]byte(""),
// 	// )
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return encryptedMessage, nil
// }

func (lsd *LSD) encryptNotification(token string, message []byte) ([]byte, error) {
	key := sha256.Sum256([]byte(token))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, message, nil), nil
}
