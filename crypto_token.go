package lsd

import (
	"errors"
)

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
