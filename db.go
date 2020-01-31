package lsd

import (
	"encoding/base64"

	bolt "go.etcd.io/bbolt"
)

const sessionsBucket = "sessions"
const keypairsBucket = "keypairs"

func openBoltConnection(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0600, nil)
}

func (lsd *LSD) registerUserSession(userID, sessionID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		// defer tx.Rollback()
		b, err := tx.CreateBucketIfNotExists([]byte(sessionsBucket))
		if err != nil {
			return err
		}

		if err = b.Put([]byte(userID), []byte(sessionID)); err != nil {
			return err
		}

		// return tx.Commit()
		return nil
	})
}

func (lsd *LSD) getUserSession(userID string) (string, error) {
	sessionID := ""

	if err := lsd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionsBucket))
		data := b.Get([]byte(userID))
		sessionID = string(data)
		// return tx.Commit()
		return nil
	}); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (lsd *LSD) clearUserSession(userID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		// defer tx.Rollback()
		b := tx.Bucket([]byte(sessionsBucket))
		if err := b.Delete([]byte(userID)); err != nil {
			return err
		}
		// return tx.Commit()
		return nil
	})
}

func (lsd *LSD) savePrivateKey(userID string, privateKey []byte) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		// defer tx.Rollback()
		b, err := tx.CreateBucketIfNotExists([]byte(keypairsBucket))
		if err != nil {
			return err
		}

		private64 := base64.StdEncoding.EncodeToString(privateKey)

		if err := b.Put([]byte(userID), []byte(private64)); err != nil {
			return err
		}

		// return tx.Commit()
		return nil
	})
}

func (lsd *LSD) getPrivateKey(userID string) ([]byte, error) {
	var privateKey []byte
	if err := lsd.db.View(func(tx *bolt.Tx) error {
		// defer tx.Rollback()
		b := tx.Bucket([]byte(keypairsBucket))
		payload := string(b.Get([]byte(userID)))
		var err error
		privateKey, err = base64.StdEncoding.DecodeString(payload)
		if err != nil {
			return err
		}

		// return tx.Commit()
		return nil
	}); err != nil {
		return nil, err
	}

	return privateKey, nil
}

func (lsd *LSD) clearKeyPair(userID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		// defer tx.Rollback()
		b := tx.Bucket([]byte(keypairsBucket))
		if err := b.Delete([]byte(userID)); err != nil {
			return err
		}
		// return tx.Commit()
		return nil
	})
}
