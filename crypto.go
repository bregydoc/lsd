package lsd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/asn1"
	"encoding/pem"
)

func (lsd *LSD) generateNewKeyPair() ([]byte, []byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	pKey := key.Public()

	var pKeyBlock = &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	privateKey := bytes.NewBuffer([]byte{})
	publicKey := bytes.NewBuffer([]byte{})

	if err = pem.Encode(privateKey, pKeyBlock); err != nil {
		return nil, nil, err
	}

	asn1Bytes, err := asn1.Marshal(pKey)
	if err != nil {
		return nil, nil, err
	}

	var pemKey = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: asn1Bytes,
	}

	if err = pem.Encode(publicKey, pemKey); err != nil {
		return nil, nil, err
	}

	return publicKey.Bytes(), privateKey.Bytes(), nil
}
