package lsd

import bolt "go.etcd.io/bbolt"

type LSD struct {
	db *bolt.DB
}
