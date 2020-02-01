package lsd

import (
	"crypto/rand"
	"fmt"
)

func (lsd *LSD) generateNewToken() (string, error){
	b := make([]byte, 8)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}
