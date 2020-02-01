package lsd

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func (lsd *LSD) encryptNotification(privateKey []byte, message string) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	private, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	encryptedMessage, err := rsa.SignPKCS1v15(rand.Reader, private, crypto.SHA256, []byte(message))
	// encryptedMessage, err := rsa.EncryptOAEP(
	// 	lsd.hash,
	// 	lsd.random,
	// 	pub,
	// 	[]byte(message),
	// 	[]byte(""),
	// )

	if err != nil {
		return nil, err
	}

	return encryptedMessage, nil
}
