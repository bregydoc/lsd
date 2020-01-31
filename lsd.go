package lsd

import (
	bolt "go.etcd.io/bbolt"
	"gopkg.in/olahol/melody.v1"
)

type LSD struct {
	db *bolt.DB
	secure bool

	sessionsMap map[string]*melody.Session
}

func NewLSD(pathDB string, secure bool) (*LSD, error) {
	db, err := openBoltConnection(pathDB)
	if err != nil {
		return nil, err
	}

	lsd := &LSD{
		db:          db,
		secure:      secure,
		sessionsMap: map[string]*melody.Session{},
	}

	return lsd, nil
}

func (lsd *LSD) RunWSService(wsPort string) error {
	return lsd.launchClientWSServer(wsPort)
}