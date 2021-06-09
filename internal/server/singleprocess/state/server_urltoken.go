package state

import (
	bolt "go.etcd.io/bbolt"
)

var (
	serverURLToken = []byte("url_token")
)

// ServerURLTokenSet writes the server URL token.
func (s *State) ServerURLTokenSet(token string) error {
	return s.db.Update(func(dbTxn *bolt.Tx) error {
		return dbTxn.Bucket(serverConfigBucket).Put(serverURLToken, []byte(token))
	})
}

// ServerURLTokenGet gets the server URL token.
func (s *State) ServerURLTokenGet() (string, error) {
	var result string
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		result = string(dbTxn.Bucket(serverConfigBucket).Get(serverURLToken))
		return nil
	})

	return result, err
}
