package lsd

import bolt "go.etcd.io/bbolt"

const sessionsBucket = "sessions"

func openBoltConnection(path string) (*bolt.DB, error) {
	return bolt.Open(path, 0666, nil)
}

func (lsd *LSD) registerNewClientSession(clientID, sessionID string) error {
	return lsd.db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte(sessionsBucket))
		if err != nil {
			return tx.Rollback()
		}

		if err =b.Put([]byte(clientID), []byte(sessionID)); err != nil {
			return tx.Rollback()
		}

		return tx.Commit()
	})
}

func (lsd *LSD) getClientSession(clientID string) (string, error) {
	sessionID := ""

	if err := lsd.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(sessionsBucket))
		data := b.Get([]byte(clientID))
		sessionID = string(data)
		return tx.Commit()
	}); err != nil {
		return "", err
	}

	return sessionID, nil
}