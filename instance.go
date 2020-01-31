package lsd

import bolt "go.etcd.io/bbolt"

func ConnectBoltDB(path string) (*bolt.DB, error) {
	return openBoltConnection(path)
}