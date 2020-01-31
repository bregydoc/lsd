package lsd

import (
	"encoding/base64"
	"errors"
	"strings"

	bolt "go.etcd.io/bbolt"
)

const sessionsBucket = "sessions"
const keypairsBucket = "keypairs"

func openBoltConnection(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0666, nil)
}

func (lsd *LSD) registerUserSession(userID, sessionID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		defer tx.Rollback()
		b, err := tx.CreateBucketIfNotExists([]byte(sessionsBucket))
		if err != nil {
			return err
		}

		if err =b.Put([]byte(userID), []byte(sessionID)); err != nil {
			return err
		}

		return tx.Commit()
	})
}

func (lsd *LSD) getUserSession(userID string) (string, error) {
	sessionID := ""

	if err := lsd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionsBucket))
		data := b.Get([]byte(userID))
		sessionID = string(data)
		return tx.Commit()
	}); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (lsd *LSD) clearUserSession(userID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		defer tx.Rollback()
		b := tx.Bucket([]byte(sessionsBucket))
		if err := b.Delete([]byte(userID)); err != nil {
			return err
		}
		return tx.Commit()
	})
}

func (lsd *LSD) saveKeyPair(userID string, publicKey []byte, privateKey []byte) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		defer tx.Rollback()
		b, err := tx.CreateBucketIfNotExists([]byte(keypairsBucket))
		if err != nil {
			return err
		}

		public64 := base64.StdEncoding.EncodeToString(publicKey)
		private64 := base64.StdEncoding.EncodeToString(privateKey)

		payload := public64 + ":" + private64

		if err := b.Put([]byte(userID), []byte(payload)); err != nil {
			return err
		}

		return tx.Commit()
	})
}


func (lsd *LSD) getKeyPair(userID string) ([]byte, []byte, error) {
	var publicKey, privateKey []byte
	if err := lsd.db.View(func(tx *bolt.Tx) error {
		defer tx.Rollback()
		b := tx.Bucket([]byte(keypairsBucket))
		payload := string(b.Get([]byte(userID)))
		cuts := strings.Split(payload, ":")
		if len(cuts) != 2 {
			return errors.New("invalid payload on store db.go:89")
		}
		public64 := cuts[0]
		private64 := cuts[1]

		var err error
		publicKey, err = base64.StdEncoding.DecodeString(public64)
		if err != nil {
			return err
		}
		privateKey, err = base64.StdEncoding.DecodeString(private64)
		if err != nil {
			return err
		}

		return tx.Commit()
	}); err != nil {
		return nil, nil, err
	}

	return publicKey, privateKey, nil
}


func (lsd *LSD) clearKeyPair(userID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		defer tx.Rollback()
		b := tx.Bucket([]byte(keypairsBucket))
		if err := b.Delete([]byte(userID)); err != nil {
			return err
		}
		return tx.Commit()
	})
}