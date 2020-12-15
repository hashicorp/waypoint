package state

import (
	"github.com/boltdb/bolt"
)

var (
	serverAPIToken = []byte("token")
)

// ServerAPITokenSet writes the server ID.
func (s *State) ServerAPITokenSet(id string) error {
	return s.db.Update(func(dbTxn *bolt.Tx) error {
		return dbTxn.Bucket(serverConfigBucket).Put(serverAPIToken, []byte(id))
	})
}

// ServerAPITokenGet gets the server ID.
func (s *State) ServerAPITokenGet() (string, error) {
	var result string
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		result = string(dbTxn.Bucket(serverConfigBucket).Get(serverAPIToken))
		return nil
	})

	return result, err
}
