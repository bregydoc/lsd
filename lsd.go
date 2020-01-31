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
