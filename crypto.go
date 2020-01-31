package lsd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
)

func (lsd *LSD) generateNewKeyPair() ([]byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	var pKeyBlock = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	privateKey := bytes.NewBuffer([]byte{})


	if err = pem.Encode(privateKey, pKeyBlock); err != nil {
		return nil, err
	}


	return privateKey.Bytes(), nil
}
