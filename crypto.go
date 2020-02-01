package lsd

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

const keySize = 2048

func (lsd *LSD) generateNewKeyPair() ([]byte, error) {
	key, err := rsa.GenerateKey(rand.Reader, keySize)
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

func (lsd *LSD) ifPublicKeyMatchWithUserID(userID string, publicKey []byte) error {
	private, err := lsd.getPrivateKey(userID)
	if err != nil {
		return err
	}

	public, err := lsd.publicKeyBytesFromPrivateKeyBytes(private)
	if err != nil {
		return err
	}

	if !bytes.Equal(public, publicKey) {
		return errors.New("invalid public key, it isn't equal to our saved key")
	}

	return nil
}

func (lsd *LSD) ifTokenMatchWithUserID(userID, token string) error {
	savedToken, err := lsd.getToken(userID)
	if err != nil {
		return err
	}

	if savedToken != token {
		return errors.New("invalid token, it not match with ours")
	}

	return nil
}
