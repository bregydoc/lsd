package lsd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"

	log "github.com/sirupsen/logrus"
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

func (lsd *LSD) privateKeyFromBytes(data []byte) (*rsa.PrivateKey, error) {
	pKeyBlock, _ := pem.Decode(data)
	log.Info("string(pKeyBlock.Bytes)", string(data))
	return x509.ParsePKCS1PrivateKey(pKeyBlock.Bytes)
}

func (lsd *LSD) publicKeyBytesFromPrivateKeyBytes(data []byte) ([]byte, error) {
	privateKey, err := lsd.privateKeyFromBytes(data)
	if err != nil {
		return nil, err
	}

	var pKeyBlock = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
	}

	publicKey := bytes.NewBuffer([]byte{})
	if err = pem.Encode(publicKey, pKeyBlock); err != nil {
		return nil, err
	}

	return publicKey.Bytes(), nil
}

// func (lsd *LSD) ifPublicKeyMatchWithPrivateKey(privateKey []byte, publicKey string)