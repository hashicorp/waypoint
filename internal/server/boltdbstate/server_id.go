// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package boltdbstate

import (
	"context"
	bolt "go.etcd.io/bbolt"
)

var (
	serverIdKey = []byte("id")
)

// ServerIdSet writes the server ID.
func (s *State) ServerIdSet(ctx context.Context, id string) error {
	return s.db.Update(func(dbTxn *bolt.Tx) error {
		return dbTxn.Bucket(serverConfigBucket).Put(serverIdKey, []byte(id))
	})
}

// ServerIdGet gets the server ID.
func (s *State) ServerIdGet(ctx context.Context) (string, error) {
	var result string
	err := s.db.View(func(dbTxn *bolt.Tx) error {
		result = string(dbTxn.Bucket(serverConfigBucket).Get(serverIdKey))
		return nil
	})

	return result, err
}
