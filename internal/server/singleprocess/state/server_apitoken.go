package state

import (
	"github.com/boltdb/bolt"
)

var (
	serverAPIToken = []byte("api_token")
)

// ServerAPITokenSet writes the server API token.
func (s *State) ServerAPITokenSet(token string) error {
	return s.db.Update(func(dbTxn *bolt.Tx) error {
		return dbTxn.Bucket(serverConfigBucket).Put(serverAPIToken, []byte(token))
	})
}

// ServerAPITokenGet gets the server API token.
func (s *State) ServerAPITokenGet() (string, error) {
	var result string
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		result = string(dbTxn.Bucket(serverConfigBucket).Get(serverAPIToken))
		return nil
	})

	return result, err
}
